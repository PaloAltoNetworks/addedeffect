package lombric

import (
	"testing"
)

func Test_stringInSlice(t *testing.T) {
	type args struct {
		str  string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test string in slice",
			args{
				"a",
				[]string{"a", "b", "c"},
			},
			true,
		},
		{
			"test string not in slice",
			args{
				"z",
				[]string{"a", "b", "c"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.str, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkRequired(t *testing.T) {
	type args struct {
		failFunc func()
		keys     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test failure",
			args{
				func() {},
				[]string{"a", "b"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRequired(tt.args.failFunc, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("checkRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkAllowedValues(t *testing.T) {
	type args struct {
		failFunc      func()
		allowedValues map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test failure",
			args{
				func() {},
				map[string][]string{"a": []string{"1", "2"}},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAllowedValues(tt.args.failFunc, tt.args.allowedValues); (err != nil) != tt.wantErr {
				t.Errorf("checkAllowedValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
