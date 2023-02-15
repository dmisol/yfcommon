package model

type TgLog struct {
	Tid    int64    `json:"tid"`
	Tokens []string `json:"tokens,omitempty"`
}
