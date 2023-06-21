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

func (s Server) GetOrigin() string { // origin = protocolo (http, https, ect..) + hostname (nombre de maquina o IP) + puerto. Ver https://developer.mozilla.org/en-US/docs/Glossary/Origin
	const protocol = "http"
	host := s.Hostname + ":" + s.Port // ver host vs hostname => https://stackoverflow.com/a/13673410/903998. host = hostname + puerto
	return protocol + "://" + host    // más sobre partes de una URL => https://blog.hubspot.com/marketing/parts-url
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
