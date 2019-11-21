package model

type Classes struct {
	Id       uint64  `json:"id" xorm:"notnull pk autoincr BIGINT"`
	AppsId   uint    `json:"apps_id" xorm:"not null comment('app id') INT"`
	Uid      uint64  `json:"uid" xorm:"BIGINT notnull comment('归属用户id')"`
	Name     string  `json:"name" xorm:"VARCHAR(16) notnull comment('分组分类名称')"`
	CreateAt string  `json:"create_at" xorm:"created not null TIMESTAMP"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (c Classes) TableName() string {
	return "classes"
}
