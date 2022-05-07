package registrar

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/nuclearblock/archgregator/node"

	"github.com/nuclearblock/archgregator/modules/telemetry"

	"github.com/nuclearblock/archgregator/logging"

	"github.com/nuclearblock/archgregator/types/config"

	"github.com/nuclearblock/archgregator/modules"
	"github.com/nuclearblock/archgregator/modules/messages"

	"github.com/nuclearblock/archgregator/database"
)

// Context represents the context of the modules registrar
type Context struct {
	archgregatorConfig config.Config
	SDKConfig          *sdk.Config
	EncodingConfig     *params.EncodingConfig
	Database           database.Database
	Proxy              node.Node
	Logger             logging.Logger
}

// NewContext allows to build a new Context instance
func NewContext(
	parsingConfig config.Config, sdkConfig *sdk.Config, encodingConfig *params.EncodingConfig,
	database database.Database, proxy node.Node, logger logging.Logger,
) Context {
	return Context{
		archgregatorConfig: parsingConfig,
		SDKConfig:          sdkConfig,
		EncodingConfig:     encodingConfig,
		Database:           database,
		Proxy:              proxy,
		Logger:             logger,
	}
}

// Registrar represents a modules registrar. This allows to build a list of modules that can later be used by
// specifying their names inside the TOML configuration file.
type Registrar interface {
	BuildModules(context Context) modules.Modules
}

// ------------------------------------------------------------------------------------------------------------------

var (
	_ Registrar = &EmptyRegistrar{}
)

// EmptyRegistrar represents a Registrar which does not register any custom module
type EmptyRegistrar struct{}

// BuildModules implements Registrar
func (*EmptyRegistrar) BuildModules(_ Context) modules.Modules {
	return nil
}

// ------------------------------------------------------------------------------------------------------------------

var (
	_ Registrar = &DefaultRegistrar{}
)

// DefaultRegistrar represents a registrar that allows to handle the default Archgregator modules
type DefaultRegistrar struct {
	parser messages.MessageAddressesParser
}

// NewDefaultRegistrar builds a new DefaultRegistrar
func NewDefaultRegistrar(parser messages.MessageAddressesParser) *DefaultRegistrar {
	return &DefaultRegistrar{
		parser: parser,
	}
}

// BuildModules implements Registrar
func (r *DefaultRegistrar) BuildModules(ctx Context) modules.Modules {
	return modules.Modules{
		messages.NewModule(r.parser, ctx.EncodingConfig.Marshaler, ctx.Database),
		telemetry.NewModule(ctx.archgregatorConfig),
	}
}

// ------------------------------------------------------------------------------------------------------------------

// GetModules returns the list of module implementations based on the given module names.
// For each module name that is specified but not found, a warning log is printed.
func GetModules(mods modules.Modules, names []string, logger logging.Logger) []modules.Module {
	var modulesImpls []modules.Module
	for _, name := range names {
		module, found := mods.FindByName(name)
		if found {
			modulesImpls = append(modulesImpls, module)
		} else {
			logger.Error("Module is required but not registered. Be sure to register it using registrar.RegisterModule", "module", name)
		}
	}
	return modulesImpls
}
