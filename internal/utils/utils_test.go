package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrls(t *testing.T) {
	cases := []struct {
		url string
		ok  bool
	}{{
		url: "https://crates.io/crates/mongodb",
		ok:  true,
	},
	}
	for _, test := range cases {
		t.Run(test.url, func(t *testing.T) {
			err, ok := IsReachable(test.url)
			assert.NoError(t, err)
			assert.True(t, ok)
		})
	}
}
