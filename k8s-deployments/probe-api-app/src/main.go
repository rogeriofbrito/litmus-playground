package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/probe", func(c echo.Context) error {
		timeoutMsParam := c.QueryParam("timeoutMs")
		envs := c.QueryParam("envs")
		image := c.QueryParam("image")

		timeoutMs, err := strconv.ParseInt(timeoutMsParam, 10, 64)
		if err != nil {
			fmt.Println("invalid timeout param")
			return c.NoContent(http.StatusInternalServerError)
		}

		uuid := uuid.New()
		podUniqueId := strings.Split(uuid.String(), "-")[0]
		podName := fmt.Sprintf("k6-probe-%s", podUniqueId)

		kubectlEnvs := getKubectlEnvs(envs)

		// TODO: set resource limits in pod
		runCommand := fmt.Sprintf("kubectl run %s -n probe-api-app %s --image-pull-policy=IfNotPresent --image=%s --restart=Never --rm -it", podName, kubectlEnvs, image)
		deleteCommand := fmt.Sprintf("kubectl delete pod %s -n probe-api-app --force=true", podName)

		runCmd := exec.Command("/bin/sh", "-c", runCommand)
		deleteCmd := exec.Command("/bin/sh", "-c", deleteCommand)

		timeout := time.After(time.Duration(timeoutMs) * time.Millisecond)
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
