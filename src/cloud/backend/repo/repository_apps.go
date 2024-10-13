package repo

import (
	"fmt"
	"ocelot/backend/tools"
)

// TODO test: is app already existing
func (r *AppRepositoryImpl) CreateApp(maintainer, app string) error {
	_, err := DB.Exec("INSERT INTO apps (maintainer, app, active_tag) VALUES (?, ?, ?)", maintainer, app, -1)
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

	app, err := r.GetApp(appId)
	if err != nil {
		return err
	}
	if app.ActiveTagId == -1 {
		tagId, err := r.GetTagId(appId, tag)
		if err != nil {
			return err
		}
		_, err = DB.Exec("UPDATE apps SET active_tag = ? WHERE app_id = ?", tagId, appId)
		if err != nil {
			return fmt.Errorf("failed to update active tag")
		}
	}

	return nil
}

// TODO Should actually be hidden to outside. I think it would be better to expose interfaces to the outside, while using the implementations internally.
func (r *AppRepositoryImpl) GetApp(appId int) (*tools.App, error) {
	var maintainer, appName string
	var activeTagId int
	err := DB.QueryRow("SELECT maintainer, app, active_tag FROM apps WHERE app_id = ?", appId).Scan(&maintainer, &appName, &activeTagId)
	if err != nil {
		return nil, fmt.Errorf("TODO2")
	}

	app := &tools.App{maintainer, appName, appId, "", -1, false}
	activeTag, err := r.getTag(activeTagId)
	if err != nil {
		// TODO if tag not found, then it becomes an empty string
		return app, nil
	}
	// TODO Add null check?
	app.ActiveTagName = activeTag.Name
	app.ActiveTagId = activeTagId
	return app, nil
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

func (r *AppRepositoryImpl) ListApps() ([]tools.App, error) {
	rows, err := DB.Query("SELECT maintainer, app, app_id, active_tag FROM apps")
	if err != nil {
		Logger.Error("Failed to fetch app list: %v", err)
		return nil, fmt.Errorf("failed to fetch app list")
	}
	defer rows.Close()

	var result []tools.App
	for rows.Next() {
		var maintainer, app string
		var appId, activeTagId int
		if err = rows.Scan(&maintainer, &app, &appId, &activeTagId); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, tools.App{maintainer, app, appId, "", activeTagId, false})
	}

	// TODO A little performance hungry. Maybe make one query for all apps and then put it together in memory.
	for i := range result {
		activeTag, err := r.getTag(result[i].ActiveTagId)
		if err != nil {
			activeTag = &Tag{"", result[i].ActiveTagId, -1}
		}
		result[i].ActiveTagName = activeTag.Name
		result[i].ActiveTagId = activeTag.TagId
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

// TODO If the ID is not found, the query is hanging indefinitely, which is bad. I want to return an error in this case.
// TODO make Tag a pointer and return nil in case of error?
func (r *AppRepositoryImpl) getTag(tagId int) (*Tag, error) {
	var tagName string
	var appId int
	err := DB.QueryRow("SELECT tag, app_id FROM tags WHERE tag_id = ?", tagId).Scan(&tagName, &appId)
	if err != nil {
		return nil, fmt.Errorf("TODO4")
	}
	return &Tag{tagName, tagId, appId}, nil
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
		tag := Tag{tagName, tagId, appId}
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
	tag, err := r.getTag(tagId)
	if err != nil {
		return err
	}

	app, err := r.GetApp(tag.AppId)
	if err != nil {
		return err
	}

	_, err = DB.Exec("DELETE FROM tags WHERE tag_id = ?", tagId)
	if err != nil {
		return fmt.Errorf("TODO8")
	}

	if app.ActiveTagId == tagId {
		_, err = DB.Exec("UPDATE apps SET active_tag = ? WHERE app_id = ?", -1, tag.AppId)
		if err != nil {
			return fmt.Errorf("TODO9")
		}
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

func (r *AppRepositoryImpl) SetActiveTag(appId, tagId int) error {
	tag, err := r.getTag(tagId)
	if err != nil {
		return err
	}
	if tag.AppId != appId {
		// TODO Logger
		return fmt.Errorf("tag is not of the app")
	}

	_, err = DB.Exec("UPDATE apps SET active_tag = ? WHERE app_id = ?", tagId, appId)
	if err != nil {
		return fmt.Errorf("TODO10")
	}
	return nil
}

// TODO Add error logs
