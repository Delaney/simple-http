package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK!")
}

type MathRequest struct {
	Op    string  `json:"op"`
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}

type MathResponse struct {
	Error  string  `json:"error"`
	Result float64 `json:"result"`
}

func mathHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode
	dec := json.NewDecoder(r.Body)
	req := &MathRequest{}

	if err := dec.Decode(req); err != nil {
		log.Printf("Error: Bad JSON: %s", err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	// Validate
	if !strings.Contains("+-*/", req.Op) {
		log.Printf("Error: Invalid operator: %s", req.Op)
		http.Error(w, "Invalid operator", http.StatusBadRequest)
		return
	}

	// Compute
	resp := &MathResponse{}
	switch req.Op {
	case "+":
		resp.Result = req.Left + req.Right
	case "-":
		resp.Result = req.Left - req.Right
	case "*":
		resp.Result = req.Left * req.Right
	case "/":
		if req.Right == 0.0 {
			resp.Error = "Division by 0"
		} else {
			resp.Result = req.Left / req.Right
		}
	default:
		resp.Error = fmt.Sprintf("Unknown operation: %s", req.Op)
	}

	// Return
	w.Header().Set("content-type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Printf("Error: Returning: %s", err)
		http.Error(w, "Returning", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/math", mathHandler)

	addr := ":8080"
	log.Printf("Running server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
