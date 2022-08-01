package main

import (
	"fmt"

	"github.com/longuan/magicbox/internal/mongobox"
)

func main() {
	mongodPath := "/home/longanliu/code/github.com/mongodb/mongo-v4.2/mongod" // 使用环境变量中的mongod

	repl, err := mongobox.NewReplicaSet(mongodPath, "rs0", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	repl.PrettyPrint()
}
