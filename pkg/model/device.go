package model

type Device struct {
	DeviceId string `json:"device"`
	Addr     string `json:"addr,omitempty"`

	Private *interface{} `json:"private,omitempty"`
}
