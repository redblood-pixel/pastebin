//go:build unit
// +build unit

package hash

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		correct  bool
	}{
		{"simple pass1 cost 14", "abcdef123", "$2a$14$XAJ4PjZu0Kl2N7l6OdmmteRwmdEFuSIAB3yLicr99IlmnhYVm1sWi", true},
		{"simple pass2 cost 14", "aaafff", "$2a$14$5cpIbHuXbEUSSC3unc8qPO6ow44yzIWFbkV2HnNR13RSCD6w0jZHe", true},
		{"jojo cost 9", "whitesnakemegagodpuchiiloveyou", "$2a$09$H4w/1c/QpKdBdZeD1nI49.6hNBs4P6HQ2/9nFZdiqfVwjoxpOQT.2", true},
		{"jojo cost 14", "whitesnakemegagodpuchiiloveyou", "$2a$14$yCZt1k3IFjDg73HSJNt1se3Xa/WdoOz.IUI9rzLv5FtLEp996FgeW", true},
		{"simple pass1 cost 9", "jQtYref8", "$2a$09$6dEtrQ9ip/KQZc8XG42Ao.dlxiit7V0CcStX8vnLXOtRNdNsYqLPe", true},
		{"simple pass1 cost 9", "jQtYref8", "$2a$09$6dEtrQ9ip/KQZc8XG42Ao.dlxiit7V0CcStX8vnLXOtRNdNsYqLee", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.correct != CheckPassword(tt.password, tt.hash) {
				t.Errorf("password %s not a valid for %s", tt.password, tt.hash)
			}
		})
	}
}
