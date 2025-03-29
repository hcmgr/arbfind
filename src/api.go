///////////////////////////
// api.go
// ------
// Defines the REST api that is accessed by our frontend.
///////////////////////////

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getArbsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	arbs, err := db.readArbs()
	if err != nil {
		fmt.Println("Failed to get arbs from db")
	}

	fmt.Println(arbs)

	// set header
	w.Header().Set("Content-Type", "application/json")

	// write header and payload
	arbsRawBytes, _ := json.Marshal(arbs)
	bytesWritten, err := w.Write(arbsRawBytes)
	if err != nil {
		http.Error(w, "Error writing payload", http.StatusInternalServerError)
	}

	fmt.Println("Payload bytes:", bytesWritten)
}

func startAPIServer() {
	// bind endpoint handlers
	http.HandleFunc("/getArbs", getArbsHandler)

	port := config.BackendPort
	addr := "localhost:" + strconv.Itoa(port)

	fmt.Println("Starting API server at http://" + addr)

	// start API server
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		panic(err)
	}
}
