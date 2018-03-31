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
	ReadChanSize      int
	WriteBufferSize   int
	WriteChanSize     int
	EnableCompression bool
}

// Websocket is the interface of channel based websocket.
type Websocket interface {

	// Reads returns a channel where the incoming messages are published.
	// If nothing pumps the Read() while it is full, new messages will be
	// discarded.
	//
	// You can configure the size of the read chan in Config.
	// The default is 64 messages.
	Read() chan []byte

	// Write write the given []byte in to the websocket.
	// If the other side of the websocket cannot get all messages
	// while the internal write channel is full, new messages will
	// be discarded.
	//
	// You can configure the size of the write chan in Config.
	// The default is 64 messages.
	Write([]byte)

	// Done returns a channel that will return when the connection
	// if closed.
	//
	// The content will be nil for clean disconnection or
	// the error that caused the disconnection. If nothing pumps the
	// Done() channel, the message will be discarded.
	//
	// If nothing pumps the Done() chan, the message will be discarded.
	Done() chan error

	// Close closes the webbsocket.
	//
	// Closing the websocket a second time has no effect.
	// A closed Websocket cannot be reused.
	Close() error
}

type ws struct {
	conn      *websocket.Conn
	readChan  chan []byte
	writeChan chan []byte
	doneChan  chan error
	cancel    context.CancelFunc
	config    Config
}

// Connect connects to the url and returns a Websocket.
func Connect(ctx context.Context, url string, config Config) (Websocket, *http.Response, error) {

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

	s, err := Accept(ctx, conn, config)

	return s, resp, err
}

// Accept handles an already connect *websocket.Conn and returns a Websocket.
func Accept(ctx context.Context, conn *websocket.Conn, config Config) (Websocket, error) {

	if config.PongWait == 0 {
		config.PongWait = 30 * time.Second
	}
	if config.WriteWait == 0 {
		config.WriteWait = 10 * time.Second
	}
	if config.PingPeriod == 0 {
		config.PingPeriod = 15 * time.Second
	}
	if config.WriteChanSize == 0 {
		config.WriteChanSize = 64
	}
	if config.ReadChanSize == 0 {
		config.ReadChanSize = 64
	}

	if err := conn.SetReadDeadline(time.Now().Add(config.PongWait)); err != nil {
		return nil, err
	}

	subCtx, cancel := context.WithCancel(ctx)

	s := &ws{
		conn:      conn,
		readChan:  make(chan []byte, config.ReadChanSize),
		writeChan: make(chan []byte, config.WriteChanSize),
		doneChan:  make(chan error, 1),
		cancel:    cancel,
		config:    config,
	}

	s.conn.SetCloseHandler(func(code int, text string) error {
		return s.Close()
	})

	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(s.config.PongWait))
	})

	go s.readPump(subCtx)
	go s.writePump(subCtx)

	return s, nil
}

// Write is part of the the Websocket interface implementation.
func (s *ws) Write(data []byte) {

	select {
	case s.writeChan <- data:
	default:
	}
}

// Read is part of the the Websocket interface implementation.
func (s *ws) Read() chan []byte {

	return s.readChan
}

// Done is part of the the Websocket interface implementation.
func (s *ws) Done() chan error {

	return s.doneChan
}

// Close is part of the the Websocket interface implementation.
func (s *ws) Close() error {

	s.cancel()
	return nil
}

func (s *ws) readPump(ctx context.Context) {

	var err error
	var msg []byte
	var msgType int

	for {
		if msgType, msg, err = s.conn.ReadMessage(); err != nil {
			s.done(err)
			return
		}

		switch msgType {

		case websocket.TextMessage, websocket.BinaryMessage:
			select {
			case s.readChan <- msg:
			default:
			}

		case websocket.CloseMessage:
			return
		}
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
				s.done(err)
				return
			}

		case <-ticker.C:

			s.conn.SetWriteDeadline(time.Now().Add(s.config.WriteWait)) // nolint: errcheck
			if err = s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.done(err)
				return
			}

		case <-ctx.Done():

			s.done(
				s.conn.WriteControl(
					websocket.CloseMessage,
					[]byte{},
					time.Now().Add(1*time.Second),
				),
			)

			_ = s.conn.Close()

			return
		}
	}
}

func (s *ws) done(err error) {

	select {
	case s.doneChan <- err:
	default:
	}
}
