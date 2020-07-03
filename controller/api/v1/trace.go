package v1

import (
	"bytes"
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/logging"
	"github.com/fabric-app/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary  传感器数据溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.SensorSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/sensors  [POST]
func Sensors(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.SensorSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceble", "query",
		[]string{"s", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, res)
}

// @Summary  图片信息溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.PicSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/pictures  [POST]
func Pictures(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.PicSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceble", "query",
		[]string{"p", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, res)
}

// @Summary  农事数据溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.FarmSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/farm  [POST]
func Farms(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.FarmSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceble", "query",
		[]string{"f", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, res)
}

// @Summary 图片下载
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.PictureSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/downloadPic  [POST]
func DownloadPic(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.PictureSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Invalid paras in json")
	}
	file, err := models.GetPicFile(reqInfo.Point, reqInfo.Date, reqInfo.Name)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "Get picture file failed.")
	}
	defer file.Close()

	buf := bytes.Buffer{}
	size, err := buf.ReadFrom(file)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "File buffer create failed.")
		return
	}
	logging.Debug("File load success,size:", size)

	appG.C.Writer.Header().Add("Content-Type", "application/octet-stream")
	appG.C.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Name())
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, buf.Bytes())
}

// @Summary  链上信息检验
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.VerifySwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/verifier  [POST]
func Verifier(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.VerifySwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}

	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, res)
}