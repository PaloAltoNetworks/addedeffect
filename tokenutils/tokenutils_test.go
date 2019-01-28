package tokenutils

import (
	"errors"
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

func TestExtractAPIAndNamespace(t *testing.T) {
	type args struct {
		jwttoken string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "Test token without api and namespace",
			args: args{
				jwttoken: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwicXVvdGEiOjMsImRhdGEiOnsiYWNjb3VudCI6ImFwb211eCIsImVtYWlsIjoiYWRtaW5AYXBvbXV4LmNvbSIsImlkIjoiNWI0OTBlY2M3ZGRmMWY3NWFiODRlN2IxIiwib3JnYW5pemF0aW9uIjoiYXBvbXV4IiwicmVhbG0iOiJ2aW5jZSJ9LCJhdWQiOiJhcG9yZXRvLmNvbSIsImV4cCI6MTU0NzY4MzIxOCwiaWF0IjoxNTQ3NTkzMjE4LCJpc3MiOiJtaWRnYXJkLmFwb211eC5jb20iLCJzdWIiOiJhcG9tdXgifQ.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "Test token with api and namespace",
			args: args{
				jwttoken: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwicXVvdGEiOjQsImRhdGEiOnsiYWNjb3VudCI6InZhcnVuIiwiZW1haWwiOiJ2YXJrc0BhcG9yZXRvLmNvbSIsImlkIjoiNWI1OGNiMDY4ODZkYjEwMDAxMGJkODNkIiwib3JnYW5pemF0aW9uIjoidmFydW4iLCJyZWFsbSI6InZpbmNlIn0sIm9wYXF1ZSI6eyJtYWNoaW5lIjoicHVsc2FyIiwibmFtZXNwYWNlIjoiL3ZhcnVuIn0sImFwaSI6Imh0dHBzOi8vYXBpLnNhbmRib3guYXBvcmV0by51cyIsImF1ZCI6InNhbmRib3guYXBvcmV0by51cyIsImV4cCI6MTU0ODc0MjI1MywiaWF0IjoxNTQ4NjUyMjUzLCJpc3MiOiJzYW5kYm94LmFwb3JldG8udXMiLCJzdWIiOiJ2YXJ1biJ9.cfkQH7ybikC3uNTR9GpbG-UIxxtHq_u3aDH4uCc7YdKZ5t1CGvefaSBW9a2TG51jQ8ilWd8sEqafKUmkbko_0w",
			},
			want:    "https://api.sandbox.aporeto.us",
			want1:   "/varun",
			wantErr: false,
		},
		{
			name: "Test token without opaque",
			args: args{
				jwttoken: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoidmFydW4iLCJlbWFpbCI6InZhcmtzQGFwb3JldG8uY29tIiwiaWQiOiI1YjU4Y2IwNjg4NmRiMTAwMDEwYmQ4M2QiLCJvcmdhbml6YXRpb24iOiJ2YXJ1biIsInJlYWxtIjoidmluY2UifSwiYXBpIjoiaHR0cHM6Ly9hcGkuc2FuZGJveC5hcG9yZXRvLnVzIiwiYXVkIjoic2FuZGJveC5hcG9yZXRvLnVzIiwiZXhwIjoxNTQ4Nzg3MTA0LCJpYXQiOjE1NDg2OTcxMDQsImlzcyI6InNhbmRib3guYXBvcmV0by51cyIsInN1YiI6InZhcnVuIn0.7mqvGj1s40J_sATfE8TiEDV0KxOj87V8E4pAAf05isa5NMKM1DpToZeo1jzg9dHaUoHZvpiZXgtTkHby5gfT_Q",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "Test token without ns",
			args: args{
				jwttoken: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoidmFydW4iLCJlbWFpbCI6InZhcmtzQGFwb3JldG8uY29tIiwiaWQiOiI1YjU4Y2IwNjg4NmRiMTAwMDEwYmQ4M2QiLCJvcmdhbml6YXRpb24iOiJ2YXJ1biIsInJlYWxtIjoidmluY2UifSwib3BhcXVlIjp7Imp1bmsiOiJ0cnVlIn0sImFwaSI6Imh0dHBzOi8vYXBpLnNhbmRib3guYXBvcmV0by51cyIsImF1ZCI6InNhbmRib3guYXBvcmV0by51cyIsImV4cCI6MTU0ODc4NzE1OCwiaWF0IjoxNTQ4Njk3MTU4LCJpc3MiOiJzYW5kYm94LmFwb3JldG8udXMiLCJzdWIiOiJ2YXJ1biJ9.sKxocX22rCofvfjIbdSISsmKn4D7RYk5entp2zblLfFdvmpj9TBJh1F69rCrDiZattf30rLjlds90_-7n5Eiyw",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "Invalid token",
			args: args{
				jwttoken: "eyJhbGciOiJFUzI1NiIsInWFsbSI6IlNlIiwicXVvdGEiOiIzIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTQ3NjgzMjE4LCJpYXQiOjE1NDc1OTMyMTgsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.N7B-X3rRcySodn0q4u1NUAVFIEtjnZEYJGidAFSwflyAhpqchRmm6P_waaVBcGhnRNhsIUayuJjeMpXccYFrWA",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ExtractAPIAndNamespace(tt.args.jwttoken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAPIAndNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractAPIAndNamespace() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ExtractAPIAndNamespace() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
