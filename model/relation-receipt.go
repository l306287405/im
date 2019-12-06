package model

import (
	"encoding/json"
)

const(
	RECEIPTS_STATUS_IS_RECEIVED=iota
	RECEIPTS_STATUS_IS_
)

type Receipts struct {
	Id  uint64 `json:"id"`
	Mid  uint64 `json:"mid"`
	Type byte   `json:"type"`	//0是单聊 1是群聊
	Status byte `json:"status"`
}

func (r Receipts) Marshal() []byte{
	b,_:=json.Marshal(r)
	return b
}
func (r Receipts) UnMarshal(v []byte) {
	json.Unmarshal(v,r)
}