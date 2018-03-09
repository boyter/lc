package parsers

import (
	"io/ioutil"
	"path/filepath"
	// "strings"
	"testing"
)

func TestCleanText(t *testing.T) {
	actual := cleanText("ToLower")
	expected := "tolower"

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = cleanText("   ToLower999$%")
	expected = " tolower999 "

	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestLoadDatabase(t *testing.T) {
	actual := loadDatabase()

	if len(actual) == 0 {
		t.Errorf("Expected database to not be empty")
	}
}

func TestWalkDirectory(t *testing.T) {
	actual := walkDirectory("../examples/identifier/", [][]LicenseMatch{})

	if len(actual) != 3 {
		t.Errorf("Expected 3 results for directory")
	}
}

func TestProcessFile(t *testing.T) {
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

// This is slow but ensures that things work as we expect for fuzzy matching
// func TestProcessFileLicensesFuzzy(t *testing.T) {
// 	files, _ := ioutil.ReadDir("../examples/licenses/")

// 	correctFuzzy := float64(0.00)
// 	correctKeywords := float64(0.00)

// 	for _, file := range files {
// 		actualFuzzy := processFile("../examples/licenses/", file.Name(), []LicenseMatch{})

// 		if len(actualFuzzy.LicenseGuesses) != 0 {
// 			if strings.Replace(file.Name(), ".json", "", 1) == actualFuzzy.LicenseGuesses[0].LicenseId {
// 				correctFuzzy++
// 			}
// 		}

// 		content := readFile(filepath.Join("../examples/licenses/", file.Name()))
// 		actualKeywords := keywordGuessLicense(string(content), loadDatabase())

// 		if len(actualKeywords) != 0 {
// 			if strings.Replace(file.Name(), ".json", "", 1) == actualKeywords[0].LicenseId {
// 				correctKeywords++
// 			}
// 		}
// 	}

// 	totalPercentageFuzzy := (correctFuzzy / float64(len(files))) * 100
// 	if totalPercentageFuzzy < 0.95 {
// 		t.Errorf("Not enough guesses correct fuzzy %f, %d, %f", correctFuzzy, len(files), totalPercentageFuzzy)
// 	}

// 	totalPercentageKeywords := (correctKeywords / float64(len(files))) * 100
// 	if totalPercentageKeywords < 0.95 {
// 		t.Errorf("Not enough guesses correct keywords %f, %d, %f", correctKeywords, len(files), totalPercentageKeywords)
// 	}
// }

func TestProcessFileLicensesTop10(t *testing.T) {
	// https://www.blackducksoftware.com/top-open-source-licenses
	files := []string{"MIT", "GPL-2.0", "Apache-2.0", "GPL-3.0-only", "ISC", "Artistic-2.0", "LGPL-2.1", "LGPL-3.0", "EPL-2.0", "MS-PL"}

	for _, file := range files {
		content, _ := ioutil.ReadFile(filepath.Join("../examples/licenses/", file+".json"))
		actual := keywordGuessLicense(string(content), loadDatabase())

		if len(actual) == 0 {
			t.Errorf("Expected some guesses for %s", file)
		}

		found := false
		for _, license := range actual {
			if license.LicenseId == file {
				found = true
			}
		}

		if found == false {
			t.Errorf("Expected license to be found %s", file)
		}
	}
}

// func TestRegressionIssue36(t *testing.T) {
// 	content := `Copyright (c) 1995-2002 Rik Faith <rikfaith@gmail.com>

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

// 	result := guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "MIT" {
// 		t.Errorf("Should be MIT")
// 	}

// 	content = `Copyright (c) 1995-2002 Rik Faith <rikfaith@gmail.com>
// Copyright (c) 2002-2018 Aleksey Cheusov <vle@gmx.net>

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

// 	result = guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "MIT" {
// 		t.Errorf("Should be MIT")
// 	}
// }

// func TestRegressionIssue41(t *testing.T) {
// 	content := `Copyright (C) 2010-2014 Jonas Borgström <jonas@borgstrom.se>
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:

//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//  3. The name of the author may not be used to endorse or promote
//     products derived from this software without specific prior
//     written permission.

// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS
// OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
// IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

// 	result := guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "BSD-2-Clause" {
// 		t.Errorf("Should be BSD-2-Clause was %s", result[0].LicenseId)
// 	}
// }

// func TestRegressionIssue40(t *testing.T) {
// 	// Ideas... trim out the copyright for the text
// 	content := `Copyright (c) 2014-2016 Lazaros Koromilas <lostd@2f30.org>
// Copyright (c) 2014-2016 Dimitris Papastamos <sin@2f30.org>
// Copyright (c) 2016-2018 Arun Prakash Jana <engineerarun@gmail.com>
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:

// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

// 	result := guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "BSD-2-Clause" {
// 		t.Errorf("Should be BSD-2-Clause was %s", result[0].LicenseId)
// 	}
// }

// func TestRegressionIssue39(t *testing.T) {
// 	content := `Copyright (c) 1989-1994
// 	The Regents of the University of California.  All rights reserved.
// Copyright (c) 1997 Christos Zoulas.  All rights reserved.
// Copyright (c) 1997-2005
// 	Herbert Xu <herbert@gondor.apana.org.au>.  All rights reserved.

// This code is derived from software contributed to Berkeley by Kenneth Almquist.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ''AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

// mksignames.c:

// This file is not directly linked with dash.  However, its output is.

// Copyright (C) 1992 Free Software Foundation, Inc.

// This file is part of GNU Bash, the Bourne Again SHell.

// Bash is free software; you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free
// Software Foundation; either version 2, or (at your option) any later
// version.

// Bash is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
// for more details.

// You should have received a copy of the GNU General Public License with
// your Debian GNU/Linux system, in /usr/share/common-licenses/GPL, or with the
// Debian GNU/Linux hello source package as the file COPYING.  If not,
// write to the Free Software Foundation, Inc., 59 Temple Place, Suite 330,
// Boston, MA 02111 USA.`

// 	result := guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "BSD-2-Clause" {
// 		t.Errorf("Should be BSD-2-Clause was %s", result[0].LicenseId)
// 	}
// }

// func TestRegressionIssue38(t *testing.T) {
// 	// Ideas... trimming the copyright works here
// 	content := `Copyright 2001-2009 Jean-Marc Valin, Timothy B. Terriberry,
//                     CSIRO, and other contributors

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:

// - Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.

// - Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// ''AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE FOUNDATION OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
// EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

// 	result := guessLicense(content, deepGuess, loadDatabase())

// 	if result[0].LicenseId != "BSD-2-Clause" {
// 		t.Errorf("Should be BSD-2-Clause was %s", result[0].LicenseId)
// 	}
// }

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
