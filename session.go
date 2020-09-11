package tdclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tada-team/tdproto"

	"github.com/pkg/errors"
)

type apiResp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type Session struct {
	Timeout  time.Duration
	logger   *log.Logger
	server   url.URL
	token    string
	features *tdproto.Features
}

func NewSession(server string, verbose bool) (Session, error) {
	s := Session{
		Timeout: 10 * time.Second,
	}

	s.logger = log.New(os.Stdout, "tdclient: ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	if !verbose {
		s.logger.SetOutput(ioutil.Discard)
	}

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
		if _, err := s.doGet("/features.json", s.features); err != nil {
			return s.features, err
		}
	}
	return s.features, nil
}

func (s *Session) SetToken(token string) { s.token = token }

func (s Session) Ping() error {
	resp := new(struct {
		apiResp
		Result string `json:"result"`
	})
	_, err := s.doGet("/api/v4/ping", resp)
	return err
}

func (s Session) Me(teamUid string) (tdproto.Contact, error) {
	resp := new(struct {
		apiResp
		Result tdproto.Team `json:"result"`
	})

	b, err := s.doGet("/api/v4/teams/"+teamUid, resp)
	if err != nil {
		return tdproto.Contact{}, err
	}

	if err := json.Unmarshal(b, resp); err != nil {
		return tdproto.Contact{}, errors.Wrap(err, "unmarshall fail")
	}

	if !resp.Ok {
		return tdproto.Contact{}, errors.New(resp.Error)
	}

	return resp.Result.Me, nil
}

func (s Session) httpClient() *http.Client {
	return &http.Client{
		Timeout: s.Timeout,
		//Transport: &http.Transport{
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//},
	}
}

func (s Session) url(path string) string {
	s.server.Path = path
	return s.server.String()
}

func (s Session) doGet(path string, v interface{}) ([]byte, error) {
	client := s.httpClient()

	req, err := http.NewRequest("GET", s.url(path), nil)
	if err != nil {
		return []byte{}, errors.Wrap(err, "new request fail")
	}

	if s.token != "" {
		req.Header.Set("token", s.token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, errors.Wrap(err, "client do fail")
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, errors.Wrap(err, "read body fail")
	}

	if resp.StatusCode != 200 {
		return respData, errors.Wrapf(err, "status code: %d %s", resp.StatusCode, string(respData))
	}

	if err := json.Unmarshal(respData, &v); err != nil {
		return respData, errors.Wrapf(err, "unmarshal fail on: %s", string(respData))
	}

	return respData, nil
}
