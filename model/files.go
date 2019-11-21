package model

type Files struct {
	Id       uint64  `json:"id" xorm:"not null pk autoincr BIGINT"`
	AppsId   uint    `json:"apps_id" xorm:"INT not null comment('app id')"`
	Mime     string  `json:"mime" xorm:"VARCHAR(11) notnull comment('文件类型')"`
	Size     int     `json:"size" xorm:"INT comment('文件尺寸')"`
	Sha1     string  `json:"sha1" xorm:"BINARY(20) notnull comment('文件sha1值') index"`
	Status   byte    `json:"status" xorm:"TINYINT notnull default 1 comment('文件状态 0:停用 1:正常')"`
	CreateAt string  `json:"create_at" xorm:"created not null TIMESTAMP"`
	UpdateAt *string `json:"update_at" xorm:"updated TIMESTAMP"`
}

func (f Files) TableName() string {
	return "files"
}
