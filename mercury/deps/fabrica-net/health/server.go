package health

import (
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*khttp.Server
}

func NewServer(addr string) *Server {
	s := khttp.NewServer(khttp.Address(addr))
	s.HandlePrefix("/q/", openapiv2.NewHandler())
	s.Handle("/metrics", promhttp.Handler())
	s.HandleFunc("/health", func(w khttp.ResponseWriter, r *khttp.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return &Server{s}
}
