package schema

//用户表
type User struct {
	ID         uint    `json:"id"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Password   string  `json:"password"`
	balance    float32 `json:"balance"`
	CreatedOn  uint    `json:"created_on"`
	ModifiedOn uint    `json:"modified_on"`
	DeletedOn  uint    `json:"deleted_on"`
	Secret     string  `json:"secret"`
}

//注册
type RegSwag struct {
	Username string `json:"username" binding:"required"`  //用户名
	Identity string `json:"identity" binding:"required"`  // identity
	Password string `json:"password"  binding:"required"` //密码
}

//登录
type AuthSwag struct {
	Name     string `json:"name"`     //登录邮箱
	Password string `json:"password"` //登录密码
}

//修改密码
type PasswordSwag struct {
	OldPassword string `json:"old_password"` //旧密码
	NewPassword string `json:"new_password"` //新密码
}

// 修改用户信息
type CurrentUserSwag struct {
	Email   string `json:"email"` //用户名
	Phone   string `json:"phone"` //电话
	Address string `json:"address"`
	Header  string `json:"header"`
}
