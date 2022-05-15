package main

import (
	"os"

	"github.com/nuclearblock/archgregator/cmd"
	"github.com/nuclearblock/archgregator/cmd/parse/types"
	"github.com/nuclearblock/archgregator/types/config"
)

func main() {

	parseCfg := types.NewConfig().WithEncodingConfigBuilder(config.MakeEncodingConfig)

	config := cmd.NewConfig("archgregator").WithParseConfig(parseCfg)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
