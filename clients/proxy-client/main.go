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

func checkProxy(ctx context.Context, proxyURL, targetURL string) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL(proxyURL)),
		},
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %s\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing HTTP request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		successMetric.Inc()
	} else {
		log.Printf("Proxy check failed. Status code: %d\n", resp.StatusCode)
	}
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

	server := &http.Server{
		Addr: ":9090",
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		log.Println("Received signal, shutting down gracefully...")
		cancel()
		server.Shutdown(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Ticker routine shutting down...")
				return
			case <-ticker.C:
				checkProxy(ctx, proxyURL, targetURL)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(server.ListenAndServe())
	}()

	log.Println("Waiting for the HTTP server to finish...")
	wg.Wait()
	log.Println("All goroutines shut down. Exiting.")
}
