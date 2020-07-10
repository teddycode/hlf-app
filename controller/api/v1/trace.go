package v1

import (
	"encoding/json"
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	DATA_TYPE_SENSOR = 0
	DATA_TYPE_PIC    = 1
	DATA_TYPE_FARM   = 2
)

// @Summary  传感器数据溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.SensorSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/sensor  [POST]
func Sensors(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.SensorSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
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
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/picture  [POST]
func Pictures(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.PicSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
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
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/farmData  [POST]
func Farms(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.FarmSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
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
//func DownloadPic(c *gin.Context) {
//	appG := app.Gin{C: c}
//	var reqInfo schema.PictureSwag
//	err := c.BindJSON(&reqInfo)
//	if err != nil {
//		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Invalid paras in json")
//	}
//	file, err := transh.GetPicFile(reqInfo.Point, reqInfo.Date, reqInfo.Name)
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "Get picture file failed.")
//	}
//	defer file.Close()
//
//	buf := bytes.Buffer{}
//	size, err := buf.ReadFrom(file)
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "File buffer create failed.")
//		return
//	}
//	logging.Debug("File load success,size:", size)
//
//	appG.C.Writer.Header().Add("Content-Type", "application/octet-stream")
//	appG.C.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Name())
//	appG.Response(http.StatusOK, e.SUCCESS, buf.Bytes())
//}

// @Summary  链上信息检验
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.VerifySwag   true "输入交易哈希，返回交易内容（包含文件内容哈希值）"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/verify  [POST]
func Verifier(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.VerifySwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	res, err := BCS.QueryTxByID(reqInfo.Hash, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_LEDGER_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, res)
}

// @Summary 数据上链接口
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.UploadSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/trace/upload  [POST]
func Uploader(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.UploadSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Bind json error.")
		return
	}
	// check data format
	if reqInfo.Type == "s" { //  sensor data
		var sensor schema.BCSensor
		err := json.Unmarshal([]byte(reqInfo.Raw), &sensor)
		if err != nil {
			appG.Response(http.StatusOK, e.INVALID_PARAMS, "Unmarshal json error.")
			return
		}
	} else if reqInfo.Type == "p" { // pics data
		var pic schema.BCPic
		err := json.Unmarshal([]byte(reqInfo.Raw), &pic)
		if err != nil {
			appG.Response(http.StatusOK, e.INVALID_PARAMS, "Unmarshal json error.")
			return
		}
	} else {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Data type not found.")
		return
	}
	txID, err := BCS.InvokeCC("traceable", "add",
		[][]byte{[]byte(reqInfo.Type), []byte(reqInfo.Point), []byte(reqInfo.Raw)}, setting.Peers)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_INVOKE_FAILED, "Chaincode traceable invoke failed.")
		return
	}
	id, _ := models.NewTx(&models.Transaction{
		Timestamp: int(time.Now().Unix()),
		Type:      reqInfo.Type,
		Hash:      string(txID),
		Point:     reqInfo.Point,
	})
	appG.Response(http.StatusOK, e.SUCCESS, id)
}
