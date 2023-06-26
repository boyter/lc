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

func TestContainsString(t *testing.T) {
	type args struct {
		ids []string
		lst []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				ids: []string{},
				lst: []string{},
			},
			want: false,
		},
		{
			name: "empty lst",
			args: args{
				ids: []string{"a", "b", "c"},
				lst: []string{},
			},
			want: false,
		},
		{
			name: "contains a",
			args: args{
				ids: []string{"a", "b", "c"},
				lst: []string{"a"},
			},
			want: true,
		},
		{
			name: "contains b",
			args: args{
				ids: []string{"a", "b", "c"},
				lst: []string{"b"},
			},
			want: true,
		},
		{
			name: "contains multiple",
			args: args{
				ids: []string{"a", "b", "c"},
				lst: []string{"a", "b", "c"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.ids, tt.args.lst); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
