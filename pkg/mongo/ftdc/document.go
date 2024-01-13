package ftdc

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	mongoFTDC "github.com/mongodb/ftdc"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/longuan/magicbox/pkg/mongo/status"
	"github.com/longuan/magicbox/pkg/utils"
)

var (
	FlatFTDCDocumentFields = []string{}
)

type FTDCDocument struct {
	Start            time.Time                    `bson:"start"`
	End              time.Time                    `bson:"end"`
	ServerStatus     status.ServerStatus          `bson:"serverStatus"`
	ReplSetGetStatus status.ReplSetGetStatus      `bson:"replSetGetStatus"`
	OplogRsCollStats status.LocalOplogRsCollStats `bson:"local.oplog.rs.stats"`
	SystemMetrics    status.SystemMetrics         `bson:"systemMetrics"`
}

func init() {
	for _, field := range utils.GetExportedFields(reflect.TypeOf(FTDCDocument{})) {
		FlatFTDCDocumentFields = append(FlatFTDCDocumentFields, utils.GetFieldNamesRecursive(field)...)
	}
}

// ParseFile 解析metric file，把文件内容转换成一条条的文档
// TODO: 把转换过程中的错误暴露出去
func ParseFile(metricFile string) (<-chan *FTDCDocument, error) {
	file, err := os.Open(metricFile)
	if err != nil {
		return nil, fmt.Errorf("open %s error %s", metricFile, err)
	}

	retChan := make(chan *FTDCDocument)

	go func() {
		defer close(retChan)

		iter := mongoFTDC.ReadStructuredMetrics(context.Background(), file)
		defer iter.Close()

		for iter.Next() {
			rawData, err := iter.Document().MarshalBSON()
			if err != nil {
				panic(err)
			}
			var doc FTDCDocument
			err = bson.Unmarshal(rawData, &doc)
			if err != nil {
				panic(err)
			}
			retChan <- &doc
		}
		if iter.Err() != nil {
			panic(iter.Err())
		}
	}()

	return retChan, nil
}
