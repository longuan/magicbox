package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
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

func (d dataPoints) getSlice(from, to int64) dataPoints {
	startIndex := sort.Search(len(d), func(i int) bool {
		return d[i][1] >= from
	})
	endIndex := sort.Search(len(d), func(i int) bool {
		return d[i][1] >= to
	})
	return d[startIndex:endIndex]
}

func (q *queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBody := &queryRequestBody{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		fmt.Println("decode request body err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("receive request for", r.Method, r.URL, reqBody)

	startTime, err := time.ParseInLocation(time.RFC3339Nano, reqBody.Range.From, time.UTC)
	if err != nil {
		fmt.Printf("parse from time %s err %s\n", reqBody.Range.From, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endTime, err := time.ParseInLocation(time.RFC3339Nano, reqBody.Range.To, time.UTC)
	if err != nil {
		fmt.Printf("parse to time %s err %s\n", reqBody.Range.To, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		allResponse := make(queryResponse, 0)

		for _, target := range reqBody.Targets {
			submetrics, ok := target.Payload["submetric"]
			if !ok {
				continue
			}
			for _, submetric := range submetrics {
				allPoints, ok := allMetricDataPoints[submetric]
				if !ok {
					continue
				}

				item := queryResponseItem{
					Target:     allPoints.Target,
					Datapoints: allPoints.Datapoints.getSlice(startTime.UnixMilli(), endTime.UnixMilli()),
				}
				allResponse = append(allResponse, item)
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
