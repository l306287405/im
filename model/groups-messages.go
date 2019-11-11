package model

import "encoding/json"

type GroupsMessages struct {
	Id         uint64  `json:"id" xorm:"not null pk autoincr BIGINT"`
	AppsId     uint    `json:"apps_id" xorm:"INT not null comment('app id')"`
	From       uint64  `json:"from" xorm:"BIGINT notnull comment('拨号用户id') index"`
	To     	   uint64  `json:"to" xorm:"BIGINT notnull comment('目标id') index"`
	Text       string  `json:"text" xorm:"TEXT comment('消息内容')"`
	TextType   byte    `json:"text_type" xorm:"TINYINT notnull default 1 comment('消息类型 1:文字 2:图片 3:语音')"`
	Status     byte    `json:"status" xorm:"TINYINT notnull default 1 comment('消息状态 -1:删除 0:撤回 1:正常')"`
	CreateAt   string  `json:"create_at,omitempty" xorm:"created notnull TIMESTAMP index"`
	UpdateAt   *string `json:"update_at,omitempty" xorm:"updated TIMESTAMP"`
	ErrCode    *string `json:"err_code,omitempty" xorm:"-"`
	ErrMsg     *string `json:"err_msg,omitempty" xorm:"-"`
}

func (m GroupsMessages) TableName() string {
	return "gourps_messages"
}

func (m *GroupsMessages) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *GroupsMessages) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}
