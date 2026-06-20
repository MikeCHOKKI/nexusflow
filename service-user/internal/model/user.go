package model

// User représente un utilisateur tel que stocké en base.
type User struct {
	ID        string
	Email     string
	Password  string // hash bcrypt
	Name      string
	Role      string
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}
