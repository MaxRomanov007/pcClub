package models

// UserRole represents the user_roles table.
type UserRole struct {
	UserRoleID int64  `db:"user_role_id" json:"user_role_id"`
	Name       string `db:"name" json:"name"`
}

// User represents the users table.
type User struct {
	UserID              int64     `db:"user_id" json:"user_id"`
	UserRoleID          int64     `db:"user_role_id" json:"user_role_id"`
	Email               string    `db:"email" json:"email"`
	Password            []byte    `db:"password" json:"password"`
	RefreshTokenVersion int64     `db:"refresh_token_version" json:"refresh_token_version"`
	Balance             float64   `db:"balance" json:"balance"`
	UserRole            *UserRole `db:"-" json:"user_role"`
}

type UserData struct {
	UserID  int64   `db:"user_id" json:"user_id"`
	Email   string  `db:"email" json:"email"`
	Balance float64 `db:"balance" json:"balance"`
}
