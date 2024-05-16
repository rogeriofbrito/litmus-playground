package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/probe", func(c echo.Context) error {
		envs := c.QueryParam("envs")
		image := c.QueryParam("image")

		uuid := uuid.New()
		podUniqueId := strings.Split(uuid.String(), "-")[0]
		podName := fmt.Sprintf("k6-probe-%s", podUniqueId)

		kubectlEnvs := getKubectlEnvs(envs)

		command := fmt.Sprintf("kubectl run %s -n probe-api-app %s --image-pull-policy=IfNotPresent --image=%s --restart=Never --rm -it", podName, kubectlEnvs, image)

		cmd := exec.Command("/bin/sh", "-c", command)

		if err := cmd.Run(); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
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
