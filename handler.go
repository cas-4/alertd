package alertd

import (
	"net/http"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// Response payload returned by the GraphQL request below
type Response struct {
	Data struct {
		NewPosition struct {
			ID int `json:"id"`
		} `json:"newPosition"`
	} `json:"data"`
}

// Request payload by the user
type AlertInput struct {
	Latitude       float64 `json:"latitude" binding:"required"`
	Longitude      float64 `json:"longitude" binding:"required"`
	MovingActivity string  `json:"movingActivity" binding:"required"`
	Login          string  `json:"login" binding:"required"`
}

// Create new alert using the backend GraphQL url
func NewAlert(c *gin.Context) {
	var input AlertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cfg, _ := GetConfig()

	// GraphQL mutation payload
	payload := map[string]interface{}{
		"query": `mutation NewPosition($input: PositionInput!) { 
        newPosition(input: $input) { 
            id
        } 
    }`,
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"latitude":       input.Latitude,
				"longitude":      input.Longitude,
				"movingActivity": input.MovingActivity,
			},
		},
	}

	response, err := MakeHttpRequest(cfg.String("backend.url"), payload, input.Login)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Prepare and send a new HTTP request
func MakeHttpRequest(endpoint string, body map[string]interface{}, auth string) (*Response, error) {
	jsonBody, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint, bodyReader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read response body: %s", err.Error()))
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP request returned a status %d and response `%s`", res.StatusCode, resBody))
	}

	var responseBody Response
	if err := json.Unmarshal(resBody, &responseBody); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal response body: %s", err))
	}

	return &responseBody, nil

}
