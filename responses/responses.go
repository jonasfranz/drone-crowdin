package responses

import (
	"encoding/xml"
	"fmt"
)

type Error struct {
	XMLName xml.Name `xml:"error"`
	Code    int      `xml:"code"`
	Message string   `xml:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error from crowdin: %s (error code %d)", e.Message, e.Code)
}
