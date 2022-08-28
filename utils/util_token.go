package utils

import (
	"math/rand"
	"modules/models"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type UtilToken struct{}

// 랜덤스트링 생성
func (u UtilToken) RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	rstring := func(n int) string {
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}(n)
	return time.Now().Format("20060102150405") + rstring
}

// 랜덤 정수 생성
func (u UtilToken) RandInt(n int) string {
	rand.Seed(time.Now().UnixNano())
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = strconv.Itoa(rand.Intn(9))
	}
	return strings.Join(arr, "")
}

// 토큰 생성
func (u UtilToken) Create(user *models.UserVo, userAdditionalVo []*models.UserAdditionalVo, expireTime time.Time) string {

	m := map[string]string{}
	username, _ := Crypto{}.Decrypt(user.UserName)

	m["UserKey"] = user.UserKey
	m["UserName"] = username
	for _, item := range userAdditionalVo {
		v, _ := Crypto{}.Decrypt(item.TypeValue)
		m[item.UserType] = v
	}
	m["ExpireDate"] = expireTime.Format("20060102")

	var param []string
	for k, v := range m {
		param = append(param, k+"="+v)
	}

	enc, err := Crypto{}.Encrypt(strings.Join(param, ";"))
	if err != nil {
		return ""
	}
	return enc
}

// 토큰 파싱
func (u UtilToken) Parse(enc string) *TokenInfoVo {

	tokenInfoVo := &TokenInfoVo{}
	enc, err := Crypto{}.Decrypt(enc)
	if err != nil {
		return nil
	}

	arr := strings.Split(enc, ";")
	target := reflect.ValueOf(tokenInfoVo)
	elem := target.Elem()

	for _, item := range arr {
		kv := strings.Split(item, "=")
		k := kv[0]
		v := kv[1]
		elem.FieldByName(k).SetString(v)
	}
	return tokenInfoVo
}
