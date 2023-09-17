package model

type User struct {
	BaseModel
	Address   string `gorm:"column:address;type:varchar(50);not null;default '' comment '用户钱包地址'" json:"address"`
	Username  string `gorm:"column:username;type:varchar(80);not null;default '' comment '用户名'" json:"username"`
	Avatar    string `gorm:"column:avatar;type:varchar(2046);not null;default '' comment '头像'" json:"avatar"`
	TwitterID int
	Twitter   Twitter
}

type Twitter struct {
	BaseModel
	TWUserID    string `gorm:"column:tw_user_id;type:varchar(40);not null;default '' comment '用户的twitter账号id'" json:"tw_user_id"`
	TWUsername  string `gorm:"column:tw_username;type:varchar(40);not null;default '' comment '用户的twitter账号名称'" json:"tw_username"`
	OAuthToken  string `gorm:"column:oauth_token;type:varchar(100);not null;default '' comment '用户的oauthtoken'" json:"oauth_token"`
	OAuthSecret string `gorm:"column:oauth_secret;type:varchar(100);not null;default '' comment '用户的oauthsecret'" json:"oauth_scret"`
}
