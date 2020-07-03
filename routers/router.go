package routers

import (
	"github.com/fabric-app/controller/api/v1"
	_ "github.com/fabric-app/docs"
	"github.com/fabric-app/middleware"
	"github.com/fabric-app/middleware/jwt"
	"github.com/fabric-app/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())      //日志
	r.Use(middleware.Cors()) // 跨域请求
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode) //设置运行模式

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //api注释文档

	apiv1 := r.Group("/api/v1")
	//用户登录
	apiv1.POST("user/login", v1.Auth)

	apiv1.Use(jwt.JWT()) //令牌 验证中间件
	{
		/* 用户管理类 */
		//用户注册
		apiv1.POST("user/register", v1.Reg)
		//当前用户信息查询
		apiv1.GET("user/current", v1.CurrentUser)
		//刷新token
		apiv1.GET("user/refresh", v1.RefreshToken)
		//用户登出
		apiv1.POST("user/logout", v1.Logout)
		//登录用户修改密码
		apiv1.POST("user/password", v1.Password)
		//登录用户修改个人信息
		apiv1.POST("user/update", v1.ModifyUser)
		//用户农事数据上传
		apiv1.POST("user/record", v1.Record)
		// 查询农事操作类型
		apiv1.GET("user/operType", v1.Operations)
		// 更换用户头像
		apiv1.GET("user/setHeader", v1.SetHeader)
		//用户头像获取
		apiv1.GET("user/getHeader", v1.GetHeader)
		// 用户注销接
		apiv1.POST("user/revoke", v1.Revoker)

		/* 区块链监控类 */
		// 获取当前区块高度、信息数、总交易数、活跃节点数
		apiv1.GET("bcs/chainInfo", v1.BcInfo)
		//条件查询交易数
		apiv1.POST("bcs/transactions", v1.Transactions)
		//查询所有采集点及其信息数量
		apiv1.GET("bcs/points", v1.Points)

		/*溯源查询类*/
		//传感器数据溯源
		apiv1.POST("trace/sensor", v1.Sensors)
		// 图片信息溯源
		apiv1.POST("trace/picture", v1.Pictures)
		// 农事数据溯源
		apiv1.POST("trace/farmData", v1.Farms)
		// 下载溯源图片
		apiv1.POST("trace/downloadPic", v1.DownloadPic)
		// 通过哈希值校验链上信息
		apiv1.POST("trace/verify", v1.Verifier)
		// 接收数据并且上链
		apiv1.POST("trace/upload", v1.Uploader)
	}

	return r
}
