package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("addr", ":9000", "http service address")

	httpReqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code.",
		},
		[]string{"code"},
	)

	httpReqDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration distribution.",
		Buckets: prometheus.ExponentialBuckets(0.05, 2, 5),
	})
)

func main() {
	flag.Parse()
	log.Printf("Listening on address %s...\n", *addr)

	http.Handle("/", http.HandlerFunc(handler))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	codeStr := getRandomValue(query.Get("status"), "200")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		// TODO(klw0): Handle gracefully.
		log.Fatal(err)
	}

	delayStr := getRandomValue(query.Get("delay"), "0s")
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		// TODO(klw0): Handle gracefully.
		log.Fatal(err)
	}

	httpReqs.WithLabelValues(codeStr).Inc()
	httpReqDuration.Observe(delay.Seconds())

	time.Sleep(delay)

	w.WriteHeader(code)
	resp := fmt.Sprintf("%s %s %s @ %s", req.Method, req.URL.Path, codeStr, delayStr)
	fmt.Fprintln(w, resp)
	log.Println(resp)
}

// Returns a random value based on percentages provided in `csv` (which has the
// form `v1:p1,v2:p2,...,vN:pN`). Unspecified percentage space is allocated to
// `fallback`.
func getRandomValue(csv string, fallback string) string {
	if len(csv) == 0 {
		return fallback
	}

	rnd := rand.Float64()
	probability := float64(0)
	// TODO(klw0): This approach ignores total percents > 100. Does it matter
	// beyond our ability to print a warning?
	for _, pair := range strings.Split(csv, ",") {
		if !strings.Contains(pair, ":") {
			log.Printf("warning: Ignoring invalid value-percent pair %s\n", pair)
			continue
		}

		parts := strings.Split(pair, ":")
		percent, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Printf("warning: Ignoring invalid value-percent pair %s\n", pair)
			continue
		}

		probability = probability + (percent / 100.0)
		if rnd < probability {
			return parts[0]
		}
	}

	return fallback
}
