package model

type User struct {
	Id      string             `json:"id"` // email
	Blocked string             `json:"blocked,omitempty"`
	Tid     string             `json:"tid,omitempty"`
	Name    string             `json:"name,omitempty"`
	Devices map[string]*Device `json:"devices"` // by DeviceId
}
