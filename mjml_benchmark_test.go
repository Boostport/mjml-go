package mjml

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

type testTemplate struct {
	file   string
	input  string
	output string
}

func getTestTemplates() ([]testTemplate, error) {

	var result []testTemplate

	files := []string{
		"black-friday",
		"one-page",
		"reactivation-email",
		"real-estate",
		"recast",
		"receipt-email",
	}

	for _, file := range files {
		input, err := ioutil.ReadFile("testdata/" + file + ".mjml")

		if err != nil {
			return result, fmt.Errorf("error reading input %s.mjml: %w", file, err)
		}

		output, err := ioutil.ReadFile("testdata/" + file + ".html")

		if err != nil {
			return result, fmt.Errorf("error reading output %s.mjml: %w", file, err)
		}

		result = append(result, testTemplate{
			file:   file,
			input:  string(input),
			output: string(output),
		})
	}

	return result, nil
}

func BenchmarkNodeJS(b *testing.B) {

	serverBin := "node-test-server/server"

	cmd := exec.Command(serverBin)

	err := cmd.Start()

	if err != nil {
		b.Fatalf("Error starting Node.js reference server: %s", err)
	}

	defer cmd.Process.Kill()

	nodeAddr := "http://localhost:8888"

	client := &http.Client{}

	var tries int

	for {
		_, err := client.Get(nodeAddr)

		if err == nil {
			break
		}

		if tries > 5 {
			b.Fatalf("Timed out while waiting for Node.js reference server to be available: %s", err)
		}

		time.Sleep(1 * time.Second)

		tries++
	}

	testCases, err := getTestTemplates()

	if err != nil {
		b.Fatalf("Error getting test cases: %s", err)
	}

	type request struct {
		MJML    string            `json:"mjml"`
		Options map[string]string `json:"options"`
	}

	for _, testCase := range testCases {

		b.Run(testCase.file, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := request{
					MJML: testCase.input,
					Options: map[string]string{
						"validationLevel": "skip",
					},
				}
				jsonRequest, err := json.Marshal(r)

				if err != nil {
					b.Fatalf("Error marshaling json request: %s", err)
				}

				result, err := client.Post(nodeAddr, "application/json", bytes.NewReader(jsonRequest))

				if err != nil {
					b.Fatalf("Error posting request to Node.js server: %s", err)
				}

				var decoded jsonResult

				body, err := ioutil.ReadAll(result.Body)

				if err != nil {
					b.Fatalf("Error reading response body: %s", err)
				}

				err = json.Unmarshal(body, &decoded)

				if err != nil {
					b.Fatalf("Error decoding response; %s", err)
				}

				if decoded.Error != nil {
					b.Fatalf("Error converting input to HTML: %s", err)
				}

				if decoded.HTML != testCase.output {
					b.Fatalf("Output for input (%s.mjml) does not match expected response", testCase.file)
				}
			}
		})
	}
}

func BenchmarkMJMLGo(b *testing.B) {
	testCases, err := getTestTemplates()

	if err != nil {
		b.Fatalf("Error getting test cases: %s", err)
	}

	if err != nil {
		b.Fatalf("Error setting max workers: %s", err)
	}

	for _, testCase := range testCases {
		b.Run(testCase.file, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				result, err := ToHTML(context.Background(), testCase.input, WithValidationLevel(Skip))

				if err != nil {
					b.Fatalf("Error converting input to HTML: %s", err)
				}

				if result != testCase.output {
					b.Fatalf("Output for input (%s.mjml) does not match expected response", testCase.file)
				}
			}
		})
	}
}
