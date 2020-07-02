// +build selinux,linux

package setrans

import (
	"strings"
	"testing"
)

func check(t *testing.T, err error, wanterr, got, want string) {
	t.Helper()

	if err != nil {
		if wanterr == "" {
			t.Fatalf("got error %q when no error expected", err)
		}
		if !strings.Contains(err.Error(), wanterr) {
			t.Fatalf("got error %q does not match expected err %q", err, wanterr)
		}
		return
	}

	if got != want {
		t.Errorf("translated %q does not match expected %q", got, want)
	}
}

// These tests require an MLS policy and will not work with the default
// Redhat targeted policy
func TestTranslation(t *testing.T) {
	t.Run("Test TransToRaw", func(t *testing.T) {
		tests := []struct {
			orig, want, err, name string
		}{
			{
				orig: "staff_u:staff_r:staff_t:SystemLow-SystemHigh",
				want: "staff_u:staff_r:staff_t:s0-s15:c0.c1023",
				name: "test a valid TransToRaw",
			},
			{
				orig: "staff_u:staff_r:staff_t:FooLow-FooHigh",
				want: "staff_u:staff_r:staff_t:FooLow-FooHigh",
				name: "test an invalid original context to TransToRaw",
			},
			{
				orig: "FooLow-FooHigh",
				name: "test an incomplete selinux context",
				err:  "failed to read from response",
			},
		}

		conn, err := New()
		if err != nil {
			t.Fatalf("failed to connect to mcstrans: %q", err)
		}
		defer conn.Close()

		for _, tt := range tests {
			got, err := conn.TransToRaw(tt.orig)
			check(t, err, tt.err, got, tt.want)
		}
	})

	t.Run("Test RawToColor", func(t *testing.T) {
		tests := []struct {
			orig, want, err, name string
		}{
			{
				orig: "staff_u:staff_r:staff_t:s0-s15:c0.c1023",
				want: "#000000 #ffffff #000000 #ffffff #000000 #ffffff #000000 #ffffff",
				name: "test a valid RawToColor",
			},
		}

		conn, err := New()
		if err != nil {
			t.Fatalf("failed to connect to mcstrans: %q", err)
		}
		defer conn.Close()

		for _, tt := range tests {
			got, err := conn.RawToColor(tt.orig)
			check(t, err, tt.err, got, tt.want)
		}
	})

	t.Run("Test RawToTrans", func(t *testing.T) {
		tests := []struct {
			orig, want, err, name string
		}{
			{
				orig: "staff_u:staff_r:staff_t:s0-s15:c0.c1023",
				want: "staff_u:staff_r:staff_t:SystemLow-SystemHigh",
				name: "test a valid RawToTrans",
			},
		}

		conn, err := New()
		if err != nil {
			t.Fatalf("failed to connect to mcstrans: %q", err)
		}
		defer conn.Close()

		for _, tt := range tests {
			got, err := conn.RawToTrans(tt.orig)
			check(t, err, tt.err, got, tt.want)
		}
	})
}