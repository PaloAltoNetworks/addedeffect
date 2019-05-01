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
	"reflect"
	"testing"

	"go.aporeto.io/elemental"
	testmodel "go.aporeto.io/elemental/test/model"
)

func TestParseAudience(t *testing.T) {
	type args struct {
		audString string
	}
	tests := []struct {
		name    string
		args    args
		want    AudiencesList
		wantErr bool
	}{
		{
			"valid single aud string",
			args{
				"aud:retrieve:lists:/a/b",
			},
			AudiencesList{
				Audience{
					Operations: []string{"retrieve"},
					Identities: []string{"lists"},
					Namespaces: []string{"/a/b"},
				},
			},
			false,
		},
		{
			"valid multiple aud string",
			args{
				"aud:retrieve:lists:/a/b;aud:create:tasks:/a/c",
			},
			AudiencesList{
				Audience{
					Operations: []string{"retrieve"},
					Identities: []string{"lists"},
					Namespaces: []string{"/a/b"},
				},
				Audience{
					Operations: []string{"create"},
					Identities: []string{"tasks"},
					Namespaces: []string{"/a/c"},
				},
			},
			false,
		},
		{
			"valid single composite aud string",
			args{
				"aud:retrieve,retrieve-many:lists,tasks:/a/b,/b/c",
			},
			AudiencesList{
				Audience{
					Operations: []string{"retrieve", "retrieve-many"},
					Identities: []string{"lists", "tasks"},
					Namespaces: []string{"/a/b", "/b/c"},
				},
			},
			false,
		},
		{
			"valid multiple composite aud string",
			args{
				"aud:retrieve,retrieve-many:lists,tasks:/a/b,/b/c;aud:create,delete:users:*",
			},
			AudiencesList{
				Audience{
					Operations: []string{"retrieve", "retrieve-many"},
					Identities: []string{"lists", "tasks"},
					Namespaces: []string{"/a/b", "/b/c"},
				},
				Audience{
					Operations: []string{"create", "delete"},
					Identities: []string{"users"},
					Namespaces: []string{AudienceAny},
				},
			},
			false,
		},
		{
			"valid full any",
			args{
				"aud:*:*:*",
			},
			AudiencesList{
				Audience{
					Operations: []string{AudienceAny},
					Identities: []string{AudienceAny},
					Namespaces: []string{AudienceAny},
				},
			},
			false,
		},
		{
			"invalid operation",
			args{
				"aud:nothing:lists:*",
			},
			nil,
			true,
		},
		{
			"invalid identity",
			args{
				"aud:create:weird:*",
			},
			nil,
			true,
		},
		{
			"invalid single aud string",
			args{
				"retrieve,:lists",
			},
			nil,
			false, // TODO: switch to false when workaround is gone
		},
		{
			"invalid multiple aud string",
			args{
				"aud:retrieve,retrieve-many:lists,tasks:/a/b,/b/c;retrieve,:lists",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAudience(tt.args.audString, testmodel.Manager())
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudiencesList_String(t *testing.T) {
	tests := []struct {
		name string
		a    AudiencesList
		want string
	}{
		{
			"simple",
			AudiencesList{
				Audience{
					Operations: []string{"op"},
					Identities: []string{"ident"},
					Namespaces: []string{"/ns"},
				},
				Audience{
					Operations: []string{"op1", "op2"},
					Identities: []string{"ident1", "ident2"},
					Namespaces: []string{"/ns1", "/ns2"},
				},
			},
			"aud:op:ident:/ns;aud:op1,op2:ident1,ident2:/ns1,/ns2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("AudiencesList.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudience_String(t *testing.T) {
	type fields struct {
		Operations []string
		Identities []string
		Namespaces []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple",
			fields{
				[]string{"op"},
				[]string{"ident"},
				[]string{"/ns"},
			},
			"aud:op:ident:/ns",
		},
		{
			"composite",
			fields{
				[]string{"op1", "op2"},
				[]string{"ident1", "ident2"},
				[]string{"/ns1", "/ns2"},
			},
			"aud:op1,op2:ident1,ident2:/ns1,/ns2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Audience{
				Operations: tt.fields.Operations,
				Identities: tt.fields.Identities,
				Namespaces: tt.fields.Namespaces,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("Audience.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudience_Verify(t *testing.T) {
	type fields struct {
		Operations []string
		Identities []string
		Namespaces []string
	}
	type args struct {
		operation elemental.Operation
		identity  elemental.Identity
		namespace string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Valid
		{
			"valid full",
			fields{
				Operations: []string{"create"},
				Identities: []string{"lists"},
				Namespaces: []string{"/ns"},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid full floating operation",
			fields{
				Operations: []string{AudienceAny},
				Identities: []string{"lists"},
				Namespaces: []string{"/ns"},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid floating identities",
			fields{
				Operations: []string{"create"},
				Identities: []string{AudienceAny},
				Namespaces: []string{"/ns"},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid floating namespace",
			fields{
				Operations: []string{"create"},
				Identities: []string{"lists"},
				Namespaces: []string{AudienceAny},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid floating operation and identity",
			fields{
				Operations: []string{AudienceAny},
				Identities: []string{AudienceAny},
				Namespaces: []string{"/ns"},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid floating operation and namespace",
			fields{
				Operations: []string{AudienceAny},
				Identities: []string{"lists"},
				Namespaces: []string{AudienceAny},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"valid floating identity and namespace",
			fields{
				Operations: []string{"create"},
				Identities: []string{AudienceAny},
				Namespaces: []string{AudienceAny},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},

		// Invalid
		{
			"invalid operation",
			fields{
				Operations: []string{"retrieve"},
				Identities: []string{AudienceAny},
				Namespaces: []string{AudienceAny},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			false,
		},
		{
			"invalid identity",
			fields{
				Operations: []string{AudienceAny},
				Identities: []string{"task"},
				Namespaces: []string{AudienceAny},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			false,
		},
		{
			"invalid ns",
			fields{
				Operations: []string{AudienceAny},
				Identities: []string{AudienceAny},
				Namespaces: []string{"/not-ns"},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Audience{
				Operations: tt.fields.Operations,
				Identities: tt.fields.Identities,
				Namespaces: tt.fields.Namespaces,
			}
			if got := a.Verify(tt.args.operation, tt.args.identity, tt.args.namespace); got != tt.want {
				t.Errorf("Audience.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudiencesList_Verify(t *testing.T) {
	type args struct {
		operation elemental.Operation
		identity  elemental.Identity
		namespace string
	}
	tests := []struct {
		name string
		a    AudiencesList
		args args
		want bool
	}{
		{
			"all valid",
			AudiencesList{
				Audience{
					Operations: []string{AudienceAny},
					Identities: []string{AudienceAny},
					Namespaces: []string{AudienceAny},
				},
				Audience{
					Operations: []string{AudienceAny},
					Identities: []string{AudienceAny},
					Namespaces: []string{AudienceAny},
				},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"first valid, second invalid",
			AudiencesList{
				Audience{
					Operations: []string{AudienceAny},
					Identities: []string{AudienceAny},
					Namespaces: []string{AudienceAny},
				},
				Audience{
					Operations: []string{"you"},
					Identities: []string{"shall not"},
					Namespaces: []string{"pass"},
				},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"first invalid, second valid",
			AudiencesList{
				Audience{
					Operations: []string{"you"},
					Identities: []string{"(maybe)"},
					Namespaces: []string{"pass"},
				},
				Audience{
					Operations: []string{AudienceAny},
					Identities: []string{AudienceAny},
					Namespaces: []string{AudienceAny},
				},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			true,
		},
		{
			"all invalid",
			AudiencesList{
				Audience{
					Operations: []string{"you"},
					Identities: []string{"shall"},
					Namespaces: []string{"pass"},
				},
				Audience{
					Operations: []string{"vous ne"},
					Identities: []string{"passerez"},
					Namespaces: []string{"pas"},
				},
			},
			args{
				elemental.OperationCreate,
				testmodel.ListIdentity,
				"/ns",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Verify(tt.args.operation, tt.args.identity, tt.args.namespace); got != tt.want {
				t.Errorf("AudiencesList.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsecureAudience(t *testing.T) {

	tokenValidWithAudience := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImF1ZDpjcmVhdGU6dGFza3M6L2FudG9pbmUvem9uZSIsImV4cCI6MTU0ODkxMzUxNiwiaWF0IjoxNTQ4ODIzNTE2LCJpc3MiOiJodHRwczovLzEyNy4wLjAuMTo0NDQzIiwic3ViIjoiYXBvbXV4In0.zk3BKj0X8e9KpdCnDVbxdZacsmtQiE9Con4AcCESp1SmSVcYxgA010-ro9HEiKwFgsl4SfosD6UKpTGnoADDcA`
	tokenValidWithoutAudience := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImV4cCI6MTU0ODkxMTMxMSwiaWF0IjoxNTQ4ODIxMzExLCJpc3MiOiJodHRwczovLzEyNy4wLjAuMTo0NDQzIiwic3ViIjoiYXBvbXV4In0.Tzqiuj1N2ti3GjLqmvd_VUJQSM3IXKZZSjvTMwpgroiQwDkoeGHNZmm4BU9UiyID6wEgqwdTorYU846B6G88hQ`
	tokenValidWithInvalidAudience := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImF1ZDpjcmVhdGU6aW52YWxpZCIsImV4cCI6MTU0ODkxMTQxNiwiaWF0IjoxNTQ4ODIxNDE2LCJpc3MiOiJodHRwczovLzEyNy4wLjAuMTo0NDQzIiwic3ViIjoiYXBvbXV4In0.evd75guuQyPR14TCl6oOSeFjJj-SASG-_qb0Yv-pvfVpldg_eXf3M5xibUuOXRZ62Ipzx39p7qUPHCOr7IodoA`
	tokenValidWithIgnoredAudience := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YjQ5MGVjYzdkZGYxZjc1YWI4NGU3YjEiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImNhdCIsImV4cCI6MTU0ODkxMTQ3NywiaWF0IjoxNTQ4ODIxNDc3LCJpc3MiOiJodHRwczovLzEyNy4wLjAuMTo0NDQzIiwic3ViIjoiYXBvbXV4In0.K39uVmt4f59AIvx1ZT6eG4ula2blkbwQxO5yw-p4kh_jBkqWguDoa7ilqppZf7tLyy8IAw_KzNm_7fTKVBwWRA`
	tokenInvalid := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.INVALID.K39uVmt4f59AIvx1ZT6eG4ula2blkbwQxO5yw-p4kh_jBkqWguDoa7ilqppZf7tLyy8IAw_KzNm_7fTKVBwWRA`

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    AudiencesList
		wantErr bool
	}{
		{
			"valid with audience",
			args{
				tokenValidWithAudience,
			},
			AudiencesList{
				Audience{
					Operations: []string{"create"},
					Identities: []string{"tasks"},
					Namespaces: []string{"/antoine/zone"},
				},
			},
			false,
		},
		{
			"valid with no audience",
			args{
				tokenValidWithoutAudience,
			},
			nil,
			false,
		},
		{
			"valid with invalid audience",
			args{
				tokenValidWithInvalidAudience,
			},
			nil,
			true,
		},
		{
			"invalid token",
			args{
				tokenInvalid,
			},
			nil,
			true,
		},

		{
			"[backward compat] valid with ignored audience",
			args{
				tokenValidWithIgnoredAudience,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnsecureAudience(tt.args.token, testmodel.Manager())
			if (err != nil) != tt.wantErr {
				t.Errorf("UnsecureAudience() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnsecureAudience() = %v, want %v", got, tt.want)
			}
		})
	}
}
