package types

type LoginFormData struct {
	Username string `form:"username" binding:"required,min=4,max=16"`
	Password string `form:"password" binding:"required,min=5,max=128"`
}
