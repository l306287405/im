package model

type Apps struct {
	Id        uint    `json:"id" xorm:"notnull pk autoincr INT"`
	KeyId     string  `json:"key_id" xorm:"VARCHAR(11) notnull comment('授权id') unique(key_id_secret)"`
	KeySecret string  `json:"key_secret" xorm:"VARCHAR(11) notnull comment('授权密钥') unique(key_id_secret)"`
	Name      string  `json:"name" xorm:"VARCHAR(11) notnull comment('应用名称')"`
	Token     *string `json:"token" xorm:"VARCHAR(255) comment('授权令牌') index"`
	Status    byte    `json:"status" xorm:"TINYINT notnull default 1 comment('状态 0:过期 1:正常')"`
	CreatedAt string  `json:"created_at" xorm:"TIMESTAMP created notnull"`
	UpdatedAt *string `json:"updated_at" xorm:"TIMESTAMP updated"`
}

func (app *Apps) TableName() string {
	return "apps"
}
