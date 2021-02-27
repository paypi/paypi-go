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

type GqlClient interface {
	MakeRequest(query GqlQuery, respData interface{}) error
}

type gqlClientImpl struct {
	HttpClient http.Client
	ApiUrl     string
}

func New(apiUrl string) GqlClient {
	return gqlClientImpl{
		ApiUrl:     apiUrl,
		HttpClient: http.Client{Timeout: time.Second * 5},
	}
}

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

func (g gqlClientImpl) MakeRequest(query GqlQuery, respData interface{}) error {
	jsonData := map[string]interface{}{
		"query":     query.Query,
		"variables": query.Variables,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", g.ApiUrl, bytes.NewBuffer(jsonValue))

	response, err := g.HttpClient.Do(request)
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
