package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ChatWithOllama(url string, model string, prompt string, messages interface{}) (string, error) {
	reqData := make(map[string]interface{})
	reqData["prompt"] = prompt
	reqData["model"] = model
	reqData["messages"] = messages

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Response status code is " + strconv.Itoa(resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	respMsg := ""

	respJsonDataList := bytes.Split(respBody, []byte("\n"))
	for _, respJsonData := range respJsonDataList {
		respJson := make(map[string]interface{})
		err := json.Unmarshal(respJsonData, &respJson)
		if err != nil {
			return "", err
		}

		if respJson["done"].(bool) {
			break
		}

		respMsg += respJson["message"].(map[string]interface{})["content"].(string)
	}

	// DeepSeek model will always output this, but we do not need...
	if len(respMsg) > len("<think>") && respMsg[:len("<think>")] == "<think>" {
		splitRes := strings.Split(respMsg, "</think>\n")
		if len(splitRes) > 1 {
			respMsg = splitRes[1]
		}
	}

	return respMsg, nil
}

func ChatWithOllamaText(url string, model string, prompt string, message string) (string, error) {
	messages := []map[string]string{
		{
			"role":    "system",
			"content": prompt,
		},
		{
			"role":    "user",
			"content": message,
		},
	}

	return ChatWithOllama(url, model, prompt, messages)
}
