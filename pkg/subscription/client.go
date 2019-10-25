package subscription

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	trialPayload struct {
		Domain string `json:"domain"`
		Email  string `json:"email"`
	}

	checkPayload struct {
		Domain string `json:"domain"`
		Key    string `json:"key"`
	}

	permintResponse struct {
		Domain        string
		Key           string
		Expires       time.Time
		IsTrial       bool `json:"isTrial"`
		LimitMaxUsers uint `json:"limitMaxUsers"`
	}

	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
)

const (
	keyLength = 64

	// permitBaseURL       = "https://permit.crust.tech/"
	permitBaseURL       = "http://localhost:8000/"
	permitTrialEndpoint = permitBaseURL + "trial"
	permitCheckEndpoint = permitBaseURL + "check"
)

var (
	logger = zap.NewNop()
)

func check(ctx context.Context, domain, key string) (sub *permintResponse, err error) {
	if len(key) == 0 {
		return nil, errors.New("key not set")
	} else if len(key) != keyLength {
		return nil, fmt.Errorf("invalid key length (%d chars)", len(key))
	}

	return send(ctx, permitCheckEndpoint, checkPayload{Domain: domain, Key: key})
}

func trial(ctx context.Context, domain, email string) (sub *permintResponse, err error) {
	return send(ctx, permitTrialEndpoint, trialPayload{Email: email, Domain: domain})
}

func send(ctx context.Context, ep string, payload interface{}) (sub *permintResponse, err error) {
	var (
		rsp *http.Response
		req *http.Request
		log = logger.With(
			zap.String("endpoint", ep),
			zap.Any("payload", payload))
	)

	if req, err = buildRequest(ep, payload); err != nil {
		log.Error("could not build request", zap.Error(err))
		return
	}

	if rsp, err = http.DefaultClient.Do(req); err != nil {
		log.Error("could not send request", zap.Error(err))
		return
	}

	defer rsp.Body.Close()
	return processResponse(rsp)
}

func buildRequest(ep string, payload interface{}) (*http.Request, error) {
	var buf = &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPost, ep, buf)
}

func processResponse(rsp *http.Response) (p *permintResponse, err error) {
	if rsp.StatusCode != http.StatusOK {
		buf, _ := ioutil.ReadAll(rsp.Body)
		logger.Debug("response from subscription server",
			zap.Int("status-code", rsp.StatusCode),
			zap.String("status", rsp.Status),
			zap.String("response", string(buf)),
		)

		switch rsp.StatusCode {
		case http.StatusBadRequest:
			return nil, errors.New("bad request")
		case http.StatusNotFound:
			return nil, errors.New("subscription key not found")
		case http.StatusInternalServerError:
			return nil, errors.New("subscription server error")
		case http.StatusUnauthorized:

			return nil, errors.New("subscription key invalid")
		}
	}

	p = &permintResponse{}
	return p, json.NewDecoder(rsp.Body).Decode(&p)
}
