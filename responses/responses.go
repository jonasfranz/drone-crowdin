package responses

import (
	"encoding/xml"
	"fmt"
)

// Error is a crowdin error message
type Error struct {
	XMLName xml.Name `xml:"error"`
	Code    int      `xml:"code"`
	Message string   `xml:"message"`
}

// Error implements the error interface to handle like an error
func (e *Error) Error() string {
	return fmt.Sprintf("Error from crowdin: %s (error code %d)", e.Message, e.Code)
}

// Success is a crowdin success message
type Success struct {
	XMLName xml.Name `xml:"success"`
	Stats   []File   `xml:"stats>file"`
}

// File represents the status of an uploaded file
type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Status  string   `xml:"status,attr"`
}
