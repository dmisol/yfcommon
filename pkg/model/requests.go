package model

// called tokenized
type TokenReq struct {
	Since    int64  `json:"since"` // GMT
	Until    int64  `json:"until"` // GMT
	Timezone string `json:"tz,omitempty"`

	Addr string `json:"addr,omitempty"`

	DeviceId string            `json:"device"`
	Devices  map[string]string `json:"devices"` // DeviceId-s by name, like ["gate"]"xxxxxxx",["door"]"yyyyyyyy"
}

type TokenResp struct {
	Web string `json:"web,omitempty"` // link for guests' browser
	Tg  string `json:"tg,omitempty"`  // invitation to tg
}
