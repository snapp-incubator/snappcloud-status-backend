package querier

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
	"go.uber.org/zap"
)

func (q *querier) Query() {
	var wg sync.WaitGroup
	wg.Add(len(q.states))

	now := time.Now().UnixMilli()
	seconds := float64(now) / 1000
	timestamp := fmt.Sprintf("%.3f", seconds)
	fmt.Println(timestamp)

	for index := 0; index < len(q.states); index++ {
		go func(index int) {
			defer wg.Done()

			//
			for _, region := range []models.Region{models.Teh1, models.Teh2, models.SnappGroup} {
				go func(region models.Region) {
					// operator will do the query operation
					operator := func(query string, badRes models.Status) bool {
						result := q.queryThanos(region, query, timestamp)
						if result != successfulResult {
							var status models.Status
							if result == badResult {
								status = badRes
							} else {
								status = models.Unknown
							}

							q.mutex.Lock()
							q.states[index].status[region] = status
							q.mutex.Unlock()
							return false
						}
						return true
					}

					// first check outage query and then check disruption query
					if !operator(q.states[index].config.Queries.Outage, models.Outage) {
						return
					} else if !operator(q.states[index].config.Queries.Disruption, models.Disruption) {
						return
					}

					q.mutex.Lock()
					q.states[index].status[region] = models.Operational
					q.mutex.Unlock()
				}(region)
			}
		}(index)
	}

	wg.Wait()
}

type result uint8

const (
	errorResult      result = 0 // invalid http request
	timeoutResult    result = 1 // timeout
	badResult        result = 2 // checks equals to false
	successfulResult result = 3 // ok
)

func (q *querier) queryThanos(region models.Region, query string, timestamp string) result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChannel := make(chan result)
	var baseURL string

	switch region {
	case models.Teh1:
		baseURL = q.config.ThanosFrontends.Teh1
	case models.Teh2:
		baseURL = q.config.ThanosFrontends.Teh2
	case models.SnappGroup:
		baseURL = q.config.ThanosFrontends.SnappGroup
	}

	go func() {
		requestURL, _ := url.Parse(baseURL)
		requestURL.RawQuery = url.Values{
			"query":            []string{query},
			"dedup":            []string{"true"},
			"partial_response": []string{"false"},
			"time":             []string{timestamp},
		}.Encode()

		request, _ := http.NewRequest("GET", requestURL.String(), nil)
		request = request.WithContext(ctx)

		client := http.Client{Timeout: q.config.RequestTimeout}
		response, err := client.Do(request)
		if err != nil {
			resultChannel <- errorResult
			return
		} else if response.StatusCode/100 != 2 {
			resultChannel <- errorResult
			return
		}
		defer response.Body.Close()

		rawBody, _ := io.ReadAll(response.Body)
		body := &struct {
			Status string `json:"status"`
			Data   struct {
				ResultType string `json:"resultType"`
				Result     []struct {
					Metric map[string]string `json:"metric"`
					Value  []any             `json:"value"`
				} `json:"result"`
			} `json:"data"`
		}{}

		if err = json.Unmarshal(rawBody, body); err != nil {
			q.loggger.Error("Error unmarshaling query result", zap.Error(err))
			return
		}

		for index := 0; index < len(body.Data.Result); index++ {
			if len(body.Data.Result[index].Value) != 2 {
				resultChannel <- badResult
				q.loggger.Error("Invalid result value length")
				return
			}

			if value, ok := body.Data.Result[index].Value[1].(string); ok {
				if value == "1" {
					continue
				}

				resultChannel <- badResult
				return
			}

			resultChannel <- badResult
			q.loggger.Error("Invalid result value type, it should be 0 or 1")
			return
		}

		resultChannel <- successfulResult
	}()

	select {
	case result := <-resultChannel:
		return result

	// on timeout
	case <-time.After(q.config.RequestTimeout):
		return timeoutResult
	}
}
