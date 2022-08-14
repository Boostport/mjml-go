package mjml

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"
)

func TestToHTML(t *testing.T) {

	files := []string{
		"black-friday",
		"one-page",
		"reactivation-email",
		"real-estate",
		"recast",
		"receipt-email",
	}

	for _, file := range files {
		contents, err := os.ReadFile("testdata/" + file + ".mjml")

		if err != nil {
			t.Fatalf("Error reading file: %s", file)
		}

		result, err := ToHTML(context.Background(), string(contents), WithValidationLevel(Skip)) // Skip validation because the templates are not fully mjml4-compliant

		if err != nil {
			t.Errorf("Error converting mjml file (%s) to html: %s", file, err)
		}

		expected, err := os.ReadFile("testdata/" + file + ".html")

		if err != nil {
			t.Fatalf("Error reading expected html: %s", err)
		}

		if result != string(expected) {
			t.Errorf("Compiled HTML for %s does not match expected html", file)
		}
	}
}

func TestConcurrency(t *testing.T) {

	files := []string{
		"black-friday",
		"one-page",
		"reactivation-email",
		"real-estate",
		"recast",
		"receipt-email",
	}

	type testCase struct {
		file     string
		input    string
		expected string
	}
	var testCases []testCase

	for _, file := range files {
		input, err := os.ReadFile("testdata/" + file + ".mjml")

		if err != nil {
			t.Fatalf("Error reading input test data (%s.mjml): %s", file, err)
		}

		expected, err := os.ReadFile("testdata/" + file + ".html")

		if err != nil {
			t.Fatalf("Error reading expected test data (%s.html): %s", file, err)
		}

		testCases = append(testCases, testCase{
			file:     file,
			input:    string(input),
			expected: string(expected),
		})
	}

	numTestCases := big.NewInt(int64(len(testCases)))

	var wg sync.WaitGroup

	numGoRoutines := 200

	errs := make(chan error, numGoRoutines)

	for i := 0; i < numGoRoutines; i++ {

		wg.Add(1)

		go func(run int) {

			defer wg.Done()

			testCaseIndex, err := rand.Int(rand.Reader, numTestCases)

			if err != nil {
				errs <- fmt.Errorf("error selecting test case for run %d: %w", run, err)
				return
			}

			testCase := testCases[testCaseIndex.Int64()]

			result, err := ToHTML(context.Background(), testCase.input, WithValidationLevel(Skip))

			if err != nil {
				errs <- fmt.Errorf("error converting input to HTML for run %d using %s.mjml as input: %w", run, testCase.file, err)
				return
			}

			if result != testCase.expected {
				errs <- fmt.Errorf("html result does not match expected result for run %d", run)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Errorf("Error running ToHTML concurrently: %s", err.Error())
		}
	}
}

func TestSetMaxWorkers(t *testing.T) {
	files := []string{
		"black-friday",
		"one-page",
		"reactivation-email",
		"real-estate",
		"recast",
		"receipt-email",
	}

	type testCase struct {
		file     string
		input    string
		expected string
	}
	var testCases []testCase

	for _, file := range files {
		input, err := os.ReadFile("testdata/" + file + ".mjml")

		if err != nil {
			t.Fatalf("Error reading input test data (%s.mjml): %s", file, err)
		}

		expected, err := os.ReadFile("testdata/" + file + ".html")

		if err != nil {
			t.Fatalf("Error reading expected test data (%s.html): %s", file, err)
		}

		testCases = append(testCases, testCase{
			file:     file,
			input:    string(input),
			expected: string(expected),
		})
	}

	numTestCases := big.NewInt(int64(len(testCases)))

	var wg sync.WaitGroup

	numGoRoutines := 200

	errs := make(chan error, numGoRoutines)

	for i := 0; i < numGoRoutines; i++ {

		wg.Add(1)

		go func(run int) {

			defer wg.Done()

			testCaseIndex, err := rand.Int(rand.Reader, numTestCases)

			if err != nil {
				errs <- fmt.Errorf("error selecting test case for run %d: %w", run, err)
				return
			}

			testCase := testCases[testCaseIndex.Int64()]

			result, err := ToHTML(context.Background(), testCase.input, WithValidationLevel(Skip))

			if err != nil {
				errs <- fmt.Errorf("error converting input to HTML for run %d using %s.mjml as input: %w", run, testCase.file, err)
				return
			}

			if result != testCase.expected {
				errs <- fmt.Errorf("html result does not match expected result for run %d", run)
				return
			}
		}(i)
	}

	// SetMaxWorkers randomly
	for i := 0; i < 10; i++ {

		wg.Add(1)

		delay, err := rand.Int(rand.Reader, big.NewInt(int64(50)))

		if err != nil {
			t.Fatalf("Error generating random delay: %s", err)
		}

		duration := time.Duration(delay.Int64()) * time.Millisecond

		time.Sleep(duration)

		go func() {

			defer wg.Done()

			randWorkers, err := rand.Int(rand.Reader, big.NewInt(int64(199)))

			if err != nil {
				errs <- fmt.Errorf("error generating random number of workers: %w", err)
				return
			}

			numWorkers := randWorkers.Int64() + 1 // Generated number needs to be between 1 and 200

			err = SetMaxWorkers(int32(numWorkers))

			if err != nil {
				errs <- fmt.Errorf("error setting max workers: %w", err)
				return
			}
		}()

	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Errorf("Error running ToHTML concurrently: %s", err.Error())
		}
	}
}
