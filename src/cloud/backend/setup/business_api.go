package setup

import (
	"github.com/gorilla/mux"
	"ocelot/backend/security"
	"ocelot/backend/tools"
)

type BusinessModule struct {
	appInitializer *ApplicationInitializer
}

func (b *BusinessModule) InitializeApplication() {
	b.appInitializer.InitializeApplicationInternally()
}

func ProvideBusinessModule(router *mux.Router, config *tools.GlobalConfig, securityModule *security.SecurityModule) BusinessModule {
	appInitializer := ProvideAppInitializer(router, config, securityModule)
	return BusinessModule{&appInitializer}
}
