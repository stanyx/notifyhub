package intistelecom

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	log "github.com/taudelta/nanolog"
)

const apiUrl = "https://go.intistele.com/external"

type Client struct {
	Client *http.Client
	Login  string
	Key    string
}

func (cl *Client) GetSignature(params map[string]string, timestamp string) string {

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var tokens []string
	for _, k := range keys {
		tokens = append(tokens, fmt.Sprintf("%v", params[k]))
	}

	tokens = append(tokens, timestamp)
	tokens = append(tokens, cl.Key)

	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(tokens, ""))))
}

func (cl *Client) GetTimestamp() (string, error) {

	resp, err := cl.Client.Get(apiUrl + "/get/timestamp.php")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	timestamp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(timestamp), nil
}

type ErrorResult struct {
	Error int `json:"error"`
}

type SmsSendResult struct {
	Error    string  `json:"error"`
	IdSms    string  `json:"id_sms"`
	Cost     float64 `json:"cost"`
	CountSms int     `json:"count_sms"`
}

func (cl *Client) SendSms(phone, message string) (*SmsSendResult, error) {

	timestamp, err := cl.GetTimestamp()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", apiUrl+"/get/send.php", nil)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"login":  cl.Login,
		"phone":  phone,
		"text":   message,
		"sender": "INFO",
	}

	signature := cl.GetSignature(params, timestamp)

	query := request.URL.Query()
	query.Add("timestamp", timestamp)
	query.Add("signature", signature)
	for k, v := range params {
		query.Add(k, v)
	}

	request.URL.RawQuery = query.Encode()

	resp, err := cl.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug().Println("sms send response: ", string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("sms send response error: %s", string(body))
	}

	var testError ErrorResult
	err = json.Unmarshal(body, &testError)
	if err != nil {
		return nil, err
	}

	if testError.Error != 0 {
		return nil, fmt.Errorf("sms send response error: %d", testError.Error)
	}

	normalizedPhone := strings.Replace(phone, "+", "", -1)

	resultData := make(map[string]SmsSendResult)
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		return nil, err
	}

	result, ok := resultData[normalizedPhone]
	if !ok {
		return nil, fmt.Errorf("sms send response error: %s", string(body))
	}

	if result.Error != "" && result.Error != "0" {
		return nil, fmt.Errorf("sms send response error: %s", result.Error)
	}

	return &result, nil
}
