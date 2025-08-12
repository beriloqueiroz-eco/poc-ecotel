package internal

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	if h.IsEnd == 1 {
		c.String(http.StatusOK, "Response from Service %s, is the end", h.ServiceName)
		return
	}
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", h.ServiceUrlTo+"/hello", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error calling Service: %v", err)
		return
	}
	defer resp.Body.Close()

	time.Sleep(time.Second * time.Duration(h.Delay))

	bodyBytes, _ := io.ReadAll(resp.Body)
	c.String(resp.StatusCode, "Response from Service %s: %s", h.ServiceUrlTo, string(bodyBytes))

}
