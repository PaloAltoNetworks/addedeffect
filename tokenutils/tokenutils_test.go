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
