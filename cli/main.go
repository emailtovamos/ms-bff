package main

import (
	"flag"
	"github.com/Graphmasters/workshop-bff/bff"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	// "time"

	// jaegercfg "github.com/uber/jaeger-client-go/config"
	// jaegerlog "github.com/uber/jaeger-client-go/log"
	// "github.com/uber/jaeger-lib/metrics"
	// "github.com/pkg/errors"
	// "os"
)

func main() {
	grpcAddr := flag.String("address-game-service-grpc", ":8082", "The GRPC server address")
	serverAddr := flag.String("address-http", ":8081", "The HTTP server address")

	// opentracingAgentUrl := flag.String("opentracing-agent-url", "localhost:6831", "UDP host:port of the remote tracing agent to send traces to.")
	// opentracingSamplingType := flag.String("opentracing-sampling-type", "rateLimiting", "trace sampling type to use: const, probabilistic, rateLimiting or remote. See http://opentracing.io/ for detailed documentation.")
	// opentracingSamplingParam := flag.Float64("opentracing-sampling-param", 3.0, "trace sampling parameter for the chosen sampling type. See http://opentracing.io/ for detailed documentation.")

	flag.Parse()

	// serviceName := os.Getenv("SERVICE_NAME")
	// if serviceName == "" {
	// 	log.Fatal().Msg("SERVICE_NAME environment var not set or empty.")
	// }

	// setup opentracing
	// opentracingCloseFunc, err := setupOpentracing(*opentracingAgentUrl, *opentracingSamplingType, *opentracingSamplingParam, serviceName)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to set up opentracing")
	// }
	// defer func() {
	// 	if err := opentracingCloseFunc(); err != nil {
	// 		log.Error().Err(err).Msg("failed to close opentracing")
	// 	}
	// }()

	gameClient, err := bff.NewGrpcGameServiceClient(*grpcAddr)
	gr := bff.NewGameResource(gameClient)

	router := mux.NewRouter()
	router.HandleFunc("/v1/quiplash/addPlayer", gr.AddPlayer).Methods(http.MethodGet)
	router.HandleFunc("/v1/quiplash/start", gr.Start).Methods(http.MethodGet)
	router.HandleFunc("/v1/quiplash/getQuestion", gr.GetQuestion).Methods(http.MethodGet)
	router.HandleFunc("/v1/quiplash/postAnswer", gr.PostAnswer).Methods(http.MethodPost)
	router.HandleFunc("/v1/quiplash/getVotePair", gr.GetVotePair).Methods(http.MethodGet)
	router.HandleFunc("/v1/quiplash/castVote", gr.CastVote).Methods(http.MethodGet)
	router.HandleFunc("/v1/quiplash/voteResult", gr.VoteResult).Methods(http.MethodGet)

	err = http.ListenAndServe(*serverAddr, router)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not init server")
	}

	log.Info().Msgf("Started HTTP-Service at [%s]", *serverAddr)
}

// setupOpentracing sets up the configuration of the opentracing communication with the jaeger agent.
// func setupOpentracing(agentUrl, samplingType string, samplingParam float64, serviceName string) (closeFunc func() error, err error) {
// 	// recommended configuration for production
// 	cfg := jaegercfg.Configuration{
// 		Reporter: &jaegercfg.ReporterConfig{
// 			BufferFlushInterval: time.Second,
// 			LocalAgentHostPort:  agentUrl,
// 		},
// 		Sampler: &jaegercfg.SamplerConfig{
// 			Type:  samplingType,
// 			Param: samplingParam,
// 		},
// 	}

// 	// example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
// 	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
// 	// frameworks
// 	jLogger := jaegerlog.StdLogger
// 	jMetricsFactory := metrics.NullFactory

// 	// initialize tracer with a logger and a metrics factory
// 	closer, err := cfg.InitGlobalTracer(
// 		serviceName,
// 		jaegercfg.Logger(jLogger),
// 		jaegercfg.Metrics(jMetricsFactory),
// 	)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "could not initialize tracer")
// 	}

// 	return closer.Close, nil
// }