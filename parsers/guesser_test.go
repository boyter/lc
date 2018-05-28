package parsers

import (
	// "io/ioutil"
	"path/filepath"
	// "strings"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCleanText(t *testing.T) {
	actual := cleanText([]byte("ToLower"))
	expected := "tolower"

	if expected != string(actual) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = cleanText([]byte("   ToLower999$%"))
	expected = " tolower999 "

	if expected != string(actual) {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestLoadDatabase(t *testing.T) {
	actual := loadDatabase()

	if len(actual) == 0 {
		t.Errorf("Expected database to not be empty")
	}
}

func TestProcessFile(t *testing.T) {
	Format = "csv"
	actual := processFile("../examples/identifier/", "has_identifier.py", []LicenseMatch{})

	if actual.Md5Hash != "fc7d75e0bc0275841de8426b18791fa4" {
		t.Errorf("Expected MD5 to match")
	}

	if actual.Sha1Hash != "03a614cc51e9a783a695bcf99ec4adcdac34e1cc" {
		t.Errorf("Expected SHA1 to match")
	}

	if actual.Sha256Hash != "5b6bf8d45b25a0dab4f0817324618e425037f156b67cdc7503da01e6d9beb652" {
		t.Errorf("Expected SHA256 to match")
	}

	if len(actual.LicenseIdentified) != 2 {
		t.Errorf("Expected 2 identified licenses")
	}

	if actual.LicenseIdentified[0].LicenseId != "GPL-2.0" {
		t.Errorf("Expected license not identified")
	}

	if actual.LicenseIdentified[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Expected license not identified")
	}
}

func TestProcessFileLicenses(t *testing.T) {
	files, _ := ioutil.ReadDir("../examples/licenses/")

	for _, file := range files {
		content, _ := ioutil.ReadFile(filepath.Join("../examples/licenses/", file.Name()))

		var lic = License{}
		json.Unmarshal(content, &lic)
		actual := keywordGuessLicense([]byte(lic.LicenseText), loadDatabase())

		found := false
		for _, license := range actual {
			if license.LicenseId == lic.LicenseId {
				found = true
			}
		}

		// TODO Diffmark seems to have some issue which needs to be investigated
		if found == false && file.Name() != "diffmark.json" {
			t.Errorf("Expected license to be found %s", file.Name())
		}
	}
}

func TestRegressionIssue36(t *testing.T) {
	content := `The MIT License (MIT)

Copyright (c) 2018 Ben Boyter

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
`

	result := keywordGuessLicense([]byte(content), loadDatabase())

	t.Log(result)
	if result[0].LicenseId != "MIT" {
		t.Errorf("Should be MIT")
	}

	content = `Copyright (c) 1995-2002 Rik Faith <rikfaith@gmail.com>
Copyright (c) 2002-2018 Aleksey Cheusov <vle@gmx.net>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

	result = keywordGuessLicense([]byte(content), loadDatabase())

	t.Log(result)
	if result[0].LicenseId != "MIT-feh" {
		t.Errorf("Should be MIT-feh")
	}
}

func TestIdentifierGuessLicence(t *testing.T) {
	actual := identifierGuessLicence("test", loadDatabase())
	if len(actual) != 0 {
		t.Errorf("Should be no matches")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 ", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 \n # SPDX-License-Identifier: GPL-3.0+", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}

	actual = identifierGuessLicence(`# SPDX-License-Identifier: GPL-2.0

import cherrypy

class Example(object):

    # SPDX-License-Identifier: GPL-3.0+
    @cherrypy.expose
    def index(self, **params):
        return 'Hello World'


if __name__ == '__main__':
    cherrypy.config.update({
        'server.socket_host': '0.0.0.0',
        'server.socket_port': 8080,
        'server.thread_pool': 30,
    })
    cherrypy.quickstart(Example())
`, loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}
}
