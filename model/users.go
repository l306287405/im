package model

type Users struct {
	Id       uint64  `json:"id" xorm:"not null pk autoincr BIGINT"`
	Account  string  `json:"mobile" xorm:"not null comment('账号') unique VARCHAR(11)"`
	Password string  `json:"password" xorm:"not null comment('密码') VARCHAR(255)"`
	Nickname string  `json:"nickname" xorm:"not null comment('昵称') VARCHAR(16)"`
	Token    *string `json:"token" xorm:"comment('令牌') VARCHAR(255)"`
	Status   byte    `json:"status" xorm:"not null default 1 comment('状态 -1:删除 0:禁用 1:正常') TINYINT(4)"`
	CreateAt string  `json:"create_at" xorm:"created not null DATETIME"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (u Users) TableName() string {
	return "users"
}
