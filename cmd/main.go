// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aryasoni98/wooak/internal/adapters/data/ssh_config_file"
	"github.com/aryasoni98/wooak/internal/logger"

	"github.com/aryasoni98/wooak/internal/adapters/ui"
	aiDomain "github.com/aryasoni98/wooak/internal/core/domain/ai"
	securityDomain "github.com/aryasoni98/wooak/internal/core/domain/security"
	"github.com/aryasoni98/wooak/internal/core/services"
	aiService "github.com/aryasoni98/wooak/internal/core/services/ai"
	"github.com/aryasoni98/wooak/internal/core/services/monitoring"
	securityService "github.com/aryasoni98/wooak/internal/core/services/security"
	"github.com/spf13/cobra"
)

var (
	version   = "develop"
	gitCommit = "unknown"
)

func main() {
	log, err := logger.New("WOOAK")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//nolint:errcheck // log.Sync may return an error which is safe to ignore here
	defer log.Sync()

	// Initialize monitoring service
	monitoringService := monitoring.NewMonitoringService()
	monitoringService.Start()
	defer monitoringService.Stop()

	// Start HTTP server for metrics endpoint
	go func() {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; version=0.0.4")
			fmt.Fprint(w, monitoringService.GetMetrics().ToPrometheusFormat())
		})

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			healthSummary := monitoringService.GetHealthMonitor().GetHealthSummary(ctx)
			w.Header().Set("Content-Type", "application/json")
			
			// Simple health status response
			if healthSummary["overall_status"] == monitoring.Healthy {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			
			fmt.Fprintf(w, `{"status":"%v","timestamp":"%v"}`, 
				healthSummary["overall_status"], 
				healthSummary["last_checked"])
		})

		log.Infow("Starting metrics server", "port", 9090)
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Errorw("Metrics server failed", "error", err)
		}
	}()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorw("failed to get user home directory", "error", err)
		//nolint:gocritic // exitAfterDefer: ensure immediate exit on unrecoverable error
		os.Exit(1)
	}
	sshConfigFile := filepath.Join(home, ".ssh", "config")
	metaDataFile := filepath.Join(home, ".wooak", "metadata.json")

	serverRepo := ssh_config_file.NewRepository(log, sshConfigFile, metaDataFile)
	
	// Set monitoring for repository
	if repo, ok := serverRepo.(*ssh_config_file.Repository); ok {
		repo.SetMonitoring(monitoringService)
	}
	
	serverService := services.NewServerService(log, serverRepo)

	// Initialize security service
	securityPolicy := securityDomain.DefaultSecurityPolicy()
	securitySvc := securityService.NewSecurityService(securityPolicy)

	// Initialize AI service
	aiConfig := aiDomain.DefaultAIConfig()
	aiSvc := aiService.NewAIService(aiConfig)
	
	// Set monitoring for AI service
	aiSvc.SetMonitoring(monitoringService)

	tui := ui.NewTUI(log, serverService, securitySvc, aiSvc, version, gitCommit)

	rootCmd := &cobra.Command{
		Use:   ui.AppName,
		Short: "Wooak SSH server picker TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run()
		},
	}
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
