package models

import "time"

type TokenVo struct {
	UserKey    string    `xorm:"varchar(50) 'USER_KEY'"`
	Token      string    `xorm:"varchar(500) 'TOKEN'"`
	ExpireDate time.Time `xorm:"datetime 'EXPIRE_DATE'"`
}

func (_ TokenVo) TableName() string {
	return "TBL_TOKEN"
}
