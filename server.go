package grafana_json

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Server implements a generic frameworks for the Grafana simpleJson API datasource
type Server struct {
	handler Handler
	port    int
}

// Handler implement the business logic of the Grafana API datasource so that
// Server can be limited to providing the generic search/query framework
type Handler interface {
	Endpoints() Endpoints
}

// Endpoints contains the functions that implements each of the SimpleJson endpoints
type Endpoints struct {
	Search      func() []string
	Query       func(target string, args *TimeSeriesQueryArgs) (*QueryResponse, error)
	TableQuery  func(target string, args *TableQueryArgs) (*TableQueryResponse, error)
	Annotations func(name, query string, args *AnnotationRequestArgs) ([]Annotation, error)
}

// Create creates a Server object for the specified Handler.
// Panics if handler is nil
func Create(handler Handler, port int) *Server {
	if handler == nil {
		panic("nil handler specified")
	}
	return &Server{handler: handler, port: port}
}

// Run the API Server
func (server *Server) Run() error {
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())
	r.HandleFunc("/", server.hello)
	if server.handler.Endpoints().Search != nil {
		r.HandleFunc("/search", server.search).Methods(http.MethodPost)
	}
	if server.handler.Endpoints().Query != nil || server.handler.Endpoints().TableQuery != nil {
		r.HandleFunc("/query", server.query).Methods(http.MethodPost)
	}
	if server.handler.Endpoints().Annotations != nil {
		r.HandleFunc("/annotations", server.annotations).Methods(http.MethodPost, http.MethodOptions)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", server.port), r)
}

// Prometheus metrics
var (
	httpDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "grafana_api_duration_seconds",
		Help: "Grafana API duration of HTTP requests.",
	}, []string{"path"})
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}
