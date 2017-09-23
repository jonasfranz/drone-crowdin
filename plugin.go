package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/JonasFranzDEV/drone-crowdin/responses"
	"golang.org/x/net/html/charset"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type (
	// Files is a mapping between the crowdin path and the real file path
	Files map[string]string

	// Config stores the credentials for the crowdin API
	Config struct {
		Key        string
		Identifier string
	}

	// Plugin represents the drone-crowdin plugin including config and file-mapping.
	Plugin struct {
		Config Config
		Files  Files
		Branch string
	}
)

// ToURL returns the API-endpoint including identifier and API-KEY
func (c Config) ToURL() string {
	return fmt.Sprintf("https://api.crowdin.com/api/project/%s/update-file?key=%s", c.Identifier, c.Key)
}

// Exec starts the plugin and updates the crowdin translation by uploading files from the files map
func (p Plugin) Exec() error {
	if len(p.Files) > 20 {
		return fmt.Errorf("20 files max are allowed to upload. %d files given", len(p.Files))
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for crowdinPath, path := range p.Files {
		var err error
		var file *os.File
		if file, err = os.Open(path); err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fmt.Sprintf("files[%s]", crowdinPath), crowdinPath)
		if err != nil {
			return err
		}
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}
	// Adding branch if it is not ignored
	if p.Branch != "" {
		writer.WriteField("branch", p.Branch)
	}
	if err := writer.Close(); err != nil {
		return err
	}
	var req *http.Request
	var err error
	if req, err = http.NewRequest("POST", p.Config.ToURL(), body); err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body = &bytes.Buffer{}
	if _, err := body.ReadFrom(resp.Body); err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var errResponse = new(responses.Error)
		decoder := xml.NewDecoder(body)
		decoder.CharsetReader = charset.NewReaderLabel
		if err := decoder.Decode(&errResponse); err != nil {
			return err
		}
		return errResponse
	}
	var success = new(responses.Success)
	decoder := xml.NewDecoder(body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&success); err != nil {
		return err
	}
	for _, file := range success.Stats {
		fmt.Printf("%s: %s\n", file.Name, file.Status)
	}
	return nil
}
