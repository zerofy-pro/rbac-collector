package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zerofy.pro/rbac-collector/src/client/k8s"
	"zerofy.pro/rbac-collector/src/collector"
	"zerofy.pro/rbac-collector/src/config"
	"zerofy.pro/rbac-collector/src/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	clientset, err := k8s.NewInClusterClient()
	if err != nil {
		log.Warn().Err(err).Msg("Could not create in-cluster Kubernetes client, falling back to kubeconfig")
		clientset, err = k8s.NewKubeConfigClient()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Kubernetes client from kubeconfig")
		}
		log.Info().Msg("Successfully created K8s client using kubeconfig")
	} else {
		log.Info().Msg("Successfully created in-cluster K8s client")
	}

	dataCollector := collector.New(clientset, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(cfg.CollectionInterval)
		defer ticker.Stop()

		log.Info().Str("interval", cfg.CollectionInterval.String()).Msg("Starting collector")
		for {
			select {
			case <-ticker.C:
				log.Info().Msg("Running collection cycle")
				if err := dataCollector.CollectAndLog(context.Background()); err != nil {
					log.Error().Err(err).Msg("Error during collection cycle")
				}
			case <-ctx.Done():
				log.Info().Msg("Collector loop stopping")
				return
			}
		}
	}()

	// Wait for a shutdown signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Info().Msg("Shutdown signal received, stopping collector...")
	cancel()

	time.Sleep(2 * time.Second)
	log.Info().Msg("Worker shut down")
}
