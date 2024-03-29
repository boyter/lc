package processor

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLicenceDetector_DetectAll(t *testing.T) {
	l := NewLicenceDetector(true)

	for _, li := range spdxLicenseIds {
		detected := l.Detect(li, fmt.Sprintf("Valid-License-Identifier: %s", li))
		if detected[0].LicenseId != li {
			t.Errorf("expected %s got %s", li, detected[0].LicenseId)
		}
	}
}

func TestLicenceDetector_Detect(t *testing.T) {
	type fields struct {
		UseFullDatabase bool
	}
	type args struct {
		filename string
		content  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []IdentifiedLicense
	}{
		{
			name:   "spdx identifier gpl multiple",
			fields: fields{true},
			args: args{
				filename: "GPL-2.0",
				content:  spdxIdentifierGpl,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "GPL-2.0",
					ScorePercentage: 100,
				},
				{
					LicenseId:       "GPL-2.0-only",
					ScorePercentage: 100,
				},
				{
					LicenseId:       "GPL-2.0+",
					ScorePercentage: 100,
				},
				{
					LicenseId:       "GPL-2.0-or-later",
					ScorePercentage: 100,
				},
			},
		},
		{
			name:   "spdx identifier mit",
			fields: fields{true},
			args: args{
				filename: "MIT",
				content:  spdxIdentifierMit,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "MIT",
					ScorePercentage: 100,
				},
			},
		},
		{
			name:   "spdx identifier lgpl multiple",
			fields: fields{true},
			args: args{
				filename: "LGPL-2.1",
				content:  spdxIdentifierLgpl,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "LGPL-2.1",
					ScorePercentage: 100,
				},
				{
					LicenseId:       "LGPL-2.1+",
					ScorePercentage: 100,
				},
			},
		},
		{
			name:   "spdx identifier bsd3 singular",
			fields: fields{true},
			args: args{
				filename: "BSD-3-Clause",
				content:  spdxIdentifierBsd3,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "BSD-3-Clause",
					ScorePercentage: 100,
				},
			},
		},
		{
			name:   "spdx identifier bsd3 duplicate",
			fields: fields{true},
			args: args{
				filename: "BSD-3-Clause",
				content:  spdxIdentifierBsd3Duplicate,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "BSD-3-Clause",
					ScorePercentage: 100,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector(tt.fields.UseFullDatabase)
			if got := l.Detect(tt.args.filename, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

var spdxIdentifierGpl = "Valid-License-Identifier: GPL-2.0\nValid-License-Identifier: GPL-2.0-only\nValid-License-Identifier: GPL-2.0+\nValid-License-Identifier: GPL-2.0-or-later\nSPDX-URL: https://spdx.org/licenses/GPL-2.0.html\nUsage-Guide:\n  To use this license in source code, put one of the following SPDX\n  tag/value pairs into a comment according to the placement\n  guidelines in the licensing rules documentation.\n  For 'GNU General Public License (GPL) version 2 only' use:\n    SPDX-License-Identifier: GPL-2.0\n  or\n    SPDX-License-Identifier: GPL-2.0-only\n  For 'GNU General Public License (GPL) version 2 or any later version' use:\n    SPDX-License-Identifier: GPL-2.0+\n  or\n    SPDX-License-Identifier: GPL-2.0-or-later\nLicense-Text:\n\n\t\t    GNU GENERAL PUBLIC LICENSE\n\t\t       Version 2, June 1991\n\n Copyright (C) 1989, 1991 Free Software Foundation, Inc.\n                       51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA\n Everyone is permitted to copy and distribute verbatim copies\n of this license document, but changing it is not allowed.\n\n\t\t\t    Preamble\n\n  The licenses for most software are designed to take away your\nfreedom to share and change it.  By contrast, the GNU General Public\nLicense is intended to guarantee your freedom to share and change free\nsoftware--to make sure the software is free for all its users.  This\nGeneral Public License applies to most of the Free Software\nFoundation's software and to any other program whose authors commit to\nusing it.  (Some other Free Software Foundation software is covered by\nthe GNU Library General Public License instead.)  You can apply it to\nyour programs, too.\n\n  When we speak of free software, we are referring to freedom, not\nprice.  Our General Public Licenses are designed to make sure that you\nhave the freedom to distribute copies of free software (and charge for\nthis service if you wish), that you receive source code or can get it\nif you want it, that you can change the software or use pieces of it\nin new free programs; and that you know you can do these things.\n\n  To protect your rights, we need to make restrictions that forbid\nanyone to deny you these rights or to ask you to surrender the rights.\nThese restrictions translate to certain responsibilities for you if you\ndistribute copies of the software, or if you modify it.\n\n  For example, if you distribute copies of such a program, whether\ngratis or for a fee, you must give the recipients all the rights that\nyou have.  You must make sure that they, too, receive or can get the\nsource code.  And you must show them these terms so they know their\nrights.\n\n  We protect your rights with two steps: (1) copyright the software, and\n(2) offer you this license which gives you legal permission to copy,\ndistribute and/or modify the software.\n\n  Also, for each author's protection and ours, we want to make certain\nthat everyone understands that there is no warranty for this free\nsoftware.  If the software is modified by someone else and passed on, we\nwant its recipients to know that what they have is not the original, so\nthat any problems introduced by others will not reflect on the original\nauthors' reputations.\n\n  Finally, any free program is threatened constantly by software\npatents.  We wish to avoid the danger that redistributors of a free\nprogram will individually obtain patent licenses, in effect making the\nprogram proprietary.  To prevent this, we have made it clear that any\npatent must be licensed for everyone's free use or not licensed at all.\n\n  The precise terms and conditions for copying, distribution and\nmodification follow.\n\n\n\t\t    GNU GENERAL PUBLIC LICENSE\n   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION\n\n  0. This License applies to any program or other work which contains\na notice placed by the copyright holder saying it may be distributed\nunder the terms of this General Public License.  The \"Program\", below,\nrefers to any such program or work, and a \"work based on the Program\"\nmeans either the Program or any derivative work under copyright law:\nthat is to say, a work containing the Program or a portion of it,\neither verbatim or with modifications and/or translated into another\nlanguage.  (Hereinafter, translation is included without limitation in\nthe term \"modification\".)  Each licensee is addressed as \"you\".\n\nActivities other than copying, distribution and modification are not\ncovered by this License; they are outside its scope.  The act of\nrunning the Program is not restricted, and the output from the Program\nis covered only if its contents constitute a work based on the\nProgram (independent of having been made by running the Program).\nWhether that is true depends on what the Program does.\n\n  1. You may copy and distribute verbatim copies of the Program's\nsource code as you receive it, in any medium, provided that you\nconspicuously and appropriately publish on each copy an appropriate\ncopyright notice and disclaimer of warranty; keep intact all the\nnotices that refer to this License and to the absence of any warranty;\nand give any other recipients of the Program a copy of this License\nalong with the Program.\n\nYou may charge a fee for the physical act of transferring a copy, and\nyou may at your option offer warranty protection in exchange for a fee.\n\n  2. You may modify your copy or copies of the Program or any portion\nof it, thus forming a work based on the Program, and copy and\ndistribute such modifications or work under the terms of Section 1\nabove, provided that you also meet all of these conditions:\n\n    a) You must cause the modified files to carry prominent notices\n    stating that you changed the files and the date of any change.\n\n    b) You must cause any work that you distribute or publish, that in\n    whole or in part contains or is derived from the Program or any\n    part thereof, to be licensed as a whole at no charge to all third\n    parties under the terms of this License.\n\n    c) If the modified program normally reads commands interactively\n    when run, you must cause it, when started running for such\n    interactive use in the most ordinary way, to print or display an\n    announcement including an appropriate copyright notice and a\n    notice that there is no warranty (or else, saying that you provide\n    a warranty) and that users may redistribute the program under\n    these conditions, and telling the user how to view a copy of this\n    License.  (Exception: if the Program itself is interactive but\n    does not normally print such an announcement, your work based on\n    the Program is not required to print an announcement.)\n\n\nThese requirements apply to the modified work as a whole.  If\nidentifiable sections of that work are not derived from the Program,\nand can be reasonably considered independent and separate works in\nthemselves, then this License, and its terms, do not apply to those\nsections when you distribute them as separate works.  But when you\ndistribute the same sections as part of a whole which is a work based\non the Program, the distribution of the whole must be on the terms of\nthis License, whose permissions for other licensees extend to the\nentire whole, and thus to each and every part regardless of who wrote it.\n\nThus, it is not the intent of this section to claim rights or contest\nyour rights to work written entirely by you; rather, the intent is to\nexercise the right to control the distribution of derivative or\ncollective works based on the Program.\n\nIn addition, mere aggregation of another work not based on the Program\nwith the Program (or with a work based on the Program) on a volume of\na storage or distribution medium does not bring the other work under\nthe scope of this License.\n\n  3. You may copy and distribute the Program (or a work based on it,\nunder Section 2) in object code or executable form under the terms of\nSections 1 and 2 above provided that you also do one of the following:\n\n    a) Accompany it with the complete corresponding machine-readable\n    source code, which must be distributed under the terms of Sections\n    1 and 2 above on a medium customarily used for software interchange; or,\n\n    b) Accompany it with a written offer, valid for at least three\n    years, to give any third party, for a charge no more than your\n    cost of physically performing source distribution, a complete\n    machine-readable copy of the corresponding source code, to be\n    distributed under the terms of Sections 1 and 2 above on a medium\n    customarily used for software interchange; or,\n\n    c) Accompany it with the information you received as to the offer\n    to distribute corresponding source code.  (This alternative is\n    allowed only for noncommercial distribution and only if you\n    received the program in object code or executable form with such\n    an offer, in accord with Subsection b above.)\n\nThe source code for a work means the preferred form of the work for\nmaking modifications to it.  For an executable work, complete source\ncode means all the source code for all modules it contains, plus any\nassociated interface definition files, plus the scripts used to\ncontrol compilation and installation of the executable.  However, as a\nspecial exception, the source code distributed need not include\nanything that is normally distributed (in either source or binary\nform) with the major components (compiler, kernel, and so on) of the\noperating system on which the executable runs, unless that component\nitself accompanies the executable.\n\nIf distribution of executable or object code is made by offering\naccess to copy from a designated place, then offering equivalent\naccess to copy the source code from the same place counts as\ndistribution of the source code, even though third parties are not\ncompelled to copy the source along with the object code.\n\n\n  4. You may not copy, modify, sublicense, or distribute the Program\nexcept as expressly provided under this License.  Any attempt\notherwise to copy, modify, sublicense or distribute the Program is\nvoid, and will automatically terminate your rights under this License.\nHowever, parties who have received copies, or rights, from you under\nthis License will not have their licenses terminated so long as such\nparties remain in full compliance.\n\n  5. You are not required to accept this License, since you have not\nsigned it.  However, nothing else grants you permission to modify or\ndistribute the Program or its derivative works.  These actions are\nprohibited by law if you do not accept this License.  Therefore, by\nmodifying or distributing the Program (or any work based on the\nProgram), you indicate your acceptance of this License to do so, and\nall its terms and conditions for copying, distributing or modifying\nthe Program or works based on it.\n\n  6. Each time you redistribute the Program (or any work based on the\nProgram), the recipient automatically receives a license from the\noriginal licensor to copy, distribute or modify the Program subject to\nthese terms and conditions.  You may not impose any further\nrestrictions on the recipients' exercise of the rights granted herein.\nYou are not responsible for enforcing compliance by third parties to\nthis License.\n\n  7. If, as a consequence of a court judgment or allegation of patent\ninfringement or for any other reason (not limited to patent issues),\nconditions are imposed on you (whether by court order, agreement or\notherwise) that contradict the conditions of this License, they do not\nexcuse you from the conditions of this License.  If you cannot\ndistribute so as to satisfy simultaneously your obligations under this\nLicense and any other pertinent obligations, then as a consequence you\nmay not distribute the Program at all.  For example, if a patent\nlicense would not permit royalty-free redistribution of the Program by\nall those who receive copies directly or indirectly through you, then\nthe only way you could satisfy both it and this License would be to\nrefrain entirely from distribution of the Program.\n\nIf any portion of this section is held invalid or unenforceable under\nany particular circumstance, the balance of the section is intended to\napply and the section as a whole is intended to apply in other\ncircumstances.\n\nIt is not the purpose of this section to induce you to infringe any\npatents or other property right claims or to contest validity of any\nsuch claims; this section has the sole purpose of protecting the\nintegrity of the free software distribution system, which is\nimplemented by public license practices.  Many people have made\ngenerous contributions to the wide range of software distributed\nthrough that system in reliance on consistent application of that\nsystem; it is up to the author/donor to decide if he or she is willing\nto distribute software through any other system and a licensee cannot\nimpose that choice.\n\nThis section is intended to make thoroughly clear what is believed to\nbe a consequence of the rest of this License.\n\n\n  8. If the distribution and/or use of the Program is restricted in\ncertain countries either by patents or by copyrighted interfaces, the\noriginal copyright holder who places the Program under this License\nmay add an explicit geographical distribution limitation excluding\nthose countries, so that distribution is permitted only in or among\ncountries not thus excluded.  In such case, this License incorporates\nthe limitation as if written in the body of this License.\n\n  9. The Free Software Foundation may publish revised and/or new versions\nof the General Public License from time to time.  Such new versions will\nbe similar in spirit to the present version, but may differ in detail to\naddress new problems or concerns.\n\nEach version is given a distinguishing version number.  If the Program\nspecifies a version number of this License which applies to it and \"any\nlater version\", you have the option of following the terms and conditions\neither of that version or of any later version published by the Free\nSoftware Foundation.  If the Program does not specify a version number of\nthis License, you may choose any version ever published by the Free Software\nFoundation.\n\n  10. If you wish to incorporate parts of the Program into other free\nprograms whose distribution conditions are different, write to the author\nto ask for permission.  For software which is copyrighted by the Free\nSoftware Foundation, write to the Free Software Foundation; we sometimes\nmake exceptions for this.  Our decision will be guided by the two goals\nof preserving the free status of all derivatives of our free software and\nof promoting the sharing and reuse of software generally.\n\n\t\t\t    NO WARRANTY\n\n  11. BECAUSE THE PROGRAM IS LICENSED FREE OF CHARGE, THERE IS NO WARRANTY\nFOR THE PROGRAM, TO THE EXTENT PERMITTED BY APPLICABLE LAW.  EXCEPT WHEN\nOTHERWISE STATED IN WRITING THE COPYRIGHT HOLDERS AND/OR OTHER PARTIES\nPROVIDE THE PROGRAM \"AS IS\" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESSED\nOR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF\nMERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.  THE ENTIRE RISK AS\nTO THE QUALITY AND PERFORMANCE OF THE PROGRAM IS WITH YOU.  SHOULD THE\nPROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF ALL NECESSARY SERVICING,\nREPAIR OR CORRECTION.\n\n  12. IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING\nWILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MAY MODIFY AND/OR\nREDISTRIBUTE THE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES,\nINCLUDING ANY GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING\nOUT OF THE USE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED\nTO LOSS OF DATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY\nYOU OR THIRD PARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER\nPROGRAMS), EVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE\nPOSSIBILITY OF SUCH DAMAGES.\n\n\t\t     END OF TERMS AND CONDITIONS\n\n\n\t    How to Apply These Terms to Your New Programs\n\n  If you develop a new program, and you want it to be of the greatest\npossible use to the public, the best way to achieve this is to make it\nfree software which everyone can redistribute and change under these terms.\n\n  To do so, attach the following notices to the program.  It is safest\nto attach them to the start of each source file to most effectively\nconvey the exclusion of warranty; and each file should have at least\nthe \"copyright\" line and a pointer to where the full notice is found.\n\n    <one line to give the program's name and a brief idea of what it does.>\n    Copyright (C) <year>  <name of author>\n\n    This program is free software; you can redistribute it and/or modify\n    it under the terms of the GNU General Public License as published by\n    the Free Software Foundation; either version 2 of the License, or\n    (at your option) any later version.\n\n    This program is distributed in the hope that it will be useful,\n    but WITHOUT ANY WARRANTY; without even the implied warranty of\n    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n    GNU General Public License for more details.\n\n    You should have received a copy of the GNU General Public License\n    along with this program; if not, write to the Free Software\n    Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA\n\n\nAlso add information on how to contact you by electronic and paper mail.\n\nIf the program is interactive, make it output a short notice like this\nwhen it starts in an interactive mode:\n\n    Gnomovision version 69, Copyright (C) year name of author\n    Gnomovision comes with ABSOLUTELY NO WARRANTY; for details type `show w'.\n    This is free software, and you are welcome to redistribute it\n    under certain conditions; type `show c' for details.\n\nThe hypothetical commands `show w' and `show c' should show the appropriate\nparts of the General Public License.  Of course, the commands you use may\nbe called something other than `show w' and `show c'; they could even be\nmouse-clicks or menu items--whatever suits your program.\n\nYou should also get your employer (if you work as a programmer) or your\nschool, if any, to sign a \"copyright disclaimer\" for the program, if\nnecessary.  Here is a sample; alter the names:\n\n  Yoyodyne, Inc., hereby disclaims all copyright interest in the program\n  `Gnomovision' (which makes passes at compilers) written by James Hacker.\n\n  <signature of Ty Coon>, 1 April 1989\n  Ty Coon, President of Vice\n\nThis General Public License does not permit incorporating your program into\nproprietary programs.  If your program is a subroutine library, you may\nconsider it more useful to permit linking proprietary applications with the\nlibrary.  If this is what you want to do, use the GNU Library General\nPublic License instead of this License."
var spdxIdentifierMit = "Valid-License-Identifier: MIT\nSPDX-URL: https://spdx.org/licenses/MIT.html\nUsage-Guide:\n  To use the MIT License put the following SPDX tag/value pair into a\n  comment according to the placement guidelines in the licensing rules\n  documentation:\n    SPDX-License-Identifier: MIT\nLicense-Text:\n\nMIT License\n\nCopyright (c) <year> <copyright holders>\n\nPermission is hereby granted, free of charge, to any person obtaining a\ncopy of this software and associated documentation files (the \"Software\"),\nto deal in the Software without restriction, including without limitation\nthe rights to use, copy, modify, merge, publish, distribute, sublicense,\nand/or sell copies of the Software, and to permit persons to whom the\nSoftware is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in\nall copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING\nFROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER\nDEALINGS IN THE SOFTWARE."
var spdxIdentifierLgpl = "Valid-License-Identifier: LGPL-2.1\nValid-License-Identifier: LGPL-2.1+\nSPDX-URL: https://spdx.org/licenses/LGPL-2.1.html\nUsage-Guide:\n  To use this license in source code, put one of the following SPDX\n  tag/value pairs into a comment according to the placement\n  guidelines in the licensing rules documentation.\n  For 'GNU Lesser General Public License (LGPL) version 2.1 only' use:\n    SPDX-License-Identifier: LGPL-2.1\n  For 'GNU Lesser General Public License (LGPL) version 2.1 or any later\n  version' use:\n    SPDX-License-Identifier: LGPL-2.1+\nLicense-Text:\n\nGNU LESSER GENERAL PUBLIC LICENSE\nVersion 2.1, February 1999\n\nCopyright (C) 1991, 1999 Free Software Foundation, Inc.\n51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA\n\nEveryone is permitted to copy and distribute verbatim copies of this\nlicense document, but changing it is not allowed.\n\n[This is the first released version of the Lesser GPL. It also counts as\nthe successor of the GNU Library Public License, version 2, hence the\nversion number 2.1.]\n\nPreamble\n\nThe licenses for most software are designed to take away your freedom to\nshare and change it. By contrast, the GNU General Public Licenses are\nintended to guarantee your freedom to share and change free software--to\nmake sure the software is free for all its users.\n\nThis license, the Lesser General Public License, applies to some specially\ndesignated software packages--typically libraries--of the Free Software\nFoundation and other authors who decide to use it. You can use it too, but\nwe suggest you first think carefully about whether this license or the\nordinary General Public License is the better strategy to use in any\nparticular case, based on the explanations below.\n\nWhen we speak of free software, we are referring to freedom of use, not\nprice. Our General Public Licenses are designed to make sure that you have\nthe freedom to distribute copies of free software (and charge for this\nservice if you wish); that you receive source code or can get it if you\nwant it; that you can change the software and use pieces of it in new free\nprograms; and that you are informed that you can do these things.\n\nTo protect your rights, we need to make restrictions that forbid\ndistributors to deny you these rights or to ask you to surrender these\nrights. These restrictions translate to certain responsibilities for you if\nyou distribute copies of the library or if you modify it.\n\nFor example, if you distribute copies of the library, whether gratis or for\na fee, you must give the recipients all the rights that we gave you. You\nmust make sure that they, too, receive or can get the source code. If you\nlink other code with the library, you must provide complete object files to\nthe recipients, so that they can relink them with the library after making\nchanges to the library and recompiling it. And you must show them these\nterms so they know their rights.\n\nWe protect your rights with a two-step method: (1) we copyright the\nlibrary, and (2) we offer you this license, which gives you legal\npermission to copy, distribute and/or modify the library.\n\nTo protect each distributor, we want to make it very clear that there is no\nwarranty for the free library. Also, if the library is modified by someone\nelse and passed on, the recipients should know that what they have is not\nthe original version, so that the original author's reputation will not be\naffected by problems that might be introduced by others.\n\nFinally, software patents pose a constant threat to the existence of any\nfree program. We wish to make sure that a company cannot effectively\nrestrict the users of a free program by obtaining a restrictive license\nfrom a patent holder. Therefore, we insist that any patent license obtained\nfor a version of the library must be consistent with the full freedom of\nuse specified in this license.\n\nMost GNU software, including some libraries, is covered by the ordinary GNU\nGeneral Public License. This license, the GNU Lesser General Public\nLicense, applies to certain designated libraries, and is quite different\nfrom the ordinary General Public License. We use this license for certain\nlibraries in order to permit linking those libraries into non-free\nprograms.\n\nWhen a program is linked with a library, whether statically or using a\nshared library, the combination of the two is legally speaking a combined\nwork, a derivative of the original library. The ordinary General Public\nLicense therefore permits such linking only if the entire combination fits\nits criteria of freedom. The Lesser General Public License permits more lax\ncriteria for linking other code with the library.\n\nWe call this license the \"Lesser\" General Public License because it does\nLess to protect the user's freedom than the ordinary General Public\nLicense. It also provides other free software developers Less of an\nadvantage over competing non-free programs. These disadvantages are the\nreason we use the ordinary General Public License for many\nlibraries. However, the Lesser license provides advantages in certain\nspecial circumstances.\n\nFor example, on rare occasions, there may be a special need to encourage\nthe widest possible use of a certain library, so that it becomes a de-facto\nstandard. To achieve this, non-free programs must be allowed to use the\nlibrary. A more frequent case is that a free library does the same job as\nwidely used non-free libraries. In this case, there is little to gain by\nlimiting the free library to free software only, so we use the Lesser\nGeneral Public License.\n\nIn other cases, permission to use a particular library in non-free programs\nenables a greater number of people to use a large body of free\nsoftware. For example, permission to use the GNU C Library in non-free\nprograms enables many more people to use the whole GNU operating system, as\nwell as its variant, the GNU/Linux operating system.\n\nAlthough the Lesser General Public License is Less protective of the users'\nfreedom, it does ensure that the user of a program that is linked with the\nLibrary has the freedom and the wherewithal to run that program using a\nmodified version of the Library.\n\nThe precise terms and conditions for copying, distribution and modification\nfollow. Pay close attention to the difference between a \"work based on the\nlibrary\" and a \"work that uses the library\". The former contains code\nderived from the library, whereas the latter must be combined with the\nlibrary in order to run.\n\nTERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION\n\n0. This License Agreement applies to any software library or other program\n   which contains a notice placed by the copyright holder or other\n   authorized party saying it may be distributed under the terms of this\n   Lesser General Public License (also called \"this License\"). Each\n   licensee is addressed as \"you\".\n\n   A \"library\" means a collection of software functions and/or data\n   prepared so as to be conveniently linked with application programs\n   (which use some of those functions and data) to form executables.\n\n   The \"Library\", below, refers to any such software library or work which\n   has been distributed under these terms. A \"work based on the Library\"\n   means either the Library or any derivative work under copyright law:\n   that is to say, a work containing the Library or a portion of it, either\n   verbatim or with modifications and/or translated straightforwardly into\n   another language. (Hereinafter, translation is included without\n   limitation in the term \"modification\".)\n\n   \"Source code\" for a work means the preferred form of the work for making\n   modifications to it. For a library, complete source code means all the\n   source code for all modules it contains, plus any associated interface\n   definition files, plus the scripts used to control compilation and\n   installation of the library.\n\n    Activities other than copying, distribution and modification are not\n    covered by this License; they are outside its scope. The act of running\n    a program using the Library is not restricted, and output from such a\n    program is covered only if its contents constitute a work based on the\n    Library (independent of the use of the Library in a tool for writing\n    it). Whether that is true depends on what the Library does and what the\n    program that uses the Library does.\n\n1. You may copy and distribute verbatim copies of the Library's complete\n   source code as you receive it, in any medium, provided that you\n   conspicuously and appropriately publish on each copy an appropriate\n   copyright notice and disclaimer of warranty; keep intact all the notices\n   that refer to this License and to the absence of any warranty; and\n   distribute a copy of this License along with the Library.\n\n   You may charge a fee for the physical act of transferring a copy, and\n   you may at your option offer warranty protection in exchange for a fee.\n\n2. You may modify your copy or copies of the Library or any portion of it,\n   thus forming a work based on the Library, and copy and distribute such\n   modifications or work under the terms of Section 1 above, provided that\n   you also meet all of these conditions:\n\n   a) The modified work must itself be a software library.\n\n   b) You must cause the files modified to carry prominent notices stating\n      that you changed the files and the date of any change.\n\n   c) You must cause the whole of the work to be licensed at no charge to\n      all third parties under the terms of this License.\n\n   d) If a facility in the modified Library refers to a function or a table\n      of data to be supplied by an application program that uses the\n      facility, other than as an argument passed when the facility is\n      invoked, then you must make a good faith effort to ensure that, in\n      the event an application does not supply such function or table, the\n      facility still operates, and performs whatever part of its purpose\n      remains meaningful.\n\n   (For example, a function in a library to compute square roots has a\n    purpose that is entirely well-defined independent of the\n    application. Therefore, Subsection 2d requires that any\n    application-supplied function or table used by this function must be\n    optional: if the application does not supply it, the square root\n    function must still compute square roots.)\n\n   These requirements apply to the modified work as a whole. If\n   identifiable sections of that work are not derived from the Library, and\n   can be reasonably considered independent and separate works in\n   themselves, then this License, and its terms, do not apply to those\n   sections when you distribute them as separate works. But when you\n   distribute the same sections as part of a whole which is a work based on\n   the Library, the distribution of the whole must be on the terms of this\n   License, whose permissions for other licensees extend to the entire\n   whole, and thus to each and every part regardless of who wrote it.\n\n   Thus, it is not the intent of this section to claim rights or contest\n   your rights to work written entirely by you; rather, the intent is to\n   exercise the right to control the distribution of derivative or\n   collective works based on the Library.\n\n   In addition, mere aggregation of another work not based on the Library\n   with the Library (or with a work based on the Library) on a volume of a\n   storage or distribution medium does not bring the other work under the\n   scope of this License.\n\n3. You may opt to apply the terms of the ordinary GNU General Public\n   License instead of this License to a given copy of the Library. To do\n   this, you must alter all the notices that refer to this License, so that\n   they refer to the ordinary GNU General Public License, version 2,\n   instead of to this License. (If a newer version than version 2 of the\n   ordinary GNU General Public License has appeared, then you can specify\n   that version instead if you wish.) Do not make any other change in these\n   notices.\n\n   Once this change is made in a given copy, it is irreversible for that\n   copy, so the ordinary GNU General Public License applies to all\n   subsequent copies and derivative works made from that copy.\n\n   This option is useful when you wish to copy part of the code of the\n   Library into a program that is not a library.\n\n4. You may copy and distribute the Library (or a portion or derivative of\n   it, under Section 2) in object code or executable form under the terms\n   of Sections 1 and 2 above provided that you accompany it with the\n   complete corresponding machine-readable source code, which must be\n   distributed under the terms of Sections 1 and 2 above on a medium\n   customarily used for software interchange.\n\n   If distribution of object code is made by offering access to copy from a\n   designated place, then offering equivalent access to copy the source\n   code from the same place satisfies the requirement to distribute the\n   source code, even though third parties are not compelled to copy the\n   source along with the object code.\n\n5. A program that contains no derivative of any portion of the Library, but\n   is designed to work with the Library by being compiled or linked with\n   it, is called a \"work that uses the Library\". Such a work, in isolation,\n   is not a derivative work of the Library, and therefore falls outside the\n   scope of this License.\n\n   However, linking a \"work that uses the Library\" with the Library creates\n   an executable that is a derivative of the Library (because it contains\n   portions of the Library), rather than a \"work that uses the\n   library\". The executable is therefore covered by this License. Section 6\n   states terms for distribution of such executables.\n\n   When a \"work that uses the Library\" uses material from a header file\n   that is part of the Library, the object code for the work may be a\n   derivative work of the Library even though the source code is\n   not. Whether this is true is especially significant if the work can be\n   linked without the Library, or if the work is itself a library. The\n   threshold for this to be true is not precisely defined by law.\n\n   If such an object file uses only numerical parameters, data structure\n   layouts and accessors, and small macros and small inline functions (ten\n   lines or less in length), then the use of the object file is\n   unrestricted, regardless of whether it is legally a derivative\n   work. (Executables containing this object code plus portions of the\n   Library will still fall under Section 6.)\n\n   Otherwise, if the work is a derivative of the Library, you may\n   distribute the object code for the work under the terms of Section\n   6. Any executables containing that work also fall under Section 6,\n   whether or not they are linked directly with the Library itself.\n\n6. As an exception to the Sections above, you may also combine or link a\n   \"work that uses the Library\" with the Library to produce a work\n   containing portions of the Library, and distribute that work under terms\n   of your choice, provided that the terms permit modification of the work\n   for the customer's own use and reverse engineering for debugging such\n   modifications.\n\n   You must give prominent notice with each copy of the work that the\n   Library is used in it and that the Library and its use are covered by\n   this License. You must supply a copy of this License. If the work during\n   execution displays copyright notices, you must include the copyright\n   notice for the Library among them, as well as a reference directing the\n   user to the copy of this License. Also, you must do one of these things:\n\n   a) Accompany the work with the complete corresponding machine-readable\n      source code for the Library including whatever changes were used in\n      the work (which must be distributed under Sections 1 and 2 above);\n      and, if the work is an executable linked with the Library, with the\n      complete machine-readable \"work that uses the Library\", as object\n      code and/or source code, so that the user can modify the Library and\n      then relink to produce a modified executable containing the modified\n      Library. (It is understood that the user who changes the contents of\n      definitions files in the Library will not necessarily be able to\n      recompile the application to use the modified definitions.)\n\n   b) Use a suitable shared library mechanism for linking with the\n      Library. A suitable mechanism is one that (1) uses at run time a copy\n      of the library already present on the user's computer system, rather\n      than copying library functions into the executable, and (2) will\n      operate properly with a modified version of the library, if the user\n      installs one, as long as the modified version is interface-compatible\n      with the version that the work was made with.\n\n   c) Accompany the work with a written offer, valid for at least three\n      years, to give the same user the materials specified in Subsection\n      6a, above, for a charge no more than the cost of performing this\n      distribution.\n\n   d) If distribution of the work is made by offering access to copy from a\n      designated place, offer equivalent access to copy the above specified\n      materials from the same place.\n\n   e) Verify that the user has already received a copy of these materials\n      or that you have already sent this user a copy.\n\n   For an executable, the required form of the \"work that uses the Library\"\n   must include any data and utility programs needed for reproducing the\n   executable from it. However, as a special exception, the materials to be\n   distributed need not include anything that is normally distributed (in\n   either source or binary form) with the major components (compiler,\n   kernel, and so on) of the operating system on which the executable runs,\n   unless that component itself accompanies the executable.\n\n   It may happen that this requirement contradicts the license restrictions\n   of other proprietary libraries that do not normally accompany the\n   operating system. Such a contradiction means you cannot use both them\n   and the Library together in an executable that you distribute.\n\n7. You may place library facilities that are a work based on the Library\n   side-by-side in a single library together with other library facilities\n   not covered by this License, and distribute such a combined library,\n   provided that the separate distribution of the work based on the Library\n   and of the other library facilities is otherwise permitted, and provided\n   that you do these two things:\n\n   a) Accompany the combined library with a copy of the same work based on\n      the Library, uncombined with any other library facilities. This must\n      be distributed under the terms of the Sections above.\n\n   b) Give prominent notice with the combined library of the fact that part\n      of it is a work based on the Library, and explaining where to find\n      the accompanying uncombined form of the same work.\n\n8. You may not copy, modify, sublicense, link with, or distribute the\n   Library except as expressly provided under this License. Any attempt\n   otherwise to copy, modify, sublicense, link with, or distribute the\n   Library is void, and will automatically terminate your rights under this\n   License. However, parties who have received copies, or rights, from you\n   under this License will not have their licenses terminated so long as\n   such parties remain in full compliance.\n\n9. You are not required to accept this License, since you have not signed\n   it. However, nothing else grants you permission to modify or distribute\n   the Library or its derivative works. These actions are prohibited by law\n   if you do not accept this License. Therefore, by modifying or\n   distributing the Library (or any work based on the Library), you\n   indicate your acceptance of this License to do so, and all its terms and\n   conditions for copying, distributing or modifying the Library or works\n   based on it.\n\n10. Each time you redistribute the Library (or any work based on the\n    Library), the recipient automatically receives a license from the\n    original licensor to copy, distribute, link with or modify the Library\n    subject to these terms and conditions. You may not impose any further\n    restrictions on the recipients' exercise of the rights granted\n    herein. You are not responsible for enforcing compliance by third\n    parties with this License.\n\n11. If, as a consequence of a court judgment or allegation of patent\n    infringement or for any other reason (not limited to patent issues),\n    conditions are imposed on you (whether by court order, agreement or\n    otherwise) that contradict the conditions of this License, they do not\n    excuse you from the conditions of this License. If you cannot\n    distribute so as to satisfy simultaneously your obligations under this\n    License and any other pertinent obligations, then as a consequence you\n    may not distribute the Library at all. For example, if a patent license\n    would not permit royalty-free redistribution of the Library by all\n    those who receive copies directly or indirectly through you, then the\n    only way you could satisfy both it and this License would be to refrain\n    entirely from distribution of the Library.\n\n    If any portion of this section is held invalid or unenforceable under\n    any particular circumstance, the balance of the section is intended to\n    apply, and the section as a whole is intended to apply in other\n    circumstances.\n\n    It is not the purpose of this section to induce you to infringe any\n    patents or other property right claims or to contest validity of any\n    such claims; this section has the sole purpose of protecting the\n    integrity of the free software distribution system which is implemented\n    by public license practices. Many people have made generous\n    contributions to the wide range of software distributed through that\n    system in reliance on consistent application of that system; it is up\n    to the author/donor to decide if he or she is willing to distribute\n    software through any other system and a licensee cannot impose that\n    choice.\n\n    This section is intended to make thoroughly clear what is believed to\n    be a consequence of the rest of this License.\n\n12. If the distribution and/or use of the Library is restricted in certain\n    countries either by patents or by copyrighted interfaces, the original\n    copyright holder who places the Library under this License may add an\n    explicit geographical distribution limitation excluding those\n    countries, so that distribution is permitted only in or among countries\n    not thus excluded. In such case, this License incorporates the\n    limitation as if written in the body of this License.\n\n13. The Free Software Foundation may publish revised and/or new versions of\n    the Lesser General Public License from time to time. Such new versions\n    will be similar in spirit to the present version, but may differ in\n    detail to address new problems or concerns.\n\n    Each version is given a distinguishing version number. If the Library\n    specifies a version number of this License which applies to it and \"any\n    later version\", you have the option of following the terms and\n    conditions either of that version or of any later version published by\n    the Free Software Foundation. If the Library does not specify a license\n    version number, you may choose any version ever published by the Free\n    Software Foundation.\n\n14. If you wish to incorporate parts of the Library into other free\n    programs whose distribution conditions are incompatible with these,\n    write to the author to ask for permission. For software which is\n    copyrighted by the Free Software Foundation, write to the Free Software\n    Foundation; we sometimes make exceptions for this. Our decision will be\n    guided by the two goals of preserving the free status of all\n    derivatives of our free software and of promoting the sharing and reuse\n    of software generally.\n\nNO WARRANTY\n\n15. BECAUSE THE LIBRARY IS LICENSED FREE OF CHARGE, THERE IS NO WARRANTY\n    FOR THE LIBRARY, TO THE EXTENT PERMITTED BY APPLICABLE LAW. EXCEPT WHEN\n    OTHERWISE STATED IN WRITING THE COPYRIGHT HOLDERS AND/OR OTHER PARTIES\n    PROVIDE THE LIBRARY \"AS IS\" WITHOUT WARRANTY OF ANY KIND, EITHER\n    EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED\n    WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE. THE\n    ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE LIBRARY IS WITH\n    YOU. SHOULD THE LIBRARY PROVE DEFECTIVE, YOU ASSUME THE COST OF ALL\n    NECESSARY SERVICING, REPAIR OR CORRECTION.\n\n16. IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING\n    WILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MAY MODIFY AND/OR\n    REDISTRIBUTE THE LIBRARY AS PERMITTED ABOVE, BE LIABLE TO YOU FOR\n    DAMAGES, INCLUDING ANY GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL\n    DAMAGES ARISING OUT OF THE USE OR INABILITY TO USE THE LIBRARY\n    (INCLUDING BUT NOT LIMITED TO LOSS OF DATA OR DATA BEING RENDERED\n    INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD PARTIES OR A FAILURE OF\n    THE LIBRARY TO OPERATE WITH ANY OTHER SOFTWARE), EVEN IF SUCH HOLDER OR\n    OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF SUCH DAMAGES.\n\nEND OF TERMS AND CONDITIONS\n\nHow to Apply These Terms to Your New Libraries\n\nIf you develop a new library, and you want it to be of the greatest\npossible use to the public, we recommend making it free software that\neveryone can redistribute and change. You can do so by permitting\nredistribution under these terms (or, alternatively, under the terms of the\nordinary General Public License).\n\nTo apply these terms, attach the following notices to the library. It is\nsafest to attach them to the start of each source file to most effectively\nconvey the exclusion of warranty; and each file should have at least the\n\"copyright\" line and a pointer to where the full notice is found.\n\none line to give the library's name and an idea of what it does.\nCopyright (C) year name of author\n\nThis library is free software; you can redistribute it and/or modify it\nunder the terms of the GNU Lesser General Public License as published by\nthe Free Software Foundation; either version 2.1 of the License, or (at\nyour option) any later version.\n\nThis library is distributed in the hope that it will be useful, but WITHOUT\nANY WARRANTY; without even the implied warranty of MERCHANTABILITY or\nFITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License\nfor more details.\n\nYou should have received a copy of the GNU Lesser General Public License\nalong with this library; if not, write to the Free Software Foundation,\nInc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA Also add\ninformation on how to contact you by electronic and paper mail.\n\nYou should also get your employer (if you work as a programmer) or your\nschool, if any, to sign a \"copyright disclaimer\" for the library, if\nnecessary. Here is a sample; alter the names:\n\nYoyodyne, Inc., hereby disclaims all copyright interest in\nthe library `Frob' (a library for tweaking knobs) written\nby James Random Hacker.\n\nsignature of Ty Coon, 1 April 1990\nTy Coon, President of Vice\nThat's all there is to it!"
var spdxIdentifierBsd3 = "Valid-License-Identifier: BSD-3-Clause\nSPDX-URL: https://spdx.org/licenses/BSD-3-Clause.html\nUsage-Guide:\n  To use the BSD 3-clause \"New\" or \"Revised\" License put the following SPDX\n  tag/value pair into a comment according to the placement guidelines in\n  the licensing rules documentation:\n    SPDX-License-Identifier: BSD-3-Clause\nLicense-Text:\n\nCopyright (c) <year> <owner> . All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are met:\n\n1. Redistributions of source code must retain the above copyright notice,\n   this list of conditions and the following disclaimer.\n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in the\n   documentation and/or other materials provided with the distribution.\n\n3. Neither the name of the copyright holder nor the names of its\n   contributors may be used to endorse or promote products derived from this\n   software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\"\nAND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE\nIMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE\nARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE\nLIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR\nCONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF\nSUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS\nINTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN\nCONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)\nARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE\nPOSSIBILITY OF SUCH DAMAGE."
var spdxIdentifierBsd3Duplicate = "Valid-License-Identifier: BSD-3-Clause\nValid-License-Identifier: BSD-3-Clause\nValid-License-Identifier: BSD-3-Clause\nSPDX-URL: https://spdx.org/licenses/BSD-3-Clause.html\nUsage-Guide:\n  To use the BSD 3-clause \"New\" or \"Revised\" License put the following SPDX\n  tag/value pair into a comment according to the placement guidelines in\n  the licensing rules documentation:\n    SPDX-License-Identifier: BSD-3-Clause\nLicense-Text:\n\nCopyright (c) <year> <owner> . All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are met:\n\n1. Redistributions of source code must retain the above copyright notice,\n   this list of conditions and the following disclaimer.\n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in the\n   documentation and/or other materials provided with the distribution.\n\n3. Neither the name of the copyright holder nor the names of its\n   contributors may be used to endorse or promote products derived from this\n   software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\"\nAND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE\nIMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE\nARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE\nLIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR\nCONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF\nSUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS\nINTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN\nCONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)\nARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE\nPOSSIBILITY OF SUCH DAMAGE."

func TestLicenceDetector_DetectFilename(t *testing.T) {
	type fields struct {
		UseFullDatabase bool
	}
	type args struct {
		filename string
		content  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []IdentifiedLicense
	}{
		{
			name:   "filename spdx identifier mit",
			fields: fields{true},
			args: args{
				filename: "MIT",
				content:  mitLicense,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "MIT",
					ScorePercentage: 100,
				},
			},
		},
		{
			name:   "filename spdx identifier mit modified",
			fields: fields{true},
			args: args{
				filename: "MIT",
				content:  mitLicense2,
			},
			want: []IdentifiedLicense{
				{
					LicenseId:       "MIT",
					ScorePercentage: 99.60854968169237,
				},
			},
		},
		{
			name:   "filename spdx identifier gpl-2.0 wrong",
			fields: fields{true},
			args: args{
				filename: "GPL-2.0",
				content:  mitLicense2,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector(tt.fields.UseFullDatabase)
			if got := l.Detect(tt.args.filename, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

var mitLicense = " MIT License\n\nCopyright (c) 2018 Ben Boyter\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE."
var mitLicense2 = " MIT License\n\nCopyright (c) 2023 Some Dude\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE."

func TestLicenceDetector_keywordDetect(t *testing.T) {
	type args struct {
		content string
	}
	var tests = []struct {
		name string
		args args
		want []IdentifiedLicense
	}{
		{
			name: "MIT",
			args: args{
				content: mitLicense,
			},
			want: []IdentifiedLicense{{
				LicenseId:       "MIT",
				ScorePercentage: 17,
			}},
		},
		{
			name: "LGPL",
			args: args{
				content: spdxIdentifierLgpl,
			},
			want: []IdentifiedLicense{{
				LicenseId:       "LGPL-2.1-or-later",
				ScorePercentage: 200,
			}, {
				LicenseId:       "LGPL-2.1-only",
				ScorePercentage: 200,
			}, {
				LicenseId:       "LGPL-2.1+",
				ScorePercentage: 200,
			}, {
				LicenseId:       "LGPL-2.1",
				ScorePercentage: 200,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector(true)
			if got := l.keywordDetect(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keywordDetect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLicenceDetector_vectorDetect(t *testing.T) {
	type args struct {
		content string
	}
	var tests = []struct {
		name string
		args args
		want IdentifiedLicense
	}{
		{
			name: "MIT",
			args: args{
				content: mitLicense2,
			},
			want: IdentifiedLicense{
				LicenseId: "MIT",
			},
		},
		{
			name: "BSD-3-Clause",
			args: args{
				content: spdxIdentifierBsd3,
			},
			want: IdentifiedLicense{
				LicenseId: "BSD-3-Clause",
			},
		},
		{
			name: "GPL-2.0",
			args: args{
				content: spdxIdentifierGpl,
			},
			want: IdentifiedLicense{
				LicenseId: "GPL-2.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector(true)
			if got := l.vectorDetect(tt.args.content); got[0].LicenseId != tt.want.LicenseId {
				t.Errorf("vectorDetect() = %v, want %v", got[0].LicenseId, tt.want.LicenseId)
			}
		})
	}
}

func TestLicenceDetector_levenshteinDetect(t *testing.T) {
	type args struct {
		content string
	}
	var tests = []struct {
		name string
		args args
		want []IdentifiedLicense
	}{
		{
			name: "MIT",
			args: args{
				content: mitLicense2,
			},
			want: []IdentifiedLicense{{
				LicenseId:       "MIT",
				ScorePercentage: 17,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector(true)
			if got := l.levenshteinDetect(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("levenshteinDetect() = %v, want %v", got, tt.want)
			}
		})
	}
}
