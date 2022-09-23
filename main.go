package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "cpuidle"
	// https://www.kernel.org/doc/Documentation/cpuidle/sysfs.txt
	system_cpu_sysfs_path = "/sys/devices/system/cpu/"
)

var (
	state_time = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "time",
			Help:      "Total time spent in this idle state (in microseconds)",
		},
		[]string{"cpu", "state", "state_name"},
	)
)

func readFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return strings.TrimSuffix(string(bytes), "\n")
}

type StateMetrics struct {
	time int
	name string
}

func readStateMetrics(state_sysfs_path string) StateMetrics {
	state_time, err := strconv.Atoi(readFile(state_sysfs_path + "/time"))
	if err != nil {
		panic(err)
	}
	state_name := readFile(state_sysfs_path + "/name")
	return StateMetrics{time: state_time, name: state_name}
}

func setMetrics(cpu string, state string) {
	state_sysfs_path := system_cpu_sysfs_path + cpu + "/cpuidle/" + state
	for {
		state_metrics := readStateMetrics(state_sysfs_path)
		state_time.With(prometheus.Labels{"cpu": cpu, "state": state, "state_name": state_metrics.name}).Set(float64(state_metrics.time))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	system_cpu_sysfs, err := ioutil.ReadDir(system_cpu_sysfs_path)
	if err != nil {
		panic(err)
	}

	for _, cpu_sysfs := range system_cpu_sysfs {
		if cpu_sysfs.IsDir() {
			cpu := cpu_sysfs.Name()
			if _, err := strconv.Atoi(cpu[3:]); err == nil && cpu[:3] == "cpu" {
				cpuidle_sysfs, err := ioutil.ReadDir(system_cpu_sysfs_path + cpu + "/cpuidle/")
				if err != nil {
					panic(err)
				}
				for _, state_sysfs := range cpuidle_sysfs {
					state := state_sysfs.Name()
					go setMetrics(cpu, state)
				}
			}
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9975", nil)
}
