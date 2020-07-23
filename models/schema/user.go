package schema

//用户表
type User struct {
	ID         int    `json:"id"`
	UserName   string `json:"user_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	CreatedOn  uint   `json:"created_on"`
	ModifiedOn uint   `json:"modified_on"`
	DeletedOn  uint   `json:"deleted_on"`
	Secret     string `json:"secret"`
}

//注册
type RegSwag struct {
	Username string `json:"username" binding:"required" example:"teddy"`            //用户名
	Role     int `json:"role" example:"1"`                                       // 用户角色   //0:管理员,1：普通员工
	Identity string `json:"identity" binding:"required" example:"4522261111111111"` // identity
	Password string `json:"password"  binding:"required" example:"1234"`            //密码
}

//登录
type AuthSwag struct {
	UserName string `json:"user_name" example:"admin"`      //登录用户名
	Password string `json:"password" example:"lzawt_admin"` //登录密码(管理员)
}

//修改密码
type PasswordSwag struct {
	OldPassword string `json:"old_password" example:"1234"`   //旧密码
	NewPassword string `json:"new_password" example:"123456"` //新密码
}

// 修改用户信息
type CurrentUserSwag struct {
	Email   string `json:"email" example:"123456@qq.com"` //用户名
	Phone   string `json:"phone" example:"18677337725"`   //电话
	Address string `json:"address" example:"广州市白云区xx路2号"`
	//Header  string `json:"header" example:"default"`
}

// 农事记录
type FarmRecordSwag struct {
	OperName  string `json:"oper_name"  example:"施肥"`
	StartTime string `json:"start_time" example:"1594382265"`
	EndTime   string `json:"end_time"  example:"1594382165"`
	OperType  string `json:"oper_type" example:"施肥"`
	Info      string `json:"info" example:"撒了点肥料"`
}

type RevokeSwag struct {
	UserName []string `json:"user_name" example:"teddy,teddy1"`
}

// 返回用户列表
type UserListSwag struct {
	ID        int    `json:"id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CreatedOn int64    `json:"created_on"`
	Role      int    `json:"role"`
}
