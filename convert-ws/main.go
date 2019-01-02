package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-redis/redis"
)


// LoggerMiddleware add logger and metrics
func LoggerMiddleware(inner http.HandlerFunc, name string, histogram *prometheus.HistogramVec, counter *prometheus.CounterVec) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		inner.ServeHTTP(w, r)

		time := time.Since(start)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time,
		)

		histogram.WithLabelValues(r.RequestURI).Observe(time.Seconds())
		if counter != nil {
			counter.WithLabelValues(r.RequestURI).Inc()
		}
	})
}

type server struct {
	nc *nats.Conn
	ec *nats.EncodedConn
	client *redis.Client
}

func main() {

	var s server
		
	///////////////////////////////// Redis Connection ////////////////////////////////
	ruri := os.Getenv("REDIS_URI")
	s.client = redis.NewClient(&redis.Options{
		Addr:     ruri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	///////////////////////////////// Nats Connection ////////////////////////////////
	var err error
	uri := os.Getenv("NATS_URI")

	for i := 0; i < 5; i++ {
		nc, err := nats.Connect(uri)
		if err == nil {
			s.nc = nc
			s.ec, err = nats.NewEncodedConn(s.nc, nats.JSON_ENCODER)
			if err != nil {
				log.Fatal("Error establishing connection to NATS encoded:", err)
			}
			break
		}

		log.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	log.Println("Connected to NATS at:", s.nc.ConnectedUrl())
	s.ec.Subscribe("users", s.eventReceive)

	///////////////////////////////// Http Connection ////////////////////////////////
	router := mux.NewRouter().StrictSlash(true)

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "webservice_uri_duration_seconds",
		Help: "Time to respond",
	}, []string{"uri"})

	promCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "webservice_count",
		Help: "counter for api",
	}, []string{"uri"})

	/// Root
	var handlerStatus http.Handler
	handlerStatus = LoggerMiddleware(handlerStatusFunc, "root", histogram, nil)
	router.
		Methods("GET").
		Path("/").
		Name("root").
		Handler(handlerStatus)

	var handlerHealth http.Handler
	handlerHealth = LoggerMiddleware(handlerHealthFunc, "health", histogram, nil)
	router.
		Methods("GET").
		Path("/healthz").
		Name("health").
		Handler(handlerHealth)

	// add prometheus
	prometheus.Register(histogram)
	prometheus.Register(promCounter)
	router.Methods("GET").Path("/metrics").Name("Metrics").Handler(promhttp.Handler())

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Printf("Start server...")
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
