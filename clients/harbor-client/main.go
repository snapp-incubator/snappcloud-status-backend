package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/containers/image/docker"
	"github.com/containers/image/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	imagePullCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "image_pulled_total",
		Help: "Total number of images pulled successfully.",
	})
)

func init() {
	prometheus.MustRegister(imagePullCounter)
}

func pullImage(imageRef types.ImageReference) {
	imageSourcePolicyContext := &types.SystemContext{}
	ctx := context.Background()

	imageSource, err := imageRef.NewImageSource(ctx, imageSourcePolicyContext)
	if err != nil {
		log.Println("Error creating image source:", err)
		return
	}
	defer imageSource.Close()

	imagePullCounter.Inc()
	fmt.Println("Image pulled successfully!")
}

func main() {
	imageRef, err := docker.ParseReference("//core.harbor.domain/test-project/busybox:test-tag") // The "//" at the first of the harbor URL is necessary.
	if err != nil {
		log.Fatal("Error parsing image reference:", err)
	}

	// Register the metrics handler
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server in a goroutine
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Initial image pull
	pullImage(imageRef)

	// Schedule image pull every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		pullImage(imageRef)
	}
}
