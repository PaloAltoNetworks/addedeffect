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

	"github.com/stretchr/testify/assert"
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
		io.Copy(&buf, oreader)
		oout <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, ereader)
		eout <- buf.String()
	}()
	wg.Wait()

	// Execute function
	f()

	owriter.Close()
	ewriter.Close()
	// Return captures
	return <-oout, <-eout
}

func TestConfigureWithOptions(t *testing.T) {
	type args struct {
		level           string
		format          string
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
			name: "no file logging",
			args: args{
				"info",
				"json",
				"",
				false,
				false,
			},
			iterations: 20,
			want:       "",
		},
		{
			name: "file only",
			args: args{
				"info",
				"json",
				"/tmp/some-log-file",
				true,
				false,
			},
			iterations: 20,
			want:       "",
		},
		{
			name: "tee file logging and stdout with wraps in files",
			args: args{
				"info",
				"json",
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

				ConfigureWithOptions(tt.args.level, tt.args.format, tt.args.file, tt.args.fileOnly, tt.args.prettyTimestamp)

				for i := 0; i < tt.iterations; i++ {
					buf := fmt.Sprintf("%80d - hello", i)
					minBytesPrinted += len(buf)
					zap.L().Info(buf)
				}
			})

			// On fileOnly case, nothing is expected on stderr as well.
			if tt.args.fileOnly {
				minBytesPrinted = 0
			}

			// validate nothing is printed on stdout
			// assert.Equal(t, 0, len(ro))
			// validate we have printed more than minBytesPrinted on stderr
			// assert.LessOrEqual(t, minBytesPrinted, len(re))

			if tt.args.file != "" {

				numFilesExpected := 1
				if minBytesPrinted > logFileSizeDefault*1024*1024 {
					numFilesExpected += logFileNumBackups
				}

				// Wait for one second as file may not have been deleted.
				time.Sleep(time.Second)

				files, err := filepath.Glob(tt.args.file + "*")
				assert.ObjectsAreEqual(err, nil)

				// logging to files tests wraparound. we should have
				assert.Equal(t, numFilesExpected, len(files))

				for _, f := range files {

					fi, err := os.Stat(f)
					assert.ObjectsAreEqual(err, nil)

					// If we printed into a file, the file should be less than wrapped size.
					assert.Less(t, fi.Size(), int64(logFileSizeDefault*1024*1024))

					err = os.Remove(f)
					assert.ObjectsAreEqual(err, nil)
				}
			}
		})
	}
}
