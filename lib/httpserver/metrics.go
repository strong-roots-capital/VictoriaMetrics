package httpserver

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/buildinfo"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/memory"
	"github.com/VictoriaMetrics/metrics"
)

func writePrometheusMetrics(w io.Writer) {
	metrics.WritePrometheus(w, true)

	fmt.Fprintf(w, "vm_app_version{version=%q} 1\n", buildinfo.Version)
	fmt.Fprintf(w, "vm_allowed_memory_bytes %d\n", memory.Allowed())

	// Export start time and uptime in seconds
	fmt.Fprintf(w, "vm_app_start_timestamp %d\n", startTime.Unix())
	fmt.Fprintf(w, "vm_app_uptime_seconds %d\n", int(time.Since(startTime).Seconds()))

	// TODO: export other interesting stuff.

	// Export flags as metrics.
	flag.VisitAll(func(f *flag.Flag) {
		lname := strings.ToLower(f.Name)
		value := f.Value.String()
		if strings.Contains(lname, "pass") || strings.Contains(lname, "key") || strings.Contains(lname, "secret") {
			// Do not expose passwords and keys to prometheus.
			value = "secret"
		}
		fmt.Fprintf(w, "flag{name=%q, value=%q} 1\n", f.Name, value)
	})
}

var startTime = time.Now()
