package utils

type allowedUsers map[string]string

func AllowedUsers() allowedUsers {
	return allowedUsers{
		"test_user": "test_password",
	}
}
