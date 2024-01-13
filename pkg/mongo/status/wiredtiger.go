package status

// ServerStatusWiredTiger 变量的命名与__wt_connection_stats保持一致
type ServerStatusWiredTiger struct {
	Uri         string                           `bson:"uri"`
	Async       ServerStatusWiredTigerAsync      `bson:"async"`
	Connection  ServerStatusWiredTigerConnection `bson:"connection"`
	DataHandle  ServerStatusWiredTigerDataHandle `bson:"data-handle"`
	Lock        ServerStatusWiredTigerLock       `bson:"lock"`
	Transaction ServerStatusWiredTigerTxn        `bson:"transaction"`
}

type ServerStatusWiredTigerAsync struct {
	Async_cur_queue  int64 `bson:"current work queue length"`
	Async_max_queue  int64 `bson:"maximum work queue length"`
	Async_alloc_race int64 `bson:"number of allocation state races"`
	Async_flush      int64 `bson:"number of flush calls"`
	Async_alloc_view int64 `bson:"number of operation slots viewed for allocation"`
	Async_full       int64 `bson:"number of times operation allocation failed"`
	Async_nowork     int64 `bson:"number of times worker found no work"`
	Async_op_alloc   int64 `bson:"total allocations"`
	Async_op_compact int64 `bson:"total compact calls"`
	Async_op_insert  int64 `bson:"total insert calls"`
	Async_op_remove  int64 `bson:"total remove calls"`
	Async_op_search  int64 `bson:"total search calls"`
	Async_op_update  int64 `bson:"total update calls"`
}

type ServerStatusWiredTigerConnection struct {
	Cond_auto_wait_reset int64 `bson:"auto adjusting condition resets"`
	Cond_auto_wait       int64 `bson:"auto adjusting condition wait calls"`
	Time_travel          int64 `bson:"detected system time went backwards"`
	File_open            int64 `bson:"files currently open"`
	Buckets_dh           int64 `bson:"hash bucket array size for data handles"`
	Buckets              int64 `bson:"hash bucket array size general"`
	Memory_allocation    int64 `bson:"memory allocations"`
	Memory_free          int64 `bson:"memory frees"`
	Memory_grow          int64 `bson:"memory re-allocations"`
	Cond_wait            int64 `bson:"pthread mutex condition wait calls"`
	Rwlock_read          int64 `bson:"pthread mutex shared lock read-lock calls"`
	Rwlock_write         int64 `bson:"pthread mutex shared lock write-lock calls"`
	Fsync_io             int64 `bson:"total fsync I/Os"`
	Read_io              int64 `bson:"total read I/Os"`
	Write_io             int64 `bson:"total write I/Os"`
}

type ServerStatusWiredTigerDataHandle struct {
	Dh_conn_handle_size  int64 `bson:"connection data handle size"`
	Dh_conn_handle_count int64 `bson:"connection data handles currently active"`
	Dh_sweep_ref         int64 `bson:"connection sweep candidate became referenced"`
	Dh_sweep_close       int64 `bson:"connection sweep dhandles closed"`
	Dh_sweep_remove      int64 `bson:"connection sweep dhandles removed from hash list"`
	Dh_sweep_tod         int64 `bson:"connection sweep time-of-death sets"`
	Dh_sweeps            int64 `bson:"connection sweeps"`
	Dh_sweep_skip_ckpt   int64 `bson:"connection sweeps skipped due to checkpoint gathering handles"`
	Dh_session_handles   int64 `bson:"session dhandles swept"`
	Dh_session_sweeps    int64 `bson:"session sweep attempts"`
}

type ServerStatusWiredTigerLock struct {
	Lock_checkpoint_count                   int64 `bson:"checkpoint lock acquisitions"`
	Lock_checkpoint_wait_application        int64 `bson:"checkpoint lock application thread wait time (usecs)"`
	Lock_checkpoint_wait_internal           int64 `bson:"checkpoint lock internal thread wait time (usecs)"`
	Lock_dhandle_wait_application           int64 `bson:"dhandle lock application thread time waiting (usecs)"`
	Lock_dhandle_wait_internal              int64 `bson:"dhandle lock internal thread time waiting (usecs)"`
	Lock_dhandle_read_count                 int64 `bson:"dhandle read lock acquisitions"`
	Lock_dhandle_write_count                int64 `bson:"dhandle write lock acquisitions"`
	Lock_durable_timestamp_wait_application int64 `bson:"durable timestamp queue lock application thread time waiting (usecs)"`
	Lock_durable_timestamp_wait_internal    int64 `bson:"durable timestamp queue lock internal thread time waiting (usecs)"`
	Lock_durable_timestamp_read_count       int64 `bson:"durable timestamp queue read lock acquisitions"`
	Lock_durable_timestamp_write_count      int64 `bson:"durable timestamp queue write lock acquisitions"`
	Lock_metadata_count                     int64 `bson:"metadata lock acquisitions"`
	Lock_metadata_wait_application          int64 `bson:"metadata lock application thread wait time (usecs)"`
	Lock_metadata_wait_internal             int64 `bson:"metadata lock internal thread wait time (usecs)"`
	Lock_read_timestamp_wait_application    int64 `bson:"read timestamp queue lock application thread time waiting (usecs)"`
	Lock_read_timestamp_wait_internal       int64 `bson:"read timestamp queue lock internal thread time waiting (usecs)"`
	Lock_read_timestamp_read_count          int64 `bson:"read timestamp queue read lock acquisitions"`
	Lock_read_timestamp_write_count         int64 `bson:"read timestamp queue write lock acquisitions"`
	Lock_schema_count                       int64 `bson:"schema lock acquisitions"`
	Lock_schema_wait_application            int64 `bson:"schema lock application thread wait time (usecs)"`
	Lock_schema_wait_internal               int64 `bson:"schema lock internal thread wait time (usecs)"`
	Lock_table_wait_application             int64 `bson:"table lock application thread time waiting for the table lock (usecs)"`
	Lock_table_wait_internal                int64 `bson:"table lock internal thread time waiting for the table lock (usecs)"`
	Lock_table_read_count                   int64 `bson:"table read lock acquisitions"`
	Lock_table_write_count                  int64 `bson:"table write lock acquisitions"`
	Lock_txn_global_wait_application        int64 `bson:"txn global lock application thread time waiting (usecs)"`
	Lock_txn_global_wait_internal           int64 `bson:"txn global lock internal thread time waiting (usecs)"`
	Lock_txn_global_read_count              int64 `bson:"txn global read lock acquisitions"`
	Lock_txn_global_write_count             int64 `bson:"txn global write lock acquisitions"`
}

type ServerStatusWiredTigerTxn struct {
	Txn_prepared_updates_count             int64 `bson:"Number of prepared updates"`
	Txn_prepared_updates_lookaside_inserts int64 `bson:"Number of prepared updates added to cache overflow"`
	Txn_durable_queue_walked               int64 `bson:"durable timestamp queue entries walked"`
	Txn_durable_queue_empty                int64 `bson:"durable timestamp queue insert to empty"`
	Txn_durable_queue_head                 int64 `bson:"durable timestamp queue inserts to head"`
	Txn_durable_queue_inserts              int64 `bson:"durable timestamp queue inserts total"`
	Txn_durable_queue_len                  int64 `bson:"durable timestamp queue length"`
	Txn_snapshots_created                  int64 `bson:"number of named snapshots created"`
	Txn_snapshots_dropped                  int64 `bson:"number of named snapshots dropped"`
	Txn_prepare                            int64 `bson:"prepared transactions"`
	Txn_prepare_commit                     int64 `bson:"prepared transactions committed"`
	Txn_prepare_active                     int64 `bson:"prepared transactions currently active"`
	Txn_prepare_rollback                   int64 `bson:"prepared transactions rolled back"`
	Txn_query_ts                           int64 `bson:"query timestamp calls"`
	Txn_read_queue_walked                  int64 `bson:"read timestamp queue entries walked"`
	Txn_read_queue_empty                   int64 `bson:"read timestamp queue insert to empty"`
	Txn_read_queue_head                    int64 `bson:"read timestamp queue inserts to head"`
	Txn_read_queue_inserts                 int64 `bson:"read timestamp queue inserts total"`
	Txn_read_queue_len                     int64 `bson:"read timestamp queue length"`
	Txn_rollback_to_stable                 int64 `bson:"rollback to stable calls"`
	Txn_rollback_upd_aborted               int64 `bson:"rollback to stable updates aborted"`
	Txn_rollback_las_removed               int64 `bson:"rollback to stable updates removed from cache overflow"`
	Txn_set_ts                             int64 `bson:"set timestamp calls"`
	Txn_set_ts_durable                     int64 `bson:"set timestamp durable calls"`
	Txn_set_ts_durable_upd                 int64 `bson:"set timestamp durable updates"`
	Txn_set_ts_oldest                      int64 `bson:"set timestamp oldest calls"`
	Txn_set_ts_oldest_upd                  int64 `bson:"set timestamp oldest updates"`
	Txn_set_ts_stable                      int64 `bson:"set timestamp stable calls"`
	Txn_set_ts_stable_upd                  int64 `bson:"set timestamp stable updates"`
	Txn_begin                              int64 `bson:"transaction begins"`
	Txn_checkpoint_running                 int64 `bson:"transaction checkpoint currently running"`
	Txn_checkpoint_generation              int64 `bson:"transaction checkpoint generation"`
	Txn_checkpoint_time_max                int64 `bson:"transaction checkpoint max time (msecs)"`
	Txn_checkpoint_time_min                int64 `bson:"transaction checkpoint min time (msecs)"`
	Txn_checkpoint_handle_duration         int64 `bson:"transaction checkpoint most recent duration for gathering all handles (usecs)"`
	Txn_checkpoint_handle_duration_apply   int64 `bson:"transaction checkpoint most recent duration for gathering applied handles (usecs)"`
	Txn_checkpoint_handle_duration_skip    int64 `bson:"transaction checkpoint most recent duration for gathering skipped handles (usecs)"`
	Txn_checkpoint_handle_applied          int64 `bson:"transaction checkpoint most recent handles applied"`
	Txn_checkpoint_handle_skipped          int64 `bson:"transaction checkpoint most recent handles skipped"`
	Txn_checkpoint_handle_walked           int64 `bson:"transaction checkpoint most recent handles walked"`
	Txn_checkpoint_time_recent             int64 `bson:"transaction checkpoint most recent time (msecs)"`
	Txn_checkpoint_scrub_target            int64 `bson:"transaction checkpoint scrub dirty target"`
	Txn_checkpoint_scrub_time              int64 `bson:"transaction checkpoint scrub time (msecs)"`
	Txn_checkpoint_time_total              int64 `bson:"transaction checkpoint total time (msecs)"`
	Txn_checkpoint                         int64 `bson:"transaction checkpoints"`
	Txn_checkpoint_skipped                 int64 `bson:"transaction checkpoints skipped because database was clean"`
	Txn_fail_cache                         int64 `bson:"transaction failures due to cache overflow"`
	Txn_checkpoint_fsync_post              int64 `bson:"transaction fsync calls for checkpoint after allocating the transaction ID"`
	Txn_checkpoint_fsync_post_duration     int64 `bson:"transaction fsync duration for checkpoint after allocating the transaction ID (usecs)"`
	Txn_pinned_range                       int64 `bson:"transaction range of IDs currently pinned"`
	Txn_pinned_checkpoint_range            int64 `bson:"transaction range of IDs currently pinned by a checkpoint"`
	Txn_pinned_snapshot_range              int64 `bson:"transaction range of IDs currently pinned by named snapshots"`
	Txn_pinned_timestamp                   int64 `bson:"transaction range of timestamps currently pinned"`
	Txn_pinned_timestamp_checkpoint        int64 `bson:"transaction range of timestamps pinned by a checkpoint"`
	Txn_pinned_timestamp_reader            int64 `bson:"transaction range of timestamps pinned by the oldest active read timestamp"`
	Txn_pinned_timestamp_oldest            int64 `bson:"transaction range of timestamps pinned by the oldest timestamp"`
	Txn_timestamp_oldest_active_read       int64 `bson:"transaction read timestamp of the oldest active reader"`
	Txn_sync                               int64 `bson:"transaction sync calls"`
	Txn_commit                             int64 `bson:"transactions committed"`
	Txn_rollback                           int64 `bson:"transactions rolled back"`
	Txn_update_conflict                    int64 `bson:"update conflicts"`
}
