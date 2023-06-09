package magma

import (
	"fmt"

	"github.com/candostyavuz/magma/pkg/magma"

	"github.com/spf13/cobra"
)

var genevmCmd = &cobra.Command{
	Use:     "gen-evm-spec [supported-apis-file] | Flags: [--chain-name] , [--chain-idx], [--imports]",
	Aliases: []string{"genevm"},
	Short:   "Generates a valid spec file from a list of supported api calls",
	Long: `Generates a valid spec file from a list of supported api calls.
	Currently, the only supported input format for the spec file is txt file.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("File is: ", args[0])
		imports, err := cmd.Flags().GetStringArray("imports")
		fmt.Println("Imported specs: ", imports)
		if err != nil {
			return err
		}
		chainName, err := cmd.Flags().GetString("chain-name")
		if err != nil {
			return err
		}
		chainIdx, err := cmd.Flags().GetString("chain-idx")
		if err != nil {
			return err
		}

		err = magma.GenerateSpec(args[0], chainName, chainIdx, imports)
		return err
	},
}

func init() {
	genevmCmd.Flags().String("chain-name", "", "Chain Name")
	genevmCmd.Flags().String("chain-idx", "", "Chain Index")
	genevmCmd.Flags().StringArray("imports", nil, "Imports for this spec")
	genevmCmd.MarkFlagRequired("chain-name")
	genevmCmd.MarkFlagRequired("chain-idx")
	rootCmd.AddCommand(genevmCmd)
}
