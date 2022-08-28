package apis

import (
	"log"
	"modules/models"
	"modules/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiUser struct {
}

// 회원가입처리
func (api *ApiUser) UserJoin(x *utils.XormHelper) gin.HandlerFunc {

	var err error
	return func(c *gin.Context) {
		// 고객의 정보를 저장한다.
		// 이메일 이름 아이디 휴대폰번호등 암호화필요
		userId, has := c.GetPostForm("userid")
		if !has {
			log.Println("아이디를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "아이디를 전송하지 않았습니다.",
			})
			return
		}

		username, has := c.GetPostForm("username")
		if !has {
			log.Println("아이디를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "아이디를 전송하지 않았습니다.",
			})
			return
		}

		pwd, has := c.GetPostForm("pwd")
		if !has {
			log.Println("비밀번호를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "비밀번호를 전송하지 않았습니다.",
			})
			return
		}

		hp, has := c.GetPostForm("hp")
		if !has {
			log.Println("휴대폰번호를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "휴대폰번호를 전송하지 않았습니다.",
			})
			return
		}

		email, has := c.GetPostForm("email")
		if !has {
			log.Println("이메일을 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "이메일을 전송하지 않았습니다.",
			})
			return
		}

		nickname, has := c.GetPostForm("nickname")
		if !has {
			log.Println("닉네임을 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "닉네임을 전송하지 않았습니다.",
			})
			return
		}

		username, err = utils.Crypto{}.Encrypt(username)

		userId, err = utils.Crypto{}.Encrypt(userId)
		hp, err = utils.Crypto{}.Encrypt(hp)
		email, err = utils.Crypto{}.Encrypt(email)
		nickname, err = utils.Crypto{}.Encrypt(nickname)

		userVo := &models.UserVo{
			UserKey:  utils.UtilToken{}.RandString(16),
			UserName: username,
			Pwd:      utils.CryptoSha256{}.Encrypt(pwd),
			RegDate:  time.Now(),
		}

		// 이전에 해당 휴대폰번호로 로그인 하던 사람의 로그인 정보를 지워준다.
		x.DeleteAdditionalUser(&models.UserAdditionalVo{
			UserType:  "Hp",
			TypeValue: hp,
		})

		// 유저 아이디로 가능한 항목들 저장
		err = x.InsertAdditionalUser([]*models.UserAdditionalVo{
			{
				UserKey:   userVo.UserKey,
				UserType:  "UserId",
				TypeValue: userId,
			}, {
				UserKey:   userVo.UserKey,
				UserType:  "Hp",
				TypeValue: hp,
			}, {
				UserKey:   userVo.UserKey,
				UserType:  "Email",
				TypeValue: email,
			}, {
				UserKey:   userVo.UserKey,
				UserType:  "Nickname",
				TypeValue: nickname,
			},
		})

		if err != nil && err.Error() == "exists data" {
			c.JSON(409, gin.H{
				"code":    409,
				"message": "서비스가 원활하지 않습니다.",
			})
		}

		if err = x.InsertUser(userVo); err != nil {
			c.JSON(200, gin.H{
				"code":    500,
				"message": "서비스가 원활하지 않습니다.",
			})
			return
		}

		c.JSON(200, gin.H{
			"code":    200,
			"message": "",
		})
	}
}

// 고객정보 취득
func (api *ApiUser) UserGet(x *utils.XormHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		// header token으로 유저의 정보를 알아낸다.
		// 토큰값으로 유저 고유키를 가져온다.
		// 유저 고유키를 가지고 고객정보 조회한다.
		// 고객정보 전송
		var err error
		token := c.GetHeader("Token")
		userVoParsed := utils.UtilToken{}.Parse(token)

		userVo, _ := x.SelectUser(userVoParsed.UserKey)
		userAdditionalVos, _ := x.ListAdditionalUser(userVoParsed.UserKey)

		m := map[string]interface{}{}
		defer func() {
			for k, _ := range m {
				delete(m, k)
			}
		}()

		log.Println("models.UserVo-", userVo, userVoParsed)

		username, _ := utils.Crypto{}.Decrypt(userVo.UserName)

		m["UserName"] = username
		m["UserKey"] = userVoParsed.UserKey
		for _, item := range userAdditionalVos {
			m[item.UserType], _ = utils.Crypto{}.Decrypt(item.TypeValue)
		}
		if err != nil {
			log.Println("user vo is nil")
		}
		c.JSON(200, gin.H(m))
	}
}

// 필수정보 중복확인
func (api *ApiUser) Dup(x *utils.XormHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		//var err error
		usertype := c.Param("usertype")
		typevalue := c.Param("typevalue")

		typevalue, _ = utils.Crypto{}.Encrypt(typevalue)

		if vo, _ := x.GetAdditionalUser(usertype, typevalue); vo != nil {
			c.JSON(409, gin.H{
				"code":    409,
				"message": "중복된 값이 있습니다.",
			})
		} else {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "",
			})
		}
	}
}
func (api *ApiUser) Modify(x *utils.XormHelper) gin.HandlerFunc {
	return func(c *gin.Context) {

		var err error
		// 고객의 정보를 수정
		// 이메일 이름 아이디 휴대폰번호등 암호화필요
		hp, has := c.GetPostForm("hp")
		if !has {
			log.Println("휴대폰번호를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "휴대폰번호를 전송하지 않았습니다.",
			})
			return
		}

		email, has := c.GetPostForm("email")
		if !has {
			log.Println("이메일을 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "이메일을 전송하지 않았습니다.",
			})
			return
		}

		nickname, has := c.GetPostForm("nickname")
		if !has {
			log.Println("닉네임을 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "닉네임을 전송하지 않았습니다.",
			})
			return
		}
		tokenVo := utils.UtilToken{}.Parse(c.GetHeader("token"))

		hp, _ = utils.Crypto{}.Encrypt(hp)
		email, _ = utils.Crypto{}.Encrypt(email)
		nickname, _ = utils.Crypto{}.Encrypt(nickname)

		// 휴대폰은 나중에 인증한사람만의 고유한 로그인 아이디로 활용됳수 있으므로 데이타를 지우고 인서트
		x.DeleteAdditionalUser(&models.UserAdditionalVo{
			UserKey:   tokenVo.UserKey,
			UserType:  "Hp",
			TypeValue: hp,
		})
		x.InsertAdditionalUser([]*models.UserAdditionalVo{
			{
				UserKey:   tokenVo.UserKey,
				UserType:  "Hp",
				TypeValue: hp,
			},
		})

		// 그외의 값은 업데이트로 처리
		err = x.UpdateAdditionalUser([]*models.UserAdditionalVo{
			{
				UserKey:   tokenVo.UserKey,
				UserType:  "Email",
				TypeValue: email,
			}, {
				UserKey:   tokenVo.UserKey,
				UserType:  "Nickname",
				TypeValue: nickname,
			},
		})

		if err != nil && err.Error() == "exists data" {
			c.JSON(409, gin.H{
				"code":    409,
				"message": "서비스가 원활하지 않습니다.",
			})
		}

		// 현재토큰 새로 생성
		userVo, _ := x.SelectUser(tokenVo.UserKey)
		userAdditionalVos, _ := x.ListAdditionalUser(tokenVo.UserKey)
		expireTime := time.Now().AddDate(1, 0, 0)
		newToken := utils.UtilToken{}.Create(userVo, userAdditionalVos, expireTime)

		// 토큰 삭제
		x.DeleteToken(c.GetHeader("token"))

		// 토큰 저장
		x.InsertToken(&models.TokenVo{
			UserKey:    tokenVo.UserKey,
			Token:      newToken,
			ExpireDate: expireTime,
		})

		c.Writer.Header().Set("Token", newToken)

		c.JSON(200, gin.H{
			"code":    200,
			"message": "",
		})
	}
}
func (api *ApiUser) ChangePassword(x *utils.XormHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 휴대폰번호 인증 완료 후 비밀번호 재설정
		// 비(AES)
		token := c.GetHeader("Token")
		userVoParsed := utils.UtilToken{}.Parse(token)

		pwd, has := c.GetPostForm("pwd")
		if !has {
			log.Println("비밀번호를 전송하지 않았습니다.")
			c.JSON(200, gin.H{
				"code":    200,
				"message": "비밀번호를 전송하지 않았습니다.",
			})
			return
		}

		if userVoParsed.UserId == "" {
			c.Writer.Header().Set("RtnType", "1")
		} else {
			c.Writer.Header().Set("RtnType", "2")
			userAdditionalVos, _ := x.ListAdditionalUser(userVoParsed.UserKey)
			expireTime := time.Now().AddDate(1, 0, 0)
			newToken := utils.UtilToken{}.Create(&models.UserVo{
				UserKey:  userVoParsed.UserKey,
				UserName: userVoParsed.UserName,
			}, userAdditionalVos, expireTime)

			if vo, _ := x.GetTokenByUserKey(userVoParsed.UserKey); vo != nil {
				// 토큰만료 시
				if vo.ExpireDate.After(time.Now()) {
					// 삭제 후
					x.DeleteToken(vo.Token)
					// 저장
					x.InsertToken(&models.TokenVo{
						UserKey:    userVoParsed.UserKey,
						Token:      newToken,
						ExpireDate: expireTime,
					})
					c.Writer.Header().Set("Token", newToken)
				} else {
					// 토큰이 만료되지 않았다면 기존토큰 삽입
					c.Writer.Header().Set("Token", vo.Token)
				}
			} else {
				c.JSON(401, gin.H{
					"code":    401,
					"message": "",
				})
			}
		}

		// 비밀번호 변경
		x.UpdatePwd(&models.UserVo{
			UserKey: userVoParsed.UserKey,
			Pwd:     utils.CryptoSha256{}.Encrypt(pwd),
		})
		// 응답결과
		c.JSON(200, gin.H{
			"code":    200,
			"message": "",
		})
	}
}
