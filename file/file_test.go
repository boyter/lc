// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense
package file

import (
	"os"
	"strings"
	"testing"
)

func TestFindRepositoryRoot(t *testing.T) {
	// We expect this to walk back from file to cs
	curdir, _ := os.Getwd()
	root := FindRepositoryRoot(curdir)

	if strings.HasSuffix(root, "file") {
		t.Error("Expected to walk back to root")
	}
}

func TestGetExtension(t *testing.T) {
	got := GetExtension("something.c")
	expected := "c"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionNoExtension(t *testing.T) {
	got := GetExtension("something")
	expected := "something"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionMultipleDots(t *testing.T) {
	got := GetExtension(".travis.yml")
	expected := "travis.yml"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionMultipleExtensions(t *testing.T) {
	got := GetExtension("something.go.yml")
	expected := "go.yml"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionStartsWith(t *testing.T) {
	got := GetExtension(".gitignore")
	expected := ".gitignore"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionTypeScriptDefinition(t *testing.T) {
	got := GetExtension("test.d.ts")
	expected := "d.ts"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}

func TestGetExtensionSecondPass(t *testing.T) {
	got := GetExtension("test.d.ts")
	got = GetExtension(got)
	expected := "ts"

	if got != expected {
		t.Errorf("Expected %s got %s", expected, got)
	}
}
