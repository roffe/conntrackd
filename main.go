package main

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricMap map[string]*prometheus.GaugeVec

/*
cpu=0   	found=0 invalid=227 ignore=26597901 insert=0 insert_failed=0 drop=0 early_drop=0 error=0 search_restart=1925333
cpu=1   	found=0 invalid=91 ignore=26617950 insert=0 insert_failed=0 drop=0 early_drop=0 error=0 search_restart=826242
*/

func getStats() {
	out := bytes.Buffer{}
	d := time.Now().Add(5 * time.Second)

	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	cmd := exec.CommandContext(ctx, "conntrack", "-S")
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Wait()

	scanner := bufio.NewScanner(&out)

	var results []map[string]int
	for scanner.Scan() {
		res := make(map[string]int)
		for _, sa := range strings.Fields(scanner.Text()) {

			parts := strings.Split(sa, "=")
			if len(parts) == 2 {
				intVal, err := strconv.Atoi(parts[1])
				if err != nil {
					log.Println(err)
					break
				}
				res[parts[0]] = intVal
			}
		}
		results = append(results, res)
	}
	publishResults(results)
}

func publishResults(results []map[string]int) {
	for _, res := range results {
		//fmt.Println(res)
		for k, v := range res {
			if k == "cpu" {
				continue
			}
			//fmt.Println(strconv.Itoa(res["cpu"]), k, v)
			metricMap[k].With(prometheus.Labels{"cpu": strconv.Itoa(res["cpu"])}).Set(float64(v))
		}
	}
}

func init() {
	metricMap = map[string]*prometheus.GaugeVec{
		"found": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_found",
			Help: "The conntrack found",
		}, []string{"cpu"}),
		"invalid": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_invalid",
			Help: "The conntrack invalid",
		}, []string{"cpu"}),
		"ignore": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_ignore",
			Help: "The conntrack ignore",
		}, []string{"cpu"}),
		"insert": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_insert",
			Help: "The conntrack insert",
		}, []string{"cpu"}),
		"insert_failed": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_insert_failed",
			Help: "The conntrack insert_failed",
		}, []string{"cpu"}),
		"drop": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_drop",
			Help: "The conntrack drop",
		}, []string{"cpu"}),
		"early_drop": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_early_drop",
			Help: "The conntrack early_drop",
		}, []string{"cpu"}),
		"error": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_error",
			Help: "The conntrack error",
		}, []string{"cpu"}),
		"search_restart": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "conntrack_search_start",
			Help: "The conntrack search_start",
		}, []string{"cpu"}),
	}
	for _, metric := range metricMap {
		prometheus.MustRegister(metric)
	}
}

func main() {
	getStats()

	ticker := time.NewTicker(30 * time.Second)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	for {
		select {
		case <-ticker.C:
			getStats()
		}
	}
}
