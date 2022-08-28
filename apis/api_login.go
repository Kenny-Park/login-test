package apis

import (
	"log"
	"modules/models"
	"modules/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiLogin struct {
}

// 로그인
func (api *ApiLogin) Login(x *utils.XormHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 로그인아이디에 해당하는 부분을 유저고유키를 저장하고 있는 테이블에서 조회
		// 유저 고유키와 비번(SHA256)을 가지고 고객정보 조회
		// 토큰값 생성
		// 리턴값은 error code와 메시지, 고객토큰으로 처리가 가능할 것 같다
		loginId, has := c.GetPostForm("userid")
		if !has {
			log.Println("아이디를 전송하지 않았습니다.")
		}
		pwd, has := c.GetPostForm("pwd")
		if !has {
			log.Println("비밀번호를 전송하지 않았습니다.")
		}

		loginId, _ = utils.Crypto{}.Encrypt(loginId)

		// 타입별 로그인
		var userAdditionalVo *models.UserAdditionalVo
		userType := []string{"UserId", "Nickname", "Hp", "Email"}
		for _, item := range userType {
			if userAdditionalVo, _ = x.GetAdditionalUser(item, loginId); userAdditionalVo != nil {
				break
			}
		}
		if userAdditionalVo == nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "로그인에 실패하였습니다.",
			})
			return
		}

		// 비밀번호 암호화
		pwd = utils.CryptoSha256{}.Encrypt(pwd)
		var userVo *models.UserVo

		userVo, _ = x.Login(userAdditionalVo.UserKey, pwd)
		if userVo == nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "로그인에 실패하였습니다.",
			})
			return
		}

		// 유저키를 가져옴
		userAdditionalVos, _ := x.ListAdditionalUser(userVo.UserKey)
		expireTime := time.Now().AddDate(1, 0, 0)
		newToken := utils.UtilToken{}.Create(userVo, userAdditionalVos, expireTime)

		if vo, _ := x.GetTokenByUserKey(userVo.UserKey); vo != nil {
			// 토큰만료 시
			if vo.ExpireDate.After(time.Now()) {
				// 삭제 후
				x.DeleteToken(vo.Token)
				// 저장
				x.InsertToken(&models.TokenVo{
					UserKey:    userVo.UserKey,
					Token:      newToken,
					ExpireDate: expireTime,
				})
				c.Writer.Header().Set("Token", newToken)
			} else {
				// 토큰이 만료되지 않았다면 기존토큰 삽입
				c.Writer.Header().Set("Token", vo.Token)
			}
			c.JSON(200, gin.H{
				"code":    200,
				"message": "",
			})
		} else {
			// 토큰 정보가 없는 경우 저장
			if userVo != nil {
				x.InsertToken(&models.TokenVo{
					UserKey:    userVo.UserKey,
					Token:      newToken,
					ExpireDate: expireTime,
				})
				c.Writer.Header().Set("Token", newToken)
				c.JSON(200, gin.H{
					"code":    200,
					"message": "",
				})
			} else {
				c.JSON(401, gin.H{
					"code":    401,
					"message": "",
				})
			}
		}
	}
}
