package clusterauthorizer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/openshift/insights-operator/pkg/config"
)

type Configurator interface {
	Config() *config.Controller
}

type Authorizer struct {
	configurator Configurator
}

func New(configurator Configurator) *Authorizer {
	return &Authorizer{
		configurator: configurator,
	}
}

func (a *Authorizer) Authorize(req *http.Request) error {
	cfg := a.configurator.Config()
	if len(cfg.Username) > 0 || len(cfg.Password) > 0 {
		req.SetBasicAuth(cfg.Username, cfg.Password)
		return nil
	}
	if len(cfg.Token) > 0 {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		token := strings.TrimSpace(cfg.Token)
		if strings.Contains(token, "\n") || strings.Contains(token, "\r") {
			return fmt.Errorf("cluster authorization token is not valid: contains newlines")
		}
		if len(token) == 0 {
			return fmt.Errorf("cluster authorization token is empty")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
	return nil
}
