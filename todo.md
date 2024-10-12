# TODO

* no access to app should redirect to login page
* Use compose-go and docker sdk to communicate with the docker daemon
* GlobalConfig and db objects should be globally available via tools.xxx
* this should be a warning if profile is TEST: "tools/config.go:112 > Profile is: TEST"
* security
  1. in handler: getAuth -> auth { user, isAdmin }, check if handler needs admin access or just user access
  2. cloud authorization - hasAccess(user, app), if app is access, check if user should have access to it
* rename "tag" to "version"?
* Admin should be asked to change password after first login?
* To protect the cloud endpoints, do a similar approach as in the hub,see "registerProtectedRoutes", but add a distinctions for handlers that require admin role and handlers that requires user role
* origin checks should be done for all endpoints. Even if they are not secured like the login endpoint.
* hub: instead of doing "checkAuthentication" at the beginning of each protected handler, I should add a middleware doing that. E.g. "ignore unprotected paths/handlers, but protected one should be checked and the users context info should be added for subsequent requests":

```
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, err := authenticate(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "user", user)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
```

* hub: stuff like this should be hidden in the repo, because this construct is used multiple times in handlers which causes duplication:

```
    err = repo.CreateTag(user, tagUpload.App, tagUpload.Tag, tagUpload.Content)
    if err != nil {
        Logger.Error("creating tag '%s' for user '%s' failed: %v", tagUpload.App, user, err)
        http.Error(w, "invalid input", http.StatusInternalServerError)
        return
    }
```

* GUI: I dont want an start/stop/open button for each app. I want a list of apps, select one and then click on one of the buttons below the table.

* PROD deploy -> hub does not work, probably because root domain "http://ocelot-cloud.localhost/hub/registration" instead of http://localhost/hub/registration or so

* in PROD, when visiting http://localhost, you should be redirect to http://ocelot-cloud.localhost

* Extract the docker service into its own package

* Currently, there is a horizontally layered architecture in both components. Convert it to a service oriented architecture.
* refactor the apps logic:
  * I think I need more architecture, e.g. more packages, for clearer structuring
  * there should be high level unit tests, e.g. startApp(...) -> Downloading, Starting, Available -> StopApp() -> Stopping, Uninitialized
  * app assets should be stored in an sqlite db
  * integrate hub for downloading

* make proper login endpoint with real cookie etc.

* Then integrate security into the "cloud-client"?
  
  * ci-runner test cloud backend/frontend etc, do they need adaption?
  * structured logging, see my notes for specification
  * log entry: "user has a valid cookie and is allowed to access protected backend functions" -> which user?
    * a user requested the frontend resources, who? anonymous
    * user has a valid cookie and is allowed to access protected backend functions, who?
    * login logic called, by who?
  * It would be cool after an acceptance test, if I can see the logs of the backend in an git-ignored directory: test-logs/acceptance.txt,backend.txt
  * Don't get rid of dummy stacks But rather put them in the hub?
  * Not sure. Should CliEvaluator_test be re-implemented?
  * In docker compose, there should be
  * I want to install ocelot without cloning git. Simply "wget" the docker-compose.yml or a small bash script and go. Git clone, only when you want to build it from scratch.
  * Introduce structured logging
  * Also use TEST and PROD profiles in hub. Also print which is used right now.
  * frontend and acceptance have bash script which can be deleted.
  * Error in hub GUI: wrong app shows password validation stuff, although it worked to create the app. Maybe create app must reset the "submitted" flag. Also check in acceptance test that this does not happen again
  * Modernize cloud GUI. Also use new vuejs3 + ts syntax in Home.vue
  * Can this simply be replaced by a blank string? -> import.meta.env.VITE_BASE_URL

* get rid of the "method" argument in doRequest, all requests are method independent

* Make hub compatible with cockroackdb

* cookies should be hashed in database

* Password min length should be 16 chars?

* open a new connection to db for each query and close it afterwards, set timeout of 5 min or so?

* always log errors on first occurrence think I didn't do that in the repo code?

* delete client side cookie after logout (only server side right now)

* frontend does not need to be build when I want to use "npm run serve" anyway

* global config should be globally visible, not passed as ar to a submodule.

* get rid of the "internal" dirs in the backend modules. Simply use private/public methods since they are not that big yet.

* dont use "sleep" in test for waiting for services in the backend component tests, better use retries

* create an easy way with which an admin user can reset his password, when he forgets it. Access to server + some bash command to sqlite maybe? Should also be tested.

* frontend build step skipping via ci-runner flag "-f" does not work somehow?

HOST=httpx://localhost throws:
ocelot-cloud  | 2024-08-29T08:18:13Z FTL ../../home/dev/Dokumente/workspace/ocelot-cloud/src/cloud/backend/tools/config.go:46 > Failed to get host params: error when evaluating port from HOST env variable
although the port is not the problem

There is sth wrong in this log entry:

```
Starting server listening on port %!(EXTRA string=8080)

ci-runner deploy logs this:
# ocelot/backend
/usr/bin/ld: /tmp/go-link-1923660456/000011.o: in function `unixDlOpen':
/home/dev/go/pkg/mod/github.com/mattn/go-sqlite3@v1.14.22/sqlite3-binding.c:44707: warning: Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
/usr/bin/ld: /tmp/go-link-1923660456/000017.o: in function `_cgo_9c8efe9babca_C2func_getaddrinfo':
/tmp/go-build/cgo_unix_cgo.cgo2.c:58: warning: Using 'getaddrinfo' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
-> maybe take the "unusual logs search approach" and ingerate it intp the ci tool
```

get rid of these of these logs?
WARN[0000] The "PROFILE" variable is not set. Defaulting to a blank string.
WARN[0000] The "LOG_LEVEL" variable is not set. Defaulting to a blank string.
WARN[0000] /home/dev/Dokumente/workspace/ocelot-cloud/src/cloud/backend/stacks/core/ocelot-cloud/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion

Remove the ocelot-cloud-auth header, when proxying to apps behind it, so they can't steal this cookie. Should be tested by API (both statuses) and once in prod environment, maybe "PROD backend". Nginx config:

```
server {
    listen 80;
    location /api/check {
        if ($http_ocelot_cloud_auth) {
            return 400 'ocelot-cloud-auth header should have been removed by ocelot-cloud proxy, but wasn't';
        }
        return 200 'no ocelot-auth header found as desired';
    }
}
```


Cloud and hub endpoints should be available at "/api/...", right? not sure

reuse:
* helper methods (like doRequest)
* tests, especially security tests
* validation
* cookie generation
* handlers and paths using mux
  * wipeData Handler
  * checkAuth
  * login
  * deleteTag
  * deleteApp
  * downloadTag?
  * nicht register/create user, hub account braucht email, cloud nicht unbedingt
* also keep in mind, that I need roles, so maybe reuse the hub methods but at a "mustBeAdmin" argument or so. There should also be tests for correct protection of admin endpoints admins. -> maybe "/api/admins/..." and "/api/users/..." as endpoints? There should also be a loop like "make a user request to this list of endpoints and get unauthorized response due to lack of admin rights"

readBody -> da müsste wie gesagt noch das "mustBeAdmin" arg rein + "job validation" muss irgendwie ausgelagert werden da die datenstrukturen Anwendungsspezifisch sind.
Es wäre cool, wenn die ReadBody oder validate Funktionen eineige Grundtypen hätten, zB validate username/password, aber mit anderen Typen erweiterbar wären, zB

reuse "doRequest" in frontend cloud logic, also ensure for hub and cloud gui that only POST requests are used.

replace "stacks" by "apps"



currently all app configs are read at the start of the application, e.g.:
"2024-08-31T12:37:58Z DBG apps/yaml_config.go:47 > file assets/local/openproject/app.yml does not exist, providing default config instead"
that means when I download it again, I need to read it again.
actually I want to transfer all of this stuff to sqlite anyhow.

Idea: registerProtectedHandler(path, myHandler, mustbeAdmin=true/false)

Is there a bash script like run-production/development still necessary for acceptance tests?

security implementation:
* try to share logic between modules
* shared packages: sectest (for test code), secprod (for production code)
* secure the exiting endpoints, then implement new ones

stack does not exist, stack already exists, etc
you should be able to stop stacks the are being started

rename stacks to apps
make a unit test for TEST and PROD config
check that the wipe data endpoint is not opened in PROD docker container
AssertCors oder so, sollte evtl auch mit dem cloud client durchgeführt werden.

* add admin-creation at startup logic
* make proper login endpoint with real cookie etc.
* Then integrate security into the "cloud-client"?

write a test that checks whether the starting of backend fails if no valid "HOST" variable is set.

frontend build step skipping does not work somehow?

Extract the docker service into its own package? -> maybe later

for the demo and in general, once downloaded stacks should be deployable without internet connection.

"latest" keyword is forbidden since it is not stable. We want

HOST=http://localhost

HOST=httpx://localhost throws:
ocelot-cloud  | 2024-08-29T08:18:13Z FTL ../../home/dev/Dokumente/workspace/ocelot-cloud/src/cloud/backend/tools/config.go:46 > Failed to get host params: error when evaluating port from HOST env variable
although the port is not the problem

When visiting ocelot-cloud on localhost:8080

Problem: Wenn ich nur mocks in acceptance checke, dann erhalte niemals einen wirklichen Request zu einem dahinter liegenden Service.
Und PROD backend API test, kann das auch nicht leisten, weil die Container nicht zugänglich sind, wenn das backend nicht in einem container liegt.
Okay, das heißt ich extrahiere den Backend PROD Test und nutze einfach keine Mocks in Acceptance Tests.

it should be possible to stop a stack download, I should also take care that when something happens, like a sudden internet disconnection during download, then this should not trap the user in an infinite loop of "downloading stage"

add this?
cmd.Env = append(os.Environ(), "LOG_LEVEL=DEBUG")
buildProxyHandler -> there should be a check like "isOcelotHost" -> redirect to ocelot, else if "is a registeredServiceHost" redirect to docker container, else "send error with message back to GUI"

when visiting "https://my-domain.com", then the user should be redirected to "ocelot-cloud.https://my-domain.com", right?
-> I think this is a good idea. People may want to use the rootDomain for a website or so, maybe?

how to handle conflict when there are multiple instances of the same service? -> I think there should be

security concept: containers run in isolated environments and they can't see each other directly, just the frontend pages -> maybe increase that security concept in the future. For example, there is a separate container for proxying. When a new service is added, say gitlab, the the proxy is restarted and has additional membership in the docker network "gitlab-net" as well as the gitlab container itself. This way, gitlab cant see the other apps. This may be a nice feature for the future for maximum security and complete app isolation.
-> too complicated, ocelot can simply join a network at runtime like this: "docker network connect <network_name> <container_name_or_id>", no extra container needed

Maybe there should also be a mechanism which shows if there are old volumes. Or if a previously unknown image/container is trying to access and old volume. This could be a potential security hazard. E.g. I delete an app but keep the data. Another app has the same name and wants to access the volume.

There is sth wrong in this log entry:
Starting server listening on port %!(EXTRA string=8080)

Does acceptance test make a request to one of the docker containers behind it?

When visitng a server that does not exist/is not available, then this is printed by the backend:
2024/08/29 09:21:04 http: proxy error: dial tcp: lookup nginx-custom-path on 127.0.0.11:53: server misbehaving
-> 1) should be a typical backend log format, message should be clearer, where is it thrown?

Idee: Suche nach invaliden logs. d.h. zB "nimm alle logs, separiere nach '\n', ignoriere leere lines, überprüfe ob jedes line format zum regex passt, falls nicht erstelle einen report mit der message"

ci-runner deploy logs this:
# ocelot/backend
/usr/bin/ld: /tmp/go-link-1923660456/000011.o: in function `unixDlOpen':
/home/dev/go/pkg/mod/github.com/mattn/go-sqlite3@v1.14.22/sqlite3-binding.c:44707: warning: Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
/usr/bin/ld: /tmp/go-link-1923660456/000017.o: in function `_cgo_9c8efe9babca_C2func_getaddrinfo':
/tmp/go-build/cgo_unix_cgo.cgo2.c:58: warning: Using 'getaddrinfo' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
-> maybe take the "unusual logs search approach" and ingerate it intp the ci tool

Take care of logs?
WARN[0000] The "PROFILE" variable is not set. Defaulting to a blank string.
WARN[0000] The "LOG_LEVEL" variable is not set. Defaulting to a blank string.
WARN[0000] /home/dev/Dokumente/workspace/ocelot-cloud/src/cloud/backend/stacks/core/ocelot-cloud/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion

PROD deploy -> hub does not work, probably because root domain "http://ocelot-cloud.localhost/hub/registration" instead of http://localhost/hub/registration or so

ci runner: make a list of envs, and set them to empty string. I encountered that I used export "ASD=true", and that it was used by the processes below.

GUI: The tab title in the browser should be "Ocelot" and I need a favicon

add a docker container that does not start properly (immediately exits), whose download aborts or takes too long so that the user wants to abort it, or whose image is not found in dockerhub. The gui should be able to communicate these cases to the users and enable them to handle it, e.g. abort/delete app/ go back to an older version that worked etc.

There should be some kind of "failed" state, when starting a service did not work or took too long (maybe a tolerance of 5 min?)

Cloud and hub endpoints should be available at "/api/...", right? not sure

reuse:
* helper methods (like doRequest)
* tests, especially security tests
* validation
* cookie generation
* handlers and paths using mux
  * wipeData Handler
  * checkAuth
  * login
  * deleteTag
  * deleteApp
  * downloadTag?
  * nicht register/create user, hub account braucht email, cloud nicht unbedingt
* also keep in mind, that I need roles, so maybe reuse the hub methods but at a "mustBeAdmin" argument or so. There should also be tests for correct protection of admin endpoints admins. -> maybe "/api/admins/..." and "/api/users/..." as endpoints? There should also be a loop like "make a user request to this list of endpoints and get unauthorized response due to lack of admin rights"

readBody -> da müsste wie gesagt noch das "mustBeAdmin" arg rein + "job validation" muss irgendwie ausgelagert werden da die datenstrukturen Anwendungsspezifisch sind.
Es wäre cool, wenn die ReadBody oder validate Funktionen eineige Grundtypen hätten, zB validate username/password, aber mit anderen Typen erweiterbar wären, zB

reuse "doRequest" in frontend cloud logic, also ensure for hub and cloud gui that only POST requests are used.

replace "stacks" by "apps"

Get rid of the bash scripts, always use go, adapt the readme
Extract the CLI runner stuff to separate module to share with other projects. Use 0BSD license.
Hub Production System must run with CockroachDB.
EnsureSchemaVersionTable() -> version the database scheme? research about database migration tools
wrap the http server of Go, in case the API changes

for later: frequent log deletion might be an interesting measure to keep the system clean and the volumes small. Centralized approach? e.g. delete files in the logs directory older than x days?


for later: frequent log deletion might be an interesting measure to keep the system clean and the volumes small. Centralized approach? e.g. delete files in the logs directory older than x days?

code quality improvements?
* "go vet ./..."
* golangci-lint:
  * go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  * create .golangci.yml (see example below)
  * golangci-lint run

```yaml
linters:
  enable:
    - deadcode   # Finds unused code
    - unused     # Checks for unused variables, constants, functions, etc.
run:
  issues-exit-code: 1  # Exit with code 1 if any issues are found
```

* staticcheck
  * go install honnef.co/go/tools/cmd/staticcheck@latest
  * staticcheck ./...
* Use -gcflags and -asmflags to pass options to the compiler: "go build -gcflags="-m""
* put code quality measures into the CI pipeline and git hooks (before push?)
* deadcode:
  * go install github.com/tsenart/deadcode@latest
  * deadcode ./...
* unused
  * go install honnef.co/go/tools/cmd/unused@latest
  * unused ./...
* enforce style: "gofmt -s -w ."
* dependency updater -> github has a bot for that I think

* don't userAndApp data structure etc to address an item. Rather use its ID.

* upload in hub must make multiple checks:
  * no extra privileges, maybe white allowed attributes
  * networks and volume names must exactly correlated with the maintainer and app name
  * the public container must have a name which is the same as the app name
  * Do I allow custom builds via "context: ."? Seems dangerous, ask ChatGPT.