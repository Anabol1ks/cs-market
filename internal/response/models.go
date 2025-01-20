package response

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	SteamID   string `json:"steam_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarurl"`
}
