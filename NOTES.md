Detection Rules
---------------

Check the filename

    "li[cs]en[cs]e(s?)",
    "legal",
    "copy(left|right|ing)",
    "unlicense",
    "l?gpl([-_ v]?)(\\d\\.?\\d)?",
    "bsd",
    "mit",
    "apache",

If it falls into one of the above, its highly likely it is a licence,
and should be tested. Note that the name itself highly indicates the 
licence itself, with unlicense for example indicating it is the unlicnse.

Something like licence, legal, or copy(left|right|ing) needs to be checked
because while it is highly likely to have a licence we cannot be sure 
as to which licence it actually is. Its also possible that these examples 
could have multiple licenses in them. Example github.com/valkey/valkey/COPYING
	
    "",
    ".md",
    ".rst",
    ".html",
    ".txt",
	
Where the file matchs the above patterns, where it has has no extention or 
one of the others we should inspect it to see if it has a license. Its possible 
a licence exists here, but we cannot be sure. Note that its possible there are multiple
licences in the file which needs to be dealt with.

    // SPDX-License-Identifier: MIT OR Unlicense

For all other files, there are a few possibilities.
The first is that it contains a SPDX header such as the above which indicates
which license the file is under. Its also possible that the header will contain 
a full copy of another licence such as MIT, GPL or otherwise. Possibly inside a comment
or a long string declaration in the case of code. Its possible it has multiple.

