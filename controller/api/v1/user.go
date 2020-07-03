package v1

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/bcs"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/logging"
	"github.com/fabric-app/pkg/setting"
	"github.com/fabric-app/pkg/util/hash"
	"github.com/fabric-app/pkg/util/rand"
	"github.com/jinzhu/gorm"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/util"
)

const HEADER_IMAGE_PATH = "./test/header/images/"

var BCS = bcs.New(setting.BcConf, "org1", "Admin", "User1")

type auth struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type currentUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// @Summary   注册用户
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.RegSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/user/reg  [POST]
func Reg(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.RegSwag //用户表字段
	var data interface{}
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	passwdEncode := hash.EncodeMD5(reqInfo.Password)
	// register in ca
	str, ok := BCS.RegisterUser(reqInfo.Username, "org1", passwdEncode, "user")
	if !ok {
		appG.Response(http.StatusOK, e.ERROR_CA_REG_FAILED, str)
		return
	}
	// register identity in blockchain
	txID, err := BCS.InvokeCC("user", "add",
		[][]byte{[]byte(reqInfo.Username), []byte(hash.EncodeMD5(reqInfo.Identity))}, setting.Peers)
	if err != nil {
		str, _ := BCS.RevokeUser(reqInfo.Username, "org1", reqInfo.Password, "user")
		logging.Error("Invoke failed: add user identity to chain failed!. Revoke res:", str)
		appG.Response(http.StatusOK, e.ERROR_CC_INVOKE_FAILED, data)
		return
	}
	logging.Debug("chaincode invoke success! tx id:" + txID)
	if user, isExist := models.FindUserByName(reqInfo.Username); isExist != gorm.ErrRecordNotFound {
		_, err := models.DelUser(&user) // delete old user
		logging.Debug("Found old user and deleted:", err)
	}
	var newUser models.User
	newUser.Username = reqInfo.Username
	newUser.Identity = reqInfo.Identity
	newUser.Role = 1
	newUser.Password = passwdEncode //密码md5值保存
	newUser.CaSecure = passwdEncode //密码md5值保存
	newUser.Secret = rand.RandStringBytesMaskImprSrcUnsafe(5)
	newUser.CreatedOn = int(time.Now().Unix())
	newUser.ModifiedOn = int(time.Now().Unix())
	newUser.Header = "default"
	userId, isSuccess := models.NewUser(&newUser)
	if userId > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"id": userId})
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, isSuccess)
}

// @Summary   用户登录 获取token 信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.AuthSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router /api/v1/user/auth  [POST]
func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo auth
	var data string
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	res, ok := BCS.EnrollUser(reqInfo.Username,
		"org1", hash.EncodeMD5(reqInfo.Password), "user")
	if !ok {
		appG.Response(http.StatusOK, e.ERROR_CA_ENROLL_FAILED, res)
		return
	}

	user, err := models.FindUserByName(reqInfo.Username)
	if err == nil || len(user.Phone) == 0 {
		data = "First Login"
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
		"msg":   data,
	})
}

// @Summary 刷新token
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/refreshtoken  [GET]
func RefreshToken(c *gin.Context) {
	var data interface{}
	var code int
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	if Authorization == "" {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	}
	token, err := util.RefreshToken(Authorization)
	if err != nil {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": err,
		})
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}

// @Summary 用户登出
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/logout  [POST]
func Logout(c *gin.Context) {
	var data interface{}
	var code int
	appG := app.Gin{C: c}
	code = e.SUCCESS
	claims := c.MustGet("claims").(*util.Claims)
	if claims == nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST, err)
		return
	}
	user, err := models.FindUserById(id)
	if err != nil {
		code = e.ERROR_EXIST_FAIL
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": err,
		})
	}
	_, isSuccess := models.UpdateUserSecret(&user)
	if isSuccess != nil {
		code = e.ERROR_EDIT_FAIL
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": isSuccess,
		})
	}
	appG.Response(http.StatusOK, code, map[string]interface{}{
		"data": data,
	})

}

// @Summary 获取登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/currentuser   [GET]
func CurrentUser(c *gin.Context) {
	var code int
	var data interface{}
	var user models.User
	var curUser currentUser
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	//token := c.Query("token")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		user, err = models.FindUserByName(claims.Audience)
		if err != nil {
			code = e.ERROR_EXIST
		} else {
			curUser = currentUser{
				Id:       user.ID,
				Email:    user.Email,
				Role:     user.Role,
				Username: user.Username,
				Phone:    user.Phone,
				Address:  user.Address,
			}
		}
	}

	if code != e.SUCCESS {
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"data": curUser,
		})
	}
}

// @Summary 修改登录用户信息
// @Tags 	用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Param   body  body   schema.CurrentUserSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/modify   [POST]
func ModifyUser(c *gin.Context) {
	var code int
	var data interface{}
	var user models.User
	var reqInfo schema.CurrentUserSwag
	appG := app.Gin{C: c}
	code = e.SUCCESS
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	//token := c.Query("token")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		user, err = models.FindUserByName(claims.Audience)
		if err != nil {
			code = e.ERROR_EXIST
		} else {
			user.Email = reqInfo.Email
			user.Phone = reqInfo.Phone
			user.Address = reqInfo.Address
			user.Header = reqInfo.Header
			_, err := models.UpdateUserInfo(&user)
			if err != nil {
				code = e.ERROR_EXIST
			}
		}
	}

	if code != e.SUCCESS {
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"data": "success",
		})
	}

}

// @Summary 登录用户修改密码
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.PasswordSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/password   [POST]
func Password(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.PasswordSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{
			"data": "Invalid json inputs",
		})
		return
	}
	claims := c.MustGet("claims").(*util.Claims)
	if claims == nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH, map[string]interface{}{
			"data": "Auth error",
		})
		return
	}
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST, map[string]interface{}{
			"data": err.Error(),
		})
		return
	}
	user, err := models.FindUserById(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST, map[string]interface{}{
			"data": err.Error(),
		})
		return
	}
	if hash.EncodeMD5(reqInfo.OldPassword) != user.Password {
		appG.Response(http.StatusOK, e.INVALID_OLD_PASS, map[string]interface{}{
			"data": err.Error(),
		})
		return
	}
	_, isOk := models.UpdateUserNewPassword(&user, reqInfo.NewPassword)
	if isOk != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_FAIL, map[string]interface{}{
			"data": "update table failed",
		})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"data": "ok",
	})
	return
}

// @Summary 用户记录上传
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.FarmRecordSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/record   [POST]
func Record(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.FarmRecordSwag
	var code int
	var userName string
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{
			"data": "Invalid json inputs",
		})
		return
	}
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	//token := c.Query("token")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		userName = claims.Audience
	}
	if code != e.SUCCESS {
		appG.Response(http.StatusOK, e.ERROR_AUTH_NOT_PERMISSION, map[string]interface{}{
			"data": "Invalid authorization",
		})
	}
	strJson, _ := json.Marshal(reqInfo)
	txID, err := BCS.InvokeCC("traceble", "add",
		[][]byte{[]byte("f"), []byte(userName), strJson}, setting.Peers)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_INVOKE_FAILED, map[string]interface{}{
			"data": "chaincode invoke failed.",
		})
	}
	transaction := models.Transaction{
		Timestamp: int(time.Now().Unix()),
		Type:      1,
		Hash:      string(txID),
		Point:     userName,
	}
	id, err := models.NewTx(&transaction)
	if err != nil {
		logging.Error("DB Error:", err.Error())
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"data": id,
	})
	return
}

// @Summary 用户查询农事数据类型
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/operType   [GET]
func Operations(c *gin.Context) {
	appG := app.Gin{C: c}
	types, err := models.QueryFarmTypes()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST, "Empty")
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, types)
	return
}

// @Summary 用户更换头像
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/user/header   [GET]
func Headers(c *gin.Context) {
	appG := app.Gin{C: c}
	var userName string
	var code int
	f, err := c.FormFile("file")
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "Form file failed.")
		return
	}
	if f == nil || f.Size == 0 {
		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "Form file size is empty.")
		return
	}
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		userName = claims.Audience
	}
	if code != e.SUCCESS {
		appG.Response(http.StatusOK, e.ERROR_AUTH_NOT_PERMISSION, map[string]interface{}{
			"data": "Invalid authorization",
		})
	}
	path := path.Join(HEADER_IMAGE_PATH, userName, ".jpg")
	err = c.SaveUploadedFile(f, path)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_FILE_SAVE_FAILED, map[string]interface{}{
			"data": "Save file failed",
		})
		return
	}
	models.UpdateUserheader(userName,userName)  // update table
	appG.Response(http.StatusOK, e.SUCCESS, "OK")
}
