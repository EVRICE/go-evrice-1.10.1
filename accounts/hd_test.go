// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package accounts

import (
	"fmt"
	"reflect"
	"testing"
)

// Tests that HD derivation paths can be correctly parsed into our internal binary
// representation.
func TestHDPathParsing(t *testing.T) {
	tests := []struct {
		input  string
		output DerivationPath
	}{
		// Plain absolute derivation paths
		{"m/44'/1020'/0'/0", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0}},
		{"m/44'/1020'/0'/128", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 128}},
		{"m/44'/1020'/0'/0'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 0}},
		{"m/44'/1020'/0'/128'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 128}},
		{"m/2147483692/2147484668/2147483648/0", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0}},
		{"m/2147483692/2147484668/2147483648/2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 0}},

		// Plain relative derivation paths
		{"0", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0}},
		{"128", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 128}},
		{"0'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 0}},
		{"128'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 128}},
		{"2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 0}},

		// Hexadecimal absolute derivation paths
		{"m/0x2C'/0x3fc'/0x00'/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0}},
		{"m/0x2C'/0x3fc'/0x00'/0x80", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 128}},
		{"m/0x2C'/0x3fc'/0x00'/0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 0}},
		{"m/0x2C'/0x3fc'/0x00'/0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 128}},
		{"m/0x8000002C/0x800003fc/0x80000000/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0}},
		{"m/0x8000002C/0x800003fc/0x80000000/0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0x80000000 + 0}},

		// Hexadecimal relative derivation paths
		{"0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0}},
		{"0x80", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 128}},
		{"0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 0}},
		{"0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 128}},
		{"0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0, 0x80000000 + 0}},

		// Weird inputs just to ensure they work
		{"	m  /   44			'\n/\n   1020	\n\n\t'   /\n0 ' /\t\t	0", DerivationPath{0x80000000 + 44, 0x80000000 + 1020, 0x80000000 + 0, 0}},

		// Invalid derivation paths
		{"", nil},              // Empty relative derivation path
		{"m", nil},             // Empty absolute derivation path
		{"m/", nil},            // Missing last derivation component
		{"/44'/1020'/0'/0", nil}, // Absolute path without m prefix, might be user error
		{"m/2147483648'", nil}, // Overflows 32 bit integer
		{"m/-1'", nil},         // Cannot contain negative number
	}
	for i, tt := range tests {
		if path, err := ParseDerivationPath(tt.input); !reflect.DeepEqual(path, tt.output) {
			t.Errorf("test %d: parse mismatch: have %v (%v), want %v", i, path, err, tt.output)
		} else if path == nil && err == nil {
			t.Errorf("test %d: nil path and error: %v", i, err)
		}
	}
}

func testDerive(t *testing.T, next func() DerivationPath, expected []string) {
	t.Helper()
	for i, want := range expected {
		if have := next(); fmt.Sprintf("%v", have) != want {
			t.Errorf("step %d, have %v, want %v", i, have, want)
		}
	}
}

func TestHdPathIteration(t *testing.T) {
	testDerive(t, DefaultIterator(DefaultBaseDerivationPath),
		[]string{
			"m/44'/1020'/0'/0/0", "m/44'/1020'/0'/0/1",
			"m/44'/1020'/0'/0/2", "m/44'/1020'/0'/0/3",
			"m/44'/1020'/0'/0/4", "m/44'/1020'/0'/0/5",
			"m/44'/1020'/0'/0/6", "m/44'/1020'/0'/0/7",
			"m/44'/1020'/0'/0/8", "m/44'/1020'/0'/0/9",
		})

	testDerive(t, DefaultIterator(LegacyLedgerBaseDerivationPath),
		[]string{
			"m/44'/1020'/0'/0", "m/44'/1020'/0'/1",
			"m/44'/1020'/0'/2", "m/44'/1020'/0'/3",
			"m/44'/1020'/0'/4", "m/44'/1020'/0'/5",
			"m/44'/1020'/0'/6", "m/44'/1020'/0'/7",
			"m/44'/1020'/0'/8", "m/44'/1020'/0'/9",
		})

	testDerive(t, LedgerLiveIterator(DefaultBaseDerivationPath),
		[]string{
			"m/44'/1020'/0'/0/0", "m/44'/1020'/1'/0/0",
			"m/44'/1020'/2'/0/0", "m/44'/1020'/3'/0/0",
			"m/44'/1020'/4'/0/0", "m/44'/1020'/5'/0/0",
			"m/44'/1020'/6'/0/0", "m/44'/1020'/7'/0/0",
			"m/44'/1020'/8'/0/0", "m/44'/1020'/9'/0/0",
		})
}
