package apps_new

import (
	"ocelot/backend/repo"
)

// TODO Add a test/function which takes the zip bytes, extracts them locally and run "docker-compose up" on them.
// TODO in shared.utils remove the "storage" stuff: Logger.Info("zipped directory %s and stored its %v bytes in the database for integration testing", dirPath, len(buf.Bytes()))
// TODO Add logs to errors.

var hubClient HubClient // TODO should be set by the application and the tests depending on the current profile

func DownloadTag(info TagInfo) error {
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
