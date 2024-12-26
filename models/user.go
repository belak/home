package models

import "fmt"

type UserInfo struct {
	ID       int
	Username string
}

func (u *UserInfo) String() string {
	return fmt.Sprintf("UserInfo<ID=%d,Username=%s>", u.ID, u.Username)
}
