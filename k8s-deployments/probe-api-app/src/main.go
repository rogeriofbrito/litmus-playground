package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/probe", func(c echo.Context) error {
		// TODO: validate params
		duration := c.QueryParam("duration")
		envs := c.QueryParam("envs")
		image := c.QueryParam("image")

		durationParsed, err := time.ParseDuration(duration)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		marginDuration, _ := time.ParseDuration("30s")

		uuid := uuid.New()
		podUniqueId := strings.Split(uuid.String(), "-")[0]
		podName := fmt.Sprintf("k6-probe-%s", podUniqueId)

		kubectlEnvs := getKubectlEnvs(envs)
		kubectlEnvs = addKubectlEnv(kubectlEnvs, "DURATION", duration)

		// TODO: set resource limits in pod
		runCommand := fmt.Sprintf("kubectl run %s -n probe-api-app %s --image-pull-policy=IfNotPresent --image=%s --restart=Never --rm -it", podName, kubectlEnvs, image)
		deleteCommand := fmt.Sprintf("kubectl delete pod %s -n probe-api-app --force=true", podName)
		runCmd := exec.Command("/bin/sh", "-c", runCommand)
		deleteCmd := exec.Command("/bin/sh", "-c", deleteCommand)

		timeout := time.After(time.Duration(durationParsed + marginDuration))
		success := make(chan bool)
		fail := make(chan bool)

		go func() {
			fmt.Printf("launching pod %s\n", podName)
			if err := runCmd.Run(); err != nil {
				fail <- true
			}
			success <- true
		}()

		select {
		case <-success:
			fmt.Println("success")
			return c.NoContent(http.StatusOK)
		case <-fail:
			fmt.Println("fail")
			return c.NoContent(http.StatusInternalServerError)
		case <-timeout:
			fmt.Println("timeout")
			_ = deleteCmd.Run()
			return c.NoContent(http.StatusRequestTimeout)
		}
	})

	if err := e.Start(os.Getenv("PORT")); err != nil {
		panic(err)
	}
}

func getKubectlEnvs(envs string) string {
	kubectlEnvs := ""

	envsSplited := strings.Split(envs, ",")
	for _, env := range envsSplited {
		envSplited := strings.Split(env, "=")
		envName := envSplited[0]
		envValue := envSplited[1]
		kubectlEnvs = kubectlEnvs + fmt.Sprintf("--env=%s=%s ", envName, envValue)
	}

	return strings.TrimSpace(kubectlEnvs)
}

func addKubectlEnv(kubectlEnvs string, envName string, envValue string) string {
	return strings.TrimSpace(kubectlEnvs) + fmt.Sprintf(" --env=%s=%s", envName, envValue)
}
