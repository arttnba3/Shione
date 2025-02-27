package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ParseOpenRouterReply(rawReplyMsg []byte) (string, error) {
	rawReplyMsg = bytes.TrimLeft(rawReplyMsg, " \t\r\n")
	respJson := make(map[string]interface{})
	err := json.Unmarshal(rawReplyMsg, &respJson)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error parsing reply: %v, original reply data:\n%v", err.Error(), rawReplyMsg))
	}

	respChoice, ok := respJson["choices"].([]interface{})
	if !ok {
		respError, ok := respJson["error"].(map[string]interface{})
		if !ok {
			return "", errors.New(fmt.Sprintf("unexpected response format1, original reply data: %v", rawReplyMsg))
		} else {
			respErrorMsg, ok := respError["message"].(string)
			if !ok {
				return "", errors.New(fmt.Sprintf("unexpected response format2, original reply data: %v", rawReplyMsg))
			} else {
				return "", errors.New(respErrorMsg)
			}

		}
	}

	if len(respChoice) == 0 {
		return "", errors.New("no reply message in response choice array")
	}

	respMessage, ok := respChoice[0].(map[string]interface{})
	if !ok {
		return "", errors.New(fmt.Sprintf("unexpected response format3, original reply data: %v", rawReplyMsg))
	}

	respMessageMsg, ok := respMessage["message"].(map[string]interface{})
	if !ok {
		return "", errors.New(fmt.Sprintf("unexpected response format4, original reply data: %v", rawReplyMsg))
	}

	return respMessageMsg["content"].(string), nil
}

func ParseOLLAMAReply(rawReplyMsg []byte) (string, error) {
	var respMsg string

	// original OLLAMA produce a list
	respJsonDataList := bytes.Split(rawReplyMsg, []byte("\n"))
	for _, respJsonData := range respJsonDataList {
		respJson := make(map[string]interface{})
		err := json.Unmarshal(respJsonData, &respJson)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Error parsing reply: %v, original reply data:\n%v", err.Error(), respJsonData))
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

func ChatWithAI(provider string, url string, model string, prompt string, headers map[string]interface{}, maxWaitingTime time.Duration, messages interface{}) (string, error) {
	reqData := make(map[string]interface{})
	//reqData["prompt"] = prompt
	reqData["model"] = model
	reqData["messages"] = messages

	reqBodyJson, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v.(string))
		}
	}

	client := &http.Client{
		Timeout: maxWaitingTime,
	}
	resp, err := client.Do(req)
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

	var respMsg string

	switch provider {
	case "Ollama":
		respMsg, err = ParseOLLAMAReply(respBody)
		if err != nil {
			return "", err
		}
		break
	case "OpenRouter":
		respMsg, err = ParseOpenRouterReply(respBody)
		if err != nil {
			return "", err
		}
		if len(respMsg) == 0 { // sometimes it occur, just do the request again
			time.Sleep(time.Second * 1)
			return ChatWithAI(provider, url, model, prompt, headers, maxWaitingTime, messages)
		}
		break
	default:
		return "", errors.New("Unknown provider: " + provider)
	}

	return respMsg, nil
}

func ChatWithAIText(provider string, url string, model string, prompt string, headers map[string]interface{}, maxWaitingTime time.Duration, message string) (string, error) {
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

	return ChatWithAI(provider, url, model, prompt, headers, maxWaitingTime, messages)
}
