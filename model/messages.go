package model

type Messages struct {
	Id         uint64  `json:"id" xorm:"not null pk autoincr BIGINT"`
	AppsId     uint    `json:"apps_id" xorm:"INT not null comment('app id')"`
	From       uint64  `json:"from" xorm:"BIGINT notnull comment('拨号用户id') index"`
	Target     uint64  `json:"target" xorm:"BIGINT notnull comment('目标id') index"`
	TargetType byte    `json:"target_type" xorm:"TINYINT notnull default 1 comment('目标类型 1:用户 2:聊天室')"`
	Msg        string  `json:"msg" xorm:"TEXT comment('消息内容')"`
	MsgType    byte    `json:"msg_type" xorm:"TINYINT notnull default 1 comment('消息类型 1:文字 2:图片 3:语音')"`
	Status     byte    `json:"status" xorm:"TINYINT notnull default 1 comment('消息状态 -1:删除 0:撤回 1:正常')"`
	CreateAt   string  `json:"create_at" xorm:"created not null DATETIME"`
	UpdateAt   *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (m Messages) TableName() string {
	return "messages"
}
