package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type queryRequestBody struct {
	Range         queryRange         `json:"range"`
	Interval      string             `json:"interval"`
	IntervalMs    int                `json:"intervalMs"`
	MaxDataPoints int                `json:"maxDataPoints"`
	Targets       []queryTarget      `json:"targets"`
	AdhocFilters  []queryAdhocFilter `json:"adhocFilters"`
}

type queryRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type queryTarget struct {
	Target  string              `json:"target"`
	RefId   string              `json:"refId"`
	Payload map[string][]string `json:"payload"`
}

type queryAdhocFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type queryHandler struct {
}

type queryResponse []queryResponseItem

type queryResponseItem struct {
	Target     string     `json:"target"`
	Datapoints dataPoints `json:"datapoints"`
}

type dataPoints [][2]int64 // [ [metric value as a float, unix timestamp in milliseconds] ]

func (q *queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBody := &queryRequestBody{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		fmt.Println("decode request body err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("receive request for", r.Method, r.URL, reqBody)

	if r.Method == http.MethodPost {
		allResponse := make(queryResponse, 0)

		for _, target := range reqBody.Targets {
			submetrics, ok := target.Payload["submetric"]
			if !ok {
				continue
			}
			for _, submetric := range submetrics {
				item, ok := allMetricDataPoints[submetric]
				if !ok {
					continue
				}

				allResponse = append(allResponse, *item)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(allResponse)
		if err != nil {
			panic(err)
		}
	}
}
