package suites

import (
	"fmt"
	"strings"
	"time"

	"github.com/clems4ever/authelia/internal/utils"
	log "github.com/sirupsen/logrus"
)

func waitUntilServiceLogDetected(
	interval time.Duration,
	timeout time.Duration,
	dockerEnvironment *DockerEnvironment,
	service string,
	logPatterns []string) error {
	log.Debug("Waiting for service " + service + " to be ready...")
	err := utils.CheckUntil(5*time.Second, 1*time.Minute, func() (bool, error) {
		logs, err := dockerEnvironment.Logs(service, []string{"--tail", "20"})
		fmt.Printf(".")

		if err != nil {
			return false, err
		}
		for _, pattern := range logPatterns {
			if strings.Contains(logs, pattern) {
				return true, nil
			}
		}
		return false, nil
	})

	fmt.Print("\n")
	return err
}

func waitUntilAutheliaIsReady(dockerEnvironment *DockerEnvironment) error {
	log.Info("Waiting for Authelia to be ready...")

	err := waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"authelia-backend",
		[]string{"Authelia is listening on"})

	if err != nil {
		return err
	}

	err = waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"authelia-frontend",
		[]string{"You can now view web in the browser.", "Compiled with warnings", "Compiled successfully!"})

	if err != nil {
		return err
	}
	log.Info("Authelia is now ready!")

	return nil
}