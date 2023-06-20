// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"testing"
)

// TODO does not work because not aware of related licences
func TestSpdxGuesser(t *testing.T) {
	lg := SpdxIdentifier{}

	actual := lg.SpdxIdentify("test")
	if len(actual) != 0 {
		t.Errorf("Should be no matches")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0")
	if len(actual) == 0 || actual[0] != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("/* SPDX-License-Identifier: GPL-2.0 */")
	if len(actual) == 0 || actual[0] != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0 ")
	if len(actual) == 0 || actual[0] != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0 \n # SPDX-License-Identifier: GPL-3.0+")
	if len(actual) == 0 || actual[0] != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if len(actual) == 0 || actual[1] != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}

	actual = lg.SpdxIdentify(`# SPDX-License-Identifier: GPL-2.0

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
`)
	if len(actual) == 0 || actual[0] != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if len(actual) == 0 || actual[1] != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}
}

func TestSpdxGuesserMultipleOr(t *testing.T) {
	lg := SpdxIdentifier{}

	actual := lg.SpdxIdentify("# SPDX-License-Identifier: MIT OR Unlicense")
	if actual[0] != "MIT" {
		t.Errorf("Should match MIT")
	}
	if actual[1] != "Unlicense" {
		t.Errorf("Should match Unlicense")
	}
}

func TestSpdxGuesserMultipleAnd(t *testing.T) {
	lg := SpdxIdentifier{}

	actual := lg.SpdxIdentify("# SPDX-License-Identifier: MIT AND Unlicense")
	if actual[0] != "MIT" {
		t.Errorf("Should match MIT")
	}
	if actual[1] != "Unlicense" {
		t.Errorf("Should match Unlicense")
	}
}

func TestSpdxGuesserMultipleLowerCase(t *testing.T) {
	lg := SpdxIdentifier{}

	actual := lg.SpdxIdentify("# SPDX-License-Identifier: mit unlicense gpl-2.0")
	if actual[0] != "MIT" {
		t.Error("Should match MIT got", actual[0])
	}
	if actual[1] != "Unlicense" {
		t.Error("Should match Unlicense got", actual[1])
	}
	if actual[2] != "GPL-2.0" {
		t.Error("Should match GPL-2.0 got", actual[2])
	}
}

func TestSpdxGuesserDuplicates(t *testing.T) {
	lg := SpdxIdentifier{}

	actual := lg.SpdxIdentify(`# SPDX-License-Identifier: mit
# SPDX-License-Identifier: mit`)

	if len(actual) != 1 {
		t.Error("should only get 1 got", len(actual))
	}
}
