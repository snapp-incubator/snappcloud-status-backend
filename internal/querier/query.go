package querier

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
	"go.uber.org/zap"
)

func (q *querier) Query() {
	var wg sync.WaitGroup
	wg.Add(len(q.states))

	for index := 0; index < len(q.states); index++ {
		go func(index int) {
			defer wg.Done()

			//
			for _, region := range []models.Region{models.Teh1, models.Teh2, models.SnappGroup} {
				go func(region models.Region) {
					// result := make(chan int, 1)

					// first check outage query
					result := q.queryWithTimeout(region, q.states[index].config.Queries.Outage)
					if !q.checkResult(result) {
						q.mutex.Lock()
						q.states[index].status[region] = models.Outage
						q.mutex.Unlock()
						return
					}

					// then check disruption query
					result = q.queryWithTimeout(region, q.states[index].config.Queries.Disruption)
					if !q.checkResult(result) {
						q.mutex.Lock()
						q.states[index].status[region] = models.Disruption
						q.mutex.Unlock()
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

func (q *querier) queryWithTimeout(region models.Region, query string) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChannel := make(chan int)
	var requestURL string

	switch region {
	case models.Teh1:
		requestURL = q.config.ThanosFrontends.Teh1
	case models.Teh2:
		requestURL = q.config.ThanosFrontends.Teh2
	case models.SnappGroup:
		requestURL = q.config.ThanosFrontends.SnappGroup
	}

	go func() {
		request, _ := http.NewRequest("GET", requestURL, nil)
		// TODO: add required parameters
		request = request.WithContext(ctx)

		client := http.Client{Timeout: q.config.RequestTimeout}
		response, err := client.Do(request)
		if err != nil {
			resultChannel <- 0
			return
		} else if response.StatusCode/100 != 2 {
			resultChannel <- 0
			return
		}

		b, _ := io.ReadAll(response.Body)
		q.loggger.Info("response", zap.ByteString("", b))
		resultChannel <- 1
	}()

	select {
	case result := <-resultChannel:
		return result

	// on timeout
	case <-time.After(q.config.RequestTimeout):
		return 0
	}
}

func (q *querier) checkResult(result int) bool {
	// TODO: check all the result
	// TODO: change result type
	return result == 1
}
