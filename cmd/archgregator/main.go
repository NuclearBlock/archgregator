package main

import (
	"os"

	"github.com/nuclearblock/archgregator/cmd"
	"github.com/nuclearblock/archgregator/cmd/parse/types"
	"github.com/nuclearblock/archgregator/types/config"
)

func main() {

	parseCfg := types.NewConfig().WithEncodingConfigBuilder(config.MakeEncodingConfig)

	// archgregatorConfig the runner
	config := cmd.NewConfig("archgregator").WithParseConfig(parseCfg)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// getBasicManagers returns the various basic managers that are used to register the encoding to
// // support custom messages.
// // This should be edited by custom implementations if needed.
// func getBasicManagers() []module.BasicManager {
// 	return []module.BasicManager{
// 		gaiaapp.ModuleBasics,
// 		cmdxapp.ModuleBasics,
// 	}
// }
