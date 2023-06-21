package processor

func NewLicenceDetector(useFullDatabase bool) LicenceDetector {
	l := LicenceDetector{}
	l.UseFullDatabase = useFullDatabase
	return l
}

type LicenceDetector struct {
	UseFullDatabase bool
}

func (ld *LicenceDetector) Detect(filename string, content string) {
	
}
