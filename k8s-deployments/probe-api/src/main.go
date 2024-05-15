package main

import (
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/probe/k6/smoke-test-order-api", func(c echo.Context) error {
		cmd := exec.Command("/bin/bash", "-c", "kubectl run k6-probe --image-pull-policy=IfNotPresent --image=localhost:5000/k6-probe:latest --restart=Never --rm -it")

		if err := cmd.Run(); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	})

	if err := e.Start("127.0.0.1:9000"); err != nil {
		panic(err)
	}
}
