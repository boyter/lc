## What we look at

licencechecker (lc) works by comparing files against the list of licences supplied by [SPDX](https://spdx.org/).

### Detecting the license

The process works like so.

First the command line arguments are checked to see if they refer to file or a folder. The difference between the two is that if it is a folder the whole folder is first scanned to see if it can identify any files which indicate a license. This can be controled using the argument `licensefiles` which by default is set to look for filenames which contain the string license, copying or readme. For example the following files would be identfied as potentially containing a licence.

* `LICENSE`
* `LICENCE`
* `license.md`
* `COPYING.txt`
* `LICENSE-MIT`
* `COPYRIGHT`
* `UNLICENSE`
* `README.md`

These are then taken as potential root licences under which all other files would be marked against.

If there are multiple root licenses then they are treated using `OR` as through the project is dual licensed. The check for a root licence happends in every folder.

When a candidate for a license is found its contents are checked against a list of unique ngrams for licenses. If there any ngrams matching then each licence is checked using a vector space to obtain the fuzzy match value for that license. If the value is over the confidence value set which by default is 85% then it is marked as being a match. 

If no license is found using the previous step then the file is checked against all licenses using the fuzzy matching. This is because some licenses do not have any unique ngrams allowing identification. Fuzzy matching only looks at the top portion of the file where license headers are expected to exist.

At the end the matches if any are sorted by and the most likely match is returned as the matching license. As such a licence file is considered to only contain a single declared licence.

Note that only the most recent root licenses are taken into account, so if a project with a root license of MIT has a subfolder with a root license of GPL-2.0+ files in that root folder will be marked as being GPL-2.0+ and not MIT.

For individual files the file is scanned in the same way but in addtion is scanned for any SPDX indicators such as,

`SPDX-License-Identifier: GPL-3.0-only`

Which will take precedence over any fuzzy matching. The indicators must match known SPDX licenses or they will be disregarded.

When finished the license determined is based on the SDPX identifiers if present, fuzzy matching if over the confidence value and then the root license(s). If multiples match they are treated as an `AND`.

Take for example,

```
Directory              File               License                        Confidence  Size
./examples/identifier  has_identifier.py  (MIT OR GPL-3.0+) AND GPL-2.0  100.00%     428B
```

The root licences were identified as being both MIT and GPL-3.0+ however inside the code itself it has a GPL-2.0 identifier. As such the lience of the file is either `MIT AND GPL-2.0` OR `GPL-3.0+ OR GPL-2.0`. The indicator for the license like this is based on SPDX examples and as specified in the SPDX specification https://spdx.org/spdx-specification-21-web-version

### Acting on this information

Currently licencechecker (lc) does not indicate if a project may be in breech of license requirements. Assuming there is some sort of license compatibility chart in the future this would be something to add at a later point.

### Known Issues

License's are fuzzy matched so its possible someone could fork an existing license and have it be a false positive match.
The license matches are based on licenses from [SPDX](https://spdx.org/) and as such may miss some licenses.
