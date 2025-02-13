package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Host struct {
	HostId string `json:"hostid"`
	Host   string `json:"host"`
	Name   string `json:"name"`
}

type hostRequest struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  hostParams `json:"params"`
	Auth    string     `json:"auth"`
	Id      int        `json:"id"`
}

type hostParams struct {
	Output []string          `json:"output"`
	Filter map[string]string `json:"filter"`
}

type hostResponse struct {
	Result []Host `json:"result"`
}

func getHosts(apiUrl, token string) ([]Host, error) {
	requestBody := hostRequest{
		Jsonrpc: "2.0",
		Method:  "host.get",
		Params: hostParams{
			Output: []string{"hostid", "host", "name", "status"},
			Filter: map[string]string{"host": "", "status": "0", "maintenance_status": "0"},
		},
		Auth: token,
		Id:   1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var zabbixResponse hostResponse
	if err := json.NewDecoder(resp.Body).Decode(&zabbixResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return zabbixResponse.Result, nil
}
