package main

import (
	"fmt"
	"reflect"

	"github.com/longuan/magicbox/pkg/mongo/ftdc"
	"github.com/longuan/magicbox/pkg/utils"
)

var (
	// allMetricDataPoints 存储所有的指标数据在内存中
	allMetricDataPoints = map[string]*queryResponseItem{}
)

func readFile(path string) error {
	docChan, err := ftdc.ParseFile(path)
	if err != nil {
		return err
	}

	count := 0
	for doc := range docChan {
		readDoc(doc)
		count++
	}
	fmt.Println("total parsed", count, "ftdc document")
	return nil
}

func readDoc(doc *ftdc.FTDCDocument) {
	// 以ftdc document的start为监控显示时间
	docTime := doc.Start.UnixMilli()

	for _, fullPath := range ftdc.FlatFTDCDocumentFields {
		value, kind := utils.GetValueByPath(reflect.ValueOf(*doc), fullPath)

		var int64Value int64
		switch kind {
		case reflect.Bool:
			boolValue, ok := value.(bool)
			if !ok {
				panic(fmt.Sprintf("%s is not bool", kind))
			}

			if boolValue {
				int64Value = 1
			} else {
				int64Value = 0
			}
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int:
			intValue, ok := value.(int)
			if !ok {
				panic(fmt.Sprintf("%s is not int", kind))
			}
			int64Value = int64(intValue)
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			intValue, ok := value.(int64)
			if !ok {
				panic(fmt.Sprintf("%s is not int64", kind))
			}
			int64Value = intValue
		case reflect.String:
			continue
		default:
			panic(fmt.Sprintf("%s is unknown kind", kind))
		}

		item, ok := allMetricDataPoints[fullPath]
		if !ok {
			item = &queryResponseItem{
				Target:     fullPath,
				Datapoints: make(dataPoints, 0, 1000),
			}
			allMetricDataPoints[fullPath] = item
		}
		item.Datapoints = append(item.Datapoints, [2]int64{int64Value, docTime})
	}
}
