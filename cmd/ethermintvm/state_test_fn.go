package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
)

type StatetestResult struct {
	Name  string      `json:"name"`
	Pass  bool        `json:"pass"`
	Fork  string      `json:"fork"`
	Error string      `json:"error,omitempty"`
	State *state.Dump `json:"state,omitempty"`
}

func runTestFile(testFile string,
	noMemory, noStack, noStorage, noReturnData bool) error {

	config := &vm.LogConfig{
		DisableMemory:     noMemory,
		DisableStack:      noStack,
		DisableStorage:    noStorage,
		DisableReturnData: noReturnData,
	}
	tracer := vm.NewJSONLogger(config, os.Stderr)

	src, err := ioutil.ReadFile(testFile)
	if err != nil {
		return err
	}

	var tests map[string]StateTest
	if err = json.Unmarshal(src, &tests); err != nil {
		return err
	}

	cfg := vm.Config{
		Tracer: tracer,
		Debug:  true,
	}
	results := make([]StatetestResult, 0, len(tests))
	for key, test := range tests {
		for _, st := range test.Subtests() {
			// Run the test and aggregate the result
			result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
			_, err := test.Run(st, cfg)
			// print state root for evmlab tracing
			//if ctx.GlobalBool(MachineFlag.Name) && state != nil {
			//	fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%x\"}\n", state.IntermediateRoot(false))
			//}
			if err != nil {
				// Test failed, mark as so and dump any state to aid debugging
				result.Pass, result.Error = false, err.Error()
				//if ctx.GlobalBool(DumpFlag.Name) && state != nil {
				//	dump := state.RawDump(false, false, true)
				//	result.State = &dump
				//}
			}

			results = append(results, *result)
		}
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	return nil
}
