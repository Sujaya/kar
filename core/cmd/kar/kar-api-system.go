package main

/*
 * This file contains the implementation of the portion of the
 * KAR REST API related to system-level operations.
 */

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.ibm.com/solsa/kar.git/core/internal/pubsub"
	"github.ibm.com/solsa/kar.git/core/internal/runtime"
	"github.ibm.com/solsa/kar.git/core/pkg/logger"
)

// swagger:route POST /v1/system/shutdown system idSystemShutdown
//
// shutdown
//
// ### Shutdown a single KAR runtime
//
// Initiate an orderly shutdown of the target KAR runtime process.
//
//     Schemes: http
//     Responses:
//       200: response200
//
func shutdown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "OK")
	logger.Info("Invoking cancel() in response to shutdown request")
	cancel()
}

// swagger:route GET /v1/system/health system isSystemHealth
//
// health
//
// ### Health-check endpoint
//
// Returns a `200` response to indicate that the KAR runtime processes is healthy.
//
//     Schemes: http
//     Responses:
//       200: response200
//
func health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "OK")
}

// post handles a direct http request from a peer sidecar
// TODO swagger
func post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	value, _ := ioutil.ReadAll(r.Body)
	m := pubsub.Message{Value: value}
	process(m)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "OK")
}

// Returns information about a specified component, controlled by the call path
// Options are given by the cases
// Format type (text/plain vs application/json) is controlled by Accept header in call
func getInformation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	format := "text/plain"
	if r.Header.Get("Accept") == "application/json" {
		format = "application/json"
	}
	component := ps.ByName("component")
	var data string
	var err error
	switch component {
	case "sidecars", "Sidecars":
		data, err = pubsub.GetSidecars(format)
	case "actors", "Actors":
		data, err = runtime.GetAllActors(ctx, format)
	case "sidecar_actors":
		data, err = runtime.GetActors()
	default:
		http.Error(w, fmt.Sprintf("Invalid information query: %v", component), http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to acquire %v information: %v", component, err), http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", format)
		fmt.Fprint(w, data)
	}
}
