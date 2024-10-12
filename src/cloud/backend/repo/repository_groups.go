package repo

import (
	"fmt"
	"strings"
)

func (r *GroupRepositoryImpl) CreateGroup(group string) error {
	_, err := DB.Exec("INSERT INTO groups(group_name) VALUES (?)", group)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *GroupRepositoryImpl) ListGroups() ([]Group, error) {
	rows, err := DB.Query("SELECT group_id, group_name FROM groups")
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()
	var groups []Group
	for rows.Next() {
		var groupId int
		var groupName string
		if err = rows.Scan(&groupId, &groupName); err != nil {
			// TODO
			return nil, err
		}
		groups = append(groups, Group{groupId, groupName})
	}
	return groups, nil
}

func (r *GroupRepositoryImpl) DeleteGroup(group string) error {
	_, err := DB.Exec("DELETE FROM groups WHERE group_name = ?", group)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *GroupRepositoryImpl) ListAllUsers() ([]string, error) {
	rows, err := DB.Query("SELECT user_name FROM users")
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var user string
		if err = rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// TODO deletion cascading: when either the group or the user is deleted, this entry should be deleted as well.
func (r *GroupRepositoryImpl) AddUserToGroup(user, group string) error {
	userId, err := r.getUserId(user)
	if err != nil {
		// TODO
		return err
	}
	groupId, err := r.GetGroupId(group)
	if err != nil {
		// TODO
		return err
	}

	_, err = DB.Exec("INSERT INTO user_to_group(user_id, group_id) VALUES (?, ?)", userId, groupId)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (r *GroupRepositoryImpl) getUserId(user string) (int, error) {
	var userId int
	err := DB.QueryRow("Select user_id from users where user_name = ?", user).Scan(&userId)
	if err != nil {
		// TODO
		return -1, err
	}
	return userId, nil
}

func (r *GroupRepositoryImpl) GetGroupId(group string) (int, error) {
	var groupId int
	err := DB.QueryRow("Select group_id from groups where group_name = ?", group).Scan(&groupId)
	if err != nil {
		// TODO
		return -1, err
	}
	return groupId, nil
}

func (r *GroupRepositoryImpl) ListMembersOfGroup(group string) ([]string, error) {
	groupId, err := r.GetGroupId(group)
	if err != nil {
		// TODO
		return nil, err
	}
	rows, err := DB.Query("SELECT user_id FROM user_to_group WHERE group_id = ?", groupId)
	if err != nil {
		// TODO
		return nil, err
	}
	defer rows.Close()
	var userIds []int
	for rows.Next() {
		var userId int
		if err = rows.Scan(&userId); err != nil {
			// TODO
			return nil, err
		}
		userIds = append(userIds, userId)
	}

	usernames, err := r.getUsernamesByIDs(userIds)
	if err != nil {
		// TODO
		return nil, err
	}

	return usernames, nil
}

func (r *GroupRepositoryImpl) getUsernamesByIDs(ids []int) ([]string, error) {
	query := fmt.Sprintf("SELECT user_name FROM users WHERE user_id IN (%s)", strings.TrimSuffix(strings.Repeat("?,", len(ids)), ","))

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

	var usernames []string
	for rows.Next() {
		var username string
		if err = rows.Scan(&username); err != nil {
			// TODO
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func (r *GroupRepositoryImpl) RemoveUserFromGroup(user, group string) error {
	// TODO Maybe abstract to getUserAndGroupId(user, group), and reuse it in other function where the same block is used.
	userId, err := r.getUserId(user)
	if err != nil {
		// TODO
		return err
	}
	groupId, err := r.GetGroupId(group)
	if err != nil {
		// TODO
		return err
	}

	_, err = DB.Exec("DELETE FROM user_to_group WHERE user_id = ? AND group_id = ?", userId, groupId)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

// TODO Feature for later: clicking on an app should display, which group have access to it. But low-prio.
