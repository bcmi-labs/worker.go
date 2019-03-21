package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bcmi-labs/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	jobs := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jobs",
	})
	prometheus.MustRegister(jobs)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening prometheus metrics on %s/metrics\n", ":8888")
	go func() {
		log.Fatal(http.ListenAndServe(":8888", nil))
	}()

	pool := worker.Pool{
		Jobs: jobs,
	}

	pool.Run(func() {
		log.Println("If you curl http://localhost:8888/metrics | jobs now, you should see it's 1")
		time.Sleep(time.Second * 10)
		log.Println("If you curl http://localhost:8888/metrics | jobs now, you should see it's 0")
	})

	select {}
}
