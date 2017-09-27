package streamer

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
)

func pprofServer(r *http.ServeMux, pport int) {
	// Register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", pport), r); err != nil {
		fmt.Fprintf(os.Stderr, "http server failed: %v", err)
	}
}
