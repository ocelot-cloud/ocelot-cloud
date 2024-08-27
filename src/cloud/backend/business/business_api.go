package business

import (
	"github.com/gorilla/mux"
	"ocelot/backend/business/internal"
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

type BusinessModule struct {
	appInitializer *internal.ApplicationInitializer
}

func (b *BusinessModule) InitializeApplication() {
	b.appInitializer.InitializeApplicationInternally()
}

func ProvideBusinessModule(router *mux.Router, config *tools.GlobalConfig, securityModule *security.SecurityModule) BusinessModule {
	appInitializer := internal.ProvideAppInitializer(router, config, securityModule)
	return BusinessModule{&appInitializer}
}
