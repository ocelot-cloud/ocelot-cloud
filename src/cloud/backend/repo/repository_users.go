package repo

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"golang.org/x/crypto/bcrypt"
	"ocelot/backend/tools"
	"time"
)

func (r *UserRepositoryImpl) DoesAnyAdminUserExist() bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE is_admin = ?)", true).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check if there is any admin user: %v", err)
		return false
	}
	return exists
}

func (r *UserRepositoryImpl) CreateUser(user string, password string, isAdmin bool) error {
	hashedPassword, err := utils.SaltAndHash(password)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO users (user_name, hashed_password, is_admin) VALUES (?, ?, ?)", user, hashedPassword, isAdmin)
	if err != nil {
		Logger.Warn("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user")
	}
	return nil
}

// TODO shift to shared module

func (r *DatabaseRepositoryImpl) WipeDatabase() {
	_, err := DB.Exec("DELETE FROM users")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}

	_, err = DB.Exec("DELETE FROM apps")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}

	_, err = DB.Exec("DELETE FROM tags")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}

	_, err = DB.Exec("DELETE FROM groups")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}

	_, err = DB.Exec("DELETE FROM app_access")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}

	_, err = DB.Exec("DELETE FROM user_to_group")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}
}

func (r *UserRepositoryImpl) IsPasswordCorrect(user string, password string) bool {
	var hashedPassword string
	err := DB.QueryRow("SELECT hashed_password FROM users WHERE user_name = ?", user).Scan(&hashedPassword)
	if err != nil {
		Logger.Error("Failed to fetch hashed password: %v", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (r *UserRepositoryImpl) DeleteUser(user string) error {
	_, err := DB.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		Logger.Warn("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user")
	}
	return nil
}

func (r *UserRepositoryImpl) HashAndSaveCookie(user string, cookieValue string, cookieExpirationDate time.Time) error {
	hashedCookieValue, err := utils.Hash(cookieValue)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET hashed_cookie_value = ?, cookie_expiration_date = ? WHERE user_name = ?", hashedCookieValue, cookieExpirationDate.Format(time.RFC3339), user)
	if err != nil {
		Logger.Warn("Failed to update cookie of user '%s': %v", user, err)
		return fmt.Errorf("failed to update cookie")
	}
	return nil
}

func (r *UserRepositoryImpl) Logout(user string) error {
	_, err := DB.Exec("UPDATE users SET hashed_cookie_value = ?, cookie_expiration_date = ? WHERE user_name = ?", "", "", user)
	if err != nil {
		Logger.Error("Failed to delete cookie of user '%s': %v", user, err)
		return fmt.Errorf("failed to delete cookie")
	}
	return nil
}

func (r *UserRepositoryImpl) DoesUserExist(user string) bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_name = ?)", user).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check if user exists: %v", err)
		return false
	}
	return exists
}

// TODO Test if isAdmin is correct in authorization.
func (r *UserRepositoryImpl) GetUserViaCookie(cookieValue string) (*tools.Authorization, error) {
	hashedCookieValue, err := utils.Hash(cookieValue)
	if err != nil {
		return nil, err
	}

	var user string
	var isAdmin bool
	err = DB.QueryRow("SELECT user_name, is_admin FROM users WHERE hashed_cookie_value = ?", hashedCookieValue).Scan(&user, &isAdmin)
	if err != nil {
		Logger.Error("Failed to fetch user data: %v", err)
		return nil, fmt.Errorf("failed to fetch user data")
	}
	return &tools.Authorization{user, isAdmin}, nil
}

func (r *UserRepositoryImpl) ChangePassword(user string, newPassword string) error {
	hashedNewPassword, err := utils.SaltAndHash(newPassword)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET hashed_password = ? WHERE user_name = ?", hashedNewPassword, user)
	if err != nil {
		Logger.Error("Failed to update password of user '%s': %v", user, err)
		return fmt.Errorf("failed to update password")
	}
	return nil
}

func (r *UserRepositoryImpl) GenerateSecret(user string) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		Logger.Error("failed to generate secret: %v", err)
		return "", fmt.Errorf("failed to generate secret")
	}
	secret := hex.EncodeToString(randomBytes)
	_, err := DB.Exec("UPDATE users SET secret = ? WHERE user_name = ?", secret, user)
	if err != nil {
		Logger.Error("failed to to secret: %v", err)
		return "", fmt.Errorf("failed to to secret")
	}

	return secret, nil
}

func (r *UserRepositoryImpl) IsSecretCorrect(user, secret string) bool {
	var repoSecret string
	err := DB.QueryRow("SELECT secret FROM users WHERE user_name = ?", user).Scan(&repoSecret)
	if err != nil {
		Logger.Error("failed to fetch secret: %v", err)
		return false
	}
	return repoSecret == secret
}

func (r *UserRepositoryImpl) RemoveSecret(user string) error {
	_, err := DB.Exec("UPDATE users SET secret = ? WHERE user_name = ?", "", user)
	if err != nil {
		Logger.Error("failed to remove secret: %v", err)
		return fmt.Errorf("failed to remove secret")
	}
	return nil
}
