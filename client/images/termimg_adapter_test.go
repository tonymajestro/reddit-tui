package images

import (
	"testing"
)

func TestParseProtocol(t *testing.T) {
	tests := []struct {
		name     string
		protocol RenderProtocol
		wantErr  bool
	}{
		{"auto", ProtocolAuto, false},
		{"kitty", ProtocolKitty, false},
		{"iterm2", ProtocolITerm2, false},
		{"sixel", ProtocolSixel, false},
		{"halfblocks", ProtocolHalfblocks, false},
		{"invalid", RenderProtocol("not-a-protocol"), true},
	}

	for _, tt := range tests {
		_, err := parseProtocol(tt.protocol)
		if tt.wantErr && err == nil {
			t.Fatalf("%s: expected error", tt.name)
		}
		if !tt.wantErr && err != nil {
			t.Fatalf("%s: unexpected error %v", tt.name, err)
		}
	}
}
