package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestDecideUpdate(t *testing.T) {
	tests := []struct {
		name            string
		current         string
		latest          string
		fetchErr        error
		wantProceed     bool
		wantMsgContains string
	}{
		{
			name:            "dev build short-circuits regardless of latest",
			current:         "dev",
			latest:          "v1.2.3",
			fetchErr:        nil,
			wantProceed:     false,
			wantMsgContains: "dev build",
		},
		{
			name:            "dev build short-circuits even on fetch error",
			current:         "dev",
			latest:          "",
			fetchErr:        errors.New("network down"),
			wantProceed:     false,
			wantMsgContains: "dev build",
		},
		{
			name:            "current equals latest -> noop with friendly message",
			current:         "v1.2.3",
			latest:          "v1.2.3",
			fetchErr:        nil,
			wantProceed:     false,
			wantMsgContains: "v1.2.3",
		},
		{
			name:            "current older than latest -> proceed silently",
			current:         "v1.2.2",
			latest:          "v1.2.3",
			fetchErr:        nil,
			wantProceed:     true,
			wantMsgContains: "",
		},
		{
			name:            "current newer than latest -> noop (don't downgrade)",
			current:         "v1.3.0",
			latest:          "v1.2.3",
			fetchErr:        nil,
			wantProceed:     false,
			wantMsgContains: "v1.3.0",
		},
		{
			name:            "fetch error on non-dev -> proceed (network flake shouldn't block)",
			current:         "v1.2.2",
			latest:          "",
			fetchErr:        errors.New("network down"),
			wantProceed:     true,
			wantMsgContains: "",
		},
		{
			name:            "empty latest with no error -> proceed (treat as unknown)",
			current:         "v1.2.2",
			latest:          "",
			fetchErr:        nil,
			wantProceed:     true,
			wantMsgContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProceed, gotMsg := decideUpdate(tt.current, tt.latest, tt.fetchErr)
			if gotProceed != tt.wantProceed {
				t.Errorf("decideUpdate(%q, %q, %v) proceed = %v, want %v",
					tt.current, tt.latest, tt.fetchErr, gotProceed, tt.wantProceed)
			}
			if tt.wantMsgContains == "" {
				if gotMsg != "" {
					t.Errorf("decideUpdate(%q, %q, %v) message = %q, want empty",
						tt.current, tt.latest, tt.fetchErr, gotMsg)
				}
			} else if !strings.Contains(gotMsg, tt.wantMsgContains) {
				t.Errorf("decideUpdate(%q, %q, %v) message = %q, want it to contain %q",
					tt.current, tt.latest, tt.fetchErr, gotMsg, tt.wantMsgContains)
			}
		})
	}
}
