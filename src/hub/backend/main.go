package main

import (
	"context"
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"os"
)

func init() {
	Logger = shared.ProvideLogger(os.Getenv("LOG_LEVEL"))
	Logger.Info("log level set to: %s", shared.GetLogLevel())
}

func main() {
	if profile == TEST {
		Logger.Info("profile is: TEST")
	} else if profile == PROD {
		Logger.Info("profile is: PROD")
	} else {
		Logger.Fatal("unknown profile: %d", profile)
	}

	initializeDatabase()
	mux := http.NewServeMux()
	initializeHandlers(mux)

	Logger.Info("server starting on port %s", port)
	err := http.ListenAndServe(":"+port, utils.GetCorsDisablingHandler(mux))
	if err != nil {
		Logger.Fatal("Server stopped: %v", err)
	}
}

// TODO shift the initialization functions into the "setup" package?
func initializeDatabase() {
	if profile == TEST {
		InitializeDatabaseWithSource(":memory:")
		Logger.Warn("initializing database only in-memory - when this hub application stops, all data will be deleted")
	} else {
		InitializeDatabaseWithSource(databaseFile)
	}
	err := createAdminUserIfNotExistent()
	if err != nil {
		Logger.Fatal("Admin user creation failed: %v", err)
	}
}

func createAdminUserIfNotExistent() error {
	// TODO Check if admin user exists. If not take it from the env variables. If not existent, crash.
	// TODO Add tests: 1) neither admin in repo nor in envs -> crash, 2) no admin in repo, but in envs -> no crash, 3) admin in repo, but not in envs -> no crash
	return nil
}

type Route struct {
	path    string
	handler http.HandlerFunc
}

func initializeHandlers(mux *http.ServeMux) {
	unprotectedRoutes := []Route{
		{loginPath, loginHandler},
		{registrationPath, registrationHandler},
		{downloadPath, downloadHandler},
		{getTagsPath, getTagsHandler},
		{searchAppsPath, searchAppsHandler},
	}

	protectedRoutes := []Route{
		{authCheckPath, authCheckHandler},
		{tagUploadPath, tagUploadHandler},
		{tagDeletePath, tagDeleteHandler},
		{changePasswordPath, changePasswordHandler},
		{appCreationPath, appCreationHandler},
		{appGetListPath, appGetListHandler},
		{appDeletePath, appDeleteHandler},
		{deleteUserPath, userDeleteHandler},
		{logoutPath, logoutHandler},
	}

	if profile == TEST {
		Logger.Warn("opening unprotected full data wipe endpoint meant for testing only")
		unprotectedRoutes = append(unprotectedRoutes, Route{wipeDataPath, wipeDataHandler})

		sampleUser := "sample"
		err := repo.CreateUser(&RegistrationForm{sampleUser, "password", "admin@admin.com"})
		if err != nil {
			Logger.Fatal("Failed to create '%s' user: %v.", sampleUser, err)
		}
		Logger.Warn("created '%s' user with weak password for manual testing", sampleUser)
		loadSampleApp()
	}

	registerUnprotectedRoutes(mux, unprotectedRoutes)
	registerProtectedRoutes(mux, protectedRoutes)
}

func loadSampleApp() {
	sampleUser := "sampleuser"
	sampleApp := "nginxdefault"
	repo.CreateUser(&RegistrationForm{sampleUser, "password", "sample@sample.com"})
	repo.CreateApp(sampleUser, sampleApp)
	dirPath := "./assets/sampleuser_nginxdefault"
	zipBytes, err := utils.ZipDirectoryToBytes(dirPath)
	if err != nil {
		Logger.Fatal("Failed to zip directory: %v", err)
	}

	appId, err := repo.GetAppId(sampleUser, sampleApp)
	if err != nil {
		Logger.Fatal("Failed to get app ID: %v", err)
	}
	err = repo.CreateTag(appId, "0.0.1", zipBytes)
	if err != nil {
		Logger.Fatal("Failed to create sample tag: %v", err)
	}
}

// getUserFromContext Since only authenticated users are added to the context, it only works in protected handlers.
func getUserFromContext(r *http.Request) string {
	return r.Context().Value("user").(string)
}

func registerUnprotectedRoutes(mux *http.ServeMux, routes []Route) {
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}
}

func registerProtectedRoutes(mux *http.ServeMux, routes []Route) {
	for _, r := range routes {
		mux.Handle(r.path, authMiddleware(r.handler))
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := checkAuthentication(w, r)
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
