package models

type UserAdditionalVo struct {
	UserKey   string `xorm:"varchar(50) 'USER_KEY'"`
	UserType  string `xorm:"varchar(20) 'USER_TYPE'"`
	TypeValue string `xorm:"varchar(255) 'TYPE_VALUE'"`
}

func (_ UserAdditionalVo) TableName() string {
	return "TBL_ADDITIONAL_USER"
}
