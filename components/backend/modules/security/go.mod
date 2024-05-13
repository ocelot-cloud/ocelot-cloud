module ocelot/security

go 1.21.6

require (
	ocelot/tools v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
)

replace ocelot/tools => ../tools

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
)
