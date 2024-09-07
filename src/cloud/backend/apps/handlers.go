package apps

import (
	"encoding/json"
	"github.com/ocelot-cloud/shared/utils"
	"io"
	"net/http"
	"ocelot/backend/tools"
)

func appReadHandler(w http.ResponseWriter, r *http.Request) {
	stackStateInfo := stackService.getAppStateInfo()
	response := make([]tools.AppInfo, 0)
	for stackName, stackDetails := range stackStateInfo {
		response = append(response, tools.AppInfo{stackName, stackDetails.State.toString(), stackDetails.Path})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func appDeployHandler(w http.ResponseWriter, r *http.Request) {
	stackName, err := decodeSingleString(r)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err = stackService.deployApp(stackName); err != nil {
		logger.Error("Deploying stack failed: " + stackName + "\n" + err.Error() + "\n")
		http.Error(w, "Deploying stack failed: "+stackName, http.StatusInternalServerError)
		return
	}
}

func appStopHandler(w http.ResponseWriter, r *http.Request) {
	stackName, err := decodeSingleString(r)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err = stackService.stopApp(stackName); err != nil {
		logger.Warn("error when trying to stop stack, %s", err.Error())
		http.Error(w, "Stopping stack failed: "+stackName, http.StatusInternalServerError)
		return
	}
}

func decodeSingleString(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var singleString utils.SingleString
	err = json.Unmarshal(body, &singleString)
	if err != nil {
		return "", err
	}
	return singleString.Value, nil
}
