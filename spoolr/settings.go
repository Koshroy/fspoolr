package spoolr

import (
	"encoding/json"
	"io/ioutil"
)

type settings struct {
	Dirs []string `dirs`
	AppRoot string `app_root`
	Port int `port`
}

func NewSettings(filePath string) (*settings, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	s := new(settings)
	err = json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}