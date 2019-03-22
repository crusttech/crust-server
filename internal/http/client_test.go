// +build integration

package http

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/test"
)

func toError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if body == nil || err != nil {
		return errors.Errorf("unexpected response (%d, %s)", resp.StatusCode, err)
	}
	return errors.New(string(body))
}

func TestHTTPClient(t *testing.T) {
	client, err := New(&config.HTTPClient{
		Timeout: 5,
	})
	test.Assert(t, err == nil, "%+v", err)
	client.Debug(FULL)

	req, err := client.Get("https://api.scene-si.org/fortune.php")
	test.Assert(t, err == nil, "%+v", err)

	resp, err := client.Do(req)
	test.Assert(t, err == nil, "%+v", err)

	err = func() error {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			return nil
		default:
			return toError(resp)
		}
	}()

	test.Assert(t, err == nil, "Invalid response: %+v", err)
}
