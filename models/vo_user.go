package models

import "time"

type UserVo struct {
	UserKey  string    `xorm:"varchar(50) 'USER_KEY'"`
	UserName string    `xorm:"varchar(255) 'USER_NAME'"`
	Pwd      string    `xorm:"varchar(256) 'PWD'"`
	RegDate  time.Time `xorm:"datetime 'REG_DATE'"`
}

func (_ UserVo) TableName() string {
	return "TBL_USER"
}
