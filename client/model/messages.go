package model

import "encoding/json"

type Messages struct {
	Id         uint64  `json:"id" xorm:"not null pk autoincr BIGINT"`
	AppsId     uint    `json:"apps_id" xorm:"INT not null comment('app id') index"`
	From       uint64  `json:"from" xorm:"BIGINT notnull comment('拨号用户id') index"`
	To     	   uint64  `json:"to" xorm:"BIGINT notnull comment('目标id') index"`
	Text       string  `json:"text" xorm:"TEXT comment('消息内容')"`
	TextType   byte    `json:"text_type" xorm:"TINYINT notnull default 1 comment('消息类型 1:文字 2:图片 3:语音 4:视频')"`
	Status     byte    `json:"status" xorm:"TINYINT notnull default 1 comment('消息状态 -1:删除 0:撤回 1:正常')"`
	CreateAt   string  `json:"create_at,omitempty" xorm:"created not null TIMESTAMP"`
	UpdateAt   *string `json:"update_at,omitempty" xorm:"updated TIMESTAMP"`
	ErrCode    *string `json:"err_code,omitempty" xorm:"-"`
	ErrMsg     *string `json:"err_msg,omitempty" xorm:"-"`
}

const MSG_TYPE_IS_TEXT byte = 1
const MSG_TYPE_IS_IMAGE byte = 2
const MSG_TYPE_IS_AUDIO byte = 3
const MSG_TYPE_IS_VIDEO byte = 4

type Multimedia struct {
	Id   int64  `json:"id"`
	Url  string `json:"url"`
	Mime string `json:"mime"`
	Size int64  `json:"size"`
}

func (m Messages) TableName() string {
	return "messages"
}

func (u *Messages) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Messages) Unmarshal(b []byte) error {
	return json.Unmarshal(b, u)
}
