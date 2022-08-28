package utils_test

import (
	"log"
	"modules/utils"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

/*func TestEncrypt(t *testing.T) {
	s, _ := utils.Crypto{}.Encrypt("박경훈")

	if s != "yWY+LssL3LXR2qlqTT/m8g==" {
		t.Error("Wrong result")
	}
}


func TestDecrypt(t *testing.T) {
	s, _ := utils.Crypto{}.Decrypt("yWY+LssL3LXR2qlqTT/m8g==")

	if s != "박경훈" {
		t.Error("Wrong result")
	}
}

func TestCryptoSha256Encrypt(t *testing.T) {
	s := utils.CryptoSha256{}.Encrypt("dkagh")

	if s != "bafe77e32a4dce5bfd4005973286acf7c5b73d8574a6621536f408ce8c041246" {
		t.Error("Wrong result")
	}
}*/

func TestGetAdditionalUser(t *testing.T) {
	// DB 접속
	engine, err := xorm.NewEngine("mysql", "root:11111111@tcp(127.0.0.1:3306)/dev")
	if err != nil {
		log.Println(err)
	}
	defer engine.Close()

	x := &utils.XormHelper{
		Engine: engine,
	}
	x.Config()

	ss, _ := utils.Crypto{}.Encrypt("00000000")
	log.Println(ss)

	s, _ := x.GetAdditionalUser("Hp", ss)
	if s == nil {
		t.Error("Wrong result")
	}
}
