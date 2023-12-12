package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/containers/image/docker"
	"github.com/containers/image/types"
	"github.com/joho/godotenv"
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
	// Load environment variables from a file, if present
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Get the image reference from the environment variable
	imageRefStr := os.Getenv("IMAGE_REFERENCE")
	if imageRefStr == "" {
		log.Fatal("IMAGE_REFERENCE environment variable not set")
	}

	imageRef, err := docker.ParseReference(imageRefStr)
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
