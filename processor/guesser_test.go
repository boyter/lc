package processor

import "testing"

func TestSpdxGuesser(t *testing.T) {
	lg := NewLicenceGuesser(false, false)

	actual := lg.SpdxIdentify("test")
	if len(actual) != 0 {
		t.Errorf("Should be no matches")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0")
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("/* SPDX-License-Identifier: GPL-2.0 */")
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0 ")
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = lg.SpdxIdentify("# SPDX-License-Identifier: GPL-2.0 \n # SPDX-License-Identifier: GPL-3.0+")
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
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
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}
}
