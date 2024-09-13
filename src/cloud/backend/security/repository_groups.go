package security

func (r *MyRepository) CreateGroup(group string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) ListGroups() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MyRepository) DeleteGroup(group string) error {
	//TODO implement me
	panic("implement me")
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
