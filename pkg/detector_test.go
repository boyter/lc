// SPDX-License-Identifier: AGPL-3.0

package pkg

import (
	"reflect"
	"testing"
)

func TestLicenceDetector_LoadDatabase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "load",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LicenceDetector{}
			l.LoadDatabase()
			if len(l.Database) == 0 {
				t.Errorf("LicenceDetector.LoadDatabase() database is empty")
			}

			if len(l.Database[0].LicenseIds) == 0 {
				t.Errorf("LicenceDetector.LoadDatabase() database license ids is empty")
			}
		})
	}
}

func TestLicenceDetector_GuessFilenameLogic(t *testing.T) {
	type args struct {
		filename string
		content  string
	}
	tests := []struct {
		name string
		args args
		want []LicenseGuess
	}{
		{
			name: "Unlicense",
			args: args{
				filename: "unlicense",
				content:  "",
			},
			want: []LicenseGuess{
				{
					Name: "Unlicense",
				},
			},
		},
		{
			name: "LGPL-2.0",
			args: args{
				filename: "LGPL-2.0",
				content:  "",
			},
			want: []LicenseGuess{
				{
					Name: "LGPL-2.0",
				},
			},
		},
		{
			name: "MIT",
			args: args{
				filename: "license",
				content:  mitLicence,
			},
			want: []LicenseGuess{
				{
					Name: "MIT",
				},
			},
		},
		{
			name: "ISC",
			args: args{
				filename: "license",
				content:  iscLicense,
			},
			want: []LicenseGuess{
				{
					Name: "ISC",
				},
			},
		},
		{
			name: "Unlicense 2",
			args: args{
				filename: "license",
				content:  unlicense,
			},
			want: []LicenseGuess{
				{
					Name: "Unlicense",
				},
			},
		},
		{
			name: "BSD-2-Clause",
			args: args{
				filename: "BSD-2-Clause",
				content:  unlicense,
			},
			want: []LicenseGuess{
				{
					Name: "BSD-2-Clause",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector()
			if got := l.Guess(tt.args.filename, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Guess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLicenceDetector_GuessFilenameReadmeLogic(t *testing.T) {
	type args struct {
		filename string
		content  string
	}
	tests := []struct {
		name string
		args args
		want []LicenseGuess
	}{
		{
			name: "README Unlicence",
			args: args{
				filename: "README.md",
				content:  unlicense,
			},
			want: []LicenseGuess{
				{
					Name: "Unlicense",
				},
			},
		},
		{
			name: "README MIT",
			args: args{
				filename: "README.md",
				content:  mitLicence,
			},
			want: []LicenseGuess{
				{
					Name: "MIT",
				},
			},
		},
		{
			name: "README ISC",
			args: args{
				filename: "README.md",
				content:  iscLicense,
			},
			want: []LicenseGuess{
				{
					Name: "ISC",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLicenceDetector()
			if got := l.Guess(tt.args.filename, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Guess() = %v, want %v", got, tt.want)
			}
		})
	}
}
