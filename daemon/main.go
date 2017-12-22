package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // serves data on /debug/pprof/*
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/eparis/access-daemon/api"

	// load all operations (via init() in the operation
	_ "github.com/eparis/access-daemon/operations/cat"
	_ "github.com/eparis/access-daemon/operations/command"
	_ "github.com/eparis/access-daemon/operations/ip"
	_ "github.com/eparis/access-daemon/operations/journalctl"
)

func handleRoles(w http.ResponseWriter, r *http.Request) {
	roles := api.GetRoleNames()
	for _, role := range roles {
		fmt.Fprintf(w, "%s\n", role)
	}
	return
}

func handleOps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	role := api.Role(vars["role"])
	ops := api.GetOperationNames(role)
	for _, op := range ops {
		fmt.Fprintf(w, "%s\n", op)
	}
	return
}

func handleOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	role := api.Role(vars["role"])
	opName := vars["opName"]

	op, err := api.GetOperation(role, opName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = op.Go(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// func main is over in root.go with the cobra init. But the real guts of the program start here.
func mainFunc(cmd *cobra.Command, args []string) error {
	fmt.Printf("bind-addr: %s\n", bindAddr)
	cfgDir, err := filepath.Abs(cfgDir)
	if err != nil {
		return err
	}
	fmt.Printf("config-dir: %s\n", cfgDir)

	err = api.InitializeOperations(cfgDir)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	s := r.Methods("GET", "POST").Subrouter()
	if staticPath != "" {
		ops := api.GetRoleNames()
		for _, op := range ops {
			if op == "static" {
				return fmt.Errorf("Role 'static' registered, but is used to serve static files")
			}
		}
		s.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	}
	s.PathPrefix("/metrics").Handler(promhttp.Handler())
	s.HandleFunc("/{role}/{opName}", handleOperation)
	s.HandleFunc("/{role}/", handleOps)
	s.HandleFunc("/", handleRoles)
	http.Handle("/", r)

	err = http.ListenAndServe(bindAddr, promRequestHandler(handlers.CompressHandler(handlers.CombinedLoggingHandler(os.Stdout, r))))
	if err != nil {
		return err
	}
	return nil
}
