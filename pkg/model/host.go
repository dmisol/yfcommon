package model

type Host struct {
	Id      string                `json:"id"` // email
	Tid     string                `json:"tid,omitempty"`
	Devices map[int64]interface{} `json:"devices"`
}
