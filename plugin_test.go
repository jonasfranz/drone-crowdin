package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupExamplePlugin() *Plugin {
	return &Plugin{
		ExportDirectory: "demo/",
		Languages:       []string{"all"},
		DoDownload:      true,
		Files:           map[string]string{"locale_en-US.ini": "LICENSE"},
		Config: Config{
			Key:        "MYKEY",
			Identifier: "test",
		},
		Branch: "master",
	}
}

func TestConfig_ToUploadURL(t *testing.T) {
	exampleConfig := setupExamplePlugin().Config
	result := exampleConfig.ToUploadURL()
	assert.Equal(t, "https://api.crowdin.com/api/project/test/update-file?key=MYKEY", result, "ToUploadURL")
}

func TestPlugin_ToLanguageDownloadURL(t *testing.T) {
	examplePlugin := setupExamplePlugin()
	result := examplePlugin.ToLanguageDownloadURL(examplePlugin.Languages[0])
	assert.Equal(t, "https://api.crowdin.com/api/project/test/download/all.zip?key=MYKEY&branch=master", result, "ToLanguageDownloadURL")

	examplePlugin.Branch = ""
	result = examplePlugin.ToLanguageDownloadURL(examplePlugin.Languages[0])
	assert.Equal(t, "https://api.crowdin.com/api/project/test/download/all.zip?key=MYKEY", result, "ToLanguageDownloadURL")
}

func TestConfig_ToProjectURL(t *testing.T) {
	exampleConfig := setupExamplePlugin().Config
	result := exampleConfig.ToProjectURL()
	assert.Equal(t, "https://api.crowdin.com/api/project/test", result, "ToProjectURL")
}
