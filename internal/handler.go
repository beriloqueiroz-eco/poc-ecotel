package internal

import (
	"fmt"
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
	Repo         *Repository
}

type Output struct {
	Cod string `json:"cod"`
	Msg string `json:"msg"`
	Res string `json:"res"`
}

func (h *HelloHandler) Handle(c *gin.Context) {
	n := rand.Intn(3000-1000+1) + 1000
	time.Sleep(time.Millisecond * time.Duration(n))
	ecotel.Info(c.Request.Context(), "Handling request for service")
	output := Output{}
	// Consulta simples no banco para mapear no trace
	if h.Repo != nil {
		if result, err := h.Repo.SimpleQuery(c.Request.Context()); err == nil {
			ecotel.Info(c.Request.Context(), fmt.Sprintf("DB SimpleQuery result: %d", result))
			output.Cod = result
		} else {
			ecotel.Error(c.Request.Context(), fmt.Sprintf("DB SimpleQuery error: %v", err))
		}
	}
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
	output.Msg = string(bodyBytes)
	output.Res = fmt.Sprintf("Response from Service %s: %s", h.ServiceUrlTo, string(bodyBytes))
	c.JSON(http.StatusOK, output)

}
