package model

type Attachments struct {
	Id       uint64  `json:"id" xorm:"notnull pk autoincr BIGINT"`
	AppsId   uint    `json:"apps_id" xorm:"INT notnull comment('app id') index"`
	Url      string  `json:"url" xorm:"TEXT notnull comment('资源路径')"`
	Mime     string  `json:"mime" xorm:"VARCHAR(32) notnull comment('mime类型')"`
	Size     int     `json:"size" xorm:"INT notnull comment('文件大小 size/1024/1024 = mb')"`
	Sha1 	 string  `json:"sha1" xorm:"VARCHAR(64) notnull index comment('文件sha1值')"`
	Status   byte    `json:"status" xorm:"TINYINT notnull default 1 comment('状态 -1:删除 0:禁用 1:正常')"`
	CreateAt string  `json:"create_at" xorm:"created TIMESTAMP notnull"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (c Attachments) TableName() string {
	return "attachments"
}
