package main

import (
	"os"

	"github.com/nuclearblock/archgregator/cmd/parse/types"

	"github.com/nuclearblock/archgregator/modules/messages"
	"github.com/nuclearblock/archgregator/modules/registrar"

	"github.com/nuclearblock/archgregator/cmd"
)

func main() {
	// archgregatorConfig the runner
	config := cmd.NewConfig("archgregator").
		WithParseConfig(types.NewConfig().
			WithRegistrar(registrar.NewDefaultRegistrar(
				messages.CosmosMessageAddressesParser,
			)),
		)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
