package model

type UsersUsers struct {
	AppsId uint   `json:"apps_id" xorm:"not null unique(apps_uid_fid) comment('app id') INT"`
	Uid    uint64 `json:"uid" xorm:"BIGINT notnull unique(apps_uid_fid) comment('归属者用户id')"`
	Fid    uint64 `json:"fid" xorm:"BIGINT notnull unique(apps_uid_fid) comment('关联用户id')"`
	Cid    uint64 `json:"cid" xorm:"BIGINT notnull comment('好友分类id -1:黑名单 0:未分组 other:分类id')"`
}

func (uu UsersUsers) TableName() string {
	return "users_users"
}
