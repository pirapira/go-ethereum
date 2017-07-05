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

package tests

import (
	"fmt"
	"testing"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) // gasLimit > 256 bits
	st.skipLoad(`^stStackTests/shallowStackOK\.json`)          // bad hex encoding
	st.skipLoad(`^stTransactionTest/zeroSigTransa.*\.json`)    // metropolis-related
	// Expected failures:
	st.fails(`^stCodeSizeLimit/codesizeOOGInvalidSize\.json/(Frontier|Homestead)`,
		"code size limit implementation is not conditional on fork")
	st.fails(`^stCallCreateCallCodeTest/createJS_ExampleContract.json`,
		"bug in test")

	st.walk(t, stateTestDir, func(t *testing.T, test *StateTest) {
		for _, subtest := range test.Subtests() {
			subtest := subtest
			name := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
			t.Run(name, func(t *testing.T) {
				if subtest.Fork != "EIP150" {
					t.Skip("only interested in EIP150")
				}
				if subtest.Index != 6 {
					t.Skip("not interested in anything except 6")
				}
				if err := st.checkFailure(t, test.Run(subtest)); err != nil {
					t.Error(err)
				}
			})
		}
	})
}
