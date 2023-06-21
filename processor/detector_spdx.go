// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"regexp"
	"strings"
)

var fileLicenceIdentifier = "SPDX-License-Identifier:"
var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)

// these were generated based on the full list from SPDX
var spdxLicenseIds = []string{"0BSD", "AAL", "ADSL", "AFL-1.1", "AFL-1.2", "AFL-2.0", "AFL-2.1", "AFL-3.0", "AGPL-1.0-only", "AGPL-1.0-or-later", "AGPL-1.0", "AGPL-3.0-only", "AGPL-3.0-or-later", "AGPL-3.0", "AMDPLPA", "AML", "AMPAS", "ANTLR-PD", "APAFML", "APL-1.0", "APSL-1.0", "APSL-1.1", "APSL-1.2", "APSL-2.0", "Abstyles", "Adobe-2006", "Adobe-Glyph", "Afmparse", "Aladdin", "Apache-1.0", "Apache-1.1", "Apache-2.0", "Artistic-1.0-Perl", "Artistic-1.0-cl8", "Artistic-1.0", "Artistic-2.0", "BSD-1-Clause", "BSD-2-Clause-FreeBSD", "BSD-2-Clause-NetBSD", "BSD-2-Clause-Patent", "BSD-2-Clause", "BSD-3-Clause-Attribution", "BSD-3-Clause-Clear", "BSD-3-Clause-LBNL", "BSD-3-Clause-No-Nuclear-License-2014", "BSD-3-Clause-No-Nuclear-License", "BSD-3-Clause-No-Nuclear-Warranty", "BSD-3-Clause-Open-MPI", "BSD-3-Clause", "BSD-4-Clause-UC", "BSD-4-Clause", "BSD-Protection", "BSD-Source-Code", "BSL-1.0", "Bahyph", "Barr", "Beerware", "BitTorrent-1.0", "BitTorrent-1.1", "BlueOak-1.0.0", "Borceux", "CAL-1.0-Combined-Work-Exception", "CAL-1.0", "CATOSL-1.1", "CC-BY-1.0", "CC-BY-2.0", "CC-BY-2.5", "CC-BY-3.0", "CC-BY-4.0", "CC-BY-NC-1.0", "CC-BY-NC-2.0", "CC-BY-NC-2.5", "CC-BY-NC-3.0", "CC-BY-NC-4.0", "CC-BY-NC-ND-1.0", "CC-BY-NC-ND-2.0", "CC-BY-NC-ND-2.5", "CC-BY-NC-ND-3.0", "CC-BY-NC-ND-4.0", "CC-BY-NC-SA-1.0", "CC-BY-NC-SA-2.0", "CC-BY-NC-SA-2.5", "CC-BY-NC-SA-3.0", "CC-BY-NC-SA-4.0", "CC-BY-ND-1.0", "CC-BY-ND-2.0", "CC-BY-ND-2.5", "CC-BY-ND-3.0", "CC-BY-ND-4.0", "CC-BY-SA-1.0", "CC-BY-SA-2.0", "CC-BY-SA-2.5", "CC-BY-SA-3.0", "CC-BY-SA-4.0", "CC-PDDC", "CC0-1.0", "CDDL-1.0", "CDDL-1.1", "CDLA-Permissive-1.0", "CDLA-Sharing-1.0", "CECILL-1.0", "CECILL-1.1", "CECILL-2.0", "CECILL-2.1", "CECILL-B", "CECILL-C", "CERN-OHL-1.1", "CERN-OHL-1.2", "CERN-OHL-P-2.0", "CERN-OHL-S-2.0", "CERN-OHL-W-2.0", "CNRI-Jython", "CNRI-Python-GPL-Compatible", "CNRI-Python", "CPAL-1.0", "CPL-1.0", "CPOL-1.02", "CUA-OPL-1.0", "Caldera", "ClArtistic", "Condor-1.1", "Crossword", "CrystalStacker", "Cube", "D-FSL-1.0", "DOC", "DSDP", "Dotseqn", "ECL-1.0", "ECL-2.0", "EFL-1.0", "EFL-2.0", "EPL-1.0", "EPL-2.0", "EUDatagrid", "EUPL-1.0", "EUPL-1.1", "EUPL-1.2", "Entessa", "ErlPL-1.1", "Eurosym", "FSFAP", "FSFUL", "FSFULLR", "FTL", "Fair-Source-0.9", "Fair", "Frameworx-1.0", "FreeImage", "GFDL-1.1-only", "GFDL-1.1-or-later", "GFDL-1.1", "GFDL-1.2-only", "GFDL-1.2-or-later", "GFDL-1.2", "GFDL-1.3-only", "GFDL-1.3-or-later", "GFDL-1.3", "GL2PS", "GPL-1.0+", "GPL-1.0-only", "GPL-1.0-or-later", "GPL-1.0", "GPL-2.0+", "GPL-2.0-only", "GPL-2.0-or-later", "GPL-2.0-with-GCC-exception", "GPL-2.0-with-autoconf-exception", "GPL-2.0-with-bison-exception", "GPL-2.0-with-classpath-exception", "GPL-2.0-with-font-exception", "GPL-2.0", "GPL-3.0+", "GPL-3.0-only", "GPL-3.0-or-later", "GPL-3.0-with-GCC-exception", "GPL-3.0-with-autoconf-exception", "GPL-3.0", "Giftware", "Glide", "Glulxe", "HPND-sell-variant", "HPND", "HaskellReport", "Hippocratic-2.1", "IBM-pibs", "ICU", "IJG", "IPA", "IPL-1.0", "ISC", "ImageMagick", "Imlib2", "Info-ZIP", "Intel-ACPI", "Intel", "Interbase-1.0", "JPNIC", "JSON", "JasPer-2.0", "LAL-1.2", "LAL-1.3", "LGPL-2.0+", "LGPL-2.0-only", "LGPL-2.0-or-later", "LGPL-2.0", "LGPL-2.1+", "LGPL-2.1-only", "LGPL-2.1-or-later", "LGPL-2.1", "LGPL-3.0+", "LGPL-3.0-only", "LGPL-3.0-or-later", "LGPL-3.0", "LGPLLR", "LPL-1.0", "LPL-1.02", "LPPL-1.0", "LPPL-1.1", "LPPL-1.2", "LPPL-1.3a", "LPPL-1.3c", "Latex2e", "Leptonica", "LiLiQ-P-1.1", "LiLiQ-R-1.1", "LiLiQ-Rplus-1.1", "Libpng", "Linux-OpenIB", "MIT-0", "MIT-CMU", "MIT-advertising", "MIT-enna", "MIT-feh", "MIT", "MITNFA", "MPL-1.0", "MPL-1.1", "MPL-2.0-no-copyleft-exception", "MPL-2.0", "MS-PL", "MS-RL", "MTLL", "MakeIndex", "MirOS", "Motosoto", "MulanPSL-1.0", "MulanPSL-2.0", "Multics", "Mup", "NASA-1.3", "NBPL-1.0", "NCGL-UK-2.0", "NCSA", "NGPL", "NLOD-1.0", "NLPL", "NOSL", "NPL-1.0", "NPL-1.1", "NPOSL-3.0", "NRL", "NTP-0", "NTP", "Naumen", "Net-SNMP", "NetCDF", "Newsletr", "Nokia", "Noweb", "Nunit", "O-UDA-1.0", "OCCT-PL", "OCLC-2.0", "ODC-By-1.0", "ODbL-1.0", "OFL-1.0-RFN", "OFL-1.0-no-RFN", "OFL-1.0", "OFL-1.1-RFN", "OFL-1.1-no-RFN", "OFL-1.1", "OGC-1.0", "OGL-Canada-2.0", "OGL-UK-1.0", "OGL-UK-2.0", "OGL-UK-3.0", "OGTSL", "OLDAP-1.1", "OLDAP-1.2", "OLDAP-1.3", "OLDAP-1.4", "OLDAP-2.0.1", "OLDAP-2.0", "OLDAP-2.1", "OLDAP-2.2.1", "OLDAP-2.2.2", "OLDAP-2.2", "OLDAP-2.3", "OLDAP-2.4", "OLDAP-2.5", "OLDAP-2.6", "OLDAP-2.7", "OLDAP-2.8", "OML", "OPL-1.0", "OSET-PL-2.1", "OSL-1.0", "OSL-1.1", "OSL-2.0", "OSL-2.1", "OSL-3.0", "OpenSSL", "PDDL-1.0", "PHP-3.0", "PHP-3.01", "PSF-2.0", "Parity-6.0.0", "Parity-7.0.0", "Plexus", "PolyForm-Noncommercial-1.0.0", "PolyForm-Small-Business-1.0.0", "PostgreSQL", "Python-2.0", "QPL-1.0", "Qhull", "RHeCos-1.1", "RPL-1.1", "RPL-1.5", "RPSL-1.0", "RSA-MD", "RSCPL", "Rdisc", "Ruby", "SAX-PD", "SCEA", "SGI-B-1.0", "SGI-B-1.1", "SGI-B-2.0", "SHL-0.5", "SHL-0.51", "SISSL-1.2", "SISSL", "SMLNJ", "SMPPL", "SNIA", "SPL-1.0", "SSH-OpenSSH", "SSH-short", "SSPL-1.0", "SWL", "Saxpath", "Sendmail-8.23", "Sendmail", "SimPL-2.0", "Sleepycat", "Spencer-86", "Spencer-94", "Spencer-99", "StandardML-NJ", "SugarCRM-1.1.3", "TAPR-OHL-1.0", "TCL", "TCP-wrappers", "TMate", "TORQUE-1.1", "TOSL", "TU-Berlin-1.0", "TU-Berlin-2.0", "UCL-1.0", "UPL-1.0", "Unicode-DFS-2015", "Unicode-DFS-2016", "Unicode-TOU", "Unlicense", "VOSTROM", "VSL-1.0", "Vim", "W3C-19980720", "W3C-20150513", "W3C", "WTFPL", "wxWindows", "Watcom-1.0", "Wsuipa", "X11", "XFree86-1.1", "XSkat", "Xerox", "Xnet", "YPL-1.0", "YPL-1.1", "ZPL-1.1", "ZPL-2.0", "ZPL-2.1", "Zed", "Zend-2.0", "Zimbra-1.3", "Zimbra-1.4", "Zlib", "blessing", "bzip2-1.0.5", "bzip2-1.0.6", "copyleft-next-0.3.0", "copyleft-next-0.3.1", "curl", "diffmark", "dvipdfm", "eCos-2.0", "eGenix", "etalab-2.0", "gSOAP-1.3b", "gnuplot", "iMatix", "libpng-2.0", "libselinux-1.0", "libtiff", "mpich2", "psfrag", "psutils", "xinetd", "xpp", "zlib-acknowledgement"}

type SpdxDetector struct{}

// SpdxDetect will identify licenses in the text which are using the SPDX indicator
// which is reasonably cheap in terms of looking things up
func (l *SpdxDetector) SpdxDetect(content string) []string {
	// cheap check to see if there might be on in the source code
	if strings.Index(content, fileLicenceIdentifier) == -1 {
		return nil
	}

	var matchingLicenses []string
	matches := spdxLicenceRegex.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var toCheck []string
		t := strings.TrimSpace(val[1])
		if strings.Contains(val[1], " ") {
			// deal with multiple with an OR or some such
			for _, x := range strings.Split(t, " ") {
				x = strings.TrimSpace(x)
				if x != "" {
					toCheck = append(toCheck, x)
				}
			}
		} else {
			toCheck = append(toCheck, t)
		}

		for _, x := range toCheck {
			found := false
			// Check the full database because there is so little cost to do so
			for _, license := range spdxLicenseIds {
				if license == x {
					matchingLicenses = append(matchingLicenses, license)
					found = true
					// we should only ever find a single per what we are checking
					break
				}
			}

			// if we didn't find anything try using lower case because hey why not
			if !found {
				x = strings.ToLower(x)
				for _, license := range spdxLicenseIds {
					if strings.ToLower(license) == x {
						matchingLicenses = append(matchingLicenses, license)
						// we should only ever find a single per what we are checking
						break
					}
				}
			}
		}
	}

	// filter out duplicates because its possible someone put in multiple markers of the same
	var found = map[string]bool{}
	var filtered []string

	for _, lic := range matchingLicenses {
		b := found[lic]
		if !b {
			filtered = append(filtered, lic)
			found[lic] = true
		}
	}

	return filtered
}
