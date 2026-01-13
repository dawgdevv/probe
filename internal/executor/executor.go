package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dawgdevv/apitestercli/internal/assert"
	"github.com/dawgdevv/apitestercli/internal/config"
	"github.com/dawgdevv/apitestercli/pkg/models"
)

type Result struct {
	Name       string
	Passed     bool
	StatusCode int
	Error      error
	Duration   time.Duration
}

func RunTest(baseURL string, env map[string]string, test models.TestCase) Result {
	start := time.Now()

	resolvePath, err := config.SubstituteString(test.Request.Path, env)

	if err != nil {
		return Result{Name: test.Name, Passed: false, Error: err}
	}
	url := baseURL + resolvePath

	var body *bytes.Reader

	if test.Request.Body != nil {
		b, err := json.Marshal(test.Request.Body)
		if err != nil {
			return Result{Name: test.Name, Passed: false, Error: err}
		}
		body = bytes.NewReader(b)

	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(test.Request.Method, url, body)

	if err != nil {
		return Result{Name: test.Name, Passed: false, Error: err}
	}

	for k, v := range test.Request.Headers {
		val, err := config.SubstituteString(v, env)

		if err != nil {
			return Result{Name: test.Name, Passed: false, Error: err}
		}
		req.Header.Set(k, val)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Name: test.Name, Passed: false, Error: err}
	}

	defer resp.Body.Close()

	if resp.StatusCode != test.Expect.Status {
		return Result{
			Name:       test.Name,
			Passed:     false,
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("expected %d, got %d", test.Expect.Status, resp.StatusCode),
		}
	}

	if len(test.Expect.JSON) > 0 {
		if err := assert.AssertJSON(bodyBytes, test.Expect.JSON); err != nil {
			return Result{
				Name:       test.Name,
				Passed:     false,
				StatusCode: resp.StatusCode,
				Error:      err,
			}
		}
	}

	return Result{
		Name:       test.Name,
		Passed:     true,
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}
