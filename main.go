package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
	cg "github.com/woraphol-j/url-shortener/pkg/codegenerator"
	r "github.com/woraphol-j/url-shortener/pkg/repository"
	"github.com/woraphol-j/url-shortener/urlshortener"
)

const (
	defaultPort = "8080"
)

func newMongoRepository() r.Repository {
	mongoURL := os.Getenv("MONGO_URL")
	mongoDb := os.Getenv("MONGO_DATABASE")
	mongoColl := os.Getenv("MONGO_COLLECTION")
	var urlDAO = r.NewMongoRepository(mongoURL, mongoDb, mongoColl)
	return urlDAO
}

func newMysqlRepository() r.Repository {
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		logger.Panic("Mysql Port is invalid: ", err)
	}
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlUsername,
		mysqlPassword,
		mysqlHost,
		mysqlPort,
		mysqlDatabase,
	)
	var urlDAO = r.NewMySQLRepository(connStr)
	return urlDAO
}

func main() {
	var (
		addr     = os.Getenv("APP_PORT")
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	fieldKeys := []string{"method"}

	urlRepo := newMysqlRepository()

	var us urlshortener.Service
	us = urlshortener.NewService(urlRepo, cg.NewCodeGenerator())
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
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

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
