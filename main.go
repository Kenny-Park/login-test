package main

import (
	"crypto/tls"
	"log"
	"modules/apis"
	"modules/models"
	"modules/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func init() {
	http.DefaultTransport = &http.Transport{
		MaxIdleConns:       16,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: false,
	}
}

func main() {

	// DB 접속
	engine, err := xorm.NewEngine("mysql", "root:11111111@tcp(127.0.0.1:3306)/dev")
	if err != nil {
		log.Println(err)
	}
	defer engine.Close()

	// Database engine
	x := &utils.XormHelper{
		Engine: engine,
	}
	x.Config()

	// 웹서버 생성
	router := gin.New()

	// 토큰체크 미들웨어
	router.Use(tokenCheckMiddleware(x))
	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/js", "./js")

	// 화면
	router.GET("/join", func(c *gin.Context) {
		c.HTML(http.StatusOK, "join.tmpl", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	router.GET("/hpconfirm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hpconfirm.tmpl", gin.H{
			"rtnUrl": c.Query("rtnurl"),
		})
	})
	router.GET("/pwdchange", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pwdchange.tmpl", gin.H{})
	})
	router.GET("/get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "modify.tmpl", gin.H{})
	})

	// API
	apiLogin := &apis.ApiLogin{}
	apiUser := &apis.ApiUser{}
	apiScTran := &apis.ApiScTran{}

	// 인증번호 요청
	router.POST("/usr/auth/phone/req", apiScTran.TranRequest(x))

	// 인증번호 비교
	router.POST("/usr/auth/phone/conf", apiScTran.TranConfirm(x))

	// 로그인 생성
	router.POST("/usr/login", apiLogin.Login(x))

	// 회원가입 (insert)
	router.POST("/usr/ins", apiUser.UserJoin(x))

	// 고객정보 조회 (get)
	router.GET("/usr/get", apiUser.UserGet(x))

	// 중복체크 (get)
	router.GET("/usr/dup/:usertype/:typevalue", apiUser.Dup(x))

	// 회원정보 수정 (patch)
	router.PATCH("/usr/upd", apiUser.Modify(x))

	// 비밀번호 재설정 (patch)
	router.PATCH("/usr/upd/pwd", apiUser.ChangePassword(x))

	// 서버시작
	router.Run()
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"code": code,
		"message": message})
}

func tokenCheckMiddleware(x *utils.XormHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.RequestURI)
		if c.Request.RequestURI == "/usr/get" {
			token := c.GetHeader("Token")
			log.Println("token", token)
			var t *models.TokenVo
			var err error
			if t, err = x.GetToken(token); err != nil {
				respondWithError(c, 401, "token error")
			} else if !t.ExpireDate.After(time.Now()) {
				log.Println("token expired error")
				respondWithError(c, 401, "token expired error")
			} else {
				c.Writer.Header().Set("token", token)
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}
