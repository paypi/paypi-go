package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	API_URL string = "http://localhost:8080/graphql"
)

type Response struct {
	Data   interface{}
	Errors *[]struct {
		Message    string
		Extensions *struct {
			Message string
			Type    string
		}
	}
}

type GqlQuery struct {
	Query     string
	Variables map[string]interface{}
}

func MakeRequest(query GqlQuery, respData interface{}) error {
	jsonData := map[string]interface{}{
		"query":     query.Query,
		"variables": query.Variables,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 3}

	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	if resp.Errors != nil {
		if len(*resp.Errors) == 1 {
			return errors.New((*resp.Errors)[0].Message)
		}
	}

	j, err := json.Marshal(resp.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &respData)
	return err
}
