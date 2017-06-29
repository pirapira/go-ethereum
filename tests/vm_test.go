// Copyright 2014 The go-ethereum Authors
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
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm"
)

func TestVM(t *testing.T) {
	t.Parallel()
	vmt := new(testMatcher)
	vmt.fails("^vmSystemOperationsTest.json/createNameRegistrator$", "fails without parallel execution")
	vmt.skipShortMode("^vmPerformanceTest.json")
	vmt.skipShortMode("^vmInputLimits(Light)?.json")

	const traceLimit = 4000000

	vmt.walk(t, vmTestDir, func(t *testing.T, test *VMTest) {
		if err := vmt.checkFailure(t, test.Run(vm.Config{})); err != nil {
			t.Error(err)

			// Output struct logs for debugging unless there's too much execution.
			if test.json.Exec.GasLimit > traceLimit {
				return
			}
			tracer := vm.NewStructLogger(nil)
			test.Run(vm.Config{Debug: true, Tracer: tracer})
			buf := new(bytes.Buffer)
			vm.WriteTrace(buf, tracer.StructLogs())
			if buf.Len() == 0 {
				t.Log("no vm operation logs generated")
			} else {
				t.Log("vm operation log:\n" + buf.String())
			}
		}
	})
}
