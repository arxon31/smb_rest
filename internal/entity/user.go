package entity

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return ErrInvalidUser
	}
	if u.Username != "admin" || u.Password != "admin" {
		return ErrInvalidUser
	}

	return nil
}
