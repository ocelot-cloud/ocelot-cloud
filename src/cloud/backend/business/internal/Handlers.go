package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ocelot/backend/config"
)

func checkSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || cookie.Value != "valid" {
		Logger.Trace("Cookie error.")
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Trace("Cookie was okay.")
		w.WriteHeader(http.StatusOK)
	}
}

func (a *ApplicationInitializer) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<html><body>Hello</body></html>")
}

func createReadHandler(stackService StackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackStateInfo := stackService.GetStackStateInfo()
		response := make([]tools.ResponsePayloadDto, 0)
		for stackName, stackDetails := range stackStateInfo {
			response = append(response, tools.ResponsePayloadDto{stackName, stackDetails.State.String(), stackDetails.Path})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func createDeployHandler(stackService StackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackName, err := decodeStackInfo(r)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if err := stackService.DeployStack(stackName); err != nil {
			if err != nil {
				Logger.Error("Deploying stack failed: " + stackName + "\n" + err.Error() + "\n")
				http.Error(w, "Deploying stack failed: "+stackName, http.StatusInternalServerError)
			}
			return
		}
	}
}

func createStopHandler(stackService StackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}

		stackName, err := decodeStackInfo(r)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if err := stackService.StopStack(stackName); err != nil {
			if err != nil {
				Logger.Warn("error when trying to stop stack, %s", err.Error())
				http.Error(w, "Stopping stack failed: "+stackName, http.StatusInternalServerError)
			}
			return
		}
	}
}

func decodeStackInfo(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var stackInfo tools.StackInfo
	err = json.Unmarshal(body, &stackInfo)
	if err != nil {
		return "", err
	}
	return stackInfo.Name, nil
}
