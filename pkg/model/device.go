package model

type Device struct {
	DeviceId int64  `json:"device"`
	Addr     string `json:"addr,omitempty"`

	Private *interface{} `json:"private,omitempty"`
}
