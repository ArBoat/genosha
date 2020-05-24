package models

type (
	Signup struct {
		UserName string `json:"userName" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Roles    string `json:"roles" binding:"required"`
		Spell    string `json:"spell" binding:"required"`
		TimeZone int    `json:"timeZone"`
	}
	Sign2Back struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	Sign2BackResp struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	ChangePassWD struct {
		UserEmail   string `json:"userEmail" binding:"required"`
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	SendRestPassWDEmail struct {
		To string `json:"to" binding:"required"`
	}
	RestPassWD struct {
		Email       string `json:"email" binding:"required"`
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
)

type (
	BackendError struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)
