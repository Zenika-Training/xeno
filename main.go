package main

import (
	"flag"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

var (
	watchedNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "population_count",
			Help: "Watched number.",
		},
		[]string{"population"},
	)
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(watchedNumber)
}

var (
	alien,
	marine,
	settler,
	loopMillisecondTimeout int
	infectedByTurn,
	marineKillByTurn,
	alienKillByTurn float64
)

func getParam(r *http.Request, p string) string {
	keys, ok := r.URL.Query()[p]

	if !ok || len(keys) < 1 {
		log.Printf("Url Param '%s' is missing\n", p)
		return ""
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	return keys[0]
}

func getenvfloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return floatVal
}

func getenvint(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intVal
}

func sendAliens(w http.ResponseWriter, r *http.Request) {
	inc, err := strconv.Atoi(getParam(r, "alien"))
	if err == nil {
		log.Printf("Adding %d 👽 aliens\n", inc)
		alien += inc
	}
}

func sendMarines(w http.ResponseWriter, r *http.Request) {
	inc, err := strconv.Atoi(getParam(r, "marine"))
	if err == nil {
		log.Printf("Adding %d 👮 marines\n", inc)
		marine += inc
	}
}

func resetSimulation() {
	alien = 0
	settler = 20000
	marine = 0
	log.Printf("Simulation reseted 🚀\n")
}

func resetSimulationHandler(w http.ResponseWriter, r *http.Request) {
	resetSimulation()
}

func maxInt(num, max int) int {
	return int(math.Max(float64(num), float64(max)))
}

func minInt(num, min int) int {
	return int(math.Min(float64(num), float64(min)))
}

func simulate(a int, s int, m int) (int, int, int) {
	alien := maxInt(a+minInt(int(infectedByTurn), s)-int(marineKillByTurn*float64(marine)), 0)
	settler := maxInt(s-int(infectedByTurn*float64(alien)), 0)
	marine := maxInt(marine-int(alienKillByTurn*float64(alien)), 0)
	return alien, settler, marine
}

func main() {
	flag.Parse()

	infectedByTurn = getenvfloat("INFECTED_BY_TURN", 2.0)
	marineKillByTurn = getenvfloat("MARINE_KILL_BY_TURN", 0.2)
	alienKillByTurn = getenvfloat("ALIEN_KILL_BY_TURN", 0.3)
	loopMillisecondTimeout = getenvint("LOOP_MILLISECOND_TIMEOUT", 1000)
	resetSimulation()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// Periodically record some sample latencies for the three services.
	go func() {
		for {
			alien, settler, marine = simulate(alien, settler, marine)

			watchedNumber.WithLabelValues("aliens").Set(float64(alien))
			watchedNumber.WithLabelValues("marines").Set(float64(marine))
			watchedNumber.WithLabelValues("settlers").Set(float64(settler))
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/sendAliens", sendAliens)
	http.HandleFunc("/sendMarines", sendMarines)
	http.HandleFunc("/resetSimulation", resetSimulationHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
