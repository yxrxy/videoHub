package model

type User struct {
	ID           int64
	Username     string
	Password     string
	AvatarURL    string
	Token        string
	RefreshToken string
}
