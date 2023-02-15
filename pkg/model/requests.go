package model

// called with tokenized SessionData
type TokenReq struct {
	Since    int64  `json:"since"` // GMT
	Until    int64  `json:"until"` // GMT
	Timezone string `json:"tz,omitempty"`
	Device   int64  `json:"device"`
	Addr     string `json:"addr,omitempty"`
}

type TokenResp struct {
	Web string `json:"web,omitempty"` // link for guests' browser
	Tg  string `json:"tg,omitempty"`  // invitation to tg
}
