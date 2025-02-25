package terror

import (
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
