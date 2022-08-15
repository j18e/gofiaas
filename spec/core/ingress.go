package core

import (
	"encoding/json"
	"fmt"
)

type IngressHost struct {
	Host        string            `json:"host"`
	Paths       []IngressPath     `json:"paths"`
	Annotations map[string]string `json:"annotations"`
}

type IngressPath struct {
	Path string      `json:"path"`
	Port IngressPort `json:"-"`
}

func (i *IngressPath) MarshalJSON() ([]byte, error) {
	if i.Port.Name != "" {
		return json.Marshal(struct {
			Path string `json:"path"`
			Port string `json:"port"`
		}{i.Path, i.Port.Name})
	}
	return json.Marshal(struct {
		Path string `json:"path"`
		Port int    `json:"port"`
	}{i.Path, i.Port.Number})
}

func (i *IngressPath) UnmarshalJSON(bs []byte) error {
	var data struct {
		Path string      `json:"path"`
		Port interface{} `json:"port"`
	}
	if err := json.Unmarshal(bs, &data); err != nil {
		return err
	}
	if data.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	i.Path = data.Path
	switch data.Port.(type) {
	case string:
		res := data.Port.(string)
		if res == "" {
			return fmt.Errorf("parsing port: cannot be empty string")
		}
		i.Port.Name = res
	case int:
		res := data.Port.(int)
		if res == 0 {
			return fmt.Errorf("parsing port: cannot be 0")
		}
		i.Port.Number = res
	case float64:
		res := data.Port.(float64)
		if res == 0 {
			return fmt.Errorf("parsing port: cannot be 0")
		}
		i.Port.Number = int(res)
	default:
		return fmt.Errorf("parsing port: expected string or int, got %T", data.Port)
	}
	return nil
}

type IngressPort struct {
	Name   string
	Number int
}
