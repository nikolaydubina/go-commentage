package gitipc

import (
	"testing"
	"time"
)

func Test_parseBlameLine(t *testing.T) {
	tests := []struct {
		s string
		v lineDetails
	}{
		{
			s: "f718957a7977 pkg/util/node.go      (<dawnchen@google.com>             1420487061 -0800  16) ",
			v: lineDetails{commit: "f718957a7977", createdAt: time.Unix(1420487061, 0)},
		},
		{
			s: "f718957a7977 pkg/util/node.go      (<dawnchen@google.com>             1420487061 -0800  15) */",
			v: lineDetails{commit: "f718957a7977", createdAt: time.Unix(1420487061, 0)},
		},
		{
			s: "d92ee41e44b4 pkg/util/node/node.go (<wfender@google.com>              1544211079 -0800  65) // NoMatchError is a typed implementation of the error interface. It indicates a failure to get a matching Node.",
			v: lineDetails{commit: "d92ee41e44b4", createdAt: time.Unix(1544211079, 0)},
		},
		{
			s: "1a7f7c539919 (<jliggitt@redhat.com>          1475779952 -0400   1) /*",
			v: lineDetails{commit: "1a7f7c539919", createdAt: time.Unix(1475779952, 0)},
		},
	}
	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			v, err := parseBlameLine(tc.s)
			if err != nil {
				t.Errorf("%#v: got err: %s", tc, err)
			}
			if v != tc.v {
				t.Errorf("%#v: got(%v)", tc, v)
			}
		})
	}
}
