package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

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

		if strings.Compare(name,"health") != 0 {
			var serverSpan opentracing.Span
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Println(err)
				log.Println(r.Header)
			}

			// Create the span referring to the RPC client if available.
			// If wireContext == nil, a root span will be created.
			serverSpan = opentracing.StartSpan(
				name,
				ext.RPCServerOption(wireContext))

			defer serverSpan.Finish()
		}

		inner.ServeHTTP(w, r)

		if strings.Compare(name,"health") != 0 {
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

	///////////////////////////////// Zipkin Connection ////////////////////////////////
	collector, err := zipkin.NewHTTPCollector("http://tracing:9411/api/v1/spans")
	if err != nil {
		log.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "read-cqrs")

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		log.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// Explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)

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

	/// Métier - user
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
		Name("user_put").
		Handler(handlerUserUpdate)

	
	/// Métier - topic
	var handlerTopicCreate http.Handler
	handlerTopicCreate = LoggerMiddleware(s.handlerTopicCreateFunc, "topicCreate", histogram, nil)
	router.
		Methods("POST").
		Path("/topics").
		Name("topics").
		Handler(handlerTopicCreate)
	
	
	var handlerTopicUpdate http.Handler
	handlerTopicUpdate = LoggerMiddleware(s.handlerTopicUpdateFunc, "topicUpdate", histogram, nil)
	router.
		Methods("PUT").
		Path("/topics/{id}").
		Name("topic_put").
		Handler(handlerTopicUpdate)

	
	/// Métier - message
	var handlerMessageCreate http.Handler
	handlerMessageCreate = LoggerMiddleware(s.handlerMessageCreateFunc, "messageCreate", histogram, nil)
	router.
		Methods("POST").
		Path("/messages").
		Name("messages").
		Handler(handlerMessageCreate)
	
	
	var handlerMessageUpdate http.Handler
	handlerMessageUpdate = LoggerMiddleware(s.handlerMessageUpdateFunc, "messageUpdate", histogram, nil)
	router.
		Methods("PUT").
		Path("/messages/{id}").
		Name("message_put").
		Handler(handlerMessageUpdate)

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
