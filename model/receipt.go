package model

import (
	"encoding/json"
)

type Receipts struct {
	Mid  uint64 `json:"mid"`
	Type byte   `json:"type"`	//0是单聊 1是群聊
}

func (r Receipts) Marshal() []byte{
	b,_:=json.Marshal(r)
	return b
}
func (r Receipts) UnMarshal(v []byte) {
	json.Unmarshal(v,r)
}