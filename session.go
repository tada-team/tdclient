package tdclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
)

type Session struct {
	Timeout  time.Duration
	logger   *log.Logger
	server   url.URL
	token    string
	cookie   string
	features *tdproto.Features
}

func NewSession(server string) (Session, error) {
	s := Session{
		Timeout: 10 * time.Second,
	}

	s.logger = log.New(os.Stdout, "tdclient: ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	s.SetVerbose(false)

	u, err := url.Parse(server)
	if err != nil {
		return Session{}, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return Session{}, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}
	s.server = *u

	return s, nil
}

func (s *Session) Features() (*tdproto.Features, error) {
	if s.features == nil {
		if err := s.doGet("/features.json", nil, &s.features); err != nil {
			return s.features, err
		}
	}
	return s.features, nil
}

func (s *Session) SetToken(v string) {
	s.token = v
}

func (s *Session) SetCookie(v string) {
	s.cookie = v
}

func (s *Session) SetVerbose(v bool) {
	if v {
		s.logger.SetOutput(os.Stdout)
	} else {
		s.logger.SetOutput(io.Discard)
	}
}

func (s Session) httpClient() *http.Client {
	return &http.Client{
		Timeout: s.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// InsecureSkipVerify: true,
				MinVersion: tls.VersionTLS12,
			},
			ForceAttemptHTTP2: true,
		},
	}
}

func (s Session) doGet(path string, params interface{}, resp interface{}) error {
	return s.doRaw("GET", path, params, nil, resp)
}

func (s Session) doPost(path string, data, v interface{}) error {
	return s.doRaw("POST", path, nil, data, v)
}

func (s Session) doDelete(path string, resp interface{}) error {
	return s.doRaw("DELETE", path, nil, nil, resp)
}

func (s Session) doRaw(method, path string, params, data, v interface{}) error {
	client := s.httpClient()

	var u = s.server
	u.Path = path
	if params != nil {
		q := make(url.Values)
		if err := schema.NewEncoder().Encode(params, q); err != nil {
			return err
		}
		for k := range q {
			v := q.Get(k)
			if v == "" {
				delete(q, k)
			}
		}
		u.RawQuery = q.Encode()
	}
	path = u.String()

	var buf *bytes.Buffer
	if data == nil {
		s.logger.Println(method, path)
		buf = bytes.NewBuffer([]byte{})
	} else {
		s.logger.Println(method, path, debugJSON(data))
		b, err := json.Marshal(data)
		if err != nil {
			return errors.Wrap(err, "json marshal fail")
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return errors.Wrap(err, "new request fail")
	}

	if s.token != "" {
		req.Header.Set("token", s.token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client do fail")
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read body fail")
	}

	if err := JSON.Unmarshal(respData, &v); err != nil {
		return errors.Wrapf(err, "unmarshal fail on: %s", string(respData))
	}

	s.logger.Println(debugJSON(v))

	return nil
}
