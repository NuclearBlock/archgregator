package genesis

import (
	"github.com/spf13/cobra"

	parsecmdtypes "github.com/nuclearblock/archgregator/cmd/parse/types"

	nodeconfig "github.com/nuclearblock/archgregator/node/config"
	"github.com/nuclearblock/archgregator/types/utils"
)

const (
	flagPath = "genesis-file-path"
)

// NewGenesisCmd returns the Cobra command allowing to parse the genesis file
func NewGenesisCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis-file",
		Short: "Parse the genesis file",
		Long: `
Parse the genesis file only. 
Note that the modules built will NOT have access to the node as they are only supposed to deal with the genesis
file itself and not the on-chain data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the configuration
			cfg, err := parsecmdtypes.ReadConfig(parseConfig)
			if err != nil {
				return err
			}

			// Set the node to be of type None so that the node won't be built
			cfg.Node.Type = nodeconfig.TypeNone

			// Get the file path
			genesisFilePath := cfg.Parser.GenesisFilePath
			customPath, _ := cmd.Flags().GetString(flagPath)
			if customPath != "" {
				genesisFilePath = customPath
			}

			// Read the genesis file
			genDoc, err := utils.ReadGenesisFileGenesisDoc(genesisFilePath)
			if err != nil {
				return err
			}

			genState, err := utils.GetGenesisState(genDoc)
			if err != nil {
				return err
			}

			utils.HandleGenesis(genDoc, genState)

			return nil
		},
	}

	cmd.Flags().String(flagPath, "", "Path to the genesis file to be used. If empty, the path will be taken from the config file")

	return cmd
}
