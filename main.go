package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/woraphol-j/url-shortener/pkg/mongo"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "github.com/woraphol-j/url-shortener/urlshortener"
	"github.com/woraphol-j/url-shortener/urlshortener"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

const (
	defaultPort = "8080"
)

func main() {
	var (
		addr     = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		// ctx      = context.Background()
	)

	flag.Parse()

	fmt.Println("^^^^^^^", httpAddr)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// var (
	// 	cargos         = inmem.NewCargoRepository()
	// 	locations      = inmem.NewLocationRepository()
	// 	voyages        = inmem.NewVoyageRepository()
	// 	handlingEvents = inmem.NewHandlingEventRepository()
	// )

	// Configure some questionable dependencies.
	// var (
	// 	handlingEventFactory = cargo.HandlingEventFactory{
	// 		CargoRepository:    cargos,
	// 		VoyageRepository:   voyages,
	// 		LocationRepository: locations,
	// 	}
	// 	handlingEventHandler = handling.NewEventHandler(
	// 		inspection.NewService(cargos, handlingEvents, nil),
	// 	)
	// )

	// Facilitate testing by adding some cargos.
	// storeTestData(cargos)

	fieldKeys := []string{"method"}

	// var rs routing.Service
	// rs = routing.NewProxyingMiddleware(ctx, *routingServiceURL)(rs)
	if err := godotenv.Load(); err != nil {
		panic("Fail to load configuration")
	}

	mongoURL := os.Getenv("MONGO_URL")
	mongoDb := os.Getenv("MONGO_DATABASE")
	mongoColl := os.Getenv("MONGO_COLLECTION")
	var urlDAO = mongo.NewDAO(mongoURL, mongoDb, mongoColl)
	var us urlshortener.Service
	// us = urlshortener.NewService(cargos, locations, handlingEvents, rs)
	us = urlshortener.NewService(urlDAO)
	us = urlshortener.NewLoggingService(log.With(logger, "component", "url-shortener"), us)
	us = urlshortener.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "url_shortener_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "url_shortener_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		us,
	)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/", urlshortener.MakeHandler(us, httpLogger))

	// http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

// func storeTestData(r cargo.Repository) {
// 	test1 := cargo.New("FTL456", cargo.RouteSpecification{
// 		Origin:          location.AUMEL,
// 		Destination:     location.SESTO,
// 		ArrivalDeadline: time.Now().AddDate(0, 0, 7),
// 	})
// 	if err := r.Store(test1); err != nil {
// 		panic(err)
// 	}

// 	test2 := cargo.New("ABC123", cargo.RouteSpecification{
// 		Origin:          location.SESTO,
// 		Destination:     location.CNHKG,
// 		ArrivalDeadline: time.Now().AddDate(0, 0, 14),
// 	})
// 	if err := r.Store(test2); err != nil {
// 		panic(err)
// 	}
// }
