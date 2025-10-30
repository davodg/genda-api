package users

type User struct {
	Id           string  `json:"id"`
	Name         string  `json:"name" validate:"required"`
	Email        string  `json:"email" validate:"required"`
	IsAdmin      bool    `json:"is_admin"`
	Type         string  `json:"type"`
	BirthDate    string  `json:"birth_date"`
	SocialNumber *string `json:"social_number,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	PhotoUrl     *string `json:"photo_url,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

type LoginRequest struct {
	Code        string `json:"code"`
	RedirectUrl string `json:"redirect_url"`
}

type LogoutRequest struct {
	Token string `json:"accessToken"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type GetUserResponse struct {
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Users []User `json:"users"`
}
