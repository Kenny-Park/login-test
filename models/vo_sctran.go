package models

import "time"

type ScTranVo struct {
	TranKey      string    `xorm:"varchar(50) 'TRAN_KEY'"`
	Hp           string    `xorm:"varchar(20) 'HP'"`
	Num          string    `xorm:"varchar(6) 'NUM'"`
	ExpireDate   time.Time `xorm:"datetime 'EXPIRE_DATE'"`
	CertCompDate time.Time `xorm:"datetime 'CERT_COMP_DATE'"`
}

func (_ ScTranVo) TableName() string {
	return "TBL_SCTRAN"
}
