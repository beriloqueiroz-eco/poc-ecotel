package internal

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tradersclub/poc-ecotel/pkg/ecotel"
)

type HelloHandler struct {
	IsEnd        int
	Delay        int
	ServiceName  string
	ServiceUrlTo string
}

func (h *HelloHandler) Handle(c *gin.Context) {
	n := rand.Intn(3000-1000+1) + 1000
	time.Sleep(time.Millisecond * time.Duration(n))
	ecotel.Info(c.Request.Context(), "Handling request for service")
	if h.IsEnd == 1 {
		c.String(http.StatusOK, "Response from Service, is the end")
		return
	}
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", h.ServiceUrlTo+"/hello", nil)
	ecotel.Info(c.Request.Context(), "Created request for service")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error calling Service: %v", err)
		return
	}
	ecotel.Error(c.Request.Context(), "Received response from service")

	defer resp.Body.Close()

	time.Sleep(time.Second * time.Duration(h.Delay))

	bodyBytes, _ := io.ReadAll(resp.Body)
	c.String(resp.StatusCode, "Response from Service %s: %s", h.ServiceUrlTo, string(bodyBytes))

}
