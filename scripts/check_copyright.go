// Copyright 2017 Elasticsearch Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// File extensions to check
var checkExts = map[string]bool{
	".c":  true,
	".go": true,
}

// Valid copyright headers, searched for in the top five lines in each file.
var copyrightRegexps = []string{
	`Copyright`,
}

var copyrightRe = regexp.MustCompile(strings.Join(copyrightRegexps, "|"))

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No directories specified.")
		os.Exit(1)
	}

	for _, dir := range args {
		err := filepath.Walk(dir, checkCopyright)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func checkCopyright(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	if strings.Contains(path, "vendor"+string(filepath.Separator)) {
		return nil
	}
	if !checkExts[filepath.Ext(path)] {
		return nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for i := 0; scanner.Scan() && i < 5; i++ {
		if copyrightRe.MatchString(scanner.Text()) {
			return nil
		}
	}

	return fmt.Errorf("Missing copyright in %s?", path)
}
