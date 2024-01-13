package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/longuan/magicbox/pkg/mongo/ftdc"
	"github.com/longuan/magicbox/pkg/utils"
)

// http接口定义参考https://github.com/simPod/GrafanaJsonDatasource/blob/0.6.x/openapi.yaml

type metricsRequestBody struct {
	Metric  string              `json:"metric"`
	Payload map[string][]string `json:"payload"`
}

type metricsResponse []metricsResponseItem

type metricsResponseItem struct {
	Value    string           `json:"value"`
	Payloads []metricsPayload `json:"payloads"`
}

type metricsPayload struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Placeholder  string          `json:"placeholder"`
	ReloadMetric bool            `json:"reloadMetric"`
	Options      []payloadOption `json:"options"`
}

type payloadOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := &metricsRequestBody{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		fmt.Println("decode request body err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("receive request for", r.Method, r.URL, reqBody)

	if r.Method == http.MethodPost {
		res := metricsResponse{}

		for _, metric := range utils.GetExportedFields(reflect.TypeOf(ftdc.FTDCDocument{})) {
			selectOptions := make([]payloadOption, 0)
			for _, subMetric := range ftdc.FlatFTDCDocumentFields {
				if strings.HasPrefix(subMetric, metric.Name) {
					selectOptions = append(selectOptions, payloadOption{
						Label: subMetric,
						Value: subMetric,
					})
				}
			}

			metricItem := metricsResponseItem{
				Value: metric.Name,
				Payloads: []metricsPayload{
					{
						Name:         "submetric",
						Type:         "multi-select",
						Placeholder:  "子指标名称",
						ReloadMetric: false,
						Options:      selectOptions,
					},
				},
			}
			res = append(res, metricItem)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			panic(err)
		}
	}
}
