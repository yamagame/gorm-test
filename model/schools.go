package model

import (
	"fmt"

	"gorm.io/gorm"
)

type School struct {
	gorm.Model
	ULID string `gorm:"index:,uniq,length:26"`
	Name string
}

func (x School) String() string {
	return fmt.Sprintf("{ULID: %s, Name: %s}", x.ULID, x.Name)
}
