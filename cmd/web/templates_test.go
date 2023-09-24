package main

import (
	"testing"
	"time"

	"github.com.scottyfionnghall.snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	// (the want field).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 9, 24, 11, 37, 0, 0, time.UTC),
			want: "24 Sep 2023 at 11:37",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 9, 24, 11, 37, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "24 Sep 2023 at 10:37",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
