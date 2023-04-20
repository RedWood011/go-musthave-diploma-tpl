package tests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func createReqBody(t *testing.T, raw interface{}) *bytes.Buffer {
	body, err := json.Marshal(raw)

	require.NoError(t, err)

	return bytes.NewBuffer(body)
}
