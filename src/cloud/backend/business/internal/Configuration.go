package internal

import (
	"github.com/ocelot-cloud/shared"
)

var Logger = shared.ProvideLogger("info") // TODO use global logger instead
var StackFileDir string
var CoreStackFileDir = "stacks/core"
