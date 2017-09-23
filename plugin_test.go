package main

import "testing"

func TestConfig_ToURL(t *testing.T) {
	exampleConfig := &Config{Identifier: "test", Key: "MYKEY"}
	result := exampleConfig.ToUploadURL()
	if result != "https://api.crowdin.com/api/project/test/update-file?key=MYKEY" {
		t.Fatalf("ToUploadURL returns \"%s\" instead of the expected \"%s\"", result, "https://api.crowdin.com/api/project/test/update-file?key=MYKEY")
	}
}
