// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense

package processor

import "testing"

func TestVectorSpaceGuess0BSD(t *testing.T) {
	lg := NewLicenceGuesser(false, true)

	licenses := lg.VectorSpaceGuessLicence([]byte(`{
  "isDeprecatedLicenseId": false,
  "licenseText": "Copyright (C) 2006 by Rob Landley \u003crob@landley.net\u003e\nPermission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted.\nTHE SOFTWARE IS PROVIDED \"AS IS\" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.",
  "standardLicenseTemplate": "Copyright (C) 2006 by Rob Landley \u003crob@landley.net\u003e\nPermission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted.\nTHE SOFTWARE IS PROVIDED \"AS IS\" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.",
  "name": "BSD Zero Clause License",
  "licenseId": "0BSD",
  "seeAlso": [
    "http://landley.net/toybox/license.html"
  ],
  "isOsiApproved": false
}`))

	if licenses[0].LicenseId != "0BSD" {
		t.Errorf("expected 0BSD")
	}
}

func TestVectorSpaceGuessMIT(t *testing.T) {
	lg := NewLicenceGuesser(false, true)

	licenses := lg.VectorSpaceGuessLicence([]byte(`The MIT License (MIT)

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
`))

	if licenses[0].LicenseId != "MIT" {
		t.Error("expected MIT got", licenses[0].LicenseId)
	}
}

func TestVectorSpaceGuessUnlicence(t *testing.T) {
	lg := NewLicenceGuesser(false, true)

	licenses := lg.VectorSpaceGuessLicence([]byte(`This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>
`))

	if licenses[0].LicenseId != "Unlicense" {
		t.Error("expected Unlicense got", licenses[0].LicenseId)
	}
}

func TestVectorSpaceGuessJSON(t *testing.T) {
	lg := NewLicenceGuesser(false, true)

	licenses := lg.VectorSpaceGuessLicence([]byte(`JSON License
Copyright (c) 2002 JSON.org
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
The Software shall be used for Good, not Evil.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`))

	if licenses[0].LicenseId != "JSON" {
		t.Error("expected JSON got", licenses[0].LicenseId)
	}
}

func TestVectorSpaceApache2Multiple(t *testing.T) {
	samples := []string{`Apache License
==============

_Version 2.0, January 2004_
_&lt;<http://www.apache.org/licenses/>&gt;_

### Terms and Conditions for use, reproduction, and distribution

#### 1. Definitions

“License” shall mean the terms and conditions for use, reproduction, and
distribution as defined by Sections 1 through 9 of this document.

“Licensor” shall mean the copyright owner or entity authorized by the copyright
owner that is granting the License.

“Legal Entity” shall mean the union of the acting entity and all other entities
that control, are controlled by, or are under common control with that entity.
For the purposes of this definition, “control” means **(i)** the power, direct or
indirect, to cause the direction or management of such entity, whether by
contract or otherwise, or **(ii)** ownership of fifty percent (50%) or more of the
outstanding shares, or **(iii)** beneficial ownership of such entity.

“You” (or “Your”) shall mean an individual or Legal Entity exercising
permissions granted by this License.

“Source” form shall mean the preferred form for making modifications, including
but not limited to software source code, documentation source, and configuration
files.

“Object” form shall mean any form resulting from mechanical transformation or
translation of a Source form, including but not limited to compiled object code,
generated documentation, and conversions to other media types.

“Work” shall mean the work of authorship, whether in Source or Object form, made
available under the License, as indicated by a copyright notice that is included
in or attached to the work (an example is provided in the Appendix below).

“Derivative Works” shall mean any work, whether in Source or Object form, that
is based on (or derived from) the Work and for which the editorial revisions,
annotations, elaborations, or other modifications represent, as a whole, an
original work of authorship. For the purposes of this License, Derivative Works
shall not include works that remain separable from, or merely link (or bind by
name) to the interfaces of, the Work and Derivative Works thereof.

“Contribution” shall mean any work of authorship, including the original version
of the Work and any modifications or additions to that Work or Derivative Works
thereof, that is intentionally submitted to Licensor for inclusion in the Work
by the copyright owner or by an individual or Legal Entity authorized to submit
on behalf of the copyright owner. For the purposes of this definition,
“submitted” means any form of electronic, verbal, or written communication sent
to the Licensor or its representatives, including but not limited to
communication on electronic mailing lists, source code control systems, and
issue tracking systems that are managed by, or on behalf of, the Licensor for
the purpose of discussing and improving the Work, but excluding communication
that is conspicuously marked or otherwise designated in writing by the copyright
owner as “Not a Contribution.”

“Contributor” shall mean Licensor and any individual or Legal Entity on behalf
of whom a Contribution has been received by Licensor and subsequently
incorporated within the Work.

#### 2. Grant of Copyright License

Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable copyright license to reproduce, prepare Derivative Works of,
publicly display, publicly perform, sublicense, and distribute the Work and such
Derivative Works in Source or Object form.

#### 3. Grant of Patent License

Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable (except as stated in this section) patent license to make, have
made, use, offer to sell, sell, import, and otherwise transfer the Work, where
such license applies only to those patent claims licensable by such Contributor
that are necessarily infringed by their Contribution(s) alone or by combination
of their Contribution(s) with the Work to which such Contribution(s) was
submitted. If You institute patent litigation against any entity (including a
cross-claim or counterclaim in a lawsuit) alleging that the Work or a
Contribution incorporated within the Work constitutes direct or contributory
patent infringement, then any patent licenses granted to You under this License
for that Work shall terminate as of the date such litigation is filed.

#### 4. Redistribution

You may reproduce and distribute copies of the Work or Derivative Works thereof
in any medium, with or without modifications, and in Source or Object form,
provided that You meet the following conditions:

* **(a)** You must give any other recipients of the Work or Derivative Works a copy of
this License; and
* **(b)** You must cause any modified files to carry prominent notices stating that You
changed the files; and
* **(c)** You must retain, in the Source form of any Derivative Works that You distribute,
all copyright, patent, trademark, and attribution notices from the Source form
of the Work, excluding those notices that do not pertain to any part of the
Derivative Works; and
* **(d)** If the Work includes a “NOTICE” text file as part of its distribution, then any
Derivative Works that You distribute must include a readable copy of the
attribution notices contained within such NOTICE file, excluding those notices
that do not pertain to any part of the Derivative Works, in at least one of the
following places: within a NOTICE text file distributed as part of the
Derivative Works; within the Source form or documentation, if provided along
with the Derivative Works; or, within a display generated by the Derivative
Works, if and wherever such third-party notices normally appear. The contents of
the NOTICE file are for informational purposes only and do not modify the
License. You may add Your own attribution notices within Derivative Works that
You distribute, alongside or as an addendum to the NOTICE text from the Work,
provided that such additional attribution notices cannot be construed as
modifying the License.

You may add Your own copyright statement to Your modifications and may provide
additional or different license terms and conditions for use, reproduction, or
distribution of Your modifications, or for any such Derivative Works as a whole,
provided Your use, reproduction, and distribution of the Work otherwise complies
with the conditions stated in this License.

#### 5. Submission of Contributions

Unless You explicitly state otherwise, any Contribution intentionally submitted
for inclusion in the Work by You to the Licensor shall be under the terms and
conditions of this License, without any additional terms or conditions.
Notwithstanding the above, nothing herein shall supersede or modify the terms of
any separate license agreement you may have executed with Licensor regarding
such Contributions.

#### 6. Trademarks

This License does not grant permission to use the trade names, trademarks,
service marks, or product names of the Licensor, except as required for
reasonable and customary use in describing the origin of the Work and
reproducing the content of the NOTICE file.

#### 7. Disclaimer of Warranty

Unless required by applicable law or agreed to in writing, Licensor provides the
Work (and each Contributor provides its Contributions) on an “AS IS” BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied,
including, without limitation, any warranties or conditions of TITLE,
NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE. You are
solely responsible for determining the appropriateness of using or
redistributing the Work and assume any risks associated with Your exercise of
permissions under this License.

#### 8. Limitation of Liability

In no event and under no legal theory, whether in tort (including negligence),
contract, or otherwise, unless required by applicable law (such as deliberate
and grossly negligent acts) or agreed to in writing, shall any Contributor be
liable to You for damages, including any direct, indirect, special, incidental,
or consequential damages of any character arising as a result of this License or
out of the use or inability to use the Work (including but not limited to
damages for loss of goodwill, work stoppage, computer failure or malfunction, or
any and all other commercial damages or losses), even if such Contributor has
been advised of the possibility of such damages.

#### 9. Accepting Warranty or Additional Liability

While redistributing the Work or Derivative Works thereof, You may choose to
offer, and charge a fee for, acceptance of support, warranty, indemnity, or
other liability obligations and/or rights consistent with this License. However,
in accepting such obligations, You may act only on Your own behalf and on Your
sole responsibility, not on behalf of any other Contributor, and only if You
agree to indemnify, defend, and hold each Contributor harmless for any liability
incurred by, or claims asserted against, such Contributor by reason of your
accepting any such warranty or additional liability.

_END OF TERMS AND CONDITIONS_

### APPENDIX: How to apply the Apache License to your work

To apply the Apache License to your work, attach the following boilerplate
notice, with the fields enclosed by brackets [] replaced with your own
identifying information. (Don't include the brackets!) The text should be
enclosed in the appropriate comment syntax for the file format. We also
recommend that a file or class name and description of purpose be included on
the same “printed page” as the copyright notice for easier identification within
third-party archives.

    Copyright [yyyy] [name of copyright owner]

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.`,
		`h1>Apache License</h1>

<p><em>Version 2.0, January 2004</em>
<em>&amp;lt;<a href="http://www.apache.org/licenses/">http://www.apache.org/licenses/</a>&amp;gt;</em></p>

<h3>Terms and Conditions for use, reproduction, and distribution</h3>

<h4>1. Definitions</h4>

<p>“License” shall mean the terms and conditions for use, reproduction, and
distribution as defined by Sections 1 through 9 of this document.</p>

<p>“Licensor” shall mean the copyright owner or entity authorized by the copyright
owner that is granting the License.</p>

<p>“Legal Entity” shall mean the union of the acting entity and all other entities
that control, are controlled by, or are under common control with that entity.
For the purposes of this definition, “control” means <strong>(i)</strong> the power, direct or
indirect, to cause the direction or management of such entity, whether by
contract or otherwise, or <strong>(ii)</strong> ownership of fifty percent (50%) or more of the
outstanding shares, or <strong>(iii)</strong> beneficial ownership of such entity.</p>

<p>“You” (or “Your”) shall mean an individual or Legal Entity exercising
permissions granted by this License.</p>

<p>“Source” form shall mean the preferred form for making modifications, including
but not limited to software source code, documentation source, and configuration
files.</p>

<p>“Object” form shall mean any form resulting from mechanical transformation or
translation of a Source form, including but not limited to compiled object code,
generated documentation, and conversions to other media types.</p>

<p>“Work” shall mean the work of authorship, whether in Source or Object form, made
available under the License, as indicated by a copyright notice that is included
in or attached to the work (an example is provided in the Appendix below).</p>

<p>“Derivative Works” shall mean any work, whether in Source or Object form, that
is based on (or derived from) the Work and for which the editorial revisions,
annotations, elaborations, or other modifications represent, as a whole, an
original work of authorship. For the purposes of this License, Derivative Works
shall not include works that remain separable from, or merely link (or bind by
name) to the interfaces of, the Work and Derivative Works thereof.</p>

<p>“Contribution” shall mean any work of authorship, including the original version
of the Work and any modifications or additions to that Work or Derivative Works
thereof, that is intentionally submitted to Licensor for inclusion in the Work
by the copyright owner or by an individual or Legal Entity authorized to submit
on behalf of the copyright owner. For the purposes of this definition,
“submitted” means any form of electronic, verbal, or written communication sent
to the Licensor or its representatives, including but not limited to
communication on electronic mailing lists, source code control systems, and
issue tracking systems that are managed by, or on behalf of, the Licensor for
the purpose of discussing and improving the Work, but excluding communication
that is conspicuously marked or otherwise designated in writing by the copyright
owner as “Not a Contribution.”</p>

<p>“Contributor” shall mean Licensor and any individual or Legal Entity on behalf
of whom a Contribution has been received by Licensor and subsequently
incorporated within the Work.</p>

<h4>2. Grant of Copyright License</h4>

<p>Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable copyright license to reproduce, prepare Derivative Works of,
publicly display, publicly perform, sublicense, and distribute the Work and such
Derivative Works in Source or Object form.</p>

<h4>3. Grant of Patent License</h4>

<p>Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable (except as stated in this section) patent license to make, have
made, use, offer to sell, sell, import, and otherwise transfer the Work, where
such license applies only to those patent claims licensable by such Contributor
that are necessarily infringed by their Contribution(s) alone or by combination
of their Contribution(s) with the Work to which such Contribution(s) was
submitted. If You institute patent litigation against any entity (including a
cross-claim or counterclaim in a lawsuit) alleging that the Work or a
Contribution incorporated within the Work constitutes direct or contributory
patent infringement, then any patent licenses granted to You under this License
for that Work shall terminate as of the date such litigation is filed.</p>

<h4>4. Redistribution</h4>

<p>You may reproduce and distribute copies of the Work or Derivative Works thereof
in any medium, with or without modifications, and in Source or Object form,
provided that You meet the following conditions:</p>

<ul>
<li><strong>(a)</strong> You must give any other recipients of the Work or Derivative Works a copy of
this License; and</li>
<li><strong>(b)</strong> You must cause any modified files to carry prominent notices stating that You
changed the files; and</li>
<li><strong>&copy;</strong> You must retain, in the Source form of any Derivative Works that You distribute,
all copyright, patent, trademark, and attribution notices from the Source form
of the Work, excluding those notices that do not pertain to any part of the
Derivative Works; and</li>
<li><strong>(d)</strong> If the Work includes a “NOTICE” text file as part of its distribution, then any
Derivative Works that You distribute must include a readable copy of the
attribution notices contained within such NOTICE file, excluding those notices
that do not pertain to any part of the Derivative Works, in at least one of the
following places: within a NOTICE text file distributed as part of the
Derivative Works; within the Source form or documentation, if provided along
with the Derivative Works; or, within a display generated by the Derivative
Works, if and wherever such third-party notices normally appear. The contents of
the NOTICE file are for informational purposes only and do not modify the
License. You may add Your own attribution notices within Derivative Works that
You distribute, alongside or as an addendum to the NOTICE text from the Work,
provided that such additional attribution notices cannot be construed as
modifying the License.</li>
</ul>

<p>You may add Your own copyright statement to Your modifications and may provide
additional or different license terms and conditions for use, reproduction, or
distribution of Your modifications, or for any such Derivative Works as a whole,
provided Your use, reproduction, and distribution of the Work otherwise complies
with the conditions stated in this License.</p>

<h4>5. Submission of Contributions</h4>

<p>Unless You explicitly state otherwise, any Contribution intentionally submitted
for inclusion in the Work by You to the Licensor shall be under the terms and
conditions of this License, without any additional terms or conditions.
Notwithstanding the above, nothing herein shall supersede or modify the terms of
any separate license agreement you may have executed with Licensor regarding
such Contributions.</p>

<h4>6. Trademarks</h4>

<p>This License does not grant permission to use the trade names, trademarks,
service marks, or product names of the Licensor, except as required for
reasonable and customary use in describing the origin of the Work and
reproducing the content of the NOTICE file.</p>

<h4>7. Disclaimer of Warranty</h4>

<p>Unless required by applicable law or agreed to in writing, Licensor provides the
Work (and each Contributor provides its Contributions) on an “AS IS” BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied,
including, without limitation, any warranties or conditions of TITLE,
NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE. You are
solely responsible for determining the appropriateness of using or
redistributing the Work and assume any risks associated with Your exercise of
permissions under this License.</p>

<h4>8. Limitation of Liability</h4>

<p>In no event and under no legal theory, whether in tort (including negligence),
contract, or otherwise, unless required by applicable law (such as deliberate
and grossly negligent acts) or agreed to in writing, shall any Contributor be
liable to You for damages, including any direct, indirect, special, incidental,
or consequential damages of any character arising as a result of this License or
out of the use or inability to use the Work (including but not limited to
damages for loss of goodwill, work stoppage, computer failure or malfunction, or
any and all other commercial damages or losses), even if such Contributor has
been advised of the possibility of such damages.</p>

<h4>9. Accepting Warranty or Additional Liability</h4>

<p>While redistributing the Work or Derivative Works thereof, You may choose to
offer, and charge a fee for, acceptance of support, warranty, indemnity, or
other liability obligations and/or rights consistent with this License. However,
in accepting such obligations, You may act only on Your own behalf and on Your
sole responsibility, not on behalf of any other Contributor, and only if You
agree to indemnify, defend, and hold each Contributor harmless for any liability
incurred by, or claims asserted against, such Contributor by reason of your
accepting any such warranty or additional liability.</p>

<p><em>END OF TERMS AND CONDITIONS</em></p>

<h3>APPENDIX: How to apply the Apache License to your work</h3>

<p>To apply the Apache License to your work, attach the following boilerplate
notice, with the fields enclosed by brackets <code>[]</code> replaced with your own
identifying information. (Don&rsquo;t include the brackets!) The text should be
enclosed in the appropriate comment syntax for the file format. We also
recommend that a file or class name and description of purpose be included on
the same “printed page” as the copyright notice for easier identification within
third-party archives.</p>

<pre><code>Copyright [yyyy] [name of copyright owner]

Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
</code></pre>`,
		`Apache License

Version 2.0, January 2004
<http://www.apache.org/licenses/>

Terms and Conditions for use, reproduction, and distribution.

1. Definitions.

“License” shall mean the terms and conditions for use, reproduction, and
distribution as defined by Sections 1 through 9 of this document.

“Licensor” shall mean the copyright owner or entity authorized by the copyright
owner that is granting the License.

“Legal Entity” shall mean the union of the acting entity and all other entities
that control, are controlled by, or are under common control with that entity.
For the purposes of this definition, “control” means (i) the power, direct or
indirect, to cause the direction or management of such entity, whether by
contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the
outstanding shares, or (iii) beneficial ownership of such entity.

“You” (or “Your”) shall mean an individual or Legal Entity exercising
permissions granted by this License.

“Source” form shall mean the preferred form for making modifications, including
but not limited to software source code, documentation source, and configuration
files.

“Object” form shall mean any form resulting from mechanical transformation or
translation of a Source form, including but not limited to compiled object code,
generated documentation, and conversions to other media types.

“Work” shall mean the work of authorship, whether in Source or Object form, made
available under the License, as indicated by a copyright notice that is included
in or attached to the work (an example is provided in the Appendix below).

“Derivative Works” shall mean any work, whether in Source or Object form, that
is based on (or derived from) the Work and for which the editorial revisions,
annotations, elaborations, or other modifications represent, as a whole, an
original work of authorship. For the purposes of this License, Derivative Works
shall not include works that remain separable from, or merely link (or bind by
name) to the interfaces of, the Work and Derivative Works thereof.

“Contribution” shall mean any work of authorship, including the original version
of the Work and any modifications or additions to that Work or Derivative Works
thereof, that is intentionally submitted to Licensor for inclusion in the Work
by the copyright owner or by an individual or Legal Entity authorized to submit
on behalf of the copyright owner. For the purposes of this definition,
“submitted” means any form of electronic, verbal, or written communication sent
to the Licensor or its representatives, including but not limited to
communication on electronic mailing lists, source code control systems, and
issue tracking systems that are managed by, or on behalf of, the Licensor for
the purpose of discussing and improving the Work, but excluding communication
that is conspicuously marked or otherwise designated in writing by the copyright
owner as “Not a Contribution.”

“Contributor” shall mean Licensor and any individual or Legal Entity on behalf
of whom a Contribution has been received by Licensor and subsequently
incorporated within the Work.

2. Grant of Copyright License.

Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable copyright license to reproduce, prepare Derivative Works of,
publicly display, publicly perform, sublicense, and distribute the Work and such
Derivative Works in Source or Object form.

3. Grant of Patent License.

Subject to the terms and conditions of this License, each Contributor hereby
grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable (except as stated in this section) patent license to make, have
made, use, offer to sell, sell, import, and otherwise transfer the Work, where
such license applies only to those patent claims licensable by such Contributor
that are necessarily infringed by their Contribution(s) alone or by combination
of their Contribution(s) with the Work to which such Contribution(s) was
submitted. If You institute patent litigation against any entity (including a
cross-claim or counterclaim in a lawsuit) alleging that the Work or a
Contribution incorporated within the Work constitutes direct or contributory
patent infringement, then any patent licenses granted to You under this License
for that Work shall terminate as of the date such litigation is filed.

4. Redistribution.

You may reproduce and distribute copies of the Work or Derivative Works thereof
in any medium, with or without modifications, and in Source or Object form,
provided that You meet the following conditions:


(a) You must give any other recipients of the Work or Derivative Works a copy of
this License; and
(b) You must cause any modified files to carry prominent notices stating that You
changed the files; and
© You must retain, in the Source form of any Derivative Works that You distribute,
all copyright, patent, trademark, and attribution notices from the Source form
of the Work, excluding those notices that do not pertain to any part of the
Derivative Works; and
(d) If the Work includes a “NOTICE” text file as part of its distribution, then any
Derivative Works that You distribute must include a readable copy of the
attribution notices contained within such NOTICE file, excluding those notices
that do not pertain to any part of the Derivative Works, in at least one of the
following places: within a NOTICE text file distributed as part of the
Derivative Works; within the Source form or documentation, if provided along
with the Derivative Works; or, within a display generated by the Derivative
Works, if and wherever such third-party notices normally appear. The contents of
the NOTICE file are for informational purposes only and do not modify the
License. You may add Your own attribution notices within Derivative Works that
You distribute, alongside or as an addendum to the NOTICE text from the Work,
provided that such additional attribution notices cannot be construed as
modifying the License.


You may add Your own copyright statement to Your modifications and may provide
additional or different license terms and conditions for use, reproduction, or
distribution of Your modifications, or for any such Derivative Works as a whole,
provided Your use, reproduction, and distribution of the Work otherwise complies
with the conditions stated in this License.

5. Submission of Contributions.

Unless You explicitly state otherwise, any Contribution intentionally submitted
for inclusion in the Work by You to the Licensor shall be under the terms and
conditions of this License, without any additional terms or conditions.
Notwithstanding the above, nothing herein shall supersede or modify the terms of
any separate license agreement you may have executed with Licensor regarding
such Contributions.

6. Trademarks.

This License does not grant permission to use the trade names, trademarks,
service marks, or product names of the Licensor, except as required for
reasonable and customary use in describing the origin of the Work and
reproducing the content of the NOTICE file.

7. Disclaimer of Warranty.

Unless required by applicable law or agreed to in writing, Licensor provides the
Work (and each Contributor provides its Contributions) on an “AS IS” BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied,
including, without limitation, any warranties or conditions of TITLE,
NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE. You are
solely responsible for determining the appropriateness of using or
redistributing the Work and assume any risks associated with Your exercise of
permissions under this License.

8. Limitation of Liability.

In no event and under no legal theory, whether in tort (including negligence),
contract, or otherwise, unless required by applicable law (such as deliberate
and grossly negligent acts) or agreed to in writing, shall any Contributor be
liable to You for damages, including any direct, indirect, special, incidental,
or consequential damages of any character arising as a result of this License or
out of the use or inability to use the Work (including but not limited to
damages for loss of goodwill, work stoppage, computer failure or malfunction, or
any and all other commercial damages or losses), even if such Contributor has
been advised of the possibility of such damages.

9. Accepting Warranty or Additional Liability.

While redistributing the Work or Derivative Works thereof, You may choose to
offer, and charge a fee for, acceptance of support, warranty, indemnity, or
other liability obligations and/or rights consistent with this License. However,
in accepting such obligations, You may act only on Your own behalf and on Your
sole responsibility, not on behalf of any other Contributor, and only if You
agree to indemnify, defend, and hold each Contributor harmless for any liability
incurred by, or claims asserted against, such Contributor by reason of your
accepting any such warranty or additional liability.

END OF TERMS AND CONDITIONS

APPENDIX: How to apply the Apache License to your work.

To apply the Apache License to your work, attach the following boilerplate
notice, with the fields enclosed by brackets [] replaced with your own
identifying information. (Don’t include the brackets!) The text should be
enclosed in the appropriate comment syntax for the file format. We also
recommend that a file or class name and description of purpose be included on
the same “printed page” as the copyright notice for easier identification within
third-party archives.

Copyright [yyyy] [name of copyright owner]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`,
		`apache license

version 2.0, january 2004
http:/www.apache.org/licenses/
terms and conditions for use, reproduction, and distribution.
definitions.
"license" shall mean the terms and conditions for use, reproduction, and
distribution as defined by sections 1 through 9 of this document.
"licensor" shall mean the © or entity authorized by the ©
owner that is granting the license

"legal entity" shall mean the union of the acting entity and all other entities
that control, are controlled by, or are under common control with that entity.
for the purposes of this definition, "control" means (i) the power, direct or
indirect, to cause the direction or management of such entity, whether by
contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the
outstanding shares, or (iii) beneficial ownership of such entity.
"you" (or "your") shall mean an individual or legal entity exercising
permissions granted by this license

"source" form shall mean the preferred form for making modifications, including
but not limited to software source code, documentation source, and configuration
files.
"object" form shall mean any form resulting from mechanical transformation or
translation of a source form, including but not limited to compiled object code,
generated documentation, and conversions to other media types.
"work" shall mean the work of authorship, whether in source or object form, made
available under the license, as indicated by a © notice that is included
in or attached to the work (an example is provided in the appendix below).
"derivative works" shall mean any work, whether in source or object form, that
is based on (or derived from) the work and for which the editorial revisions,
annotations, elaborations, or other modifications represent, as a whole, an
original work of authorship. for the purposes of this license, derivative works
shall not include works that remain separable from, or merely link (or bind by
name) to the interfaces of, the work and derivative works thereof.
"contribution" shall mean any work of authorship, including the original version
of the work and any modifications or additions to that work or derivative works
thereof, that is intentionally submitted to licensor for inclusion in the work
by the © or by an individual or legal entity authorized to submit
on behalf of the ©. for the purposes of this definition,
"submitted" means any form of electronic, verbal, or written communication sent
to the licensor or its representatives, including but not limited to
communication on electronic mailing lists, source code control systems, and
issue tracking systems that are managed by, or on behalf of, the licensor for
the purpose of discussing and improving the work, but excluding communication
that is conspicuously marked or otherwise designated in writing by the ©
owner as "not a contribution."
"contributor" shall mean licensor and any individual or legal entity on behalf
of whom a contribution has been received by licensor and subsequently
incorporated within the work.
grant of © license

subject to the terms and conditions of this license, each contributor hereby
grants to you a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable © license to reproduce, prepare derivative works of,
publicly display, publicly perform, sublicense, and distribute the work and such
derivative works in source or object form.
grant of patent license

subject to the terms and conditions of this license, each contributor hereby
grants to you a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable (except as stated in this section) patent license to make, have
made, use, offer to sell, sell, import, and otherwise transfer the work, where
such license applies only to those patent claims licensable by such contributor
that are necessarily infringed by their contribution(s) alone or by combination
of their contribution(s) with the work to which such contribution(s) was
submitted. if you institute patent litigation against any entity (including a
cross-claim or counterclaim in a lawsuit) alleging that the work or a
contribution incorporated within the work constitutes direct or contributory
patent infringement, then any patent licenses granted to you under this license
for that work shall terminate as of the date such litigation is filed.
redistribution.
you may reproduce and distribute copies of the work or derivative works thereof
in any medium, with or without modifications, and in source or object form,
provided that you meet the following conditions:
you must give any other recipients of the work or derivative works a copy of
this license; and
you must cause any modified files to carry prominent notices stating that you
changed the files; and
© you must retain, in the source form of any derivative works that you distribute,
all ©, patent, ™, and attribution notices from the source form
of the work, excluding those notices that do not pertain to any part of the
derivative works; and
if the work includes a "notice" text file as part of its distribution, then any
derivative works that you distribute must include a readable copy of the
attribution notices contained within such notice file, excluding those notices
that do not pertain to any part of the derivative works, in at least one of the
following places: within a notice text file distributed as part of the
derivative works; within the source form or documentation, if provided along
with the derivative works; or, within a display generated by the derivative
works, if and wherever such third-party notices normally appear. the contents of
the notice file are for informational purposes only and do not modify the
license. you may add your own attribution notices within derivative works that
you distribute, alongside or as an addendum to the notice text from the work,
provided that such additional attribution notices cannot be construed as
modifying the license

you may add your own © statement to your modifications and may provide
additional or different license terms and conditions for use, reproduction, or
distribution of your modifications, or for any such derivative works as a whole,
provided your use, reproduction, and distribution of the work otherwise complies
with the conditions stated in this license

submission of contributions.
unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you to the licensor shall be under the terms and
conditions of this license, without any additional terms or conditions.
notwithstanding the above, nothing herein shall supersede or modify the terms of
any separate license agreement you may have executed with licensor regarding
such contributions.
™.
this license does not grant permission to use the trade names, ™,
service marks, or product names of the licensor, except as required for
reasonable and customary use in describing the origin of the work and
reproducing the content of the notice file.
disclaimer of warranty.
unless required by applicable law or agreed to in writing, licensor provides the
work (and each contributor provides its contributions) on an "as is" basis,
without warranties or conditions of any kind, either express or implied,
including, without limitation, any warranties or conditions of title,
non-infringement, merchantability, or fitness for a particular purpose. you are
solely responsible for determining the appropriateness of using or
redistributing the work and assume any risks associated with your exercise of
permissions under this license

limitation of liability.
in no event and under no legal theory, whether in tort (including negligence),
contract, or otherwise, unless required by applicable law (such as deliberate
and grossly negligent acts) or agreed to in writing, shall any contributor be
liable to you for damages, including any direct, indirect, special, incidental,
or consequential damages of any character arising as a result of this license or
out of the use or inability to use the work (including but not limited to
damages for loss of goodwill, work stoppage, computer failure or malfunction, or
any and all other commercial damages or losses), even if such contributor has
been advised of the possibility of such damages.
accepting warranty or additional liability.
while redistributing the work or derivative works thereof, you may choose to
offer, and charge a fee for, acceptance of support, warranty, indemnity, or
other liability obligations and/or rights consistent with this license. however,
in accepting such obligations, you may act only on your own behalf and on your
sole responsibility, not on behalf of any other contributor, and only if you
agree to indemnify, defend, and hold each contributor harmless for any liability
incurred by, or claims asserted against, such contributor by reason of your
accepting any such warranty or additional liability.
end of terms and conditions
appendix: how to apply the apache license to your work.
to apply the apache license to your work, attach the following boilerplate
notice, with the fields enclosed by brackets [] replaced with your own
identifying information. (don"t include the brackets!) the text should be
enclosed in the appropriate comment syntax for the file format. we also
recommend that a file or class name and description of purpose be included on
the same "printed page" as the © notice for easier identification within
third-party archives.
© [yyyy] [name of ©]
licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at
http:/www.apache.org/licenses/license-2.0
unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license`,
		`apache license

version 20, january 2004
http:/wwwapacheorg/licenses/
terms and conditions for use, reproduction, and distribution
definitions
"license" shall mean the terms and conditions for use, reproduction, and
distribution as defined by sections 1 through 9 of this document
"licensor" shall mean the © or entity authorized by the ©
owner that is granting the license

"legal entity" shall mean the union of the acting entity and all other entities
that control, are controlled by, or are under common control with that entity
for the purposes of this definition, "control" means (i) the power, direct or
indirect, to cause the direction or management of such entity, whether by
contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the
outstanding shares, or (iii) beneficial ownership of such entity
"you" (or "your") shall mean an individual or legal entity exercising
permissions granted by this license

"source" form shall mean the preferred form for making modifications, including
but not limited to software source code, documentation source, and configuration
files
"object" form shall mean any form resulting from mechanical transformation or
translation of a source form, including but not limited to compiled object code,
generated documentation, and conversions to other media types
"work" shall mean the work of authorship, whether in source or object form, made
available under the license, as indicated by a © notice that is included
in or attached to the work (an example is provided in the appendix below)
"derivative works" shall mean any work, whether in source or object form, that
is based on (or derived from) the work and for which the editorial revisions,
annotations, elaborations, or other modifications represent, as a whole, an
original work of authorship for the purposes of this license, derivative works
shall not include works that remain separable from, or merely link (or bind by
name) to the interfaces of, the work and derivative works thereof
"contribution" shall mean any work of authorship, including the original version
of the work and any modifications or additions to that work or derivative works
thereof, that is intentionally submitted to licensor for inclusion in the work
by the © or by an individual or legal entity authorized to submit
on behalf of the © for the purposes of this definition,
"submitted" means any form of electronic, verbal, or written communication sent
to the licensor or its representatives, including but not limited to
communication on electronic mailing lists, source code control systems, and
issue tracking systems that are managed by, or on behalf of, the licensor for
the purpose of discussing and improving the work, but excluding communication
that is conspicuously marked or otherwise designated in writing by the ©
owner as "not a contribution"
"contributor" shall mean licensor and any individual or legal entity on behalf
of whom a contribution has been received by licensor and subsequently
incorporated within the work
grant of © license

subject to the terms and conditions of this license, each contributor hereby
grants to you a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable © license to reproduce, prepare derivative works of,
publicly display, publicly perform, sublicense, and distribute the work and such
derivative works in source or object form
grant of patent license

subject to the terms and conditions of this license, each contributor hereby
grants to you a perpetual, worldwide, non-exclusive, no-charge, royalty-free,
irrevocable (except as stated in this section) patent license to make, have
made, use, offer to sell, sell, import, and otherwise transfer the work, where
such license applies only to those patent claims licensable by such contributor
that are necessarily infringed by their contribution(s) alone or by combination
of their contribution(s) with the work to which such contribution(s) was
submitted if you institute patent litigation against any entity (including a
cross-claim or counterclaim in a lawsuit) alleging that the work or a
contribution incorporated within the work constitutes direct or contributory
patent infringement, then any patent licenses granted to you under this license
for that work shall terminate as of the date such litigation is filed
redistribution
you may reproduce and distribute copies of the work or derivative works thereof
in any medium, with or without modifications, and in source or object form,
provided that you meet the following conditions:
you must give any other recipients of the work or derivative works a copy of
this license; and
you must cause any modified files to carry prominent notices stating that you
changed the files; and
all ©, patent, ™, and attribution notices from the source form
of the work, excluding those notices that do not pertain to any part of the
derivative works; and
if the work includes a "notice" text file as part of its distribution, then any
derivative works that you distribute must include a readable copy of the
attribution notices contained within such notice file, excluding those notices
that do not pertain to any part of the derivative works, in at least one of the
following places: within a notice text file distributed as part of the
derivative works; within the source form or documentation, if provided along
with the derivative works; or, within a display generated by the derivative
works, if and wherever such third-party notices normally appear the contents of
the notice file are for informational purposes only and do not modify the
license you may add your own attribution notices within derivative works that
you distribute, alongside or as an addendum to the notice text from the work,
provided that such additional attribution notices cannot be construed as
modifying the license

you may add your own © statement to your modifications and may provide
additional or different license terms and conditions for use, reproduction, or
distribution of your modifications, or for any such derivative works as a whole,
provided your use, reproduction, and distribution of the work otherwise complies
with the conditions stated in this license

submission of contributions
unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you to the licensor shall be under the terms and
conditions of this license, without any additional terms or conditions
notwithstanding the above, nothing herein shall supersede or modify the terms of
any separate license agreement you may have executed with licensor regarding
such contributions
™
this license does not grant permission to use the trade names, ™,
service marks, or product names of the licensor, except as required for
reasonable and customary use in describing the origin of the work and
reproducing the content of the notice file
disclaimer of warranty
unless required by applicable law or agreed to in writing, licensor provides the
work (and each contributor provides its contributions) on an "as is" basis,
without warranties or conditions of any kind, either express or implied,
including, without limitation, any warranties or conditions of title,
non-infringement, merchantability, or fitness for a particular purpose you are
solely responsible for determining the appropriateness of using or
redistributing the work and assume any risks associated with your exercise of
permissions under this license

limitation of liability
in no event and under no legal theory, whether in tort (including negligence),
contract, or otherwise, unless required by applicable law (such as deliberate
and grossly negligent acts) or agreed to in writing, shall any contributor be
liable to you for damages, including any direct, indirect, special, incidental,
or consequential damages of any character arising as a result of this license or
out of the use or inability to use the work (including but not limited to
damages for loss of goodwill, work stoppage, computer failure or malfunction, or
any and all other commercial damages or losses), even if such contributor has
been advised of the possibility of such damages
accepting warranty or additional liability
while redistributing the work or derivative works thereof, you may choose to
offer, and charge a fee for, acceptance of support, warranty, indemnity, or
other liability obligations and/or rights consistent with this license however,
in accepting such obligations, you may act only on your own behalf and on your
sole responsibility, not on behalf of any other contributor, and only if you
agree to indemnify, defend, and hold each contributor harmless for any liability
incurred by, or claims asserted against, such contributor by reason of your
accepting any such warranty or additional liability
end of terms and conditions
appendix: how to apply the apache license to your work
to apply the apache license to your work, attach the following boilerplate
notice, with the fields enclosed by brackets [] replaced with your own
identifying information (don"t include the brackets!) the text should be
enclosed in the appropriate comment syntax for the file format we also
recommend that a file or class name and description of purpose be included on
the same "printed page" as the © notice for easier identification within
third-party archives
licensed under the apache license, version 20 (the "license");
you may not use this file except in compliance with the license
you may obtain a copy of the license at
http:/wwwapacheorg/licenses/license-20
unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied
see the license for the specific language governing permissions and
limitations under the license`}

	lg := NewLicenceGuesser(true, true)

	for i, l := range samples {
		license := lg.VectorSpaceGuessLicence([]byte(l))

		if license[0].LicenseId != "Apache-2.0" {
			t.Error("expected Apache-2.0 got", license[0].LicenseId, "for", i)
		}
	}
}