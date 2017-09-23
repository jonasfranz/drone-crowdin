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

type Success struct {
	XMLName xml.Name `xml:"success"`
	Stats   []File   `xml:"stats>file"`
}

type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Status  string   `xml:"status,attr"`
}
