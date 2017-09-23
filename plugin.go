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
	Files map[string]string

	Config struct {
		Key        string
		Identifier string
	}

	Plugin struct {
		Config Config
		Files  Files
	}
)

func (c Config) ToURL() string {
	return fmt.Sprintf("https://api.crowdin.com/api/project/%s/update-file?key=%s", c.Identifier, c.Key)
}

func (p Plugin) Exec() error {
	p.Files = map[string]string{"locale_en-US.ini": "DOCS.md"}
	if len(p.Files) > 20 {
		return fmt.Errorf("20 files max are allowed to upload. %d files given", len(p.Files))
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for crowdin_path, path := range p.Files {
		var err error
		var file *os.File
		if file, err = os.Open(path); err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fmt.Sprintf("files[%s]", crowdin_path), crowdin_path)
		if err != nil {
			return err
		}
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
		if err = writer.Close(); err != nil {
			return err
		}
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
		var err_response = new(responses.Error)
		decoder := xml.NewDecoder(body)
		decoder.CharsetReader = charset.NewReaderLabel
		if err := decoder.Decode(&err_response); err != nil {
			return err
		}
		return err_response
	} else {
		var success = new(responses.Success)
		decoder := xml.NewDecoder(body)
		decoder.CharsetReader = charset.NewReaderLabel
		if err := decoder.Decode(&success); err != nil {
			return err
		}
		for _, file := range success.Stats {
			fmt.Printf("%s: %s\n", file.Name, file.Status)
		}
	}
	return nil
}
