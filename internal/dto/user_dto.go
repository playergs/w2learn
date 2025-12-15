package dto

type CreateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=32"`
}
