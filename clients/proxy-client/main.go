package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	successMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_check_success",
		Help: "Indicates the number of successful proxy checks",
	})
)

func init() {
	prometheus.MustRegister(successMetric)
}

func checkProxy(ctx context.Context, proxyURL, targetURL string) error {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL(proxyURL)),
		},
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		successMetric.Inc()
	}

	return nil
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func main() {
	proxyURL := os.Getenv("PROXY_URL")
	if proxyURL == "" {
		log.Fatal("PROXY_URL environment variable not set")
	}
	targetURL := "https://ifconfig.me"

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			log.Printf("Error starting Prometheus HTTP server: %s\n", err)
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		log.Println("Received signal, shutting down gracefully...")
		cancel()
	}()

	defer wg.Done()

	defer func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines shut down. Exiting.")
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	tickerFunc := func() {
		if err := checkProxy(ctx, proxyURL, targetURL); err != nil {
			log.Printf("Proxy check failed: %s\n", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		case <-ticker.C:
			tickerFunc()
		}
	}
}
