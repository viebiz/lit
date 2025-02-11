package monitoring

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestLogger_Infof(t *testing.T) {
	// Given
	expResult := map[string]interface{}{
		"level":        "info",
		"msg":          "WAAAGH!",
		"project_name": "lightning",
		"version":      "0.0.1",
		"group":        "bizgroup",
		"env":          "dev",
	}

	buf := bytes.NewBuffer(nil)
	logger := NewLoggerWithWriter(buf, WithFields(
		Field("project_name", "lightning"),
		Field("version", "0.0.1"),
		Field("group", "bizgroup"),
		Field("env", "dev"),
	))

	// When
	logger.Infof("WAAAGH!")

	// Then
	var actual map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &actual))
	if diff := cmp.Diff(expResult, actual, cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
		return key == "ts"
	})); diff != "" {
		t.Errorf("unexpected result (-want, +got) = %v", diff)
	}
}

func TestLogger_Errorf(t *testing.T) {
	// Given
	expErr := errors.New("simulated error")
	expResult := map[string]interface{}{
		"level":        "error",
		"msg":          "WAAAGH!",
		"project_name": "lightning",
		"version":      "0.0.1",
		"group":        "bizgroup",
		"env":          "dev",
		"error":        "simulated error",
	}

	buf := bytes.NewBuffer(nil)
	logger := NewLoggerWithWriter(buf, WithFields(
		Field("project_name", "lightning"),
		Field("version", "0.0.1"),
		Field("group", "bizgroup"),
		Field("env", "dev"),
	))

	// When
	logger.Errorf(expErr, "WAAAGH!")

	// Then
	var actual map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &actual))
	if diff := cmp.Diff(expResult, actual, cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
		return key == "ts"
	})); diff != "" {
		t.Errorf("unexpected result (-want, +got) = %v", diff)
	}
}
