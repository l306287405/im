package model

type Apps struct {
	Id        uint64  `json:"id" xorm:"notnull pk autoincr BIGINT"`
	KeyId     string  `json:"key_id" xorm:"VARCHAR(11) notnull comment('授权id')"`
	KeySecret string  `json:"key_secret" xorm:"VARCHAR(11) notnull comment('授权密钥')"`
	Name      string  `json:"name" xorm:"VARCHAR(11) notnull comment('应用名称')"`
	Token     *string `json:"token" xorm:"VARCHAR(11) comment('授权令牌')"`
	Status    byte    `json:"status" xorm:"TINYINT notnull default 1 comment('状态 0:过期 1:正常')"`
	CreatedAt string  `json:"created_at" xorm:"TIMESTAMP created notnull"`
	UpdatedAt *string `json:"updated_at" xorm:"TIMESTAMP updated"`
}

func (app Apps) TableName() string {
	return "apps"
}
