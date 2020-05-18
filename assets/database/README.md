The following appear to not have any unique ngrams

```
AGPL-3.0-only Ngrams 46034 Unique Ngrams 0
AGPL-3.0-or-later Ngrams 46078 Unique Ngrams 0
AGPL-3.0 Ngrams 46034 Unique Ngrams 0

Artistic-1.0 Ngrams 6706 Unique Ngrams 0

GFDL-1.1-or-later Ngrams 24654 Unique Ngrams 0
GFDL-1.1 Ngrams 24654 Unique Ngrams 0
GFDL-1.2-or-later Ngrams 27626 Unique Ngrams 0
GFDL-1.2 Ngrams 27626 Unique Ngrams 0
GFDL-1.3-or-later Ngrams 30978 Unique Ngrams 0
GFDL-1.3 Ngrams 30978 Unique Ngrams 0

GPL-1.0+ Ngrams 17582 Unique Ngrams 0
GPL-1.0-only Ngrams 17550 Unique Ngrams 0
GPL-1.0 Ngrams 17550 Unique Ngrams 0
GPL-2.0+ Ngrams 24502 Unique Ngrams 0
GPL-2.0-only Ngrams 24502 Unique Ngrams 0
GPL-2.0-or-later Ngrams 24518 Unique Ngrams 0
GPL-2.0 Ngrams 24502 Unique Ngrams 0
GPL-3.0+ Ngrams 46974 Unique Ngrams 0
GPL-3.0-only Ngrams 46978 Unique Ngrams 0
GPL-3.0 Ngrams 46974 Unique Ngrams 0

LGPL-2.0+ Ngrams 34866 Unique Ngrams 0
LGPL-2.0-only Ngrams 34866 Unique Ngrams 0
LGPL-2.0-or-later Ngrams 34882 Unique Ngrams 0
LGPL-2.0 Ngrams 34866 Unique Ngrams 0
LGPL-2.1+ Ngrams 36578 Unique Ngrams 0
LGPL-2.1-only Ngrams 36578 Unique Ngrams 0
LGPL-2.1 Ngrams 36578 Unique Ngrams 0
LGPL-3.0+ Ngrams 10502 Unique Ngrams 0
LGPL-3.0-only Ngrams 10502 Unique Ngrams 0
LGPL-3.0-or-later Ngrams 10502 Unique Ngrams 0
LGPL-3.0 Ngrams 10502 Unique Ngrams 0

MPL-2.0-no-copyleft-exception Ngrams 20926 Unique Ngrams 0
MPL-2.0 Ngrams 20926 Unique Ngrams 0

SMLNJ Ngrams 1390 Unique Ngrams 0

StandardML-NJ Ngrams 1390 Unique Ngrams 0
```

For the above we need to determine if it falls into one of the above buckets... which means we need to find
a ngram that's perhaps common to that group

So if we check the file, and we don't have a match at all, it means it might be one of the above. IE the lack of anything
indicates that it could be one of the above. 