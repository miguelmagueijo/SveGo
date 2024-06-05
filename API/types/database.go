package types

type UserJwtData struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserMeData struct {
	UserJwtData
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type RefreshTokenData struct {
	Id        string `json:"id"`
	JwtId     string `json:"jwt_id"`
	UserId    int64  `json:"user_id"`
	IsActive  bool   `json:"is_active"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
