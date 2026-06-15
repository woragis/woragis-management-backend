package middleware

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type metricsKey struct {
	method string
	route  string
	status int
}

var (
	metricsMu sync.RWMutex
	counters  = map[metricsKey]*atomic.Uint64{}
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		route := r.Pattern
		if route == "" {
			route = r.URL.Path
		}
		incMetric(r.Method, route, rec.status)
	})
}

func incMetric(method, route string, status int) {
	k := metricsKey{method: method, route: route, status: status}
	metricsMu.RLock()
	c, ok := counters[k]
	metricsMu.RUnlock()
	if ok {
		c.Add(1)
		return
	}
	metricsMu.Lock()
	defer metricsMu.Unlock()
	if c, ok = counters[k]; !ok {
		c = &atomic.Uint64{}
		counters[k] = c
	}
	c.Add(1)
}

func PrometheusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		_, _ = fmt.Fprintln(w, "# HELP http_requests_total Total HTTP requests handled by the API.")
		_, _ = fmt.Fprintln(w, "# TYPE http_requests_total counter")

		metricsMu.RLock()
		keys := make([]metricsKey, 0, len(counters))
		for k := range counters {
			keys = append(keys, k)
		}
		metricsMu.RUnlock()
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].route != keys[j].route {
				return keys[i].route < keys[j].route
			}
			if keys[i].method != keys[j].method {
				return keys[i].method < keys[j].method
			}
			return keys[i].status < keys[j].status
		})

		for _, k := range keys {
			metricsMu.RLock()
			c := counters[k]
			metricsMu.RUnlock()
			val := c.Load()
			if val == 0 {
				continue
			}
			route := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(k.route)
			_, _ = fmt.Fprintf(w, "http_requests_total{method=%q,route=%q,status=%q} %d\n",
				k.method, route, strconv.Itoa(k.status), val)
		}
	}
}
