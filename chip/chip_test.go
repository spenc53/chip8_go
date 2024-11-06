package chip

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		ins     uint16
		wantIns instruction
	}{
		{
			name: "Clear",
			ins:  0x00E0,
			wantIns: instruction{
				Ins: 0x00E0,
				I:   0x0,
				X:   0x0,
				Y:   0xE,
				NN:  0xE0,
				NNN: 0x0E0,
			},
		},
		{
			name: "Jump",
			ins:  0x1228,
			wantIns: instruction{
				Ins: 0x1228,
				I:   0x1,
				X:   0x2,
				Y:   0x2,
				N:   0x8,
				NN:  0x28,
				NNN: 0x228,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotIns := decode(tc.ins)
			if diff := cmp.Diff(gotIns, tc.wantIns); diff != "" {
				t.Errorf("Diff between instructions (-got, +want): %s", diff)
			}
		})
	}
}
