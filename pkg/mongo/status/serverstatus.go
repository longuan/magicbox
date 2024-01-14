package status

// ServerStatus 具体结构参考serverStatus命令的输出
type ServerStatus struct {
	Connections ServerStatusConnections `bson:"connections"`
	FlowControl ServerStatusFlowControl `bson:"flowControl"`
	GlobalLock  ServerStatusGlobalLock  `bson:"globalLock"`
	Tcmalloc    ServerStatusTcmalloc    `bson:"tcmalloc"`
	WiredTiger  ServerStatusWiredTiger  `bson:"wiredTiger"`
}

type ServerStatusConnections struct {
	Current      int `bson:"current"`
	Available    int `bson:"available"`
	TotalCreated int `bson:"totalCreated"`
	Active       int `bson:"active"`
}

type ServerStatusFlowControl struct {
	Enabled             bool  `bson:"enabled"`
	TargetRateLimit     int   `bson:"targetRateLimit"`
	TimeAcquiringMicros int64 `bson:"timeAcquiringMicros"`
	LocksPerOp          int   `bson:"locksPerOp"`
	SustainerRate       int   `bson:"sustainerRate"`
	IsLagged            bool  `bson:"isLagged"`
	IsLaggedCount       int   `bson:"isLaggedCount"`
	IsLaggedTimeMicros  int64 `bson:"isLaggedTimeMicros"`
}

type globalLockRWInfo struct {
	Total   int `bson:"total"`
	Readers int `bson:"readers"`
	Writers int `bson:"writers"`
}

type ServerStatusGlobalLock struct {
	TotalTime     int64            `bson:"totalTime"`
	CurrentQueue  globalLockRWInfo `bson:"currentQueue"`
	ActiveClients globalLockRWInfo `bson:"activeClients"`
}
