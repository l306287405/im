package model

type Chatrooms struct {
	Id       uint64  `json:"id" xorm:"notnull pk autoincr BIGINT"`
	Uid      *uint64 `json:"uid" xorm:"BIGINT notnull comment('归属用户id')"`
	Name     string  `json:"name" xorm:"VARCHAR(11) notnull comment('名称')"`
	Desc     string  `json:"desc" xorm:"TEXT comment('描述')"`
	MaxUsers *uint16 `json:"max_users" xorm:"INT notnull default 200 comment('最大人数 默认200 最大2000')"`
	Approval *byte   `json:"approval" xorm:"TINYINT notnull default 0 comment('入群批准 0:不需要 1:需要')"`
	Status   byte    `json:"status" xorm:"TINYINT notnull default 1 comment('状态 -1:封禁 1:正常')"`
	CreateAt string  `json:"create_at" xorm:"created not null DATETIME"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (c Chatrooms) TableName() string {
	return "chatrooms"
}
