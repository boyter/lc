package parsers

import (
	"strings"
	"testing"
)

func TestGeneratePackageVerificationCode(t *testing.T) {
	result := generatePackageVerificationCode([]FileResult{})

	if result != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("Expect generatePackageVerificationCode to have da39a3ee5e6b4b0d3255bfef95601890afd80709 but got %q", result)
	}
}

func TestGenerateDocumentNamespaceNothingSpecified(t *testing.T) {
	DocumentNamespace = ""
	result := generateDocumentNamespace()

	if strings.Contains(result, "http://spdx.org/spdxdocs/") == false {
		t.Errorf("Expect generateDocumentNamespace to have http://spdx.org/spdxdocs/ but got %q", result)
	}
}

func TestGenerateDocumentNamespaceSomethingSpecified(t *testing.T) {
	DocumentNamespace = "something"
	result := generateDocumentNamespace()

	if strings.Contains(result, "something") == false {
		t.Errorf("Expect generateDocumentNamespace to have something but got %q", result)
	}
}
