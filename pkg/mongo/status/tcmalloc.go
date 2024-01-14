package status

type ServerStatusTcmalloc struct {
	Generic  ServerStatusTcmallocGeneric  `bson:"generic"`
	Tcmalloc ServerStatusTcmallocInternal `bson:"tcmalloc"`
}

type ServerStatusTcmallocGeneric struct {
	Current_allocated_bytes int64 `bson:"current_allocated_bytes"`
	Heap_size               int64 `bson:"heap_size"`
}

type ServerStatusTcmallocInternal struct {
	Pageheap_free_bytes              int64 `bson:"pageheap_free_bytes"`
	Pageheap_unmapped_bytes          int64 `bson:"pageheap_unmapped_bytes"`
	Max_total_thread_cache_bytes     int64 `bson:"max_total_thread_cache_bytes"`
	Current_total_thread_cache_bytes int64 `bson:"current_total_thread_cache_bytes"`
	Total_free_bytes                 int64 `bson:"total_free_bytes"`
	Central_cache_free_bytes         int64 `bson:"central_cache_free_bytes"`
	Transfer_cache_free_bytes        int64 `bson:"transfer_cache_free_bytes"`
	Thread_cache_free_bytes          int64 `bson:"thread_cache_free_bytes"`
	Aggressive_memory_decommit       int64 `bson:"aggressive_memory_decommit"`
	Pageheap_committed_bytes         int64 `bson:"pageheap_committed_bytes"`
	Pageheap_scavenge_count          int64 `bson:"pageheap_scavenge_count"`
	Pageheap_commit_count            int64 `bson:"pageheap_commit_count"`
	Pageheap_total_commit_bytes      int64 `bson:"pageheap_total_commit_bytes"`
	Pageheap_decommit_count          int64 `bson:"pageheap_decommit_count"`
	Pageheap_total_decommit_bytes    int64 `bson:"pageheap_total_decommit_bytes"`
	Pageheap_reserve_count           int64 `bson:"pageheap_reserve_count"`
	Pageheap_total_reserve_bytes     int64 `bson:"pageheap_total_reserve_bytes"`
	Spinlock_total_delay_ns          int64 `bson:"spinlock_total_delay_ns"`
	Release_rate                     int64 `bson:"release_rate"`
}
