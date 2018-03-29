package wsc

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Config contains configuration for the webbsocket.
type Config struct {
	WriteWait         time.Duration
	PongWait          time.Duration
	PingPeriod        time.Duration
	TLSConfig         *tls.Config
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
}

// Websocket is the interface of channel based websocket.
type Websocket interface {
	Read() chan []byte
	Write([]byte)
	Done() chan error
	Close() error
}

type ws struct {
	conn      *websocket.Conn
	tlsConfig *tls.Config
	readChan  chan []byte
	writeChan chan []byte
	doneChan  chan error
	isClosed  bool
	cancel    context.CancelFunc
	config    Config
}

// NewWebsocket returns a new connected Websocket.
func NewWebsocket(ctx context.Context, url string, config Config) (Websocket, *http.Response, error) {

	dialer := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   config.TLSConfig,
		ReadBufferSize:    config.ReadBufferSize,
		WriteBufferSize:   config.ReadBufferSize,
		EnableCompression: config.EnableCompression,
	}

	conn, resp, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, resp, err
	}

	if config.PongWait == 0 {
		config.PongWait = 30 * time.Second
	}
	if config.WriteWait == 0 {
		config.WriteWait = 10 * time.Second
	}
	if config.PingPeriod == 0 {
		config.PingPeriod = 15 * time.Second
	}

	if err = conn.SetReadDeadline(time.Now().Add(config.PongWait)); err != nil {
		return nil, resp, err
	}

	subCtx, cancel := context.WithCancel(ctx)

	s := &ws{
		conn:      conn,
		readChan:  make(chan []byte),
		writeChan: make(chan []byte),
		doneChan:  make(chan error, 16),
		cancel:    cancel,
		config:    config,
	}

	s.conn.SetCloseHandler(s.closeHandler)
	s.conn.SetPongHandler(s.pongHandler)

	go s.readPump(subCtx)
	go s.writePump(subCtx)

	return s, resp, nil
}

func (s *ws) Write(data []byte) { s.writeChan <- data }
func (s *ws) Read() chan []byte { return s.readChan }
func (s *ws) Done() chan error  { return s.doneChan }
func (s *ws) Close() error {

	if s.isClosed {
		return nil
	}

	s.cancel()

	s.isClosed = true
	if err := s.conn.Close(); err != nil {
		s.publishDoneMessage(err)
	}

	s.publishDoneMessage(nil)

	return nil
}

func (s *ws) readPump(ctx context.Context) {

	var err error
	var message []byte

	for {
		if _, message, err = s.conn.ReadMessage(); err != nil {
			s.publishDoneMessage(err)
			return
		}

		s.readChan <- message
	}
}

func (s *ws) writePump(ctx context.Context) {

	var err error

	ticker := time.NewTicker(s.config.PingPeriod)
	defer ticker.Stop()

	for {
		select {

		case message := <-s.writeChan:

			s.conn.SetWriteDeadline(time.Now().Add(s.config.WriteWait)) // nolint: errcheck
			if err = s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				s.publishDoneMessage(err)
				return
			}

		case <-ticker.C:

			s.conn.SetWriteDeadline(time.Now().Add(s.config.WriteWait)) // nolint: errcheck
			if err = s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.publishDoneMessage(nil)
				return
			}

		case <-ctx.Done():

			s.publishDoneMessage(s.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			return
		}
	}
}

func (s *ws) publishDoneMessage(err error) {
	select {
	case s.doneChan <- err:
	default:
	}
}

func (s *ws) pongHandler(string) error {

	return s.conn.SetReadDeadline(time.Now().Add(s.config.PongWait))
}

func (s *ws) closeHandler(code int, text string) error {

	return s.Close()
}
