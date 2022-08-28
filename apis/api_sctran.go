package apis

import (
	"log"
	"modules/models"
	"modules/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiScTran struct {
}

// 문자전송요청
func (api *ApiScTran) TranRequest(x *utils.XormHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 휴대폰번호는 암호화 한다.(AES)
		// 파라메터의 유저의 휴대폰번호와 인증번호를 받아와서 현재시간과 인증문자가 만료되는 시간을 비교하여 값을 리턴한다.
		// 리턴값은 error code와 메시지면 처리가 가능할 것 같다
		hp, _ := c.GetPostForm("hp")
		userid, _ := c.GetPostForm("userid")
		join, _ := c.GetPostForm("join")
		userid, _ = utils.Crypto{}.Encrypt(userid)
		log.Println("userid", userid)
		if voUser, _ := x.GetAdditionalUser("UserId", userid); voUser == nil && join != "true" {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "없는 고객",
			})
			return
		}

		tranKey := utils.UtilToken{}.RandString(16)
		num := utils.UtilToken{}.RandInt(6)
		hp, _ = utils.Crypto{}.Encrypt(hp)

		x.InsertTran(&models.ScTranVo{
			TranKey:    tranKey,
			Hp:         hp,
			Num:        num,
			ExpireDate: time.Now().Add(5 * time.Minute),
		})

		c.JSON(200, gin.H{
			"code":    200,
			"message": "",
			// 트랜키만 내려줘야 하나, 인증번호를 확인할수 있는 상황을 만들수 없으므로 번호도 넘겨줌
			"tranKey": tranKey,
			"num":     num,
		})
	}

}

// 인증문자확인
func (api *ApiScTran) TranConfirm(x *utils.XormHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 휴대폰번호는 암호화 한다.(AES)
		// 파라메터의 유저의 휴대폰번호와 인증번호를 받아와서 현재시간과 인증문자가 만료되는 시간을 비교하여 값을 리턴한다.
		// 리턴값은 error code와 메시지면 처리가 가능할 것 같다
		hp, _ := c.GetPostForm("hp")
		num, _ := c.GetPostForm("num")
		userid, _ := c.GetPostForm("userid")
		tranKey, _ := c.GetPostForm("trankey")
		join, _ := c.GetPostForm("join")

		hp, _ = utils.Crypto{}.Encrypt(hp)
		userid, _ = utils.Crypto{}.Encrypt(userid)
		vo, _ := x.SelectTran(hp, num, tranKey)
		tmpToken := c.GetHeader("Token")

		if voUser, _ := x.GetAdditionalUser("UserId", userid); voUser == nil && join != "true" {
			c.JSON(401, gin.H{
				"code":    403,
				"message": "없는 고객",
			})
			return
		}

		// 받아온 토큰이 없다면 임시토큰 생성
		if len(tmpToken) <= 0 || strings.EqualFold(tmpToken, "null") {
			if voAdditionalUser, _ := x.GetAdditionalUser("Hp", hp); voAdditionalUser != nil {
				voUser, _ := x.SelectUser(voAdditionalUser.UserKey)
				tmpToken = utils.UtilToken{}.Create(voUser, nil, time.Now().AddDate(0, 0, 1))
			}
		}

		if vo != nil {
			x.UpdateTranByCompleteDate(&models.ScTranVo{
				CertCompDate: time.Now(),
			})
			// Token 정보를 삽입해준다.
			c.Writer.Header().Set("Token", tmpToken)
			// 완료정보를 헤더에 삽입해준다.
			c.Writer.Header().Set("CertCompleted", time.Now().Format("20060102"))
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
