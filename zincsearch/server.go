package zincsearch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

type Server struct {
	Hostname    string
	Port        string
	Credentials Credentials
}

var DefaultServer Server = Server{
	Hostname: "localhost",
	Port:     "4080",
	Credentials: Credentials{
		Username: "admin",
		Password: "Complexpass#123",
	},
}

func (s Server) GetOrigin() string {
	const protocol = "http"
	return protocol + "://" + s.Hostname + ":" + s.Port
}

func (s Server) Version() (string, error) {
	request, err := http.NewRequest("GET", s.GetOrigin()+"/version", nil)
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(s.Credentials.Username, s.Credentials.Password)
	bytes, err := performRequestAndReadResponse(request)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s Server) ListIndexes() (string, error) {
	request, err := http.NewRequest("GET", s.GetOrigin()+"/api/index", nil)
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(s.Credentials.Username, s.Credentials.Password)

	bytes, err := performRequestAndReadResponse(request)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func performRequestAndReadResponse(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		errMsg := fmt.Sprintf("Status code is not OK: %d", response.StatusCode)
		return nil, errors.New(errMsg)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
