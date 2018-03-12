package main

import (
	"testing"
	"os/exec"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"github.com/68696c6c/goat"
	"github.com/stretchr/testify/require"
)

func TestCLIConfig(t *testing.T) {
	cmd := exec.Command("./run.sh")
	err := cmd.Start()
	requireNilError(t, err, "test script failed: ")
	defer cmd.Process.Kill()

	// Give the server time to start
	time.Sleep(2 * time.Second)

	response, err := http.Get("http://127.0.0.1:8080/ping")
	requireNilError(t, err, "request failed: ")
	defer response.Body.Close()
	assert.Equal(t, http.StatusOK, response.StatusCode, "unexpected status code")

	b, err := ioutil.ReadAll(response.Body)
	require.Nil(t, err, "failed to read body")

	var r goat.Response
	err = json.Unmarshal(b, &r)
	require.Nil(t, err, "failed to parse body")
	assert.Equal(t, "pong", r.Message, "unexpected response")
}

func requireNilError(t *testing.T, err error, msg string) {
	if err != nil {
		msg += err.Error()
	}
	require.Nil(t, err, msg)
}
