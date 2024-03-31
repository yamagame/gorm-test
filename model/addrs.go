package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Addr struct {
	gorm.Model
	ULID string `gorm:"index:,uniq,length:26"`
	City string
}

func (x Addr) String() string {
	return fmt.Sprintf("{ULID: %s, Name: %s}", x.ULID, x.City)
}
