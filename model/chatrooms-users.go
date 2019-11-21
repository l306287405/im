package model

const ROOM_USER_ROLE_IS_OWNER byte = 0
const ROOM_USER_ROLE_IS_ADMIN byte = 1
const ROOM_USER_ROLE_IS_MEMBER byte = 2

type ChatroomsUsers struct {
	AppsId uint   `json:"apps_id" xorm:"not null unique(apps_roomid_uid) comment('app id') INT"`
	RoomId uint64 `json:"room_id" xorm:"BIGINT notnull unique(apps_roomid_uid) comment('聊天室id')"`
	Uid    uint64 `json:"uid" xorm:"BIGINT notnull unique(apps_roomid_uid) comment('用户id')"`
	Role   byte	  `json:"role" xorm:"TINYINT notnull default 2 comment('角色 0:拥有者 1:管理员 2:成员')"`
	Status int8   `json:"status" xorm:"TINYINT notnull default 1 comment('状态 -1:软删 0:待审核 1:正常')"`
	JoinedAt *string `json:"joined_at" xorm:"TIMESTAMP"`
	CreateAt string  `json:"create_at" xorm:"created TIMESTAMP not null"`
}

func (uc ChatroomsUsers) TableName() string {
	return "chatrooms_users"
}