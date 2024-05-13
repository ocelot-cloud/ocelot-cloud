module ocelot/backend

go 1.21.6

require (
	github.com/gorilla/mux v1.8.1
	ocelot/business v0.0.0
	ocelot/security v0.0.0
	ocelot/tools v0.0.0
)

replace (
	ocelot/business => ./modules/business
	ocelot/component-tests => ./modules/component-tests
	ocelot/security => ./modules/security
	ocelot/tools => ./modules/tools
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
