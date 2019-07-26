package model

type UsersChatrooms struct {
	AppsId uint   `json:"apps_id" xorm:"not null unique(apps_uid_roomid) comment('app id') INT"`
	Uid    uint64 `json:"uid" xorm:"BIGINT notnull unique(apps_uid_roomid) comment('归属者用户id')"`
	RoomId uint64 `json:"room_id" xorm:"BIGINT notnull unique(apps_uid_roomid) comment('聊天室id')"`
}

func (uc UsersChatrooms) TableName() string {
	return "users_chatrooms"
}
