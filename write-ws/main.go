package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

)

type server struct {
	nc *nats.Conn
	ec *nats.EncodedConn
	db *gorm.DB
}

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

func main() {

	var s server
	var err error

	///////////////////////////////// Nats Connection ////////////////////////////////
	
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

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	fmt.Println("Connected to NATS at:", s.nc.ConnectedUrl())


	///////////////////////////////// MySQL Connection ////////////////////////////////
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_NAME")
	dbUser := os.Getenv("MYSQL_USER")

	for i := 0; i < 10; i++ {
		db, err := gorm.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
		if err == nil {
			s.db = db
			break
		}

		fmt.Println("Waiting before connecting to MySQL at:", dbHost)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to MySQL:", err)
	}
	
  	defer s.db.Close()

	s.testCreateTables()

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

	/// Métier
	var handlerUserCreate http.Handler
	handlerUserCreate = LoggerMiddleware(s.handlerUserCreateFunc, "userCreate", histogram, nil)
	router.
		Methods("POST").
		Path("/users").
		Name("users").
		Handler(handlerUserCreate)
	
	
	var handlerUserUpdate http.Handler
	handlerUserUpdate = LoggerMiddleware(s.handlerUserUpdateFunc, "userUpdate", histogram, nil)
	router.
		Methods("PUT").
		Path("/users/{id}").
		Name("users").
		Handler(handlerUserUpdate)

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
