package querier

import (
	"sync"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (q *querier) Query() {
	var wg sync.WaitGroup
	wg.Add(len(q.states))

	for index := 0; index < len(q.states); index++ {
		go func(index int) {
			defer wg.Done()

			_ = map[models.Region]models.Status{
				models.Teh1:       models.Unknown,
				models.Teh2:       models.Unknown,
				models.SnappGroup: models.Unknown,
			}

			// TODO: do operation over regions
			// TODO: handle timeout and set to unknown

		}(index)
	}

	wg.Wait()
}
