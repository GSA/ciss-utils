package sm

import "time"

type Secret struct {
	ID       string    `json:"ID,omitempty"`
	Name     string    `json:"name,omitempty"`
	Type     string    `json:"type,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
	Username string    `json:"username,omitempty"`
	Value    string    `json:"value,omitempty"`
}
