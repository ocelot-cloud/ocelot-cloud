package apps

import (
	"encoding/json"
	"io"
	"net/http"
	"ocelot/backend/tools"
)

func createReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
		return
	}

	stackStateInfo := stackService.getAppStateInfo()
	response := make([]tools.ResponsePayloadDto, 0)
	for stackName, stackDetails := range stackStateInfo {
		response = append(response, tools.ResponsePayloadDto{stackName, stackDetails.State.toString(), stackDetails.Path})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createDeployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	stackName, err := decodeStackInfo(r)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err = stackService.deployApp(stackName); err != nil {
		if err != nil { // TODO condition is always true
			logger.Error("Deploying stack failed: " + stackName + "\n" + err.Error() + "\n")
			http.Error(w, "Deploying stack failed: "+stackName, http.StatusInternalServerError)
		}
		return
	}
}

func createStopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	stackName, err := decodeStackInfo(r)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err = stackService.stopApp(stackName); err != nil {
		if err != nil {
			logger.Warn("error when trying to stop stack, %s", err.Error())
			http.Error(w, "Stopping stack failed: "+stackName, http.StatusInternalServerError)
		}
		return
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
