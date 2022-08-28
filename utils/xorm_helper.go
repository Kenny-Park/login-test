package utils

import (
	"errors"
	"log"
	"modules/models"
	"time"

	"xorm.io/xorm"
)

type XormHelper struct {
	Engine *xorm.Engine
}

func (x *XormHelper) Config() {
	x.Engine.DB().SetMaxIdleConns(16)
	x.Engine.DB().SetMaxOpenConns(8)

	tz, _ := time.LoadLocation("Asia/Seoul")

	// 타임존 설정
	x.Engine.TZLocation = tz
	x.Engine.DatabaseTZ = tz
}

// 토큰정보 저장
func (x *XormHelper) InsertToken(vo *models.TokenVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Insert(vo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 토큰정보 가져오기
func (x *XormHelper) GetToken(token string) (*models.TokenVo, error) {
	tokenVo := new(models.TokenVo)
	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("TOKEN", "EXPIRE_DATE").Where("TOKEN = ?", token).Get(tokenVo)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}

	return tokenVo, nil
}

func (x *XormHelper) GetTokenByUserKey(userKey string) (*models.TokenVo, error) {
	tokenVo := new(models.TokenVo)
	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("TOKEN", "EXPIRE_DATE").Where("USER_KEY = ?", userKey).Get(tokenVo)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}

	return tokenVo, nil
}

// 토큰정보 가져오기
func (x *XormHelper) DeleteToken(token string) error {
	s := x.Engine.NewSession()
	defer s.Close()

	rows, err := s.Delete(&models.TokenVo{Token: token})
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.New("no data")
	}
	return nil
}

// 고객정보 저장
func (x *XormHelper) InsertUser(vo *models.UserVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Insert(vo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 비밀번호 변경
func (x *XormHelper) UpdatePwd(vo *models.UserVo) error {
	s := x.Engine.NewSession()
	defer s.Close()

	rows, err := s.Cols("pwd").Where("user_key = ?", vo.UserKey).Update(vo)
	if err != nil {
		log.Println(err)
	}
	if rows <= 0 {
		return errors.New("no data")
	}
	return nil
}

// 고객정보 수정
func (x *XormHelper) UpdateAdditionalUser(vos []*models.UserAdditionalVo) error {
	s := x.Engine.NewSession()
	defer s.Close()

	var err error
	s.Begin()

	var has bool
	for _, item := range vos {
		has, err = s.Where("USER_KEY != ? AND USER_TYPE = ? AND TYPE_VALUE = ?", item.UserKey, item.UserType, item.TypeValue).Exist(item)
		if has {
			s.Rollback()
			return errors.New("exists data")
		}
		if err != nil {
			s.Rollback()
			return err
		}
	}
	for _, item := range vos {
		_, err := s.Cols("TYPE_VALUE").Where("USER_KEY = ? AND USER_TYPE = ?", item.UserKey, item.UserType).Update(item)
		if err != nil {
			s.Rollback()
			return err
		}
	}
	s.Commit()
	return nil
}

// 고객정보 저장
// 저장하는 순간에도 중복된값이 있다면 에러출력
func (x *XormHelper) InsertAdditionalUser(vos []*models.UserAdditionalVo) error {
	s := x.Engine.NewSession()
	defer s.Close()

	var err error
	s.Begin()

	var has bool
	for _, item := range vos {
		has, err = s.Where("USER_TYPE = ? AND TYPE_VALUE = ?", item.UserType, item.TypeValue).Exist(item)
		if has {
			s.Rollback()
			return errors.New("exists data")
		}
		if err != nil {
			s.Rollback()
			return err
		}
	}
	for _, item := range vos {
		_, err = s.Insert(item)
		if err != nil {
			s.Rollback()
			return err
		}
	}

	s.Commit()
	return nil
}

// 고객부가정보 조회
func (x *XormHelper) GetAdditionalUser(userType string, typeValue string) (*models.UserAdditionalVo, error) {
	userAdditionalVo := new(models.UserAdditionalVo)

	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("USER_KEY", "USER_TYPE", "TYPE_VALUE").
		Where("USER_TYPE = ? AND TYPE_VALUE = ? ", userType, typeValue).
		Get(userAdditionalVo)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}

	return userAdditionalVo, nil
}

func (x *XormHelper) ListAdditionalUserByValue(typeValue string) ([]*models.UserAdditionalVo, error) {
	var userAdditionalVos []*models.UserAdditionalVo

	s := x.Engine.NewSession()
	defer s.Close()

	err := s.Cols("USER_KEY").
		Where("TYPE_VALUE = ?", typeValue).
		Find(&userAdditionalVos)

	if err != nil {
		return nil, err
	}

	return userAdditionalVos, nil
}

// 고객부가정보 조회
func (x *XormHelper) ListAdditionalUser(userKey string) ([]*models.UserAdditionalVo, error) {
	var userAdditionalVos []*models.UserAdditionalVo

	s := x.Engine.NewSession()
	defer s.Close()

	err := s.Cols("USER_KEY", "USER_TYPE", "TYPE_VALUE").
		Where("USER_KEY = ?", userKey).
		Find(&userAdditionalVos)

	if err != nil {
		return nil, err
	}

	return userAdditionalVos, nil
}

// 고객정보 저장
func (x *XormHelper) DeleteAdditionalUser(vo *models.UserAdditionalVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Delete(vo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 고객정보 조회
func (x *XormHelper) SelectUser(userKey string) (*models.UserVo, error) {
	user := new(models.UserVo)

	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("USER_KEY", "USER_NAME", "PWD").Where("USER_KEY = ?", userKey).Get(user)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}

	return user, nil
}

// 고객 로그인
func (x *XormHelper) Login(userKey string, pwd string) (*models.UserVo, error) {
	user := new(models.UserVo)

	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("USER_KEY", "USER_NAME").Where("USER_KEY = ? and PWD = ?", userKey, pwd).Get(user)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}

	return user, nil
}

// 인증번호 삽입
func (x *XormHelper) InsertTran(vo *models.ScTranVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Insert(vo)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// 인증번호 삽입
func (x *XormHelper) UpdateTranByCompleteDate(vo *models.ScTranVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Cols("CERT_COMP_DATE").Where("TRAN_KEY = ?", vo.TranKey).Update(vo)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// 인증번호 제거
func (x *XormHelper) DelTran(vo *models.ScTranVo) error {
	s := x.Engine.NewSession()
	defer s.Close()
	_, err := s.Cols("SEND_YN").Where("TRAN_KEY = ?", vo.TranKey).Update(vo)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// 인증번호 조회
func (x *XormHelper) SelectTran(hp string, num string, tran_key string) (*models.ScTranVo, error) {
	tran := new(models.ScTranVo)

	s := x.Engine.NewSession()
	defer s.Close()

	has, err := s.Cols("HP", "NUM", "EXPIRE_DATE").
		Where("TRAN_KEY = ? and HP = ? and NUM = ?", tran_key, hp, num).
		Desc("EXPIRE_DATE").Get(tran)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !has {
		return nil, errors.New("no data")
	}
	return tran, nil
}
