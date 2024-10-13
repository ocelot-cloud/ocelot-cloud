package repo

import (
	"fmt"
)

// TODO test: is app already existing
func (r *AppRepositoryImpl) CreateApp(maintainer, app string) error {
	_, err := DB.Exec("INSERT INTO apps (maintainer, app) VALUES (?, ?)", maintainer, app)
	if err != nil {
		Logger.Error("Failed to create app: %v", err)
		return fmt.Errorf("failed to create app")
	}
	return nil
}

func (r *AppRepositoryImpl) CreateTag(appId int, tag string, blob []byte) error {
	_, err := DB.Exec("INSERT INTO tags (app_id, tag, blob) VALUES (?, ?, ?)", appId, tag, blob)
	if err != nil {
		return fmt.Errorf("failed to create tag")
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

func (r *AppRepositoryImpl) ListApps() ([]App, error) {
	rows, err := DB.Query("SELECT maintainer, app, app_id FROM apps")
	if err != nil {
		Logger.Error("Failed to fetch app list: %v", err)
		return nil, fmt.Errorf("failed to fetch app list")
	}
	defer rows.Close()

	var result []App
	for rows.Next() {
		var maintainer, app string
		var appId int
		if err = rows.Scan(&maintainer, &app, &appId); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, App{Maintainer: maintainer, Name: app, AppId: appId})
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}
func (r *AppRepositoryImpl) ListTagsOfApp(appId int) ([]Tag, error) {
	rows, err := DB.Query("SELECT tag, tag_id FROM tags WHERE app_id = ?", appId)
	if err != nil {
		Logger.Error("Failed to fetch tag list: %v", err)
		return nil, fmt.Errorf("failed to fetch tag list")
	}
	defer rows.Close()

	var result []Tag
	for rows.Next() {
		var tagName string
		var tagId int
		if err = rows.Scan(&tagName, &tagId); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		tag := Tag{tagName, tagId}
		result = append(result, tag)
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

func (r *AppRepositoryImpl) LoadTagBlob(tagId int) ([]byte, error) {
	var blob []byte
	err := DB.QueryRow("SELECT blob FROM tags WHERE tag_id = ?", tagId).Scan(&blob)
	if err != nil {
		return nil, fmt.Errorf("TODO3")
	}
	return blob, nil
}

func (r *AppRepositoryImpl) DeleteApp(appId int) error {
	_, err := DB.Exec("DELETE FROM apps WHERE app_id = ?", appId)
	if err != nil {
		return fmt.Errorf("TODO6")
	}
	return nil
}

func (r *AppRepositoryImpl) DeleteTag(tagId int) error {
	_, err := DB.Exec("DELETE FROM tags WHERE tag_id = ?", tagId)
	if err != nil {
		return fmt.Errorf("TODO8")
	}
	return nil
}

func (r *AppRepositoryImpl) GetTagId(appId int, tag string) (int, error) {
	var tagId int
	err := DB.QueryRow("SELECT tag_id FROM tags WHERE app_id = ? AND tag = ?", appId, tag).Scan(&tagId)
	if err != nil {
		return -1, err
	}
	return tagId, nil
}

// TODO Add error logs
