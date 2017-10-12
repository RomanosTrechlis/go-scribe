package profiling

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
)

// Serve serves a profiling web server.
func Serve(port int) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	// Register pprof handlers
	http.HandleFunc("/debug/pprof", pprof.Index)
	http.HandleFunc("/debug/pprof/cmdline/", pprof.Cmdline)
	http.HandleFunc("/debug/pprof/profile/", pprof.Profile)
	http.HandleFunc("/debug/pprof/symbol/", pprof.Symbol)
	http.HandleFunc("/debug/pprof/trace/", pprof.Trace)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "http server failed: %v", err)
		}
	}()
	return srv
}
