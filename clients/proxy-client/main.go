package main

import (
	"fmt"
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
	proxyURL      = os.Getenv("PROXY_URL")
	targetURL     = "https://ifconfig.me"
	successMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_check_success",
		Help: "Indicates if the proxy check was successful (1) or not (0)",
	})
	shutdownCh = make(chan struct{})
	wg         sync.WaitGroup
)

func init() {
	prometheus.MustRegister(successMetric)
}

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			fmt.Printf("Error starting Prometheus HTTP server: %s\n", err)
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go proxyCheckRoutine()

	select {
	case <-signalCh:
		close(shutdownCh)
		wg.Wait()
		fmt.Println("Shutting down gracefully.")
		os.Exit(0)
	}
}

func proxyCheckRoutine() {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			err := checkProxy()
			if err != nil {
				fmt.Printf("Proxy check failed: %s\n", err)
				successMetric.Set(0)
			} else {
				successMetric.Set(1)
			}

		case <-shutdownCh:
			return
		}
	}
}

func checkProxy() error {
	if proxyURL == "" {
		log.Fatal("PROXY_URL environment variable not set")
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(MustParseURL(proxyURL)),
		},
	}

	req, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
