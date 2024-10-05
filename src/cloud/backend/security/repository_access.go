package security

import (
	"database/sql"
	"fmt"
	"strings"
)

// TODO Replace maintainer + app (separate strings) with MaintainerAndApp in all repos. Probably wil be extended later.

func (r *AccessRepositoryImpl) GiveGroupAccessToApp(group string, app MaintainerAndApp) error {
	groupId, err := GroupRepo.getGroupId(group)
	if err != nil {
		// TODO
		return err
	}

	// TODO Should have MaintainerAndApp as argument. Maybe rename to AppInfo?
	appId, err := AppRepo.getAppId(app.Maintainer, app.App)
	if err != nil {
		// TODO
		return err
	}

	_, err = DB.Exec("INSERT INTO app_access(group_id, app_id) VALUES (?, ?)", groupId, appId)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *AccessRepositoryImpl) ListAppAccessesOfGroup(group string) ([]MaintainerAndApp, error) {
	groupId, err := GroupRepo.getGroupId(group)
	if err != nil {
		// TODO
		return nil, err
	}

	rows, err := DB.Query("SELECT app_id FROM app_access WHERE group_id = ?", groupId)
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()
	var appsIds []int
	for rows.Next() {
		var appId int
		if err = rows.Scan(&appId); err != nil {
			// TODO
			return nil, err
		}
		appsIds = append(appsIds, appId)
	}

	apps, err := r.getAppsByIDs(appsIds)
	if err != nil {
		// TODO
		return nil, err
	}

	return apps, nil
}

func (r *AccessRepositoryImpl) getAppsByIDs(ids []int) ([]MaintainerAndApp, error) {
	query := fmt.Sprintf("SELECT maintainer, app FROM apps WHERE app_id IN (%s)", strings.TrimSuffix(strings.Repeat("?,", len(ids)), ","))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()

	var apps []MaintainerAndApp
	for rows.Next() {
		var maintainer string
		var app string
		if err = rows.Scan(&maintainer, &app); err != nil {
			// TODO
			return nil, err
		}
		apps = append(apps, MaintainerAndApp{maintainer, app})
	}

	return apps, nil
}

func (r *AccessRepositoryImpl) RemoveGroupsAccessToApp(group string, app MaintainerAndApp) error {
	groupId, err := GroupRepo.getGroupId(group)
	if err != nil {
		// TODO
		return err
	}

	appId, err := AppRepo.getAppId(app.Maintainer, app.App)
	if err != nil {
		// TODO
		return err
	}

	_, err = DB.Exec("DELETE FROM app_access WHERE group_id = ? AND app_id = ?", groupId, appId)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccessRepositoryImpl) DoesUserHaveAccessToApp(user string, app MaintainerAndApp) bool {
	userId, err := GroupRepo.getUserId(user) // TODO should be in userRepo
	if err != nil {
		Logger.Info("Error getting user id: %v", err)
		return false
	}

	rows, err := DB.Query("SELECT group_id FROM user_to_group WHERE user_id = ?", userId)
	if err != nil {
		Logger.Info("Error getting group ids: %v", err)
		return false
	}
	defer rows.Close()

	var groupIds []int
	for rows.Next() {
		var groupId int
		if err = rows.Scan(&groupId); err != nil {
			return false
		}
		groupIds = append(groupIds, groupId)
	}

	if len(groupIds) == 0 {
		return false
	}

	appId, err := AppRepo.getAppId(app.Maintainer, app.App)
	if err != nil {
		return false
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(groupIds)), ",")
	query := fmt.Sprintf("SELECT 1 FROM app_access WHERE app_id = ? AND group_id IN (%s) LIMIT 1", placeholders)

	args := make([]interface{}, len(groupIds)+1)
	args[0] = appId
	for i, groupId := range groupIds {
		args[i+1] = groupId
	}

	row := DB.QueryRow(query, args...)
	var exists int
	err = row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}
