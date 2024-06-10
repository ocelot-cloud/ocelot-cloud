module component_tests

go 1.21.6

require (
	ocelot/tools v0.0.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	ocelot/business => ../business
	ocelot/security => ../security
	ocelot/tools => ../tools
)
