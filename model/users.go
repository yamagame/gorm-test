package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ULID string `gorm:"index:,uniq,length:26"`
	Name string
}

func (x User) String() string {
	return fmt.Sprintf("{ULID: %s, Name: %s}", x.ULID, x.Name)
}
