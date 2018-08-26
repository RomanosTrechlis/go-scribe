package scribe

import (
	"testing"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		path     string
		port     int
		size     int64
		mediator string

		crt string
		pk  string
		ca  string

		err bool
	}{
		// {"testdata", 1235, 111111, "", "", "", "", false}, // fails when run on travis
		{"testdata", 1234, 111111, "", "dummy", "dummy", "dummy", true},
		{"no_dir", 1234, 111111, "", "", "", "", false},
	}

	for _, tt := range tests {
		_, err := New(tt.path, tt.port, tt.size, tt.mediator, tt.crt, tt.pk, tt.ca)
		if err != nil && !tt.err {
			t.Errorf("expecting no err, got error %v", err)
		}
		if err == nil && tt.err {
			t.Error("expecting err, got no error")
		}
	}
}
