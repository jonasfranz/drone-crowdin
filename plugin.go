package main

import (
	"bytes"
	"fmt"
	"github.com/JonasFranzDEV/drone-crowdin/responses"
	"github.com/JonasFranzDEV/drone-crowdin/utils"
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
		Config          Config
		Files           Files
		Branch          string
		Languages       []string
		ExportDirectory string
		DoDownload      bool
	}
)

func (c Config) ToProjectURL() string {
	return fmt.Sprintf("https://api.crowdin.com/api/project/%s", c.Identifier)
}

// ToUploadURL returns the API-endpoint including identifier and API-KEY
func (c Config) ToUploadURL() string {
	return fmt.Sprintf("%s/update-file?key=%s", c.ToProjectURL(), c.Key)
}

func (p Plugin) ToLanguageDownloadURL(language string) string {
	if p.Branch != "" {
		return fmt.Sprintf("%s/download/%s.zip?key=%s&branch=%s", p.Config.ToProjectURL(), language, p.Config.Key, p.Branch)
	}
	return fmt.Sprintf("%s/download/%s.zip?key=%s", p.Config.ToProjectURL(), language, p.Config.Key)

}

// Exec starts the plugin and updates the crowdin translation by uploading files from the files map
func (p Plugin) Exec() error {
	client := &http.Client{}

	//SECTION: Upload
	if len(p.Files) > 20 {
		return fmt.Errorf("20 files max are allowed to upload. %d files given", len(p.Files))
	} else if len(p.Files) > 0 {
		req, err := p.buildUploadRequest()
		if err != nil {
			return fmt.Errorf("error while building upload request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body := &bytes.Buffer{}
		if _, err := body.ReadFrom(resp.Body); err != nil {
			return err
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			if e, err := responses.ParseAsError(body); err != nil {
				return err
			} else {
				return fmt.Errorf("error while uploading: %v", e)
			}
		}
		var success = new(responses.Success)
		if success, err = responses.ParseAsSuccess(body); err != nil {
			return err
		}
		for _, file := range success.Stats {
			fmt.Printf("%s: %s\n", file.Name, file.Status)
		}
	}

	//SECTION: Download
	if p.DoDownload {
		if err := p.buildTranslations(client); err != nil {
			return fmt.Errorf("error while building languages: %v", err)
		}
		for _, language := range p.Languages {
			if err := p.downloadLanguage(client, language); err != nil {
				return fmt.Errorf("error while downloading %s: %v", language, err)
			}
			fmt.Printf("Downloaded package: %s\n", language)
		}
	}
	return nil
}

func (p Plugin) buildUploadRequest() (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for crowdinPath, path := range p.Files {
		var err error
		var file *os.File
		if file, err = os.Open(path); err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fmt.Sprintf("files[%s]", crowdinPath), crowdinPath)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}
	// Adding branch if it is not ignored
	if p.Branch != "" {
		writer.WriteField("branch", p.Branch)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	var req *http.Request
	var err error
	if req, err = http.NewRequest("POST", p.Config.ToUploadURL(), body); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func (p Plugin) buildTranslations(client *http.Client) error {
	// Step 1: Export translations (aka generate server side)
	exportURL := fmt.Sprintf("%s/export?key=%s", p.Config.ToProjectURL(), p.Config.Key)
	if p.Branch != "" {
		exportURL = fmt.Sprintf("%s&branch=%s", exportURL, p.Branch)
	}
	fmt.Println(exportURL)
	if resp, err := client.Get(exportURL); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		defer resp.Body.Close()
		if e, err := responses.ParseAsError(resp.Body); err != nil {
			return err
		} else {
			return e
		}
	} else {
		defer resp.Body.Close()
	}
	return nil
}

func (p Plugin) downloadLanguage(client *http.Client, language string) error {
	fmt.Println(p.ToLanguageDownloadURL(language))
	file, err := downloadFromUrl(p.ToLanguageDownloadURL(language))
	if err != nil {
		return err
	}
	err = utils.Unzip(file.Name(), p.ExportDirectory)
	if err != nil {
		return err
	}
	err = os.Remove(file.Name())
	if err != nil {
		return err
	}
	return nil
}

func downloadFromUrl(url string) (*os.File, error) {
	output, err := os.Create("lang.zip")
	if err != nil {
		return nil, err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		os.Remove(output.Name())
		return nil, err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	return output, err
}
