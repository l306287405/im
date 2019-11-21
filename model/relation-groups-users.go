package model

type RelationGroupsUsers struct {
	ChatroomsUsers `xorm:"extends"`
	Name	string	`json:"room_name"`
}

func (RelationGroupsUsers) TableName() string{
	return "chatrooms_users"
}