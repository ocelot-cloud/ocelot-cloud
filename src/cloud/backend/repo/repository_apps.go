package repo

import (
	"fmt"
	"ocelot/backend/tools"
)

// TODO test: is app already existing
func (r *AppRepositoryImpl) CreateAppWithTag(maintainer string, app string, tag string, blob []byte) error {
	var doesAppExist bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM apps WHERE maintainer = ? AND app = ?)", maintainer, app).Scan(&doesAppExist)
	if err != nil {
		Logger.Error("Failed to check if app exists: %v", err)
		return fmt.Errorf("failed to check if app exists")
	}

	if !doesAppExist {
		_, err = DB.Exec("INSERT INTO apps (maintainer, app) VALUES (?, ?)", maintainer, app)
		if err != nil {
			Logger.Error("Failed to create app: %v", err)
			return fmt.Errorf("failed to create app")
		}
	}

	appId, err := r.GetAppId(maintainer, app)
	if err != nil {
		return fmt.Errorf("TODO4")
	}

	_, err = DB.Exec("INSERT INTO tags (app_id, tag, blob) VALUES (?, ?, ?)", appId, tag, blob)
	if err != nil {
		return fmt.Errorf("TODO5")
	}

	return nil
}

func (r *AppRepositoryImpl) GetAppId(maintainer string, app string) (int, error) {
	var appId int
	err := DB.QueryRow("SELECT app_id FROM apps WHERE maintainer = ? AND app = ?", maintainer, app).Scan(&appId)
	if err != nil {
		// TODO Make better error message. Should not be error log, since DoesTagExist uses it and is expected to not find the app any time.
		Logger.Info("TODO error: %v", err)
		return -1, err
	}
	return appId, nil
}

func (r *AppRepositoryImpl) ListApps() ([]MaintainerAndApp, error) {
	rows, err := DB.Query("SELECT maintainer, app FROM apps")
	if err != nil {
		Logger.Error("Failed to fetch app list: %v", err)
		return nil, fmt.Errorf("failed to fetch app list")
	}
	defer rows.Close()

	var result []MaintainerAndApp
	for rows.Next() {
		var maintainer, app string
		if err = rows.Scan(&maintainer, &app); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, MaintainerAndApp{Maintainer: maintainer, App: app})
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

func (r *AppRepositoryImpl) ListTagsOfApp(maintainer string, app string) ([]string, error) {
	appId, err := r.GetAppId(maintainer, app)
	if err != nil {
		return nil, fmt.Errorf("TODO1")
	}

	rows, err := DB.Query("SELECT tag FROM tags WHERE app_id = ?", appId)
	if err != nil {
		Logger.Error("Failed to fetch tag list of app: %s/%s, %v", maintainer, app, err)
		return nil, fmt.Errorf("failed to fetch tag list")
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var singleTag string
		if err = rows.Scan(&singleTag); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, singleTag)
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

func (r *AppRepositoryImpl) LoadTagBlob(maintainer, app, tag string) ([]byte, error) {
	appId, err := r.GetAppId(maintainer, app)
	if err != nil {
		return nil, fmt.Errorf("TODO2")
	}

	var blob []byte
	err = DB.QueryRow("SELECT blob FROM tags WHERE app_id = ? AND tag = ?", appId, tag).Scan(&blob)
	if err != nil {
		return nil, fmt.Errorf("TODO3")
	}

	return blob, nil
}

func (r *AppRepositoryImpl) DeleteApp(maintainer, app string) error {
	_, err := DB.Exec("DELETE FROM apps WHERE maintainer = ? AND app = ?", maintainer, app)
	if err != nil {
		return fmt.Errorf("TODO6")
	}
	return nil
}

func (r *AppRepositoryImpl) DeleteTag(maintainer, app, tag string) error {
	appId, err := r.GetAppId(maintainer, app)
	if err != nil {
		return fmt.Errorf("TODO7")
	}

	_, err = DB.Exec("DELETE FROM tags WHERE app_id = ? AND tag = ?", appId, tag)
	if err != nil {
		return fmt.Errorf("TODO8")
	}
	return nil
}

// TODO Method is not tested yet.
// TODO Add error logs
func (r *AppRepositoryImpl) DoesTagExist(tagInfo tools.TagInfo) bool {
	appId, err := r.GetAppId(tagInfo.User, tagInfo.App)
	if err != nil {
		return false
	}

	var doesTagExist bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tags WHERE app_id = ? AND tag = ?)", appId, tagInfo.Tag).Scan(&doesTagExist)
	if err != nil {
		return false
	}

	return doesTagExist
}
