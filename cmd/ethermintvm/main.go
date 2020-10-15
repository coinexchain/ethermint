package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagNoMemory     = "nomemory"
	flagNoStack      = "nostack"
	flagNoStorage    = "nostorage"
	flagNoReturnData = "noreturndata"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ethermintvm",
		Short: "EthermintVM CLI",
	}

	rootCmd.AddCommand(stateTestCmd())
	//executor := cli.PrepareBaseCmd(rootCmd, "EM", app.DefaultNodeHome)
	//rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
	//	0, "Assert registered invariants every N blocks")
	executor := cli.Executor{rootCmd, os.Exit}
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func stateTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statetest <file>",
		Short: "Executes the given state tests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			testFile := args[0]
			noMemory, _ := cmd.Flags().GetBool(flagNoMemory)
			noStack, _ := cmd.Flags().GetBool(flagNoStack)
			noStorage, _ := cmd.Flags().GetBool(flagNoStorage)
			noReturnData, _ := cmd.Flags().GetBool(flagNoReturnData)
			return runTestFile(testFile, noMemory, noStack, noStorage, noReturnData)
		},
	}
	cmd.Flags().Bool(flagNoMemory, true, "disable memory output")
	cmd.Flags().Bool(flagNoStack, true, "disable stack output")
	cmd.Flags().Bool(flagNoStorage, true, "disable storage output")
	cmd.Flags().Bool(flagNoReturnData, true, "disable return data output")
	return cmd
}
