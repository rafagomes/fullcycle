package usecase

import (
	"net/http"
	"sync"
	"time"

	"loadtester/internal/domain"
)

type LoadTester struct {
	httpClient *http.Client
}

func NewLoadTester() *LoadTester {
	return &LoadTester{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (lt *LoadTester) Run(url string, totalRequests, concurrency int) (domain.Report, error) {
	jobs := make(chan struct{}, totalRequests)
	results := make(chan domain.Result, totalRequests)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go lt.worker(url, jobs, results, &wg)
	}

	for i := 0; i < totalRequests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	close(results)

	report := domain.Report{
		TotalTime:          time.Since(startTime),
		TotalRequests:      totalRequests,
		StatusDistribution: make(map[int]int),
	}

	for res := range results {
		if res.Err != nil {
			report.ErrorCount++
			report.StatusDistribution[0]++
		} else {
			report.StatusDistribution[res.StatusCode]++
			if res.StatusCode == 200 {
				report.Status200Count++
			}
		}
	}

	return report, nil
}

func (lt *LoadTester) worker(url string, jobs <-chan struct{}, results chan<- domain.Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for range jobs {
		resp, err := lt.httpClient.Get(url)
		if err != nil {
			results <- domain.Result{StatusCode: 0, Err: err}
			continue
		}
		results <- domain.Result{StatusCode: resp.StatusCode, Err: nil}
		resp.Body.Close()
	}
}
