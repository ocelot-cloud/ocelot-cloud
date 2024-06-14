package business

import (
	"github.com/gorilla/mux"
	"github.com/ocelot-cloud/shared"
	"ocelot/business/internal"
	"ocelot/security"
	"ocelot/tools"
)

type BusinessModule struct {
	appInitializer *internal.ApplicationInitializer
}

func (b *BusinessModule) InitializeApplication() {
	b.appInitializer.InitializeApplicationInternally()
}

func ProvideBusinessModule(router *mux.Router, config *tools.GlobalConfig, securityModule *security.SecurityModule) BusinessModule {
	internal.Logger = shared.ProvideLogger()
	appInitializer := internal.ProvideAppInitializer(router, config, securityModule)
	return BusinessModule{&appInitializer}
}
