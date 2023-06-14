package processor

import (
	"testing"
)

func TestIsLicenceFile(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{
			name: "",
			val:  "mit.txt",
			want: true,
		},
		{
			name: "",
			val:  "hppc-LICENSE-ASL.txt",
			want: true,
		},
		{
			name: "",
			val:  "org.apache.lucene.codecs.Codec",
			want: false,
		},
		{
			name: "",
			val:  "org.apache.lucene.analysis.TokenFilterFactory",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLicenceFile(tt.val); got != tt.want {
				t.Errorf("licenseFileRe = %v, want %v for %v", got, tt.want, tt.val)
			}
		})
	}
}
