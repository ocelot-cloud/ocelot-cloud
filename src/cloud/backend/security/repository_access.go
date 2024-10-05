package security

import (
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
	//TODO implement me
	panic("implement me")
}
