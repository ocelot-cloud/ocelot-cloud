package internal

import (
	"github.com/ocelot-cloud/shared"
)

var Logger = shared.ProvideLogger()
var StackFileDir string
var CoreStackFileDir = "stacks/core"

// TODO GlobalConfig instance should be put here, instead of being distributed as function argument to all units.
