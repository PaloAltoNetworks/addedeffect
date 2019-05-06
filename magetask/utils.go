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

package magetask

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/imports"
)

// writeVersionsFile writes the version file from the given templateData at the given path.
func writeVersionsFile(templateData versionTemplate, versionsFilePath string) error {

	t := template.Must(template.New("versions").Funcs(template.FuncMap{
		"cap":       strings.Title,
		"short":     path.Base,
		"hasprefix": strings.HasPrefix,
		"varname": func(v string) string {
			return strings.Title(path.Base(strings.Replace(v, "-", "", -1)))
		},
	}).Parse(versionFileTemplate))

	outFile := path.Join(versionsFilePath, "versions.go")
	_, err := os.Stat(versionsFilePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(versionsFilePath, 0744); err != nil {
			panic(err)
		}
	}

	_, err = os.Stat(outFile)
	if err == nil {
		if err = os.Remove(outFile); err != nil {
			panic(err)
		}
	}

	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close() // nolint: errcheck

	buffer := &bytes.Buffer{}
	if err = t.Execute(buffer, templateData); err != nil {
		return err
	}

	data, err := imports.Process(".", buffer.Bytes(), &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	})
	if err != nil {
		fmt.Println(buffer.String())
		panic(err)
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}
