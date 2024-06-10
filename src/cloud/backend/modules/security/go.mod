module ocelot/security

go 1.21.6

require (
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
	ocelot/tools v0.0.0
)

replace ocelot/tools => ../tools

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
