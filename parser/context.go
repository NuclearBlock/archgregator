package parser

import (
	"github.com/archway-network/archway/app/params"
	"github.com/nuclearblock/archgregator/database"
	"github.com/nuclearblock/archgregator/logging"
	"github.com/nuclearblock/archgregator/node"
)

// Context represents the context that is shared among different workers
type Context struct {
	EncodingConfig *params.EncodingConfig
	Node           node.Node
	Database       database.Database
	Logger         logging.Logger
}

// NewContext builds a new Context instance
func NewContext(
	encodingConfig *params.EncodingConfig,
	proxy node.Node,
	db database.Database,
	logger logging.Logger,
) *Context {
	return &Context{
		EncodingConfig: encodingConfig,
		Node:           proxy,
		Database:       db,
		Logger:         logger,
	}
}
