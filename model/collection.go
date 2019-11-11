package model

type Collections struct {
	Id       uint64  `json:"id" xorm:"notnull pk autoincr BIGINT"`
	AppsId   uint    `json:"apps_id" xorm:"INT notnull comment('app id') index"`
	Uid      uint64  `json:"uid" xorm:"BIGINT notnull comment('归属用户id')"`
	From	 uint64	 `json:"from" xorm:"BIGINT notnull comment('来自')"`
	FromType byte	 `json:"from_type" xorm:"TINYINT notnull comment('文件状态 0:私聊用户 1:群聊房间')"`
	Text     string  `json:"text" xorm:"TEXT comment('消息内容')"`
	TextType byte    `json:"text_type" xorm:"TINYINT notnull default 1 comment('消息类型 参考messages消息类型') index"`
	CreateAt string  `json:"create_at" xorm:"created TIMESTAMP notnull"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (c Collections) TableName() string {
	return "collections"
}
