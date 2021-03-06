// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fuchsia.googlesource.com/jiri/jiritest"
)

type importTestCase struct {
	Args           []string
	Filename       string
	Exist, Want    string
	Stdout, Stderr string
	SetFlags       func()
}

func TestImport(t *testing.T) {
	tests := []importTestCase{
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = false
				flagImportOut = ""
			},
			Stderr: `wrong number of arguments`,
		},
		{
			Args:   []string{"a"},
			Stderr: `wrong number of arguments`,
		},
		{
			Args:   []string{"a", "b", "c"},
			Stderr: `wrong number of arguments`,
		},
		// Remote imports, default append behavior
		{
			SetFlags: func() {
				flagImportName = "name"
				flagImportRemoteBranch = "remotebranch"
				flagImportRoot = "root"
				flagImportOverwrite = false
				flagImportOut = ""
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="name" remote="https://github.com/new.git" remotebranch="remotebranch" root="root"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = false
				flagImportOut = ""
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = false
				flagImportOut = "file"
			},
			Args:     []string{"foo", "https://github.com/new.git"},
			Filename: `file`,
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = false
				flagImportOut = "-"
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Stdout: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = false
				flagImportOut = ""
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Exist: `<manifest>
  <imports>
    <import manifest="bar" name="manifest" remote="https://github.com/orig.git"/>
  </imports>
</manifest>
`,
			Want: `<manifest>
  <imports>
    <import manifest="bar" name="manifest" remote="https://github.com/orig.git"/>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		// Remote imports, explicit overwrite behavior
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = true
				flagImportOut = ""
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = true
				flagImportOut = "file"
			},
			Args:     []string{"foo", "https://github.com/new.git"},
			Filename: `file`,
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = true
				flagImportOut = "-"
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Stdout: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
		{
			SetFlags: func() {
				flagImportName = "manifest"
				flagImportRemoteBranch = "master"
				flagImportRoot = ""
				flagImportOverwrite = true
				flagImportOut = ""
			},
			Args: []string{"foo", "https://github.com/new.git"},
			Exist: `<manifest>
  <imports>
    <import manifest="bar" name="manifest" remote="https://github.com/orig.git"/>
  </imports>
</manifest>
`,
			Want: `<manifest>
  <imports>
    <import manifest="foo" name="manifest" remote="https://github.com/new.git"/>
  </imports>
</manifest>
`,
		},
	}

	// Temporary directory in which our jiri binary will live.
	binDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(binDir)

	for _, test := range tests {
		if err := testImport(t, test); err != nil {
			t.Errorf("%v: %v", test.Args, err)
		}
	}
}

func testImport(t *testing.T, test importTestCase) error {
	jirix, cleanup := jiritest.NewX(t)
	defer cleanup()
	// Temporary directory in which to run `jiri import`.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Return to the current working directory when done.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	// cd into a root directory in which to do the actual import.
	jiriRoot := jirix.Root
	if err := os.Chdir(jiriRoot); err != nil {
		return err
	}

	// Allow optional non-default filenames, for testing the -out option.
	filename := test.Filename
	if filename == "" {
		filename = ".jiri_manifest"
	}

	// Set up manfile for the local file import tests.  It should exist in both
	// the tmpDir (for ../manfile tests) and jiriRoot.
	for _, dir := range []string{tmpDir, jiriRoot} {
		if err := ioutil.WriteFile(filepath.Join(dir, "manfile"), nil, 0644); err != nil {
			return err
		}
	}

	// Set up an existing file if it was specified.
	if test.Exist != "" {
		if err := ioutil.WriteFile(filename, []byte(test.Exist), 0644); err != nil {
			return err
		}
	}

	// Run import and check the results.
	importCmd := func() {
		if test.SetFlags != nil {
			test.SetFlags()
		}
		err = runImport(jirix, test.Args)
	}
	stdout, _, runErr := runfunc(importCmd)
	if runErr != nil {
		return err
	}
	stderr := ""
	if err != nil {
		stderr = err.Error()
	}
	if got, want := stdout, test.Stdout; !strings.Contains(got, want) || (got != "" && want == "") {
		return fmt.Errorf("stdout got %q, want substr %q", got, want)
	}
	if got, want := stderr, test.Stderr; !strings.Contains(got, want) || (got != "" && want == "") {
		return fmt.Errorf("stderr got %q, want substr %q", got, want)
	}

	// Make sure the right file is generated.
	if test.Want != "" {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		if got, want := string(data), test.Want; got != want {
			return fmt.Errorf("GOT\n%s\nWANT\n%s", got, want)
		}
	}
	return nil
}
