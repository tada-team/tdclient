package tdclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
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
	cookie   string
	features *tdproto.Features
}

func NewSession(server string) (Session, error) {
	s := Session{
		Timeout: 10 * time.Second,
	}

	s.logger = log.New(os.Stdout, "tdclient: ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

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
		if _, err := s.doGet("/features.json", &s.features); err != nil {
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
		s.logger.SetOutput(ioutil.Discard)
	}
}

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

	if !tdproto.ValidUid(teamUid) {
		return tdproto.Contact{}, errors.New("invalid team uid")
	}

	b, err := s.doGet("/api/v4/teams/"+teamUid, resp)
	if err != nil {
		return tdproto.Contact{}, err
	}

	if err := JSON.Unmarshal(b, resp); err != nil {
		return tdproto.Contact{}, errors.Wrap(err, "unmarshall fail")
	}

	if !resp.Ok {
		return tdproto.Contact{}, errors.New(resp.Error)
	}

	return resp.Result.Me, nil
}

func (s Session) Contacts(teamUid string) ([]tdproto.Contact, error) {
	resp := new(struct {
		apiResp
		Result []tdproto.Contact `json:"result"`
	})

	if !tdproto.ValidUid(teamUid) {
		return resp.Result, errors.New("invalid team uid")
	}

	b, err := s.doGet("/api/v4/teams/"+teamUid+"/contacts/", resp)
	if err != nil {
		return resp.Result, err
	}

	if err := JSON.Unmarshal(b, resp); err != nil {
		return resp.Result, errors.Wrap(err, "unmarshall fail")
	}

	if !resp.Ok {
		return resp.Result, errors.New(resp.Error)
	}

	return resp.Result, nil
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

	path = s.url(path)
	s.logger.Println("GET", path)
	req, err := http.NewRequest("GET", path, nil)
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

	if err := JSON.Unmarshal(respData, &v); err != nil {
		return respData, errors.Wrapf(err, "unmarshal fail on: %s", string(respData))
	}

	return respData, nil
}
