package apps_new

import (
	"archive/zip"
	"fmt"
	"io"
	"ocelot/backend/repo"
	"ocelot/backend/tools"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO I am missing a test, where the app uses custom port and path, and client tries to access a specific message from there. Add a seconds nginx app for this.
// TODO Add logs to errors.

var hubClient HubClient

// TODO should directly give tagId
func DownloadTag(info tools.TagInfo) error {
	apps, err := hubClient.SearchApps(info.App)
	if err != nil {
		return err
	}

	app := (*apps)[0]
	tags, err := hubClient.GetTags(app.Id)
	if err != nil {
		return err
	}
	tag := (*tags)[0]

	tagContent, err := hubClient.DownloadTag(tag.Id)
	if err != nil {
		return err
	}
	err = repo.AppRepo.CreateApp(info.User, info.App)
	if err != nil {
		// TODO log: app was already existing, skip creation -> or maybe make a DoesAppExist check before that?
	}
	appId, err := repo.AppRepo.GetAppId(info.User, info.App)
	if err != nil {
		return err
	}
	err = repo.AppRepo.CreateTag(appId, info.Tag, *tagContent)
	if err != nil {
		return err
	}
	return nil
}

func StartContainer(appId int) error {
	app, err := repo.AppRepo.GetApp(appId)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "compose", "-p", app.Name, "up", "-d")
	err = extractTagToDirAndDeploy(app.ActiveTagId, cmd)
	if err != nil {
		return err
	}
	return nil
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, f.Mode())
			if err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(fpath), f.Mode())
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitializeAppsModule() {
	if tools.Config.UseRealHubClient {
		hubClient = NewHubClientReal()
	} else {
		hubClient = NewHubClientMock()
	}

	routes := []tools.Route{
		{"/apps/search", AppSearchHandler},
		{"/tags/list", GetTagsHandler},
		{"/tags/download", TagDownloadHandler},
		{"/apps/start", AppStartHandler},
		{"/apps/stop", AppStopHandler},
	}
	tools.RegisterRoutes(routes)
	tools.Router.HandleFunc("/api/apps/read", AppReadHandler)
}

func StopContainer(appId int) error {
	app, err := repo.AppRepo.GetApp(appId)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "compose", "-p", app.Name, "down")
	err = extractTagToDirAndDeploy(app.ActiveTagId, cmd)
	if err != nil {
		return err
	}
	return nil
}

func extractTagToDirAndDeploy(tagId int, command *exec.Cmd) error {
	tagContent, err := repo.AppRepo.LoadTagBlob(tagId)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "docker-compose")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	zipFilePath := filepath.Join(tempDir, "archive.zip")
	err = os.WriteFile(zipFilePath, tagContent, 0644)
	if err != nil {
		return err
	}

	err = unzip(zipFilePath, tempDir)
	if err != nil {
		return err
	}

	cmd := command
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
