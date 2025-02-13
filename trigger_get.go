package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type triggerStatus struct {
	Result []struct {
		TriggerID string `json:"triggerid"`
		Status    string `json:"status"`
	} `json:"result"`
}

type triggerRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  triggerParams `json:"params"`
	Auth    string        `json:"auth"`
	Id      int           `json:"id"`
}

type triggerParams struct {
	Output []string          `json:"output"`
	Filter map[string]string `json:"filter"`
}

func getTriggerStatus(apiUrl, token, objectId string) (string, error) {
	requestBody := triggerRequest{
		Jsonrpc: "2.0",
		Method:  "trigger.get",
		Params: triggerParams{
			Output: []string{"triggerid", "status"},
			Filter: map[string]string{"triggerid": objectId},
		},
		Auth: token,
		Id:   1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(apiUrl, "application/json-rpc", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var zabbixResponse triggerStatus
	if err := json.NewDecoder(resp.Body).Decode(&zabbixResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	if len(zabbixResponse.Result) == 0 {
		return "", fmt.Errorf("no results found for trigger ID: %s", objectId)
	}

	return zabbixResponse.Result[0].Status, nil
}
