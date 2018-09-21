package parsers

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCleanText(t *testing.T) {
	actual := cleanText([]byte("ToLower"))
	expected := "tolower"

	if expected != string(actual) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = cleanText([]byte("   ToLower999$%   "))
	expected = "tolower999"

	if expected != string(actual) {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestLoadDatabase(t *testing.T) {
	actual := LoadDatabase()

	if len(actual) == 0 {
		t.Errorf("Expected database to not be empty")
	}
}

func TestProcessFile(t *testing.T) {
	Format = "csv"
	actual := processFile("../examples/identifier/", "has_identifier.py", []LicenseMatch{}, []LicenseMatch{})

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
		actual := keywordGuessLicense([]byte(lic.LicenseText), LoadDatabase())

		found := false
		for _, license := range actual {
			if license.LicenseId == lic.LicenseId {
				found = true
			}
		}

		// TODO Diffmark seems to have some issue which needs to be investigated
		if found == false && file.Name() != "diffmark.json" && file.Name() != "NCSA.json" {
			t.Errorf("Expected license to be found %s", file.Name())
		}
	}
}

func TestRegression(t *testing.T) {
	content := `##The MIT License (MIT)

> Copyright (c) 2014-2015 Vicc Alexander

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in all
> copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.`

	result := keywordGuessLicense([]byte(content), LoadDatabase())

	t.Log(result)
	if result[0].LicenseId != "MIT" {
		t.Errorf("Should be MIT")
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

	result := keywordGuessLicense([]byte(content), LoadDatabase())

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

	result = keywordGuessLicense([]byte(content), LoadDatabase())

	t.Log(result)
	if result[0].LicenseId != "MIT" {
		t.Errorf("Should be MIT")
	}
}

func TestIdentifierGuessLicence(t *testing.T) {
	actual := identifierGuessLicence("test", LoadDatabase())
	if len(actual) != 0 {
		t.Errorf("Should be no matches")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0", LoadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 ", LoadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 \n # SPDX-License-Identifier: GPL-3.0+", LoadDatabase())
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
`, LoadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}
}

func TestApache2SampleHeader(t *testing.T) {
	content := `************************************
This project is available under the following license:
************************************
Copyright 2012-2016 eBusiness Information
Copyright 2016-2017 the AndroidAnnotations project

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.

This project uses CodeModel (http://codemodel.java.net/), which is
distributed under the GlassFish Dual License, which means CodeModel is
subject to the terms of either the GNU General Public License Version 2 only
("GPL") or the Common Development and Distribution License("CDDL").

You may obtain a copy of the "CDDL" License at

http://www.opensource.org/licenses/cddl1.php

As per section 3.6 ("Larger Works") of the "CDDL" License, we may create a
Larger Work by combining Covered Software with other code not governed by
the terms of this License and distribute the Larger Work as a single
product.

We are therefore allowed to distribute CodeModel without Modification as
part of AndroidAnnotations.

About AndroidAnnotations generated files: if you create a larger work that
contains files generated by AndroidAnnotations, you can distribute that work
under terms of your choice.`

	result := keywordGuessLicense([]byte(content), LoadDatabase())

	if len(result) == 0 {
		t.Errorf("Should find at least one license")
	}

	if result[0].LicenseId != "Apache-2.0" {
		t.Errorf("Should be Apache-2.0")
	}
}

func TestBsd3Clause(t *testing.T) {
	content := `Copyright (c) 2006-2015, Salvatore Sanfilippo
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
    * Neither the name of Redis nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	result := keywordGuessLicense([]byte(content), LoadDatabase())

	if len(result) == 0 {
		t.Errorf("Should find at least one license")
	}

	if result[0].LicenseId != "BSD-3-Clause" {
		t.Errorf("Should be BSD-3-Clause was " + result[0].LicenseId)
	}
}

func TestBsd3ClauseWithAPI(t *testing.T) {
	content := `Copyright (c) 2006-2015, Salvatore Sanfilippo
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
    * Neither the name of Redis nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	result := GuessLicense([]byte(content))

	if len(result) == 0 {
		t.Errorf("Should find at least one license")
	}

	if result[0].LicenseId != "BSD-3-Clause" {
		t.Errorf("Should be BSD-3-Clause was " + result[0].LicenseId)
	}
}
