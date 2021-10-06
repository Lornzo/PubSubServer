package connection

import "encoding/json"

type ClientMsg struct {
	Type int64       `json:"msgtype"`
	Data interface{} `json:"msgdata"`
}

func (thisObj *ClientMsg) DataToJson() (jsonStr string) {
	var (
		jsonByte []byte
		err      error
	)
	if jsonByte, err = json.Marshal(thisObj.Data); err != nil {
		return
	}

	jsonStr = string(jsonByte)
	return
}
