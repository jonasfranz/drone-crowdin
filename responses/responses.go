package responses

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
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

// ParseAsError parses XML to Error
func ParseAsError(body io.Reader) (*Error, error) {
	var errResponse = new(Error)
	decoder := xml.NewDecoder(body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&errResponse); err != nil {
		return nil, err
	}
	return errResponse, nil
}

// Success is a crowdin success message
type Success struct {
	XMLName xml.Name `xml:"success"`
	Stats   []File   `xml:"stats>file"`
}

// ParseAsSuccess parses XML to Success
func ParseAsSuccess(body io.Reader) (*Success, error) {
	var success = new(Success)
	decoder := xml.NewDecoder(body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&success); err != nil {
		return nil, err
	}
	return success, nil
}

// File represents the status of an uploaded file
type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Status  string   `xml:"status,attr"`
}
