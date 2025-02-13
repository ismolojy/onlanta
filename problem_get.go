package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type problemRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  problemParams `json:"params"`
	Auth    string        `json:"auth"`
	Id      int           `json:"id"`
}

type problemParams struct {
	Output  []string          `json:"output"`
	HostIds string            `json:"hostids"`
	Filter  map[string]string `json:"filter"`
}

type problemResponse struct {
	Result []problem `json:"result"`
}

type problem struct {
	Name     string `json:"name"`
	ObjectId string `json:"objectid"`
}

func getProblems(apiUrl, token string, host Host) ([]problem, error) {
	requestBody := problemRequest{
		Jsonrpc: "2.0",
		Method:  "problem.get",
		Params: problemParams{
			Output:  []string{"objectid", "name"},
			HostIds: host.HostId,
			Filter:  map[string]string{"host": "", "status": "0"},
		},
		Auth: token,
		Id:   1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(apiUrl, "application/json-rpc", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var zabbixResponse problemResponse
	if err := json.NewDecoder(resp.Body).Decode(&zabbixResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return zabbixResponse.Result, nil
}
