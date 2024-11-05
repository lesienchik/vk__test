package models

import "github.com/golang-jwt/jwt/v5"

/*
Здесь описываются структуры, которые используются для обработки и передачи данных.
*/

type UserRegReq struct { // При регистрации пользователя.
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthReq struct { // При аутентификации пользователя.
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthClaims struct { // Для создания jwt-токенов аутентифицированного пользователя.
	Id int `json:"id"`
	jwt.RegisteredClaims
}

type UserAccessResp struct { // Для отдачи access-токена в теле ответа.
	Token string `json:"access_token"`
}

type UserChangePassReq struct { // При смене пароля пользователя.
	Id          int    `json:"-"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserResetPassReq struct {
	Id    int    `json:"-"`
	Email string `json:"email"`
}

type UserConfirmResetPassReq struct {
	Email       string `json:"-"`
	NewPassword string `json:"new_password"`
}
