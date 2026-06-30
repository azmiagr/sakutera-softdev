package mlclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/azmiagr/sakutera-softdev/model"
)

type Interface interface {
	Forecast(req model.MLForecastRequest) (*model.MLForecastResponse, error)
}

type mlClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func Init() Interface {
	timeoutSec, _ := strconv.Atoi(os.Getenv("ML_REQUEST_TIMEOUT_SECONDS"))
	if timeoutSec == 0 {
		timeoutSec = 15
	}
	return &mlClient{
		baseURL:    os.Getenv("ML_SERVICE_URL"),
		token:      os.Getenv("ML_SERVICE_TOKEN"),
		httpClient: &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (c *mlClient) Forecast(req model.MLForecastRequest) (*model.MLForecastResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/forecast", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Service-Token", c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	var result model.MLForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
