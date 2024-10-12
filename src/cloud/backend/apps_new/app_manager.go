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

// TODO Add logs to errors.

var hubClient HubClient

func DownloadTag(info tools.TagInfo) error {
	tagContent, err := hubClient.DownloadTag(info)
	if err != nil {
		return err
	}
	err = repo.AppRepo.CreateAppWithTag(info.User, info.App, info.Tag, *tagContent)
	if err != nil {
		return err
	}
	return nil
}

func StartContainer(info tools.TagInfo) error {
	cmd := exec.Command("docker", "compose", "-p", info.App, "up", "-d")
	err := extractTagToDir(info, cmd)
	if err != nil {
		return err
	}
	return nil
}

func extractTagToDir(info tools.TagInfo, command *exec.Cmd) error {
	tagContent, err := repo.AppRepo.LoadTagBlob(info.User, info.App, info.Tag)
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

func StopContainer(info tools.TagInfo) error {
	cmd := exec.Command("docker", "compose", "-p", info.App, "down")
	err := extractTagToDir(info, cmd)
	if err != nil {
		return err
	}
	return nil
}
