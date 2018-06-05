package magetask

import (
	"reflect"
	"testing"
)

func Test_prependPathToExcludes(t *testing.T) {
	type args struct {
		exclude []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test Excludes Empty input",
			args: args{
				exclude: []string{},
			},
			want: []string{},
		},
		{
			name: "Test Excludes one value",
			args: args{
				exclude: []string{"somepackage"},
			},
			want: []string{"github.com/aporeto-inc/addedeffect/magetask/somepackage"},
		},
		{
			name: "Test Excludes some values",
			args: args{
				exclude: []string{
					"somepackage",
					"someotherpackage",
				},
			},
			want: []string{
				"github.com/aporeto-inc/addedeffect/magetask/somepackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prependPathToExcludes(tt.args.exclude); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prependPathToExcludes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prunePackagesToExclude(t *testing.T) {
	type args struct {
		packages []string
		exclude  []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		{
			name: "Empty exclude list",
			args: args{
				packages: []string{
					"github.com/aporeto-inc/addedeffect/magetask/somepackage",
					"github.com/aporeto-inc/addedeffect/magetask/somepackage/subpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
				},
				exclude: []string{},
			},
			wantRet: []string{
				"github.com/aporeto-inc/addedeffect/magetask/somepackage",
				"github.com/aporeto-inc/addedeffect/magetask/somepackage/subpackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
			},
		},
		{
			name: "Somepackage exclude list",
			args: args{
				packages: []string{
					"github.com/aporeto-inc/addedeffect/magetask/somepackage",
					"github.com/aporeto-inc/addedeffect/magetask/somepackage/subpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
				},
				exclude: []string{"github.com/aporeto-inc/addedeffect/magetask/somepackage"},
			},
			wantRet: []string{
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
			},
		},
		{
			name: "Some exclude  which shouldnt remove any entry",
			args: args{
				packages: []string{
					"github.com/aporeto-inc/addedeffect/magetask/somepackage",
					"github.com/aporeto-inc/addedeffect/magetask/somepackage/subpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
					"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
				},
				exclude: []string{"github.com/aporeto-inc/addedeffect/magetask/some"},
			},
			wantRet: []string{
				"github.com/aporeto-inc/addedeffect/magetask/somepackage",
				"github.com/aporeto-inc/addedeffect/magetask/somepackage/subpackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage",
				"github.com/aporeto-inc/addedeffect/magetask/someotherpackage/subpackage",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := prunePackagesToExclude(tt.args.packages, tt.args.exclude); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("prunePackagesToExclude() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
