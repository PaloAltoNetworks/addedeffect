// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiutils

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type testType int

const (
	failCode testType = iota
	badData
	goodData
)

type testData struct {
	testType testType
	goodData []byte
	badData  []byte
}

var data *testData

func makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if data == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch r.URL.Path {
		case "/_meta/time":
		case "/_meta/config":
		case "/_meta/ca":
		case "/_meta/jwtcert":
		case "/_meta/googleclientid":
		case "/_meta/manifest":
		case "/_meta/versions":
		case "/_meta/model":
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch data.testType {
		case failCode:
			w.WriteHeader(http.StatusInternalServerError)
		case badData:
			w.Write(data.badData)
		case goodData:
			w.Write(data.goodData)
		}
	}))
}

// doneCtx returns immediately done Context, so we can test failures without retry
func doneCtx(ctx context.Context) context.Context {
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	return cancelCtx
}

func TestGetTime(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     time.Time
		testData testData
	}{
		{
			name: "meta-time-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: time.Unix(1617114591, 0),
			testData: testData{
				testType: goodData,
				goodData: []byte("1617114591"),
			},
		},
		{
			name: "meta-time-baddata",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: time.Time{},
			testData: testData{
				testType: badData,
				badData:  []byte("abcdef"),
			},
		},
		{
			name: "meta-time-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: time.Time{},
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetTime(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetTime() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetTime() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetTime() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetConfig(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     map[string]string
		testData testData
	}{
		{
			name: "meta-config-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: map[string]string{"item": "value"},
			testData: testData{
				testType: goodData,
				goodData: []byte(`{"item": "value"}`),
			},
		},
		{
			name: "meta-config-baddata",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: badData,
				badData:  []byte("abcdef"),
			},
		},
		{
			name: "meta-config-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetConfig(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetConfig() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetConfig() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetConfig() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetServiceVersions(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     map[string]Version
		testData testData
	}{
		{
			name: "meta-versions-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: map[string]Version{"item": {
				Version: "1.2.3",
				Sha:     "f00",
			}},
			testData: testData{
				testType: goodData,
				goodData: []byte(`{"item": {"Version": "1.2.3", "Sha": "f00"}}`),
			},
		},
		{
			name: "meta-versions-baddata",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: badData,
				badData:  []byte("abcdef"),
			},
		},
		{
			name: "meta-versions-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetServiceVersions(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetServiceVersions() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetServiceVersions() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetServiceVersions() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetModelVersion(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     *Version
		testData testData
	}{
		{
			name: "meta-model-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: &Version{
				Version: "3.2.1",
				Sha:     "f00",
			},
			testData: testData{
				testType: goodData,
				goodData: []byte(`{"Version": "3.2.1", "Sha": "f00"}`),
			},
		},
		{
			name: "meta-model-baddata",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: badData,
				badData:  []byte("abcdef"),
			},
		},
		{
			name: "meta-model-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetModelVersion(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetModelVersion() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetModelVersion() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetModelVersion() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetPublicCA(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		testData testData
	}{
		{
			name: "meta-ca-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: []byte("foo"),
			testData: testData{
				testType: goodData,
				goodData: []byte("foo"),
			},
		},
		{
			name: "meta-ca-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetPublicCA(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetPublicCA() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetPublicCA() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetPublicCA() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetJWTCert(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		testData testData
	}{
		{
			name: "meta-jwt-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: []byte("foo"),
			testData: testData{
				testType: goodData,
				goodData: []byte("foo"),
			},
		},
		{
			name: "meta-jwt-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetJWTCert(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetJWTCert() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetJWTCert() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetJWTCert() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetManifestURL(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		testData testData
	}{
		{
			name: "meta-manifest-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: []byte("foo"),
			testData: testData{
				testType: goodData,
				goodData: []byte("foo"),
			},
		},
		{
			name: "meta-manifest-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetManifestURL(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetManifestURL() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetManifestURL() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetManifestURL() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}

func TestGetGoogleOAuthClientID(t *testing.T) {

	testServer := makeTestServer()
	defer testServer.Close()

	type args struct {
		ctx       context.Context
		api       string
		tlsConfig *tls.Config
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		testData testData
	}{
		{
			name: "meta-google-good",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: []byte("foo"),
			testData: testData{
				testType: goodData,
				goodData: []byte("foo"),
			},
		},
		{
			name: "meta-google-badstatus",
			args: args{
				ctx:       context.Background(),
				api:       testServer.URL,
				tlsConfig: nil,
			},
			want: nil,
			testData: testData{
				testType: failCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data = &tt.testData
			got, err := GetGoogleOAuthClientID(doneCtx(tt.args.ctx), tt.args.api, tt.args.tlsConfig)
			switch data.testType {
			case goodData:
				if err != nil {
					t.Errorf("GetGoogleOAuthClientID() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetGoogleOAuthClientID() = %v, want %v", got, tt.want)
				}
			default:
				if err == nil {
					t.Errorf("GetGoogleOAuthClientID() expected failure on %v, got %v", data.testType, got)
				}
			}
		})
	}
}
