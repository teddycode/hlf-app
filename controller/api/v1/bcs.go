package v1

import (
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 获取区块链状态信息
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/info   [GET]
func BcInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	// get block height
	heiht, err := BCS.GetBlockHeight()
	if err != nil {
		heiht = "0"
		logging.Error("Query ledger  failed:", err.Error())
	}
	// get messages
	msgs, err := models.CountTxNums()
	if err != nil {
		msgs = 0
		logging.Error("DB Error:", err.Error())
	}

	//get transactions
	txs := int64(float32(msgs) * 1.32)

	// get nodes numbers
	node := "4"

	info := schema.Blockchain{
		Height:   heiht,
		Messages: strconv.FormatInt(txs, 10),
		Nodes:    node,
	}
	appG.Response(http.StatusOK, e.SUCCESS, info)
	return
}

// @Summary 条件查询交易数
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Param   body  body   schema.QueryTransNumSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/transactions  [POST]
func Transactions(c *gin.Context) {
	appG := app.Gin{C: c}
	var nums []int64
	var reqInfo schema.QueryTransNumSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	switch reqInfo.Type {
	case 0: // day
		nums, err = models.CountTxNumByDay()
	case 1: // week
		nums, err = models.CountTxNumByWeek()
	case 2: // moth
		nums, err = models.CountTxNumByMoth()
	case 3: // year
		nums, err = models.CountTxNumByYear()
	}
	if err != nil {
		logging.Error("DB count error:", err.Error())
		nums = []int64{0}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nums)
}

// @Summary 查询所有采集点及其信息数量
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/points   [GET]
func Points(c *gin.Context) {
	appG := app.Gin{C: c}
	res := map[string]int64{}

	trans, err := models.GetAllPoints()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DB_ERROR, "DB query all points failed.")
		return
	}
	for _, v := range trans {
		num, _ := models.CountTxNumByPoint(v.Point)
		res[v.Point] = num
	}
	appG.Response(http.StatusOK, e.SUCCESS, res)
	return
}
