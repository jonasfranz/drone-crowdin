package responses

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var errorData = `<?xml version="1.0" encoding="ISO-8859-1"?>
<error>
  <code>3</code>
  <message>API key is not valid</message>
</error>
`
var invalidErrorData = `<?xml version="1.0" encoding="ISO-fdsfsdf-1"?>
<error>
  <code>3</code>
  <message>API key is not valid</message>
</error>
`

var successData = `<?xml version="1.0" encoding="ISO-8859-1"?>
<success>
  <stats>
     <file status="skipped" name="demo.ini"></file>
  </stats>
</success>`
var invalidSuccessData = `<?xml version="1.0" encoding="ISO-sefsgdfb-1"?>
<success>
  <stats>
     <file status="skipped" name="demo.ini"></file>
  </stats>
</success>`

func TestParseAsError(t *testing.T) {
	result, err := ParseAsError(bytes.NewBufferString(errorData))
	assert.NoError(t, err)
	assert.Equal(t, 3, result.Code, "error code")
	assert.Equal(t, "API key is not valid", result.Message, "error message")

	_, err = ParseAsSuccess(bytes.NewBufferString(invalidErrorData))
	assert.Error(t, err)
}

func TestError_Error(t *testing.T) {
	result, err := ParseAsError(bytes.NewBufferString(errorData))
	assert.NoError(t, err)
	assert.Error(t, result)
	assert.Equal(t, result.Error(), "Error from crowdin: API key is not valid (error code 3)", "error message")
}

func TestParseAsSuccess(t *testing.T) {
	result, err := ParseAsSuccess(bytes.NewBufferString(successData))
	assert.NoError(t, err)
	assert.Len(t, result.Stats, 1, "files")
	assert.Equal(t, result.Stats[0].Status, "skipped", "status of first file")
	assert.Equal(t, result.Stats[0].Name, "demo.ini", "name of first file")

	_, err = ParseAsError(bytes.NewBufferString(invalidSuccessData))
	assert.Error(t, err)
}
