package security

func (r *MyRepository) CreateGroup(group string) error {
	_, err := DB.Exec("INSERT INTO groups(group_name) VALUES (?)", group)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *MyRepository) ListGroups() ([]string, error) {
	rows, err := DB.Query("SELECT group_name FROM groups")
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()
	var groups []string
	for rows.Next() {
		var group string
		if err = rows.Scan(&group); err != nil {
			// TODO
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *MyRepository) DeleteGroup(group string) error {
	_, err := DB.Exec("DELETE FROM groups WHERE group_name = ?", group)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *MyRepository) ListAllUsers() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) AddUserToGroup(user, group string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) ListMembersOfGroup(group string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) RemoveUserFromGroup(user, group string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) GiveGroupAccessToApp(group, app string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) ListAppAccessesOfGroup(group string) ([]MaintainerAndApp, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) DoesUserHaveAccessToApp(user, maintainer, app string) bool {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) RemoveGroupsAccessToApp(group, app string) error {
	//TODO implement me
	panic("implement me")
}
