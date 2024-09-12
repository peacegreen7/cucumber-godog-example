package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

var (
	response *http.Response
	body     []byte
)

type registerPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func iSendAPostRequestToWithPayload(method, url string, payload *godog.DocString) error {
	// Parse the payload from the feature file
	var reqBody registerPayload
	err := json.Unmarshal([]byte(payload.Content), &reqBody)
	if err != nil {
		return err
	}

	// Create request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send POST request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// Read the response body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}

func theResponseCodeShouldBe(expectedCode int) error {
	if response.StatusCode != expectedCode {
		return fmt.Errorf("expected status code %d but got %d", expectedCode, response.StatusCode)
	}
	return nil
}

func theResponseShouldMatchJson(expectedJson *godog.DocString) error {
	var expected, actual map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJson.Content), &expected); err != nil {
		return err
	}
	if err := json.Unmarshal(body, &actual); err != nil {
		return err
	}

	// Use assert to check equality of the JSON response
	assert.Equal(nil, expected, actual)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with payload:$`, iSendAPostRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, theResponseShouldMatchJson)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{"features"}, // Make sure to point to the correct feature path
	}
	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
