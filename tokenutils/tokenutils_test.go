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

package tokenutils

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenUtils_Snip(t *testing.T) {

	Convey("Given have a token and and error containing the token", t, func() {

		token := "token"
		err := errors.New("your token is token")

		Convey("When I call Snip", func() {

			e := Snip(err, token)

			Convey("Then err should have the reference to token snipped", func() {
				So(e.Error(), ShouldEqual, "your [snip] is [snip]")
			})
		})
	})

	Convey("Given have a token and and error that doesn't contain the token", t, func() {

		token := "token"
		err := errors.New("your secret is secret")

		Convey("When I call Snip", func() {

			e := Snip(err, token)

			Convey("Then err should have the reference to token snipped", func() {
				So(e.Error(), ShouldEqual, "your secret is secret")
			})
		})
	})

	Convey("Given I have a token and a nil error", t, func() {

		token := "token"

		Convey("When I call Snip", func() {

			e := Snip(nil, token)

			Convey("Then err should be nil", func() {
				So(e, ShouldBeNil)
			})
		})
	})
}

func TestTokenUtils_UnsecureClaimsMap(t *testing.T) {

	Convey("Given I have a valid token", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YTZhNTUxMTdkZGYxZjIxMmY4ZWIwY2UiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTIwNjQ5MTAyLCJpYXQiOjE1MTgwNTcxMDIsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then claims should be correct", func() {
				So(claims, ShouldNotBeNil)
				So(claims["data"].(map[string]interface{})["realm"].(string), ShouldEqual, "vince")
				So(claims["sub"].(string), ShouldEqual, "apomux")
			})
		})
	})

	Convey("Given I have a token an invalid token", t, func() {

		token := "not good"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: not enough segments")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token an empty token", t, func() {

		token := ""

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: empty")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token a token with invalid base64", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.not-base64.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: illegal base64 data at input byte 3")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token a token with invalid json data", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicm9rZW46ICJqc29u.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			alg, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: invalid character 'j' after object key")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})
}

func TestJWTUtils_SigAlg(t *testing.T) {

	Convey("Given I have a valid token", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YTZhNTUxMTdkZGYxZjIxMmY4ZWIwY2UiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTIwNjQ5MTAyLCJpYXQiOjE1MTgwNTcxMDIsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I SigAlg", func() {

			alg, err := SigAlg(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then alg should be correct", func() {
				So(alg, ShouldEqual, "HS256")
			})
		})
	})

	Convey("Given I have a token an invalid token", t, func() {

		token := "not good"

		Convey("When I SigAlg", func() {

			alg, err := SigAlg(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: not enough segments")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})

	Convey("Given I have an empty token", t, func() {

		token := ""

		Convey("When I SigAlg", func() {

			alg, err := SigAlg(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: empty")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})

	Convey("Given I have a token a token with invalid base64", t, func() {

		token := "not-base-64.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YTZhNTUxMTdkZGYxZjIxMmY4ZWIwY2UiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTIwNjQ5MTAyLCJpYXQiOjE1MTgwNTcxMDIsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I SigAlg", func() {

			alg, err := SigAlg(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: illegal base64 data at input byte 3")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})

	Convey("Given I have a token a token with invalid json data", t, func() {

		token := "eyJicm9rZW46ICJqc29u.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YTZhNTUxMTdkZGYxZjIxMmY4ZWIwY2UiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTIwNjQ5MTAyLCJpYXQiOjE1MTgwNTcxMDIsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I SigAlg", func() {
			alg, err := SigAlg(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: invalid character 'j' after object key")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})
}

func TestExtractQuota(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"valid token with valid quota",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwicXVvdGEiOjMsImRhdGEiOnsiYWNjb3VudCI6ImFwb211eCIsImVtYWlsIjoiYWRtaW5AYXBvbXV4LmNvbSIsImlkIjoiNWI0OTBlY2M3ZGRmMWY3NWFiODRlN2IxIiwib3JnYW5pemF0aW9uIjoiYXBvbXV4IiwicmVhbG0iOiJ2aW5jZSJ9LCJhdWQiOiJhcG9yZXRvLmNvbSIsImV4cCI6MTU0NzY4MzIxOCwiaWF0IjoxNTQ3NTkzMjE4LCJpc3MiOiJtaWRnYXJkLmFwb211eC5jb20iLCJzdWIiOiJhcG9tdXgifQ.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA`,
			},
			3,
			false,
		},
		{
			"valid token with no quota",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTQ3NjgzMjE4LCJpYXQiOjE1NDc1OTMyMTgsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA`,
			},
			0,
			false,
		},
		{
			"valid token with invalid quota",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwicXVvdGEiOiIzIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTQ3NjgzMjE4LCJpYXQiOjE1NDc1OTMyMTgsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA`,
			},
			0,
			true,
		},
		{
			"invalid token",
			args{
				`eyJhbGciOiJFUzI1NiIsInWFsbSI6IlNlIiwicXVvdGEiOiIzIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTQ3NjgzMjE4LCJpYXQiOjE1NDc1OTMyMTgsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA`,
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractQuota(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractQuota() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractQuota() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRestrictions(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name      string
		args      args
		wantNs    string
		wantPerms []string
		wantNets  []string
		wantErr   bool
	}{
		{
			"valid token with no restrictions",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJleHAiOjE1OTAwMTUzNTIsImlhdCI6MTU4OTkyNTM1MiwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.agqImtfkfjJugJH59XfQwkasIayYtvG6tz3p84jMulfbgwZzTLzgfRDLNIcfnfqfUix_702BUJxvdlsaSsgeUg`,
			},
			"",
			nil,
			nil,
			false,
		},

		{
			"valid token with namespace restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmFtZXNwYWNlIjoiL3Rlc3QifSwiZXhwIjoxNTkwMDE1MzY2LCJpYXQiOjE1ODk5MjUzNjYsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.3xYCsLm6FuIUnYKE3GPYqrB9-GhUTDUETSuk6hWy8DqnK8DOj95AXWbJsKN7yUSpBp8CuPmbrPRAljnOSROIJg`,
			},
			"/test",
			nil,
			nil,
			false,
		},
		{
			"valid token with invalid namespace restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmFtZXNwYWNlIjoxfSwiZXhwIjoxNTkwMDE1MzY2LCJpYXQiOjE1ODk5MjUzNjYsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.dIsnGMSEy961FqXgJH-TBVw8_9VrzH_j4xcQJG4JY0--ekwNuMpLr0CyOJFj_XFuVsY-ZS8Lwj5yJCYHv7TS8Q`,
			},
			"",
			nil,
			nil,
			true,
		},

		{
			"valid token with permissions restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0Il19LCJleHAiOjE1OTAwMTUzOTIsImlhdCI6MTU4OTkyNTM5MiwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.ZgFBGxiU5F_ImSZMy8BwwyOv-1dnDw1z4zY5spUqxSKAzCrtrr_G0FomsL-w0CZOKWDv2Hs99FQPxIRRjfSwSQ`,
			},
			"",
			[]string{"@auth:role=test"},
			nil,
			false,
		},
		{
			"valid token with invalid permissions restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOiJAYXV0aDpyb2xlPXRlc3QifSwiZXhwIjoxNTkwMDE1MzkyLCJpYXQiOjE1ODk5MjUzOTIsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.zBBcocIRkBUzcHGUKQk6ywitqMFe9sQsyMRp2raVyuBW0jirIR-yUybXGtOuw8YucCq8uIvC9CPHXUPww6kGtA`,
			},
			"",
			nil,
			nil,
			true,
		},
		{
			"valid token with invalid permissions restriction item type",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0IiwgMV19LCJleHAiOjE1OTAwMTUzOTIsImlhdCI6MTU4OTkyNTM5MiwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.zBBcocIRkBUzcHGUKQk6ywitqMFe9sQsyMRp2raVyuBW0jirIR-yUybXGtOuw8YucCq8uIvC9CPHXUPww6kGtA`,
			},
			"",
			nil,
			nil,
			true,
		},

		{
			"valid token with networks restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmV0d29ya3MiOlsiMTI3LjAuMC4xLzMyIl19LCJleHAiOjE1OTAwNDMyMDUsImlhdCI6MTU4OTk1MzIwNSwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.aO-D-1e7XX18ZGGi2FoKl-lfpc8xBzShrPVHZzLHqXd22ojZEWxKOVaIh8ZVyjs8KWBu-CmCKBCi3eg875j_1A`,
			},
			"",
			nil,
			[]string{"127.0.0.1/32"},
			false,
		},
		{
			"valid token with invalid networks restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmV0d29ya3MiOiIxMjcuMC4wLjEvMzIifSwiZXhwIjoxNTkwMDQzMjA1LCJpYXQiOjE1ODk5NTMyMDUsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.dIsnGMSEy961FqXgJH-TBVw8_9VrzH_j4xcQJG4JY0--ekwNuMpLr0CyOJFj_XFuVsY-ZS8Lwj5yJCYHv7TS8Q`,
			},
			"",
			nil,
			nil,
			true,
		},
		{
			"valid token with invalid networks restriction item",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmV0d29ya3MiOlsxXX0sImV4cCI6MTU5MDA0MzIwNSwiaWF0IjoxNTg5OTUzMjA1LCJpc3MiOiJodHRwczovL2xvY2FsaG9zdDo0NDQzIiwic3ViIjoiYXBvbXV4In0.dIsnGMSEy961FqXgJH-TBVw8_9VrzH_j4xcQJG4JY0--ekwNuMpLr0CyOJFj_XFuVsY-ZS8Lwj5yJCYHv7TS8Q`,
			},
			"",
			nil,
			nil,
			true,
		},

		{
			"valid token with namespace, identities and networks restrictions",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0Il0sIm5hbWVzcGFjZSI6Ii9hcG9tdXgvY2hpbGQiLCJuZXR3b3JrcyI6WyIxMjcuMC4wLjEvMzIiXX0sImV4cCI6MTU5MDA0Mjk5OCwiaWF0IjoxNTg5OTUyOTk4LCJpc3MiOiJodHRwczovL2xvY2FsaG9zdDo0NDQzIiwic3ViIjoiYXBvbXV4In0.JQcljRWeraT2Ma2u9DpeO0ub0SLNj5jDjKMVppibsm17YH6CyNKO5pyf-Kg6SldxuJau1nf0W_V7K3sQFmqj0g`,
			},
			"/apomux/child",
			[]string{"@auth:role=test"},
			[]string{"127.0.0.1/32"},
			false,
		},
		{
			"invalid token",
			args{
				`eyJhbGciOiJFUzI1NiIsInWFsbSI6IlNlIiwicXVvdGEiOiIzIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTQ3NjgzMjE4LCJpYXQiOjE1NDc1OTMyMTgsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA`,
			},
			"",
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNs, gotPerms, gotNetworks, err := ExtractRestrictions(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractRestrictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNs != tt.wantNs {
				t.Errorf("ExtractRestrictions() gotAuthorizedNS = %v, want %v", gotNs, tt.wantNs)
			}
			if !reflect.DeepEqual(gotPerms, tt.wantPerms) {
				t.Errorf("ExtractRestrictions() gotAuthorizedPerms = %v, want %v", gotPerms, tt.wantPerms)
			}
			if !reflect.DeepEqual(gotNetworks, tt.wantNets) {
				t.Errorf("ExtractRestrictions() gotNetworks = %v, want %v", gotNetworks, tt.wantNets)
			}
		})
	}
}
