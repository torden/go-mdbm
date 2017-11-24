package mdbm

/*
#cgo CFLAGS: -I/tmp/install/include/
#cgo LDFLAGS: -L/tmp/install/lib64/ -Wl,-rpath=/tmp/install/lib64/ -lmdbm

#include <stdlib.h>
#include <stdio.h>
#include <sys/syslog.h>
#include <unistd.h>
#include <string.h>
#include <time.h>
#include <mdbm.h>
#include <mdbm_log.h>

static void get_mdbm_iter(MDBM_ITER *iter) {
	MDBM_ITER_INIT(iter);
}

static mdbm_ubig_t get_pageno_of_mdbm_iter(MDBM_ITER *iter) {
	return (*iter).m_pageno;
}

static int get_next_of_mdbm_iter(MDBM_ITER *iter) {
	return (*iter).m_next;
}

//debug
static void echo_str(const char *str) {
	fprintf(stderr,"----- echo_str -----\n");
	fprintf(stderr,"[%s]\n", str);
	fprintf(stderr,"--------------------\n");
}

// protection to raise the locking exception at times
static int set_mdbm_store_with_lock(MDBM *db, datum key, datum val, int flags) {
	int rv;
	rv = mdbm_lock(db);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store(db,key,val,flags);
	mdbm_unlock(db);
	return rv;
}

static int set_mdbm_store_r_with_lock(MDBM *db, datum *key, datum *val, int flags, MDBM_ITER *iter) {
	int rv;
	rv = mdbm_lock(db);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store_r(db,key,val,flags, iter);
	mdbm_unlock(db);
	return rv;
}

static int set_mdbm_store_str_with_lock(MDBM *db, const char *key, const char *val, int flags) {

	int rv;
	rv = mdbm_lock(db);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store_str(db, key, val, flags);
	mdbm_unlock(db);
	return rv;
}

static int set_mdbm_store_with_lock_smart(MDBM *db, datum key, datum val, int flags, int lockflags) {

	int rv;
	rv = mdbm_lock_smart(db, &key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store(db, key, val, flags);
	mdbm_unlock_smart(db, &key, lockflags);
	return rv;
}

static int set_mdbm_store_r_with_lock_smart(MDBM *db, datum *key, datum *val, int flags, int lockflags, MDBM_ITER *iter) {

	int rv;
	rv = mdbm_lock_smart(db, key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store_r(db, key, val, flags, iter);
	mdbm_unlock_smart(db, key, lockflags);
	return rv;
}

static int set_mdbm_store_with_plock(MDBM *db, datum key, datum val, int flags, int lockflags) {

	int rv;
	rv = mdbm_plock(db, &key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store(db, key, val, flags);
	mdbm_punlock(db, &key, lockflags);
	return rv;
}

static int set_mdbm_store_r_with_plock(MDBM *db, datum *key, datum *val, int flags, int lockflags, MDBM_ITER *iter) {

	int rv;
	rv = mdbm_plock(db, key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store_r(db, key, val, flags, iter);
	mdbm_punlock(db, key, lockflags);
	return rv;
}

static char *get_mdbm_fetch_with_lock(MDBM *db, datum key) {

	int rv;
	datum val;
	char *buf = NULL;

	rv = mdbm_lock(db);
	if(rv != 1) {
		return NULL;
	}

	val = mdbm_fetch(db, key);
	if (val.dptr != NULL) {
		buf = (char *)calloc(val.dsize+1, sizeof(char));
		if (buf == NULL) {
			perror("failed allocation:");
		} else {
			strncpy(buf, val.dptr, val.dsize);
		}
	}

	mdbm_unlock(db);
	return buf;
}

static char *get_mdbm_fetch_with_lock_smart(MDBM *db, datum key, int lockflags) {

	int rv;
	datum val;
	char *buf = NULL;

	rv = mdbm_lock_smart(db, &key, lockflags);
	if(rv != 1) {
		return NULL;
	}

	val = mdbm_fetch(db, key);
	if (val.dptr != NULL) {
		buf = (char *)calloc(val.dsize+1, sizeof(char));
		if (buf == NULL) {
			perror("failed allocation:");
		} else {
			strncpy(buf, val.dptr, val.dsize);
		}
	}

	mdbm_unlock_smart(db, &key, lockflags);
	return buf;
}

static char *get_mdbm_fetch_with_plock(MDBM *db, datum key, int lockflags) {

	int rv;
	datum val;
	char *buf = NULL;

	rv = mdbm_plock(db, &key, lockflags);
	if(rv != 1) {
		return NULL;
	}

	val = mdbm_fetch(db, key);
	if (val.dptr != NULL) {
		buf = (char *)calloc(val.dsize+1, sizeof(char));
		if (buf == NULL) {
			perror("failed allocation:");
		} else {
			strncpy(buf, val.dptr, val.dsize);
		}
	}

	mdbm_punlock(db, &key, lockflags);
	return buf;
}

static int get_mdbm_fetch_r_with_lock(MDBM *db, datum *key, datum *val, MDBM_ITER *iter) {

	int rv;

	rv = mdbm_lock(db);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_fetch_r(db, key, val, iter);
	mdbm_unlock(db);
	return rv;
}

static int get_mdbm_fetch_r_with_lock_smart(MDBM *db, datum *key, datum *val, int lockflags, MDBM_ITER *iter) {

	int rv;

	rv = mdbm_lock_smart(db, key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_fetch_r(db, key, val, iter);
	mdbm_unlock_smart(db, key, lockflags);
	return rv;
}

static int get_mdbm_fetch_r_with_plock(MDBM *db, datum *key, datum *val, int lockflags, MDBM_ITER *iter) {

	int rv;

	rv = mdbm_plock(db, key, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_fetch_r(db, key, val, iter);
	mdbm_punlock(db, key, lockflags);
	return rv;
}
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
)

// MDBM represents the MDBM Go's methods
type MDBM struct {
	pmdbm    *C.MDBM
	iter     C.MDBM_ITER
	locked   bool
	isopened bool
	dbmfile  string
	pdbmfile *C.char
	flags    int
	perms    int
	psize    int
	dsize    int
	mutex    sync.RWMutex
	cgomtx   sync.RWMutex

	scpagesize  uint
	minpagesize uint
	//_ struct{}
}

// Iter represents the mdbm_iter struct of the MDBM
type Iter struct { //mdbm_iter
	PageNo uint32 //mdbm_ubig_t     uint32_t
	Next   int
	//_      struct{}
}

// FetchInfo represents the mdbm_fetch_info struct of the MDBM
type FetchInfo struct { //mdbm_fetch_info
	Flags            uint32 //Entry flags
	CacheNumAccesses uint32 //Number of accesses to cache entry
	CacheAccessTime  uint32 //Last access time (LRU/LFU only)
	//_ struct{}
}

// Stats represents the mdbm_stats_t struct of the MDBM
type Stats struct { //mdbm_stats_t
	Size                uint32
	PageSize            uint32
	PageCount           uint32
	PagesUsed           uint32
	BytesUsed           uint32
	NumEntries          uint32
	MinLevel            uint32
	MaxLevel            uint32
	LargePageSize       uint32
	LargePageCount      uint32
	LargeThreshold      uint32
	LargePagesUsed      uint32
	LargeNumFreeEntries uint32
	LargeMaxFree        uint32
	LargeNumEntries     uint32
	LargeBytesUsed      uint32
	LargeMinSize        uint32
	LargeMaxSize        uint32
	CacheMode           uint32
}

// DBInfo represents the mdbm_db_info_t struct of the MDBM
type DBInfo struct { //mdbm_db_info_t
	PageSize     uint32
	NumPages     uint32
	MaxPages     uint32
	NumDirPages  uint32
	DirWidth     uint32
	MaxDirShift  uint32
	DirMinLevel  uint32
	DirMaxLevel  uint32
	DirNumNodes  uint32
	HashFunc     uint32
	HashFuncName string
	SpillSize    uint32
	CacheMode    uint32
}

// BucketStat represents the mdbm_bucket_stat_t struct of the MDBM
type BucketStat struct { //mdbm_bucket_stat_t
	NumPages     uint32
	MinBytes     uint32
	MaxBytes     uint32
	MinFreeBytes uint32
	MaxFreeBytes uint32
	SumEntries   uint64
	SumBytes     uint64
	SumFreeBytes uint64
}

// StatInfo represents the mdbm_stat_info_t struct of the MDBM
type StatInfo struct { //mdbm_stat_info_t
	Flags               int
	NumActiveEntries    uint64
	NumActiveLobEntries uint64
	SumKeyBytes         uint64
	SumLobValBytes      uint64
	SumNormalValBytes   uint64
	SumOverheadBytes    uint64
	MinEntryBytes       uint32
	MaxEntryBytes       uint32
	MinKeyBytes         uint32
	MaxKeyBytes         uint32
	MinValBytes         uint32
	MaxValBytes         uint32
	MinLobBytes         uint32
	MaxLobBytes         uint32
	MaxPageUsedSpace    uint32
	MaxDataPages        uint32
	NumFreePages        uint32
	NumActivePages      uint32
	NumNormalPages      uint32
	NumOversizedPages   uint32
	NumLobPages         uint32
	Buckets             []BucketStat
	MaxPageEntries      uint32
	MinPageEntries      uint32
}

// WindowStats represents the mdbm_window_stats_t struct of the MDBM
type WindowStats struct { //mdbm_window_stats_t
	WnumReused     uint64
	WnumRemapped   uint64
	WwindowSize    uint32
	WmaxWindowUsed uint32
}

// NewMDBM creates and returns the MDBM methods's pointer and initialized a MDBM struct
func NewMDBM() *MDBM {

	obj := &MDBM{
		locked:   false,
		dbmfile:  "",
		flags:    Create | Rdrw,
		perms:    0666,
		psize:    0,
		dsize:    0,
		isopened: false,
	}

	obj.scpagesize = uint(C.sysconf(C._SC_PAGESIZE))
	obj.minpagesize = obj.scpagesize * 2
	obj.iter = obj.GetNewIter()
	return obj
}

// convertIter returns a data of the MDBM_ITER convert to Iter
func (db *MDBM) convertIter(iter *C.MDBM_ITER) Iter {
	return Iter{
		PageNo: uint32(C.get_pageno_of_mdbm_iter(iter)),
		Next:   int(C.get_next_of_mdbm_iter(iter)),
	}
}

// convertIterToC returns a data of the Iter convert to MDBM_ITER
func (db *MDBM) convertIterToC(iter Iter) C.MDBM_ITER {

	var rv C.MDBM_ITER
	rv.m_pageno = C.uint32_t(iter.PageNo)
	rv.m_next = C.int(iter.Next)
	return rv
}

// convertFetchInfo returns a data of the mdbm_fetch_info convert to FetchInfo
func (db *MDBM) convertFetchInfo(info C.struct_mdbm_fetch_info) FetchInfo {
	return FetchInfo{
		Flags:            uint32(info.flags),
		CacheNumAccesses: uint32(info.cache_num_accesses),
		CacheAccessTime:  uint32(info.cache_access_time),
	}
}

// convertStats returns a data of the mdbm_stats_t convert to Stats
func (db *MDBM) convertStats(stats C.mdbm_stats_t) Stats {

	return Stats{
		Size:                uint32(stats.s_size),
		PageSize:            uint32(stats.s_page_size),
		PageCount:           uint32(stats.s_page_count),
		PagesUsed:           uint32(stats.s_pages_used),
		BytesUsed:           uint32(stats.s_bytes_used),
		NumEntries:          uint32(stats.s_num_entries),
		MinLevel:            uint32(stats.s_min_level),
		MaxLevel:            uint32(stats.s_max_level),
		LargePageSize:       uint32(stats.s_large_page_size),
		LargePageCount:      uint32(stats.s_large_page_count),
		LargeThreshold:      uint32(stats.s_large_threshold),
		LargePagesUsed:      uint32(stats.s_large_pages_used),
		LargeNumFreeEntries: uint32(stats.s_large_num_free_entries),
		LargeMaxFree:        uint32(stats.s_large_max_free),
		LargeNumEntries:     uint32(stats.s_large_num_entries),
		LargeBytesUsed:      uint32(stats.s_large_bytes_used),
		LargeMinSize:        uint32(stats.s_large_min_size),
		LargeMaxSize:        uint32(stats.s_large_max_size),
		CacheMode:           uint32(stats.s_cache_mode),
	}
}

// convertStatsToC returns a data of the conver Stats to mdbm_stats_t
func (db *MDBM) convertStatsToC(stats Stats) C.mdbm_stats_t {

	var rv C.mdbm_stats_t

	rv.s_size = C.mdbm_ubig_t(stats.Size)
	rv.s_page_size = C.mdbm_ubig_t(stats.PageSize)
	rv.s_page_count = C.mdbm_ubig_t(stats.PageCount)
	rv.s_pages_used = C.mdbm_ubig_t(stats.PagesUsed)
	rv.s_bytes_used = C.mdbm_ubig_t(stats.BytesUsed)
	rv.s_num_entries = C.mdbm_ubig_t(stats.NumEntries)
	rv.s_min_level = C.mdbm_ubig_t(stats.MinLevel)
	rv.s_max_level = C.mdbm_ubig_t(stats.MaxLevel)
	rv.s_large_page_size = C.mdbm_ubig_t(stats.LargePageSize)
	rv.s_large_page_count = C.mdbm_ubig_t(stats.LargePageCount)
	rv.s_large_threshold = C.mdbm_ubig_t(stats.LargeThreshold)
	rv.s_large_pages_used = C.mdbm_ubig_t(stats.LargePagesUsed)
	rv.s_large_num_free_entries = C.mdbm_ubig_t(stats.LargeNumFreeEntries)
	rv.s_large_max_free = C.mdbm_ubig_t(stats.LargeMaxFree)
	rv.s_large_num_entries = C.mdbm_ubig_t(stats.LargeNumEntries)
	rv.s_large_bytes_used = C.mdbm_ubig_t(stats.LargeBytesUsed)
	rv.s_large_min_size = C.mdbm_ubig_t(stats.LargeMinSize)
	rv.s_large_max_size = C.mdbm_ubig_t(stats.LargeMaxSize)
	rv.s_cache_mode = C.mdbm_ubig_t(stats.CacheMode)

	return rv
}

// convertDBInfo returns a data of the mdbm_db_info_t convert to DBInfo
func (db *MDBM) convertDBInfo(info C.mdbm_db_info_t) DBInfo {

	return DBInfo{
		PageSize:     uint32(info.db_page_size),
		NumPages:     uint32(info.db_num_pages),
		MaxPages:     uint32(info.db_max_pages),
		NumDirPages:  uint32(info.db_num_dir_pages),
		DirWidth:     uint32(info.db_dir_width),
		MaxDirShift:  uint32(info.db_max_dir_shift),
		DirMinLevel:  uint32(info.db_dir_min_level),
		DirMaxLevel:  uint32(info.db_dir_max_level),
		DirNumNodes:  uint32(info.db_dir_num_nodes),
		HashFunc:     uint32(info.db_hash_func),
		HashFuncName: C.GoString(info.db_hash_funcname),
		SpillSize:    uint32(info.db_spill_size),
		CacheMode:    uint32(info.db_cache_mode),
	}
}

// convertDBInfoToC returns a data of the DBInfo convert to mdbm_db_info_t
func (db *MDBM) convertDBInfoToC(info DBInfo) C.mdbm_db_info_t {

	var rv C.mdbm_db_info_t
	rv.db_page_size = C.uint32_t(info.PageSize)
	rv.db_num_pages = C.uint32_t(info.NumPages)
	rv.db_max_pages = C.uint32_t(info.MaxPages)
	rv.db_num_dir_pages = C.uint32_t(info.NumDirPages)
	rv.db_dir_width = C.uint32_t(info.DirWidth)
	rv.db_max_dir_shift = C.uint32_t(info.MaxDirShift)
	rv.db_dir_min_level = C.uint32_t(info.DirMinLevel)
	rv.db_dir_max_level = C.uint32_t(info.DirMaxLevel)
	rv.db_dir_num_nodes = C.uint32_t(info.DirNumNodes)
	rv.db_hash_func = C.uint32_t(info.HashFunc)
	rv.db_hash_funcname = C.CString(info.HashFuncName)
	rv.db_spill_size = C.uint32_t(info.SpillSize)
	rv.db_cache_mode = C.uint32_t(info.CacheMode)
	return rv
}

// convertBucketStat returns a data of the mdbm_bucket_stat_t convert to BucketStat
func (db *MDBM) convertBucketStat(bucket C.mdbm_bucket_stat_t) BucketStat {

	return BucketStat{
		NumPages:     uint32(bucket.num_pages),
		MinBytes:     uint32(bucket.min_bytes),
		MaxBytes:     uint32(bucket.max_bytes),
		MinFreeBytes: uint32(bucket.min_free_bytes),
		MaxFreeBytes: uint32(bucket.max_free_bytes),
		SumEntries:   uint64(bucket.sum_entries),
		SumBytes:     uint64(bucket.sum_bytes),
		SumFreeBytes: uint64(bucket.sum_free_bytes),
	}
}

// convertBucketStatToC returns a data of the BucketStat convert to mdbm_bucket_stat_t
func (db *MDBM) convertBucketStatToC(bucket BucketStat) C.mdbm_bucket_stat_t {

	var rv C.mdbm_bucket_stat_t
	rv.num_pages = C.uint32_t(bucket.NumPages)
	rv.min_bytes = C.uint32_t(bucket.MinBytes)
	rv.max_bytes = C.uint32_t(bucket.MaxBytes)
	rv.min_free_bytes = C.uint32_t(bucket.MinFreeBytes)
	rv.max_free_bytes = C.uint32_t(bucket.MaxFreeBytes)
	rv.sum_entries = C.uint64_t(bucket.SumEntries)
	rv.sum_bytes = C.uint64_t(bucket.SumBytes)
	rv.sum_free_bytes = C.uint64_t(bucket.SumFreeBytes)

	return rv
}

// convertStatInfo returns a data of the mdbm_stat_info_t convert to StatInfo
func (db *MDBM) convertStatInfo(info C.mdbm_stat_info_t) StatInfo {

	var buckets []BucketStat

	bucketLen := len(info.buckets)
	for i := 0; i < bucketLen; i++ {
		buckets = append(buckets, db.convertBucketStat(info.buckets[i]))
	}

	return StatInfo{
		Flags:               int(info.flags),
		NumActiveEntries:    uint64(info.num_active_entries),
		NumActiveLobEntries: uint64(info.num_active_lob_entries),
		SumKeyBytes:         uint64(info.sum_key_bytes),
		SumLobValBytes:      uint64(info.sum_lob_val_bytes),
		SumNormalValBytes:   uint64(info.sum_normal_val_bytes),
		SumOverheadBytes:    uint64(info.sum_overhead_bytes),
		MinEntryBytes:       uint32(info.min_entry_bytes),
		MaxEntryBytes:       uint32(info.max_entry_bytes),
		MinKeyBytes:         uint32(info.min_key_bytes),
		MaxKeyBytes:         uint32(info.max_key_bytes),
		MinValBytes:         uint32(info.min_val_bytes),
		MaxValBytes:         uint32(info.max_val_bytes),
		MinLobBytes:         uint32(info.min_lob_bytes),
		MaxLobBytes:         uint32(info.max_lob_bytes),
		MaxPageUsedSpace:    uint32(info.max_page_used_space),
		MaxDataPages:        uint32(info.max_data_pages),
		NumFreePages:        uint32(info.num_free_pages),
		NumActivePages:      uint32(info.num_active_pages),
		NumNormalPages:      uint32(info.num_normal_pages),
		NumOversizedPages:   uint32(info.num_oversized_pages),
		NumLobPages:         uint32(info.num_lob_pages),
		Buckets:             buckets,
		MaxPageEntries:      uint32(info.max_page_entries),
		MinPageEntries:      uint32(info.min_page_entries),
	}
}

// convertStatInfoToC returns a data of the StatInfo to mdbm_stat_info_t
// NOT COMPLETE
func (db *MDBM) convertStatInfoToC(info StatInfo) C.mdbm_stat_info_t {

	var rv C.mdbm_stat_info_t

	rv.flags = C.int(info.Flags)
	rv.num_active_entries = C.uint64_t(info.NumActiveEntries)
	rv.num_active_lob_entries = C.uint64_t(info.NumActiveLobEntries)
	rv.sum_key_bytes = C.uint64_t(info.SumKeyBytes)
	rv.sum_lob_val_bytes = C.uint64_t(info.SumLobValBytes)
	rv.sum_normal_val_bytes = C.uint64_t(info.SumNormalValBytes)
	rv.sum_overhead_bytes = C.uint64_t(info.SumOverheadBytes)
	rv.min_entry_bytes = C.uint32_t(info.MinEntryBytes)
	rv.max_entry_bytes = C.uint32_t(info.MaxEntryBytes)
	rv.min_key_bytes = C.uint32_t(info.MinKeyBytes)
	rv.max_key_bytes = C.uint32_t(info.MaxKeyBytes)
	rv.min_val_bytes = C.uint32_t(info.MinValBytes)
	rv.max_val_bytes = C.uint32_t(info.MaxValBytes)
	rv.min_lob_bytes = C.uint32_t(info.MinLobBytes)
	rv.max_lob_bytes = C.uint32_t(info.MaxLobBytes)
	rv.max_page_used_space = C.uint32_t(info.MaxPageUsedSpace)
	rv.max_data_pages = C.uint32_t(info.MaxDataPages)
	rv.num_free_pages = C.uint32_t(info.NumFreePages)
	rv.num_active_pages = C.uint32_t(info.NumActivePages)
	rv.num_normal_pages = C.uint32_t(info.NumNormalPages)
	rv.num_oversized_pages = C.uint32_t(info.NumOversizedPages)
	rv.num_lob_pages = C.uint32_t(info.NumLobPages)
	//rv.buckets = C.mdbm_buck(info.//Buckets           )
	rv.max_page_entries = C.uint32_t(info.MaxPageEntries)
	rv.min_page_entries = C.uint32_t(info.MinPageEntries)

	return rv
}

// convertWindowStat returns a data of the mdbm_window_stats_t convert to WindowStats
func (db *MDBM) convertWindowStat(ws C.mdbm_window_stats_t) WindowStats {

	return WindowStats{
		WnumReused:     uint64(ws.w_num_reused),
		WnumRemapped:   uint64(ws.w_num_remapped),
		WwindowSize:    uint32(ws.w_window_size),
		WmaxWindowUsed: uint32(ws.w_max_window_used),
	}
}

// convertWindowStatToC return a data of the WindowStats convert to mdbm_window_stats_t
func (db *MDBM) convertWindowStatToC(ws WindowStats) C.mdbm_window_stats_t {

	var rv C.mdbm_window_stats_t
	rv.w_num_reused = C.uint64_t(ws.WnumReused)
	rv.w_num_remapped = C.uint64_t(ws.WnumRemapped)
	rv.w_window_size = C.uint32_t(ws.WwindowSize)
	rv.w_max_window_used = C.uint32_t(ws.WmaxWindowUsed)
	return rv
}

func (db *MDBM) cgoRun(call func() (int, error)) (int, string, error) {

	db.cgomtx.Lock()
	defer db.cgomtx.Unlock()

	orgStdErr, orgStdOut := os.Stderr, os.Stdout
	orgCStdErr, orgCStdOut := C.stderr, C.stdout
	defer func() {
		os.Stderr, os.Stdout = orgStdErr, orgStdOut
		C.stderr, C.stdout = orgCStdErr, orgCStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed create a os.Pipe()")
	}

	defer func() {
		e := r.Close()
		if e != nil {
			err = e
		}
	}()

	cw := C.CString("w")
	defer C.free(unsafe.Pointer(cw))

	f := C.fdopen((C.int)(w.Fd()), cw)
	if f == nil {
		return 0, "", errors.New("failed call to mdbm::cgoRun(C.fdopen)")
	}
	defer C.fclose(f)

	os.Stderr, os.Stdout = w, w
	C.stderr, C.stdout = f, f

	out := make(chan []byte)
	go func() {

		var b bytes.Buffer
		_, rerr := io.Copy(&b, r)
		if rerr != nil {
			out <- []byte(fmt.Sprintf("failed can't capture cgo output : %s", rerr.Error()))
		} else {
			out <- b.Bytes()
		}
		runtime.Gosched()
	}()

	//run
	rv, err := call()

	C.fflush(f)
	cerr := w.Close()
	if cerr != nil {
		log.Println("failred close the os.Stderr, os.Stdout")
	}

	return rv, string(<-out), err
}

// convertToArByte returns a data of the any data convert ot Byte Array
// and returns error at raise the exception
func (db *MDBM) convertToArByte(obj interface{}) ([]byte, error) {

	switch obj.(type) {

	case bool:
		if obj.(bool) {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case byte:
		return []byte{obj.(byte)}, nil

	case []uint8:
		return reflect.ValueOf(obj).Bytes(), nil

	case string:
		return []byte(obj.(string)), nil

	case int:
		return []byte(strconv.FormatInt(int64(obj.(int)), 10)), nil
	case int8:
		return []byte(strconv.FormatInt(int64(obj.(int8)), 10)), nil
	case int16:
		return []byte(strconv.FormatInt(int64(obj.(int16)), 10)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(obj.(int32)), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(obj.(int64), 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(obj.(uint)), 10)), nil
	case uint16:
		return []byte(strconv.FormatUint(uint64(obj.(uint16)), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(obj.(uint32)), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(obj.(uint64), 10)), nil
	case float32:
		return []byte(fmt.Sprintf("%g", obj.(float32))), nil
	case float64:
		return []byte(fmt.Sprintf("%g", obj.(float64))), nil
	case complex64:
		return []byte(fmt.Sprintf("%g", obj.(complex64))), nil
	case complex128:
		return []byte(fmt.Sprintf("%g", obj.(complex128))), nil

	default:
		return nil, fmt.Errorf("not support type(%s)", reflect.TypeOf(obj).String())
	}
}

// convertToString returns a data of the any data convert ot Byte Array
// and returns error at raise the exception
func (db *MDBM) convertToString(obj interface{}) (string, error) {

	switch obj.(type) {

	case bool:
		if obj.(bool) {
			return "true", nil
		}
		return "false", nil

	case byte:
		return string(obj.(byte)), nil

	case []uint8:
		return reflect.ValueOf(obj).String(), nil

	case string:
		return obj.(string), nil

	case int:
		return strconv.FormatInt(int64(obj.(int)), 10), nil
	case int8:
		return strconv.FormatInt(int64(obj.(int8)), 10), nil
	case int16:
		return strconv.FormatInt(int64(obj.(int16)), 10), nil
	case int32:
		return strconv.FormatInt(int64(obj.(int32)), 10), nil
	case int64:
		return strconv.FormatInt(obj.(int64), 10), nil
	case uint:
		return strconv.FormatUint(uint64(obj.(uint)), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(obj.(uint16)), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(obj.(uint32)), 10), nil
	case uint64:
		return strconv.FormatUint(obj.(uint64), 10), nil
	case float32:
		return fmt.Sprintf("%g", obj.(float32)), nil
	case float64:
		return fmt.Sprintf("%g", obj.(float64)), nil
	case complex64:
		return fmt.Sprintf("%g", obj.(complex64)), nil
	case complex128:
		return fmt.Sprintf("%g", obj.(complex128)), nil

	default:
		return "", fmt.Errorf("not support type(%s)", reflect.TypeOf(obj).String())
	}
}

// isVersion3Above checks out what the version is 3 or higher
func (db *MDBM) isVersion3Above() error {

	dbver, err := db.GetVersion()
	if err != nil {
		return errors.Wrapf(err, "failed obtain the MDBM version")
	}

	if dbver < 3 {
		return fmt.Errorf("only supported by MDBM version 3 or higher (current version=%d)", dbver)
	}

	return nil
}

// isVersion3Above checks out what the version is 2
func (db *MDBM) isVersion2() error {

	dbver, err := db.GetVersion()
	if err != nil {
		return errors.Wrapf(err, "failed obtain the MDBM version")
	}

	if dbver != 2 {
		return fmt.Errorf("there is only a v2 implementation. v3 or higher version not currently supported (current version=%d)", dbver)
	}

	return nil
}

// GetNewIter returns a pointer of the C.MDBM_ITer struct
func (db *MDBM) GetNewIter() C.MDBM_ITER {
	var iter C.MDBM_ITER
	C.get_mdbm_iter(&iter)
	return iter
}

// DupHandle returns a pointer of the Duplicate an existing database handle.
// The advantage of dup'ing a handle over doing a separate Open is that dup's handle share the same virtual
// page mapping within the process space (saving memory).
// Threaded applications should use pthread_mutex_lock and unlock around calls to mdbm_dup_handle.
func (db *MDBM) DupHandle() (*MDBM, error) {

	if !db.isopened {
		return nil, errors.New("Not an existing database handle")
	}

	var err error

	db.mutex.Lock()
	defer db.mutex.Unlock()

	obj := &MDBM{
		locked:   db.locked,
		dbmfile:  db.dbmfile,
		flags:    db.flags,
		perms:    db.perms,
		psize:    db.psize,
		dsize:    db.dsize,
		isopened: db.isopened,
	}

	_, _, err = db.cgoRun(func() (int, error) {
		obj.iter = db.GetNewIter()
		obj.pmdbm, err = C.mdbm_dup_handle(db.pmdbm, 0)
		if err != nil {
			return -1, err
		}
		return 0, err
	})

	return obj, err
}

// GetErrNo returns the value of internally saved errno.
// Contains the value of errno that is set during some lock failures.
// Under other circumstances, GetErrNo will not return the actual value of the errno variable.
func (db *MDBM) GetErrNo() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_errno(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LogMinLevel sets the minimum logging level,Lower priority messages are discarded
func (db *MDBM) LogMinLevel(lv int) error {

	if lv < C.LOG_EMERG || lv > C.LOG_DEBUG {
		return errors.New("Not support log level")
	}

	_, _, err := db.cgoRun(func() (int, error) {
		C.mdbm_log_minlevel(C.int(1))
		return 0, nil
	})

	return err
}

// EasyOpen Creates and/or opens the MDBM database use the default options
func (db *MDBM) EasyOpen(dbmfile string, perms int) error {

	var err error

	db.dbmfile = dbmfile
	if perms > 0 {
		db.perms = perms
	}

	db.pdbmfile = C.CString(dbmfile)

	//db.LogMinLevel(LOG_INFO)
	_, _, err = db.cgoRun(func() (int, error) {
		db.pmdbm, err = C.mdbm_open(db.pdbmfile, C.int(db.flags), C.int(db.perms), C.int(db.psize), C.int(db.dsize))
		if db.pmdbm != nil {
			return 0, nil
		}

		return 0, err
	})

	if err == nil {
		db.mutex.RLock()
		{
			db.isopened = true
		}
		db.mutex.RUnlock()
	}

	return err
}

// Open Creates and/or opens the MDBM database
func (db *MDBM) Open(mdbmfn string, flags, perms, psize, dsize int) error {

	db.flags = flags
	db.perms = perms
	db.psize = psize
	db.dsize = dsize

	return db.EasyOpen(mdbmfn, perms)
}

// Sync syncs all pages to disk asynchronously. it's mapped pages are scheduled to be flushed to disk.
func (db *MDBM) Sync() (int, error) {

	var rv int
	var err error

	if db.isopened {
		rv, _, err = db.cgoRun(func() (int, error) {
			rv, err := C.mdbm_sync(db.pmdbm)
			return int(rv), err
		})
	}
	return rv, err
}

// Fsync syncs all pages to disk synchronously. it's will pages have been flushed to disk.
// The database is locked while pages are flushed.
func (db *MDBM) Fsync() (int, error) {

	var rv int
	var err error

	if db.isopened {
		rv, _, err = db.cgoRun(func() (int, error) {
			rv, err := C.mdbm_fsync(db.pmdbm)
			return int(rv), err
		})
	}

	return rv, err
}

// CloseFD closes the MDBM's underlying file descriptor.
func (db *MDBM) CloseFD() error {

	var err error

	if db.isopened {

		_, _, err = db.cgoRun(func() (int, error) {
			_, err := C.mdbm_close_fd(db.pmdbm)
			db.isopened = false
			return 0, err
		})
	}

	return err
}

// Close Closes the database after Sync
func (db *MDBM) Close() {

	if db == nil {
		return
	}

	if db.isopened {

		db.mutex.Lock()
		{
			C.mdbm_close(db.pmdbm)
			db.isopened = false
		}
		db.mutex.Unlock()
	}

	C.free(unsafe.Pointer(db.pdbmfile))
}

// EasyClose Closes the database after Sync
func (db *MDBM) EasyClose() {

	if db == nil {
		return
	}

	if db.isopened {

		rv, err := db.Sync()
		if err != nil {
			log.Printf("failed db.Sync(), rv=%d, err=%v", rv, err)
		}

		db.mutex.Lock()
		{
			C.mdbm_close(db.pmdbm)
			db.isopened = false
		}
		db.mutex.Unlock()
	}
}

// Lock locks the database for exclusive access by the caller.
// The lock is nestable, so a caller already holding
// the lock may call mdbm_lock again as long as an equal number of calls
// to Unlock are made to release the lock.
func (db *MDBM) Lock() error {

	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock(db.pmdbm)

		db.mutex.Lock()
		{
			if !db.locked && int(rv) == 1 {
				db.locked = true
			}
		}
		db.mutex.Unlock()
		return 0, err
	})

	return err
}

// Unlock unlocks the database, releasing exclusive or shared access by the caller.
// If the caller has called Lock() or LockShared() multiple times
// in a row, an equal number of unlock calls are required.
func (db *MDBM) Unlock() error {

	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_unlock(db.pmdbm)
		db.mutex.Lock()
		{
			if db.locked && int(rv) == 1 {
				db.locked = false
			}
		}
		db.mutex.Unlock()
		return 0, err
	})

	return err
}

// TryLock attempts to exclusively lock the MDBM.
func (db *MDBM) TryLock() error {

	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_trylock(db.pmdbm)
		db.mutex.Lock()
		{
			if !db.locked && int(rv) == 1 {
				db.locked = true
			}
		}
		db.mutex.Unlock()

		return 0, err
	})

	return err
}

// IsLocked returns whether or not MDBM is locked by another process or thread.
// rv 0 Database is not locked
// rv 1 Database is locked
func (db *MDBM) IsLocked() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_islocked(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LockShared Locks the database for shared access by readers, excluding access to writers.
// This is multiple-readers, one writer (MROW) locking.dwi
// The database must be opened with the mdbm.RwLocks (=C.MDBM_RW_LOCKS) flag to enable shared locks.
// Use Unlock() to release a shared lock.
func (db *MDBM) LockShared() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock_shared(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// TryLockShared locks the database for shared access by readers, excluding access to writers.
// This is the non-blocking version of LockShared()
// This is MROW locking. The database must be opened with the mdbm.RwLocks (=C.MDBM_RW_LOCKS) flag to enable shared locks.
func (db *MDBM) TryLockShared() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_trylock_shared(db.pmdbm)
		return int(rv), err
	})

	return rv, err

}

// LockReset resets the global lock ownership state of a database.
// USE THIS FUNCTION WITH EXTREME CAUTION!
func (db *MDBM) LockReset(dbmpath string) (int, error) {

	pdbmfile := C.CString(dbmpath)
	defer C.free(unsafe.Pointer(pdbmfile))

	//flags(2nd arg) Reserved for future use, and must be 0.
	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock_reset(pdbmfile, 0)
		if rv == 0 {
			db.mutex.Lock()
			{
				db.locked = false
			}
			db.mutex.Unlock()
		}
		return int(rv), err
	})

	return rv, err
}

// MyLockReset resets the global lock ownership state of a database.
// USE THIS FUNCTION WITH EXTREME CAUTION!
func (db *MDBM) MyLockReset() (int, error) {

	return db.LockReset(db.dbmfile)
}

// DeleteLockFiles removes all lockfiles associated with the MDBM file.
// USE THIS FUNCTION WITH EXTREME CAUTION!
// HINT: /tmp/.mlock-named/[PATH]
func (db *MDBM) DeleteLockFiles(dbmpath string) (int, error) {

	path := C.CString(dbmpath)
	defer C.free(unsafe.Pointer(path))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_delete_lockfiles(path)
		if rv == 0 {
			db.mutex.Lock()
			{
				db.locked = false
			}
			db.mutex.Unlock()
		}
		return int(rv), err
	})

	return rv, err
}

// ReplaceDB replaces the database currently in oldfile db with the new database in newfile.
func (db *MDBM) ReplaceDB(newfile string) error {

	newmdbmfn := C.CString(newfile)
	defer C.free(unsafe.Pointer(newmdbmfn))

	_, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_replace_db(db.pmdbm, newmdbmfn)
		return int(rv), err
	})

	return errors.Wrapf(err, newfile)
}

// ReplaceFile replaces an old database in oldfile with new database in newfile.
// oldfile is deleted, and a newfile is renamed to a oldfile.
func (db *MDBM) ReplaceFile(oldfile, newfile string) error {

	oldmdbmfn := C.CString(oldfile)
	defer C.free(unsafe.Pointer(oldmdbmfn))

	newmdbmfn := C.CString(newfile)
	defer C.free(unsafe.Pointer(newmdbmfn))

	_, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_replace_file(oldmdbmfn, newmdbmfn)
		return int(rv), err
	})

	return err
}

// GetHash returns the MDBM's hash function identifier.
func (db *MDBM) GetHash() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_hash(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// SetHash sets the hashing function for a given MDBM
func (db *MDBM) SetHash(hashid int) error {

	if hashid < HashCRC32 || hashid > MaxHash {
		return fmt.Errorf("not support hash : hashid(%d)", hashid)
	}

	_, _, err := db.cgoRun(func() (int, error) {
		fmt.Println(C.int(hashid))
		rv, err := C.mdbm_set_hash(db.pmdbm, C.int(hashid))
		return int(rv), err
	})

	return err
}

// SetSpillSize sets the size of item data value which will be put on the large-object heap rather than inline.
// The spill size can be changed at any point after the db has been created.
// However, it's a recommended practice to set the spill size at creation time.
// NOTE: The database has to be opened with the MDBM_LARGE_OBJECTS flag for spillsize to take effect.
func (db *MDBM) SetSpillSize(size int) error {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_setspillsize(db.pmdbm, C.int(size))
		return int(rv), err
	})

	if rv == -1 {
		return fmt.Errorf("failed to set spill size=%d errno=%d. Verify you enabled Large object support(-L)", size, rv)
	}

	return err
}

// GetAlignment gets the MDBM's record byte-alignment.
// Alignment mask.
// rv 0 - 8-bit alignment
// rv 1 - 16-bit alignment
// rv 3 - 32-bit alignment
// rv 7 - 64-bit alignment
func (db *MDBM) GetAlignment() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_alignment(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// SetAlignment sets a database's byte-size alignment for keys and values within a page.
// This feature is useful for hardware/memory architectures that incur a performance penalty for unaligned accesses.
// Later (2006+) i386 and x86 architectures do not need special byte alignment,
// and should use the default of 8-bit alignment.
func (db *MDBM) SetAlignment(align int) (int, error) {

	if align < Align8Bits || align > Align64Bits {
		return -1, fmt.Errorf("not support align=%d", align)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_set_alignment(db.pmdbm, C.int(align))
		return int(rv), err
	})

	return rv, err
}

//GetLimitSize gets the MDBM's size limit. returns the limit set for the size of the db
func (db *MDBM) GetLimitSize() (uint64, error) {

	var rv uint64
	_, _, err := db.cgoRun(func() (int, error) {
		size, err := C.mdbm_get_limit_size(db.pmdbm)
		rv = uint64(size)
		return 0, err
	})

	return rv, err
}

// LimitDirSize sets limit the internal page directory size to a number of pages.
// The number of pages is rounded up to a power of 2.
func (db *MDBM) LimitDirSize(pages int) error {

	if pages < 1 {
		return fmt.Errorf("the internal page directory size must be at least 1, pages=%d", pages)
	}

	_, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_limit_dir_size(db.pmdbm, C.int(pages))
		return int(rv), err
	})

	return err
}

// GetVersion returns the on-disk format version number of the MDBM.
func (db *MDBM) GetVersion() (uint32, error) {

	var rv uint32
	_, _, err := db.cgoRun(func() (int, error) {
		ver, err := C.mdbm_get_version(db.pmdbm)
		rv = uint32(ver)
		return 0, err
	})

	return rv, err
}

// GetSize returns the current MDBM's size.
func (db *MDBM) GetSize() (uint64, error) {

	var rv uint64
	_, _, err := db.cgoRun(func() (int, error) {
		size, err := C.mdbm_get_size(db.pmdbm)
		rv = uint64(size)
		return 0, err
	})

	return rv, err
}

// GetPageSize returns the MDBM's page size.
func (db *MDBM) GetPageSize() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_page_size(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// GetMagicNumber returns the magic number from the MDBM.
func (db *MDBM) GetMagicNumber() (uint32, error) {

	var rv uint32
	var magic C.uint32_t

	errno, _, err := db.cgoRun(func() (int, error) {
		no, err := C.mdbm_get_magic_number(db.pmdbm, &magic)
		rv = uint32(no)
		return int(no), err
	})

	switch errno {
	case -3:
		return rv, fmt.Errorf("Cannot read all of the magic number : %s", db.dbmfile)
	case -2:
		return rv, fmt.Errorf("File is truncated (empty) : %s", db.dbmfile)
	case -1:
		return rv, fmt.Errorf("Cannot read file : %s", db.dbmfile)
	}

	return uint32(magic), err
}

// SetWindowSize sets the window size for the MDBM.
// Windowing is typically used for a very large data store where only part of it will be mapped to memory.
// In windowing mode, pages are accessed through a "window" to the database.
func (db *MDBM) SetWindowSize(wsize uint) error {

	if wsize < db.minpagesize {
		return fmt.Errorf("wsize should be at least 2 pages, SC_PAGESIZE=%d, wsize=%d", db.minpagesize, wsize)
	}

	_, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_set_window_size(db.pmdbm, C.size_t(wsize))
		return 0, err
	})

	return err
}

// IsOwned returns whether or not MDBM is currently locked (owned) by the calling process.
// Owned MDBMs have multiple nested locks in place.
func (db *MDBM) IsOwned() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_isowned(db.pmdbm)
		return int(rv), err
	})

	/*
		switch rv {
		case 0:
			log.Println("Database is not owned")
		case 1:
			log.Println("Database is owned")
		}
	*/

	return rv, err
}

// GetLockMode return the MDBM's lock mode.
func (db *MDBM) GetLockMode() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_lockmode(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// CompressTree compresses the existing MDBM directory.
// Attempts to rebalance the directory and to compress the db to a smaller size.
func (db *MDBM) CompressTree() error {

	//mdbm_compress_tree
	_, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_compress_tree(db.pmdbm)
		return 0, err
	})

	return err
}

// Truncate truncates the MDBM to single empty page
func (db *MDBM) Truncate() error {

	_, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_truncate(db.pmdbm)
		return 0, err
	})

	return err
}

// Purge purges (removes) all entries from the MDBM.
// This does not change the MDBM's configuration or general structure.
func (db *MDBM) Purge() error {

	_, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_purge(db.pmdbm)
		return 0, err
	})

	return err
}

// Check checks the MDBM's integrity, and displays information on standard output.
func (db *MDBM) Check(level int, verbose bool) (int, string, error) {

	//leverl between 0 and 10
	//verbose 0 or 1
	var verb C.int
	if verbose {
		verb = 1
	} else {
		verb = 0
	}
	rv, out, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_check(db.pmdbm, C.int(level), verb)
		return int(rv), err
	})

	return rv, out, err
}

// CheckAllPage checks the database for errors.
// It will report same as ChkPage() for all pages in the database.
// See v2 and v3 in ChkPage() to determine if errors detected in the database.
func (db *MDBM) CheckAllPage() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_chk_all_page(db.pmdbm)
		return int(rv), err
	})

	//log.Print(rv, out, err)
	return rv, err
}

// Protect sets all database pages to protect permission.
// This function is for advanced users only.
// Users that want to use the built-in protect feature should specify Protect (=MDBM_PROTECT) in their Open() flags.
func (db *MDBM) Protect(protect int) (int, error) {

	if protect < ProtNone || protect > ProtAccess {
		return -1, fmt.Errorf("not support protect=%d", protect)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_protect(db.pmdbm, C.int(protect))
		return int(rv), err
	})

	return rv, err
}

// DumpAllPage dumps information for all pages, in version-specific format, to standard output.
func (db *MDBM) DumpAllPage() (string, error) {

	_, out, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_dump_all_page(db.pmdbm)
		return 0, err
	})

	return out, err
}

// StoreWithLock adds key and value into the current MDBM with locking
// NOTE : Update if key exists; insert if does not exist
func (db *MDBM) StoreWithLock(key interface{}, val interface{}, flags int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}
	bval, err := db.convertToArByte(val)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	v.dptr = (*C.char)(unsafe.Pointer(&bval[0]))
	v.dsize = C.int(len(bval))

	rv, out, err := db.cgoRun(func() (int, error) {
		rv := C.set_mdbm_store_with_lock(db.pmdbm, k, v, C.int(flags))
		return int(rv), nil
	})

	switch rv {
	case -1:
		return rv, errors.New(out)
	case 1:
		return rv, errors.New("Flag const:mdbm.Insert was specified, and the key already exists")
	}

	return rv, err
}

// StoreWithLockSmart adds key and value into the current MDBM with locking
// NOTE : Update if key exists; insert if does not exist
func (db *MDBM) StoreWithLockSmart(key interface{}, val interface{}, flags int, lockflags int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}
	bval, err := db.convertToArByte(val)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	v.dptr = (*C.char)(unsafe.Pointer(&bval[0]))
	v.dsize = C.int(len(bval))

	rv, out, err := db.cgoRun(func() (int, error) {
		rv := C.set_mdbm_store_with_lock_smart(db.pmdbm, k, v, C.int(flags), C.int(lockflags))
		return int(rv), nil
	})

	switch rv {
	case -1:
		return rv, errors.New(out)
	case 1:
		return rv, errors.New("Flag const:mdbm.Insert was specified, and the key already exists")
	}

	return rv, err
}

// Store stores the record specified by the key and val parameters.
func (db *MDBM) Store(key interface{}, val interface{}, flags int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}
	bval, err := db.convertToArByte(val)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	v.dptr = (*C.char)(unsafe.Pointer(&bval[0]))
	v.dsize = C.int(len(bval))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_store(db.pmdbm, k, v, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// StoreRWithLock the record specified by the key and val parameters with locking.
func (db *MDBM) StoreRWithLock(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}
	sval, err := db.convertToString(val)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	v.dptr = C.CString(sval)
	v.dsize = C.int(len(sval))
	defer C.free(unsafe.Pointer(v.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_store_r_with_lock(db.pmdbm, &k, &v, C.int(flags), iter)
		return int(rv), err
	})

	return rv, db.convertIter(iter), err
}

// StoreRWithLockSmart the record specified by the key and val parameters with locking.
func (db *MDBM) StoreRWithLockSmart(key interface{}, val interface{}, flags int, lockflags int, iter *C.MDBM_ITER) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}
	sval, err := db.convertToString(val)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	v.dptr = C.CString(sval)
	v.dsize = C.int(len(sval))
	defer C.free(unsafe.Pointer(v.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_store_r_with_lock_smart(db.pmdbm, &k, &v, C.int(flags), C.int(lockflags), iter)
		return int(rv), err
	})

	return rv, db.convertIter(iter), err
}

// StoreR stores the record specified by the key and val parameters.
func (db *MDBM) StoreR(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}
	sval, err := db.convertToString(val)
	if err != nil {
		return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	v.dptr = C.CString(sval)
	v.dsize = C.int(len(sval))
	defer C.free(unsafe.Pointer(v.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_store_r(db.pmdbm, &k, &v, C.int(flags), iter)
		return int(rv), err
	})

	return rv, db.convertIter(iter), err
}

// StoreStrWitchLock stores the record specified by the key and val parameters with locking
// BUG: tail \00
func (db *MDBM) StoreStrWitchLock(key interface{}, val interface{}, flags int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}
	sval, err := db.convertToString(val)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)
	v := C.CString(sval)

	defer C.free(unsafe.Pointer(k))
	defer C.free(unsafe.Pointer(v))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_store_str_with_lock(db.pmdbm, k, v, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// StoreStr stores a string into the MDBM.
// BUG: tail \00
func (db *MDBM) StoreStr(key interface{}, val interface{}, flags int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}
	sval, err := db.convertToString(val)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)
	v := C.CString(sval)
	defer C.free(unsafe.Pointer(k))
	defer C.free(unsafe.Pointer(v))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_store_str(db.pmdbm, k, v, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// FetchWithLock returns fetche the record specified by the key argument and returns a value with lock
func (db *MDBM) FetchWithLock(key interface{}) (int, string, error) {

	var retval string
	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, retval, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		v, err := C.get_mdbm_fetch_with_lock(db.pmdbm, k)

		retval = C.GoString(v)
		C.free(unsafe.Pointer(v))
		return 0, err
	})

	return rv, retval, err
}

// FetchWithLockSmart returns fetche the record specified by the key argument and returns a value with lockSmart
func (db *MDBM) FetchWithLockSmart(key interface{}, lockflags int) (int, string, error) {

	var retval string
	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, retval, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		v, err := C.get_mdbm_fetch_with_lock_smart(db.pmdbm, k, C.int(lockflags))

		retval = C.GoString(v)
		C.free(unsafe.Pointer(v))
		return 0, err
	})

	return rv, retval, err
}

// FetchWithPlock returns fetche the record specified by the key argument and returns a value with plock
func (db *MDBM) FetchWithPlock(key interface{}, lockflags int) (int, string, error) {

	var retval string
	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, retval, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		v, err := C.get_mdbm_fetch_with_plock(db.pmdbm, k, C.int(lockflags)) //flags Ignored.

		retval = C.GoString(v)
		C.free(unsafe.Pointer(v))
		return 0, err
	})

	return rv, retval, err
}

// Fetch returns fetche the record specified by the key argument and returns a value
func (db *MDBM) Fetch(key interface{}) (int, string, error) {

	var retval string
	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, retval, errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		v, err = C.mdbm_fetch(db.pmdbm, k)
		return 0, err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	return rv, retval, err
}

// FetchR returns fetche the record specified by the key argument and returns a value
func (db *MDBM) FetchR(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_fetch_r(db.pmdbm, &k, &v, iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	goiter := db.convertIter(iter)

	return rv, retval, goiter, err
}

// FetchRWithLock returns fetche the record specified by the key argument and returns a value with lock
func (db *MDBM) FetchRWithLock(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.get_mdbm_fetch_r_with_lock(db.pmdbm, &k, &v, iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	goiter := db.convertIter(iter)

	return rv, retval, goiter, err
}

// FetchRWithLockSmart returns fetche the record specified by the key argument and returns a value with lock smart
func (db *MDBM) FetchRWithLockSmart(key interface{}, lockflags int, iter *C.MDBM_ITER) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.get_mdbm_fetch_r_with_lock_smart(db.pmdbm, &k, &v, C.int(lockflags), iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	goiter := db.convertIter(iter)

	return rv, retval, goiter, err
}

// FetchRWithPlock returns fetche the record specified by the key argument and returns a value with plock
func (db *MDBM) FetchRWithPlock(key interface{}, lockflags int, iter *C.MDBM_ITER) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.get_mdbm_fetch_r_with_plock(db.pmdbm, &k, &v, C.int(lockflags), iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	goiter := db.convertIter(iter)

	return rv, retval, goiter, err
}

// FetchBuf Fetches and copies the record specified by the key argument.
func (db *MDBM) FetchBuf(key interface{}, sbuf *string) (int, string, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, errors.Wrapf(err, "failured")
	}

	var k, v, buf C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_fetch_buf(db.pmdbm, &k, &v, &buf, 0) //"0" flags Reserved for future use
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)
	*sbuf = C.GoStringN(buf.dptr, buf.dsize)

	return rv, retval, err
}

// FetchDupR fetches the next value for a key inserted via mdbm_store_r with the mdbm.InsertDup (=C.MDBM_INSERT_DUP) flag set.
func (db *MDBM) FetchDupR(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_fetch_dup_r(db.pmdbm, &k, &v, iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)
	return rv, retval, db.convertIter(iter), err
}

// FetchStr returns fetche the record specified by the key argument and returns a value
func (db *MDBM) FetchStr(key interface{}) (string, error) {

	var retval string

	skey, err := db.convertToString(key)
	if err != nil {
		return retval, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)
	defer C.free(unsafe.Pointer(k))

	_, _, err = db.cgoRun(func() (int, error) {

		val, err := C.mdbm_fetch_str(db.pmdbm, k)
		retval = C.GoString(val)
		return 0, err
	})

	return retval, err
}

// FetchInfo ...
func (db *MDBM) FetchInfo(key interface{}, sbuf *string, iter *C.MDBM_ITER) (int, string, FetchInfo, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string
	var info C.struct_mdbm_fetch_info

	err := db.isVersion3Above()
	if err != nil {
		return rv, retval, db.convertFetchInfo(info), db.convertIter(iter), err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertFetchInfo(info), db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v, buf C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_fetch_info(db.pmdbm, &k, &v, &buf, &info, iter)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)
	*sbuf = C.GoStringN(buf.dptr, buf.dsize)
	goiter := db.convertIter(iter)
	goinfo := db.convertFetchInfo(info)

	return rv, retval, goinfo, goiter, err
}

// Delete deletes a specific record
func (db *MDBM) Delete(key interface{}) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	var k C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_delete(db.pmdbm, k)
		return int(rv), err
	})

	return rv, err
}

// DeleteR deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
func (db *MDBM) DeleteR(iter Iter) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	citer := db.convertIterToC(iter)

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_delete_r(db.pmdbm, &citer)
		return int(rv), err
	})

	return rv, db.convertIter(&citer), err
}

// DeleteStr deletes a string from the MDBM.
func (db *MDBM) DeleteStr(key interface{}) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_delete_str(db.pmdbm, k)
		return int(rv), err
	})

	return rv, err
}

// First returns the first key/value pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) First() (string, string, error) {

	var kv C.kvpair
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_first(db.pmdbm)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	return key, val, nil
}

// Next returns the next key/value pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) Next() (string, string, error) {

	var kv C.kvpair
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_next(db.pmdbm)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	return key, val, nil
}

// FirstR returns the first key/value pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) FirstR(iter *Iter) (string, string, Iter, error) {

	var kv C.kvpair
	var err error
	citer := db.convertIterToC(*iter)

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_first_r(db.pmdbm, &citer)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", db.convertIter(&citer), errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	*iter = db.convertIter(&citer)

	return key, val, *iter, nil
}

// NextR Fetches the next record in an MDBM.
// Returns the next key/value pair from the db, based on the iterator.
func (db *MDBM) NextR(iter *Iter) (string, string, Iter, error) {

	var kv C.kvpair
	var err error
	citer := db.convertIterToC(*iter)

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_next_r(db.pmdbm, &citer)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", db.convertIter(&citer), errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	*iter = db.convertIter(&citer)

	return key, val, *iter, nil
}

// FirstKey Returns the first key from the database.
// The order that records are returned is not specified.
func (db *MDBM) FirstKey() (string, error) {

	var k C.datum
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_firstkey(db.pmdbm)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)

	return key, err
}

// NextKey Returns the next key pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) NextKey() (string, error) {

	var k C.datum
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_nextkey(db.pmdbm)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)

	return key, err
}

// FirstKeyR fetches the first key in an MDBM.
// Initializes the iterator, and returns the first key from the db.
// Subsequent calls to NextR() or NextKeyR() with this iterator will loop through the entire db.
func (db *MDBM) FirstKeyR(iter *Iter) (string, Iter, error) {

	var k C.datum
	var err error
	citer := db.convertIterToC(*iter)

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_firstkey_r(db.pmdbm, &citer)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", db.convertIter(&citer), errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)
	*iter = db.convertIter(&citer)

	return key, *iter, nil
}

// NextKeyR fetches the next key in an MDBM.  Returns the next key from the db.
// Subsequent calls to NextR() or NextKeyR() with this iterator
// will loop through the entire db.
func (db *MDBM) NextKeyR(iter *Iter) (string, Iter, error) {

	var k C.datum
	var err error
	citer := db.convertIterToC(*iter)

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_nextkey_r(db.pmdbm, &citer)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", db.convertIter(&citer), errors.Wrapf(err, "database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)
	*iter = db.convertIter(&citer)

	return key, *iter, nil
}

// GetCacheMode returns the current cache style of the database.
// See the cachemode parameter in SetCacheMode() for the valid values.
func (db *MDBM) GetCacheMode() (int, error) {

	err := db.isVersion3Above()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_cachemode(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// SetCacheMode sets the database as a cache
// it's must be called before data is inserted.
// Tracking metadata is stored with each entry which allows MDBM to do cache eviction via LRU, LFU, and GDSF
// (greedy-dual-size-frequency). MDBM also supports clean/dirty tracking and the application can supply a callback (see SetBackingStore())
// which is called by MDBM when a dirty entry is about to be evicted allowing
// the application to sync the entry to a backing store or perform some other type of "clean" operation.
func (db *MDBM) SetCacheMode(cachemode int) (int, error) {

	err := db.isVersion3Above()
	if err != nil {
		return -1, err
	}

	if cachemode < CacheModeNone || cachemode > CacheModeMax {
		return -1, fmt.Errorf("not support cachemode=%d", cachemode)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_set_cachemode(db.pmdbm, C.int(cachemode))
		return int(rv), err
	})

	return rv, err
}

// GetCacheModeName returns the cache mode as a string. See SetCacheMode()
func (db *MDBM) GetCacheModeName(cachemode int) (string, error) {

	err := db.isVersion3Above()
	if err != nil {
		return "", err
	}

	var retval string

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_cachemode_name(C.int(cachemode))

		retval = C.GoString(rv)
		return 0, err
	})

	return retval, err
}

// CountRecords counts the number of records in an MDBM.
func (db *MDBM) CountRecords() (uint64, error) {

	var retval uint64
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_count_records(db.pmdbm)
		retval = uint64(rv)
		return 0, err
	})

	return retval, err
}

// CountPages counts the number of pages used by an MDBM.
func (db *MDBM) CountPages() (uint32, error) {

	var retval uint32
	var err error

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_count_pages(db.pmdbm)
		retval = uint32(rv)
		return 0, err
	})

	return retval, err
}

// GetPage gets the MDBM page number for a given key.
// The key does not actually have to exist.
func (db *MDBM) GetPage(key interface{}) (uint32, error) {

	var retval uint32
	var err error

	skey, err := db.convertToString(key)
	if err != nil {
		return retval, errors.Wrapf(err, "failed")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_page(db.pmdbm, &k)
		retval = uint32(rv)
		return 0, err
	})

	return retval, err

}

// PreLoad pre-loads mdbm: Read every 4k bytes to force all pages into memory
func (db *MDBM) PreLoad() (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_preload(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LockDump returns the state of lock
func (db *MDBM) LockDump() (string, error) {

	_, out, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_lock_dump(db.pmdbm)
		return 0, err
	})

	return out, err
}

// LockPages locks MDBM data pages into memory.
// When running MDBM as root, LockPages() will expand the amount of RAM
// that can be locked to infinity using setrlimit(RLIMIT_MEMLOCK).
// When not running as root, mdbm_lock_pages will expand the amount of RAM
// that can be locked up to the maximum allowed (retrieved using getrlimit(MEMLOCK),
// and normally a very small amount), and if the MDBM is larger than the amount of RAM that can be
// locked, a warning will be logged but LockPages() will return 0 for success.
func (db *MDBM) LockPages() (int, error) {

	var rv int
	var err error

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if user.Username == "root" {

		rv, _, err = db.cgoRun(func() (int, error) {
			rv, err := C.mdbm_lock_pages(db.pmdbm)
			return int(rv), err
		})

	} else {
		rv, err = -9, fmt.Errorf("When running MDBM as root, current.Username=%s", user.Username)
	}

	return rv, err
}

// UnLockPages releases MDBM data pages from always staying in memory
// When running MDBM as root
func (db *MDBM) UnLockPages() (int, error) {

	var rv int
	var err error

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if user.Username == "root" {

		rv, _, err = db.cgoRun(func() (int, error) {
			rv, err := C.mdbm_unlock_pages(db.pmdbm)
			return int(rv), err
		})

	} else {
		rv, err = -9, fmt.Errorf("When running MDBM as root, current.Username=%s", user.Username)
	}

	return rv, err
}

// ChkPage checks the specified page for errors.
// It will print errors found on the page, including bad key size, bad val size, and bad offsets of various fields.
func (db *MDBM) ChkPage(pagenum int) (int, string, error) {

	rv, out, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_chk_page(db.pmdbm, C.int(pagenum))
		return int(rv), err
	})

	return rv, out, err
}

// ChkError checks integrity of an entry on a page.
func (db *MDBM) ChkError(pagenum int, mappedpagenum int, index int) error {

	_, out, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_chk_error(db.pmdbm, C.int(pagenum), C.int(mappedpagenum), C.int(index))
		return 0, err
	})

	if len(out) > 0 {
		return errors.Wrapf(err, out)
	}

	return err
}

// DumpPage dumps specified page's information, in version-specific format, to standard output.
func (db *MDBM) DumpPage(pno int) (string, error) {

	_, out, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_dump_page(db.pmdbm, C.int(pno))
		return 0, err
	})

	return out, err
}

// ResetStatOperations resets the stat counter and last-time performed for fetch, store, and remove operations.
func (db *MDBM) ResetStatOperations() error {

	_, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_reset_stat_operations(db.pmdbm)
		return 0, err
	})

	return err
}

// EnableStatOperations enables and disables gathering of stat counters and/or last-time performed for fetch, store, and remove operations.
func (db *MDBM) EnableStatOperations(flags int) (int, error) {

	if flags < 0 || flags > (StatsBasic|StatsTimed) {
		return -1, fmt.Errorf("not support flags=%d", flags)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_enable_stat_operations(db.pmdbm, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// GetStatCounter enables and disables gathering of stat counters and/or last-time performed for fetch, store, and remove operations.
func (db *MDBM) GetStatCounter(stype int) (int, uint64, error) {

	if stype < StatTypeFetch || stype > StatTypeMax {
		return -1, 0, fmt.Errorf("not support stype=%d", stype)
	}

	var cntvalue C.mdbm_counter_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stat_counter(db.pmdbm, C.mdbm_stat_type(stype), &cntvalue)
		return int(rv), err
	})

	return rv, uint64(cntvalue), err
}

// GetStatName gets the name of a stat.
func (db *MDBM) GetStatName(stype int) (string, error) {

	//stype between 0 (MDBM_STAT_TAG_FETCH) and 16 (MDBM_STAT_TAG_DELETE_FAILED) on mdbm 4.x
	var retval string
	_, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stat_name(C.int(stype))
		retval = C.GoString(rv)

		return 0, err
	})

	return retval, err
}

// GetStatTime Gets the last time when an type operation was performed.
func (db *MDBM) GetStatTime(stype int) (int, uint64, error) {

	if stype < StatTypeFetch || stype > StatTypeMax {
		return -1, 0, fmt.Errorf("not support stype=%d", stype)
	}

	var value C.time_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stat_time(db.pmdbm, C.mdbm_stat_type(stype), &value)

		return int(rv), err
	})

	return rv, uint64(value), err
}

// SetStatTimeFunc tells the MDBM library whether to use TSC (CPU TimeStamp Counters)
// for timing the performance of fetch, store and delete operations.
// The standard behavior of timed stat operations is to use clock_gettime(MONOTONIC)
func (db *MDBM) SetStatTimeFunc(flags int) (int, error) {

	if flags < ClockStandard || flags > ClockTsc {
		return -1, fmt.Errorf("not support flags=%d", flags)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_set_stat_time_func(db.pmdbm, C.int(flags))

		return int(rv), err
	})

	return rv, err
}

// StatAllPage returns various pieces of information, specifically
// Total mapped db size, Total number of entries, ADDRESS SPACE page efficiency,
// ADDRESS SPACE byte efficiency, PHYSICAL MEM/DISK SPACE efficiency,
// Average bytes per record, Maximum B-tree level, Minimum B-tree level,
// Minimum free bytes on page, Minimum free bytes on page post compress
func (db *MDBM) StatAllPage() (string, error) {

	err := db.isVersion2()
	if err != nil {
		return "", err
	}

	//raise the SIGABORT when V3
	_, out, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_stat_all_page(db.pmdbm)
		return 0, err
	})

	return out, err
}

// GetStats gets a a stats block with individual stat values.
func (db *MDBM) GetStats() (int, Stats, error) {

	var stats C.mdbm_stats_t
	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stats(db.pmdbm, &stats, C.sizeof_mdbm_stats_t)
		return int(rv), err
	})

	return rv, db.convertStats(stats), err
}

// GetDBInfo gets configuration information about a database
func (db *MDBM) GetDBInfo() (int, DBInfo, error) {

	var info C.mdbm_db_info_t
	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_db_info(db.pmdbm, &info)
		return int(rv), err
	})

	return rv, db.convertDBInfo(info), err
}

// GetDBStats gets overall database stats.
func (db *MDBM) GetDBStats(flags int) (int, DBInfo, StatInfo, error) {

	if flags != StatNolock && flags > IterateNolock {
		return -1, DBInfo{}, StatInfo{}, fmt.Errorf("not support flags=%d", flags)
	}

	var dbinfo C.mdbm_db_info_t
	var statinfo C.mdbm_stat_info_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_db_stats(db.pmdbm, &dbinfo, &statinfo, C.int(flags))
		return int(rv), err
	})

	return rv, db.convertDBInfo(dbinfo), db.convertStatInfo(statinfo), err
}

// GetWindowStats retrieves statistics about windowing usage on the associated database.
func (db *MDBM) GetWindowStats() (int, WindowStats, error) {

	var stats C.mdbm_window_stats_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_window_stats(db.pmdbm, &stats, C.sizeof_mdbm_window_stats_t)
		return int(rv), err
	})

	return rv, db.convertWindowStat(stats), err
}

// GetHashValue return the hash function code, get the hash value for the given key.
// See SetHash() for the list of valid hash function codes.
func (db *MDBM) GetHashValue(key interface{}, hashFunctionCode int) (uint32, error) {

	if hashFunctionCode < HashCRC32 || hashFunctionCode > MaxHash {
		return 0, fmt.Errorf("not support hashFunctionCode=%d", hashFunctionCode)
	}

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return 0, errors.Wrapf(err, "failured")
	}

	var k C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	var hashValue C.uint32_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_hash_value(k, C.int(hashFunctionCode), &hashValue)
		return int(rv), err
	})

	if rv == 0 {
		return uint32(hashValue), nil
	}

	return 0, err
}

// Plock locks a specific partition in the database for exclusive access by the caller.
// The lock is nestable, so a caller already holding the lock may call Plock() again
// as long as an equal number of calls to Punlock() are made to release the lock.
func (db *MDBM) Plock(key interface{}) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_plock(db.pmdbm, &k, C.int(0)) //flags ignored
		return int(rv), err
	})

	//rv 1 Success, partition lock was acquired
	return rv, err
}

// Punlock unlocks a specific partition in the database, releasing exclusive access by the caller.
// If the caller has called Plock() multiple times in a row, an equal number of unlock calls are required.
// See Plock() for usage.
func (db *MDBM) Punlock(key interface{}) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_punlock(db.pmdbm, &k, C.int(0)) //flags ignored
		return int(rv), err
	})

	// rv 1 Success, partition lock was released
	return rv, err
}

// TryPlock tries to locks a specific partition in the database for exclusive access by the caller.
// The lock is nestable, so a caller already holding the lock may call Plock() again
// as long as an equal number of calls to Punlock() are made to release the lock.
// See Plock() for usage.
func (db *MDBM) TryPlock(key interface{}) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_tryplock(db.pmdbm, &k, C.int(0)) //flags ignored
		return int(rv), err
	})

	return rv, err
}

// LockSmart performs either partition, shared or exclusive locking based on the locking-related flags supplied to Open()
func (db *MDBM) LockSmart(key interface{}, flags int) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock_smart(db.pmdbm, &k, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// UnLockSmart unlocks an MDBM based on the locking flags supplied to Open()
func (db *MDBM) UnLockSmart(key interface{}, flags int) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_unlock_smart(db.pmdbm, &k, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// TryLockSmart attempts to lock an MDBM based on the locking flags supplied to Open()
func (db *MDBM) TryLockSmart(key interface{}) (int, error) {

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_trylock_smart(db.pmdbm, &k, C.int(0)) //flags ignored
		return int(rv), err
	})

	return rv, err
}

// CheckResidency checks mdbm page residency: count the number of DB pages mapped into memory.
// The counts are in units of the system-page-size (typically 4k)
func (db *MDBM) CheckResidency() (int, uint32, uint32, error) {

	var pgsin, pgsout C.mdbm_ubig_t

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_check_residency(db.pmdbm, &pgsin, &pgsout)
		return int(rv), err
	})

	return rv, uint32(pgsin), uint32(pgsout), err
}

/*
func (db *MDBM) CDBDumpToFile(key interface{}, val interface{}, fnpath string, mode string) (int, *C.FILE, error) {

	rv := -1

	var fd *C.struct__IO_FILE

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, fd, errors.Wrapf(err, "failured")
	}
	bval, err := db.convertToArByte(val)
	if err != nil {
		return rv, fd, errors.Wrapf(err, "failured")
	}

	var k, v C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	v.dptr = (*C.char)(unsafe.Pointer(&bval[0]))
	v.dsize = C.int(len(bval))

	var kv C.kvpair
	kv.key = k
	kv.val = v

	fdpath := C.CString(fnpath)
	fdmode := C.CString(mode)
	defer C.free(unsafe.Pointer(fdpath))
	defer C.free(unsafe.Pointer(fdmode))

	fd, err = C.fopen(fdpath, fdmode)
	if err != nil {
		return rv, fd, errors.Wrapf(err, "fnpath=%s, mode=%s", fnpath, mode)
	}

	//defer C.fclose(fd)

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_cdbdump_to_file(kv, fd)
		return int(rv), err
	})

	return rv, fd, err

}

func (db *MDBM) CDBDumpTrailerAndClose(fd *C.FILE) (int, error) {

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_cdbdump_trailer_and_close(fd)
		return int(rv), err
	})

	return rv, err
}
*/

/*
func (db *MDBM) Clean(pagenum int) (int, error) {

	rv, out, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_clean(db.pmdbm, C.int(pagenum), C.int(0)) //flags ignored
		return int(rv), err
	})

	fmt.Println(out)

	return rv, err
}
*/
