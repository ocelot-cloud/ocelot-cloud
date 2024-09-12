package apps_new

import (
	"database/sql"
	"log"
	"ocelot/backend/security"
)

// TODO
/*
type appCreationForm struct: {src_domain, maintainer, app, tag, blob}
add app: func(appEntry) err
type appForm struct: {src_domain, maintainer, app, tag}
delete app: func(appForm) err
load app: func(appForm) (blob, error)
*/

var db *sql.DB

func InitAppRepo() {
	if security.DB == nil {
		security.InitializeDatabaseWithSource(security.DatabaseFile)
		if security.DB == nil {
			log.Fatalf("Failed to initialize app repo")
		}
	}
	db = security.DB
}
