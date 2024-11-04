package model

import "github.com/golang-jwt/jwt/v5"

type GuestToken struct {
	jwt.RegisteredClaims

	Tz        string            `json:"tz"`
	Addr      string            `json:"addr"`
	DevId     string            `json:"devid,omitempty"`
	NotifyTid int64             `json:"notify,omitempty"`
	Devices   map[string]string `json:"devices,omitempty"` // by name, like "gate"
}
