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

package logutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

// captureOutAndErr executes a function and returns captured os.Stdout and os.Stderr during this execution.
func captureOutAndErr(f func()) (o, e string) {

	// Create output and error capture pipes
	ereader, ewriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	oreader, owriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Store original stdout and stderr and restore them by defer call.
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	// Hijack stdout and stderr
	os.Stdout = owriter
	os.Stderr = ewriter

	// Setup capture funcs
	oout := make(chan string)
	eout := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, err := io.Copy(&buf, oreader)
		if err != nil {
			panic(err)
		}
		oout <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, err := io.Copy(&buf, ereader)
		if err != nil {
			panic(err)
		}
		eout <- buf.String()
	}()
	wg.Wait()

	// Execute function
	f()

	err = owriter.Close()
	if err != nil {
		panic(err)
	}
	err = ewriter.Close()
	if err != nil {
		panic(err)
	}
	// Return captures
	return <-oout, <-eout
}

func TestConfigureWithOptions(t *testing.T) {
	type args struct {
		level           string
		format          string
		service         string
		file            string
		fileOnly        bool
		prettyTimestamp bool
	}
	tests := []struct {
		name       string
		args       args
		iterations int
		want       string
	}{
		{
			name: "file only",
			args: args{
				"info",
				"stackdriver",
				"",
				"/tmp/some-log-file",
				true,
				false,
			},
			iterations: 3,
			want:       "",
		},
		{
			name: "no file logging at info level",
			args: args{
				"info",
				"json",
				"",
				"",
				false,
				true,
			},
			iterations: 6,
			want:       "",
		},
		{
			name: "service name on json",
			args: args{
				"info",
				"stackdriver",
				"some-bizzare-bad-named-service-1",
				"",
				false,
				false,
			},
			iterations: 1,
			want:       "",
		},
		{
			name: "service name on stackdriver",
			args: args{
				"info",
				"json",
				"some-bizzare-bad-named-service-2",
				"",
				false,
				false,
			},
			iterations: 1,
			want:       "",
		},
		{
			name: "tee file logging and stdout with wraps in files",
			args: args{
				"info",
				"json",
				"",
				"/tmp/some-tee-log-file",
				false,
				false,
			},
			iterations: 160 * 1024,
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			minBytesPrinted := 0
			ro, re := captureOutAndErr(func() {

				opts := []Option{OptionFile(tt.args.file, tt.args.fileOnly)}
				if tt.args.service != "" {
					opts = append(opts, OptionServiceName(tt.args.service))
				}
				if tt.args.prettyTimestamp {
					opts = append(opts, OptionPrettyTimeStamp())
				}
				Configure(tt.args.level, tt.args.format, opts...)
				for i := 0; i < tt.iterations; i++ {
					buf := fmt.Sprintf("%80d - hello", i)
					minBytesPrinted += len(buf)
					switch i % 6 {
					case 0:
						zap.L().Info(buf)
					case 1:
						zap.L().Warn(buf)
					case 2:
						zap.L().Debug(buf)
					case 3, 4, 5:
						zap.L().Error(buf)
					}
				}
			})

			// validate nothing is printed on stdout
			if ok, message := assertions.So(ro, ShouldBeBlank); !ok {
				t.Fatalf("Failure in %s:\n%s", tt.name, message)
			}
			if tt.args.fileOnly {
				// validate nothing is printed on stderr in fileOnly case
				if ok, message := assertions.So(re, ShouldBeBlank); !ok {
					t.Fatalf("Failure in %s:\n%s", tt.name, message)
				}
			} else {
				// validate we have printed more than minBytesPrinted on stderr
				if ok, message := assertions.So(minBytesPrinted, ShouldBeLessThan, len(re)); !ok {
					t.Fatalf("Failure in %s:\n%s", tt.name, message)
				}
			}

			if tt.args.service != "" {
				if ok, message := assertions.So(re, ShouldContainSubstring, tt.args.service); !ok {
					t.Fatalf("Failure in %s:\n%s", tt.name, message)
				}
			}

			if tt.args.file != "" {

				numFilesExpected := 1
				if minBytesPrinted > logFileSizeDefault*1024*1024 {
					numFilesExpected += logFileNumBackups
				}

				// Wait for one second as file may not have been deleted.
				time.Sleep(5 * time.Second)

				files, err := filepath.Glob(tt.args.file + "*")
				if ok, message := assertions.So(err, ShouldBeNil); !ok {
					t.Fatalf("Failure in %s:\n%s", tt.name, message)
				}

				// logging to files tests wraparound. we should have
				if ok, message := assertions.So(len(files), ShouldEqual, numFilesExpected); !ok {
					t.Fatalf("Failure in %s:\n%s", tt.name, message)
				}

				for _, f := range files {

					fi, err := os.Stat(f)
					if ok, message := assertions.So(err, ShouldBeNil); !ok {
						t.Fatalf("Failure in %s:\n%s", tt.name, message)
					}

					// If we printed into a file, the file should be less than wrapped size.
					if ok, message := assertions.So(fi.Size(), ShouldBeLessThan, int64(logFileSizeDefault*1024*1024)); !ok {
						t.Fatalf("Failure in %s:\n%s", tt.name, message)
					}

					err = os.Remove(f)
					if ok, message := assertions.So(err, ShouldBeNil); !ok {
						t.Fatalf("Failure in %s:\n%s", tt.name, message)
					}
				}
			}
		})
	}
}
