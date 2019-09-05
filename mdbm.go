package mdbm

// #cgo CFLAGS: -I/usr/local/mdbm/include/ -I./
// #cgo LDFLAGS: -L/usr/local/mdbm/lib64/ -Wl,-rpath,/usr/local/mdbm/lib64/ -lmdbm
// #include <mdbm-binding.h>
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
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

// MDBM represents the MDBM Go's methods
type MDBM struct {
	pmdbm    *C.MDBM
	iter     C.MDBM_ITER
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
		dbmfile:  "",
		flags:    Create | Rdrw | AnyLocks | LargeObjects,
		perms:    0666,
		psize:    0,
		dsize:    0,
		isopened: false,
	}

	obj.scpagesize = uint(C.sysconf(C._SC_PAGESIZE))
	obj.minpagesize = obj.scpagesize * 2
	obj.iter = obj.GetNewIter()

	//runtime.GOMAXPROCS(1)
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

// IsOpened returns whether a dbm opened and is writable/readable
func (db *MDBM) IsOpened() bool {

	return db.isopened
}

func (db *MDBM) checkAvailable() error {

	if false == db.IsOpened() {
		return errors.New("not yet open the MDBM or closed the MDBM")
	}

	return nil
}

func (db *MDBM) cgoRunCapture(call func() (int, error)) (int, string, error) {

	var err error

	db.cgomtx.Lock()
	defer db.cgomtx.Unlock()

	orgStdErr, err := syscall.Dup(syscall.Stderr)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Dup(STDERR)")
	}

	defer func() {
		e := syscall.Dup2(orgStdErr, syscall.Stderr)
		if e != nil {
			err = e
		}
	}()

	orgStdOut, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Dup(STDOUT)")
	}

	defer func() {
		e := syscall.Dup2(orgStdOut, syscall.Stdout)
		if e != nil {
			err = e
		}
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed create a os.Pipe()")
	}

	defer func() {
		r.Close()
	}()

	err = syscall.Dup2(int(w.Fd()), syscall.Stderr)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Dup2(STDERR)")
	}

	err = syscall.Dup2(int(w.Fd()), syscall.Stdout)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Dup2(STDOUT)")
	}

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
	runtime.Gosched()

	err = w.Close()
	if err != nil {
		log.Println("failred close the os.Stderr, os.Stdout")
	}

	err = syscall.Close(syscall.Stderr)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Close(STDERR)")
	}

	err = syscall.Close(syscall.Stderr)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed call to syscall.Close(STDOUT)")
	}

	return rv, string(<-out), err
}

func (db *MDBM) cgoRun(call func() (int, error)) (int, string, error) {

	//run
	rv, err := call()
	return rv, "", err
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
		return nil, fmt.Errorf("not support type(%v)", reflect.TypeOf(obj))
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
		return string(reflect.ValueOf(obj).Bytes()), nil

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
		return "", fmt.Errorf("not support type(%v)", reflect.TypeOf(obj))
	}
}

// isVersion3Above checks out what the version is 3 or higher
func (db *MDBM) isVersion3Above() error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return err
	}

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

// GetDBMFile returns a current dbm file path
func (db *MDBM) GetDBMFile() string {
	return db.dbmfile
}

// DupHandle returns a pointer of the Duplicate an existing database handle.
// The advantage of dup'ing a handle over doing a separate Open is that dup's handle share the same virtual
// page mapping Within the process space (saving memory).
// Threaded applications should use pthread_mutex_lock and unlock around calls to mdbm_dup_handle.
func (db *MDBM) DupHandle() (*MDBM, error) {

	var err error

	err = db.checkAvailable()
	if err != nil {
		return nil, err
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	obj := &MDBM{
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

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_errno(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LogMinLevel sets the minimum logging level,Lower priority messages are discarded
func (db *MDBM) LogMinLevel(lv C.int) error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	if lv < LogOff || lv > LogDebug {
		return fmt.Errorf("Not support log level=%d", int(lv))
	}

	_, _, err = db.cgoRun(func() (int, error) {
		C.mdbm_log_minlevel(lv)
		return 0, nil
	})

	return err
}

// LogPlugin sets the logging plug-in.
// LogToStdErr		stderr
// LogToFile		file
// LogToSyslog		syslog
func (db *MDBM) LogPlugin(plugin int) error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	var plugname string

	switch plugin {
	case LogToStdErr:
		plugname = "stderr"
	case LogToFile:
		plugname = "file"
	case LogToSysLog:
		plugname = "file"
	case LogToSkip:
		return nil

	default:
		return fmt.Errorf("Not support log plugin=%d", plugin)
	}

	pplugname := C.CString(plugname)
	defer C.free(unsafe.Pointer(pplugname))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_select_log_plugin(pplugname)
		return int(rv), err
	})

	if rv == 0 {
		return nil
	}

	return err
}

// LogToAutoFile sets the logging to file (name: /mdbm_path/mdbm_file + .log-PID)
func (db *MDBM) LogToAutoFile() (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	logpath := fmt.Sprintf("%s.log-%d", db.dbmfile, os.Getpid())
	return db.LogToFile(logpath)
}

// LogToFile sets the logging to file
func (db *MDBM) LogToFile(fnpath string) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	logpath := C.CString(fnpath)
	defer C.free(unsafe.Pointer(logpath))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_set_log_filename(logpath)
		return int(rv), err
	})

	return rv, err
}

// EasyOpen Creates and/or opens the MDBM database use the default options
func (db *MDBM) EasyOpen(dbmfile string, perms int) error {

	var err error

	if len(strings.TrimSpace(dbmfile)) < 1 {
		return errors.New("dbm file path is empty")
	}

	// if opened
	if db.isopened {
		log.Printf("Already opened db handler(=%s). will open after close the previously db handler.", db.dbmfile)
		db.EasyClose()
	}

	//protect : sigfault
	if db.flags == (db.flags|C.MDBM_O_CREAT) && db.flags == (db.flags|C.MDBM_PROTECT) {
		return errors.New("failed to open the MDBM, not support create db.flags with MDBM_PROTECT")
	}

	if db.flags == (db.flags|C.MDBM_O_ASYNC) && db.flags == (db.flags|C.MDBM_O_FSYNC) {
		return errors.New("failed to open the MDBM, not support mixed sync db.flags (MDBM_O_FSYNC, MDBM_O_ASYNC)")
	}

	if db.flags == (db.flags|C.MDBM_O_RDONLY) && db.flags == (db.flags|C.MDBM_O_WRONLY) {
		return errors.New("failed to open the MDBM, not support mixed access db.flags (MDBM_O_RDONLY, MDBM_O_WRONLY, MDBM_O_RDWR)")
	}

	db.dbmfile = dbmfile
	if perms > 0 {
		db.perms = perms
	}

	db.pdbmfile = C.CString(dbmfile)

	_, out, err := db.cgoRunCapture(func() (int, error) {
		db.pmdbm, err = C.mdbm_open(db.pdbmfile, C.int(db.flags), C.int(db.perms), C.int(db.psize), C.int(db.dsize))
		if db.pmdbm != nil {
			return 0, nil
		}

		return 0, err
	})

	if err == nil {
		db.mutex.Lock()
		db.isopened = true
		db.mutex.Unlock()
	} else {
		return errors.Wrapf(err, out)
	}

	err = db.LogMinLevel(LogOff)
	if err != nil {
		return err
	}

	return err
}

// Open Creates and/or opens the MDBM database
// mdbmfn	Name of the backing file for the database.
// flags	Specifies the open-mode for the file, usually either (MDBM_O_RDWR|MDBM_O_CREAT) or (MDBM_O_RDONLY). Flag MDBM_LARGE_OBJECTS may be used to enable large object support. Large object support can only be enabled when the database is first created. Subsequent mdbm_open calls will ignore the flag. Flag MDBM_PARTITIONED_LOCKS may be used to enable partition locking a per mdbm_open basis.
// mode	Used to set the file permissions if the file needs to be created.
// psize	Specifies the page size for the database and is set when the database is created. The minimum page size is 128. In v2, the maximum is 64K. In v3, the maximum is 16M - 64. The default, if 0 is specified, is 4096.
// presize	Specifies the initial size for the database. The database will dynamically grow as records are added, but specifying an initial size may improve efficiency. If this is not a multiple of psize, it will be increased to the next psize multiple.
//
func (db *MDBM) Open(mdbmfn string, flags int, perms int, psize int, dsize int) error {

	db.flags = flags
	db.perms = perms
	db.psize = psize
	db.dsize = dsize

	return db.EasyOpen(mdbmfn, perms)
}

// Sync syncs all pages to disk asynchronously. it's mapped pages are scheduled to be flushed to disk.
func (db *MDBM) Sync() (int, error) {

	var rv int

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}
	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_sync(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// Fsync syncs all pages to disk synchronously. it's will pages have been flushed to disk.
// The database is locked while pages are flushed.
func (db *MDBM) Fsync() (int, error) {

	var rv int

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_fsync(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// CloseFD closes the MDBM's underlying file descriptor.
func (db *MDBM) CloseFD() error {

	var err error

	err = db.checkAvailable()
	if err != nil {
		return err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		_, err := C.mdbm_close_fd(db.pmdbm)
		db.isopened = false
		return 0, err
	})

	return err
}

// Close Closes the database after Sync
func (db *MDBM) Close() {

	err := db.checkAvailable()
	if err != nil {
		return
	}

	if db.isopened {
		C.mdbm_close(db.pmdbm)
		db.isopened = false
	}

	if len(db.dbmfile) > 0 {
		C.free(unsafe.Pointer(db.pdbmfile))
	}
}

// EasyClose Closes the database after Sync
func (db *MDBM) EasyClose() {

	err := db.checkAvailable()
	if err != nil {
		return
	}

	rv, err := db.Sync()
	if err != nil {
		log.Printf("failed db.Sync(), rv=%d, err=%v", rv, err)
	}

	C.mdbm_close(db.pmdbm)

	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.isopened = false
}

// Lock locks the database for exclusive access by the caller.
// The lock is nestable, so a caller already holding
// the lock may call mdbm_lock again as long as an equal number of calls
// to Unlock are made to release the lock.
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithLock API
func (db *MDBM) Lock() (int, error) {

	var rv int

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// Unlock unlocks the database, releasing exclusive or shared access by the caller.
// If the caller has called Lock() or LockShared() multiple times
// in a row, an equal number of unlock calls are required.
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithLock API
func (db *MDBM) Unlock() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_unlock(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// TryLock attempts to exclusively lock the MDBM.
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithTryLock API
func (db *MDBM) TryLock() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_trylock(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// IsLocked returns whether or not MDBM is locked by another process or thread.
// rv 0 Database is not locked
// rv 1 Database is locked
func (db *MDBM) IsLocked() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_islocked(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LockShared Locks the database for shared access by readers, excluding access to writers.
// This is multiple-readers, one writer (MROW) locking.dwi
// The database must be opened With the mdbm.RwLocks (=C.MDBM_RW_LOCKS) flag to enable shared locks.
// Use Unlock() to release a shared lock.
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithLockShared API
func (db *MDBM) LockShared() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_lock_shared(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// TryLockShared locks the database for shared access by readers, excluding access to writers.
// This is the non-blocking version of LockShared()
// This is MROW locking. The database must be opened With the mdbm.RwLocks (=C.MDBM_RW_LOCKS) flag to enable shared locks.
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithTryLockShared API
func (db *MDBM) TryLockShared() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
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
	rv, _, err := db.cgoRunCapture(func() (int, error) {
		rv, err := C.mdbm_lock_reset(pdbmfile, 0)
		return int(rv), err
	})

	return rv, err
}

// MyLockReset resets the global lock ownership state of a database.
// USE THIS FUNCTION WITH EXTREME CAUTION!
func (db *MDBM) MyLockReset() (int, error) {
	return db.LockReset(db.dbmfile)
}

// DeleteLockFiles removes all lockfiles associated With the MDBM file.
// USE THIS FUNCTION WITH EXTREME CAUTION!
// HINT: /tmp/.mlock-named/[PATH]
func (db *MDBM) DeleteLockFiles(dbmpath string) (int, error) {

	path := C.CString(dbmpath)
	defer C.free(unsafe.Pointer(path))

	rv, _, err := db.cgoRunCapture(func() (int, error) {
		rv, err := C.mdbm_delete_lockfiles(path)
		return int(rv), err
	})

	return rv, err
}

// ReplaceDB replaces the database currently in oldfile db With the new database in newfile.
func (db *MDBM) ReplaceDB(newfile string) error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	newmdbmfn := C.CString(newfile)
	defer C.free(unsafe.Pointer(newmdbmfn))

	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_replace_db(db.pmdbm, newmdbmfn)
		return int(rv), err
	})

	return err
}

// ReplaceFile replaces an old database in oldfile With new database in newfile.
// oldfile is deleted, and a newfile is renamed to a oldfile.
func (db *MDBM) ReplaceFile(oldfile, newfile string) error {

	db.mutex.Lock()
	defer db.mutex.Unlock()

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

/*
// ReplaceBackingStore Atomically replaces an old database in oldfile with a new database in
// newfile. oldfile is deleted, and newfile is renamed to oldfile.
//
// The old database is locked (if the MDBM were opened with locking) while the
// new database is renamed from newfile to oldfile, and the old database
// is marked as having been replaced.  The marked old database causes all
// processes that have the old database open to reopen using the new database on
// their next access.
//
// Only database files of the same version may be specified for oldfile and
// newfile. For example, mixing and matching of v2 and v3 for oldfile and
// newfile is not allowed.
//
// replaceFile() may be used if the MDBM is opened with locking or without
// locking (using mdbm_open flag MDBM_OPEN_NOLOCK), and without per-access
// locking, if all accesses are read (fetches) accesses across all programs that
// open that MDBM.  If there are any write (store/delete) accesses, you must
// open the MDBM with locking, and you must lock around all operations (fetch, store, delete, iterate).
func (db *MDBM) ReplaceBackingStore(newfile string) error {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	newmdbmfn := C.CString(newfile)
	defer C.free(unsafe.Pointer(newmdbmfn))

	_, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_replace_backing_store(db.pmdbm, newmdbmfn)
		return int(rv), err
	})

	return err
}
*/

// GetHash returns the MDBM's hash function identifier.
func (db *MDBM) GetHash() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_hash(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// SetHash sets the hashing function for a given MDBM
func (db *MDBM) SetHash(hashid int) error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	if hashid < HashCRC32 || hashid > MaxHash {
		return fmt.Errorf("not support hash : hashid(%d)", hashid)
	}

	_, _, err = db.cgoRunCapture(func() (int, error) {
		fmt.Println(C.int(hashid))
		rv, err := C.mdbm_set_hash(db.pmdbm, C.int(hashid))
		return int(rv), err
	})

	return err
}

// SetSpillSize sets the size of item data value which will be put on the large-object heap rather than inline.
// The spill size can be changed at any point after the db has been created.
// However, it's a recommended practice to set the spill size at creation time.
// NOTE: The database has to be opened With the MDBM_LARGE_OBJECTS flag for spillsize to take effect.
func (db *MDBM) SetSpillSize(size int) error {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
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

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_alignment(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// SetAlignment sets a database's byte-size alignment for keys and values Within a page.
// This feature is useful for hardware/memory architectures that incur a performance penalty for unaligned accesses.
// Later (2006+) i386 and x86 architectures do not need special byte alignment,
// and should use the default of 8-bit alignment.
func (db *MDBM) SetAlignment(align int) (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	if align < Align8Bits || align > Align64Bits {
		return -1, fmt.Errorf("not support align=%d", align)
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_set_alignment(db.pmdbm, C.int(align))
		return int(rv), err
	})

	return rv, err
}

//GetLimitSize gets the MDBM's size limit. returns the limit set for the size of the db
func (db *MDBM) GetLimitSize() (uint64, error) {

	var rv uint64
	err := db.checkAvailable()
	if err != nil {
		return 0, err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		size, err := C.mdbm_get_limit_size(db.pmdbm)
		rv = uint64(size)
		return 0, err
	})

	return rv, err
}

// LimitDirSize sets limit the internal page directory size to a number of pages.
// The number of pages is rounded up to a power of 2.
func (db *MDBM) LimitDirSize(pages int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	if pages < 1 {
		return -1, fmt.Errorf("the internal page directory size must be at least 1, pages=%d", pages)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_limit_dir_size(db.pmdbm, C.int(pages))
		return int(rv), err
	})

	return rv, err
}

// GetVersion returns the on-disk format version number of the MDBM.
func (db *MDBM) GetVersion() (uint32, error) {

	var rv uint32
	err := db.checkAvailable()
	if err != nil {
		return 0, err
	}
	_, _, err = db.cgoRun(func() (int, error) {
		ver, err := C.mdbm_get_version(db.pmdbm)
		rv = uint32(ver)
		return 0, err
	})

	return rv, err
}

// GetSize returns the current MDBM's size.
func (db *MDBM) GetSize() (uint64, error) {

	var rv uint64
	err := db.checkAvailable()
	if err != nil {
		return 0, err
	}
	_, _, err = db.cgoRun(func() (int, error) {
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

	err := db.checkAvailable()
	if err != nil {
		return 0, err
	}

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
func (db *MDBM) SetWindowSize(wsize uint) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	if wsize < db.minpagesize {
		return -1, fmt.Errorf("wsize should be at least 2 pages, SC_PAGESIZE=%d, wsize=%d", db.minpagesize, wsize)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		_, err := C.mdbm_set_window_size(db.pmdbm, C.size_t(wsize))
		return 0, err
	})

	return rv, err
}

// IsOwned returns whether or not MDBM is currently locked (owned) by the calling process.
// Owned MDBMs have multiple nested locks in place.
func (db *MDBM) IsOwned() (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_lockmode(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// CompressTree compresses the existing MDBM directory.
// Attempts to rebalance the directory and to compress the db to a smaller size.
func (db *MDBM) CompressTree() error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	//mdbm_compress_tree
	_, _, err = db.cgoRun(func() (int, error) {
		_, err := C.mdbm_compress_tree(db.pmdbm)
		return 0, err
	})

	return err
}

// Truncate truncates the MDBM to single empty page
func (db *MDBM) Truncate() error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		_, err := C.mdbm_truncate(db.pmdbm)
		return 0, err
	})

	return err
}

// Purge purges (removes) all entries from the MDBM.
// This does not change the MDBM's configuration or general structure.
func (db *MDBM) Purge() error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		_, err := C.mdbm_purge(db.pmdbm)
		return 0, err
	})

	return err
}

// Check checks the MDBM's integrity, and displays information on standard output.
func (db *MDBM) Check(level int, verbose bool) (int, string, error) {

	var rv int
	var out string
	var err error

	err = db.checkAvailable()
	if err != nil {
		return -1, "", err
	}

	//level between 0 and 10
	if level < 0 || level > 10 {
		return -1, "", fmt.Errorf("not support level(=%d), required level between 1 and 10", level)
	}

	//verbose 0 or 1
	var verb C.int
	if verbose {
		verb = 1
		rv, out, err = db.cgoRunCapture(func() (int, error) {
			rv, err := C.mdbm_check(db.pmdbm, C.int(level), verb)
			return int(rv), err
		})

	} else {
		verb = 0
		rv, out, err = db.cgoRun(func() (int, error) {
			rv, err := C.mdbm_check(db.pmdbm, C.int(level), verb)
			return int(rv), err
		})

	}

	return rv, out, err
}

// CheckAllPage checks the database for errors.
// It will report same as ChkPage() for all pages in the database.
// See v2 and v3 in ChkPage() to determine if errors detected in the database.
func (db *MDBM) CheckAllPage() (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_chk_all_page(db.pmdbm)
		return int(rv), err
	})

	//log.Print(rv, out, err)
	return rv, err
}

// Protect sets all database pages to protect permission.
// This function is for advanced users only.
// Users that want to use the built-in protect feature should specify Protect (=MDBM_PROTECT) in their Open() flags.
// NOTE: RHEL is unable to set mdbm.ProtWrite without mdbm.ProtRead , so specifying mdbm.ProtWrite does not protect against reads.
func (db *MDBM) Protect(protect int) (int, error) {

	var rv int
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	if protect < ProtNone || protect > ProtAccess {
		return -1, fmt.Errorf("not support protect=%d", protect)
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_protect(db.pmdbm, C.int(protect))
		return int(rv), err
	})

	return rv, err
}

// DumpAllPage dumps information for all pages, in version-specific format, to standard output.
func (db *MDBM) DumpAllPage() (string, error) {

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	_, out, err := db.cgoRunCapture(func() (int, error) {
		_, err := C.mdbm_dump_all_page(db.pmdbm)
		return 0, err
	})

	return out, err
}

func (db *MDBM) storeWithAnyLock(key interface{}, val interface{}, flags int, lockType C.int, lockFlags C.int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
		rv, err := C.set_mdbm_store_with_lock(db.pmdbm, k, v, C.int(flags), lockType, lockFlags)
		return int(rv), err
	})

	switch rv {
	case -1:
		return rv, errors.New(out)
	case 1:
		return rv, errors.New("Flag const:mdbm.Insert was specified, and the key already exists")
	}

	return rv, err
}

// StoreWithLock adds key and value into the current MDBM With Lock()
func (db *MDBM) StoreWithLock(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeLock, lockFlagsSkip)
}

// StoreWithLockSmart adds key and value into the current MDBM With LocckSmart()
func (db *MDBM) StoreWithLockSmart(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeSmart, C.int(lockflags))
}

// StoreWithLockShared adds key and value into the current MDBM With LockShared()
func (db *MDBM) StoreWithLockShared(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeShared, lockFlagsSkip)
}

// StoreWithPlock adds key and value into the current MDBM With Plock()
func (db *MDBM) StoreWithPlock(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypePlock, C.int(lockflags))
}

// StoreWithTryLock adds key and value into the current MDBM With TryLock()
func (db *MDBM) StoreWithTryLock(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeTryLock, lockFlagsSkip)
}

// StoreWithTryLockSmart adds key and value into the current MDBM With TryLockSmart()
func (db *MDBM) StoreWithTryLockSmart(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeTrySmart, C.int(lockflags))
}

// StoreWithTryLockShared adds key and value into the current MDBM with TryShared()
func (db *MDBM) StoreWithTryLockShared(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeTryShared, lockFlagsSkip)
}

// StoreWithTryPlock adds key and value into the current MDBM With TryPlock()
func (db *MDBM) StoreWithTryPlock(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeTryPlock, C.int(lockflags))
}

// Store stores the record specified by the key and val parameters.
func (db *MDBM) Store(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeWithAnyLock(key, val, flags, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) storeRWithAnyLock(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER, lockType C.int, lockFlags C.int) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertIter(iter), err
	}

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
		rv, err := C.set_mdbm_store_r_with_lock(db.pmdbm, &k, &v, C.int(flags), iter, lockType, lockFlags)
		return int(rv), err
	})

	return rv, db.convertIter(iter), err
}

// StoreRWithLock the record specified by the key and val parameters With Lock()
func (db *MDBM) StoreRWithLock(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeLock, lockFlagsSkip)
}

// StoreRWithLockSmart the record specified by the key and val parameters With LockSmart()
func (db *MDBM) StoreRWithLockSmart(key interface{}, val interface{}, flags int, lockflags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeSmart, C.int(lockflags))
}

// StoreRWithLockShared the record specified by the key and val parameters With LockShared()
func (db *MDBM) StoreRWithLockShared(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeShared, lockFlagsSkip)
}

// StoreRWithPlock the record specified by the key and val parameters With Plock()
func (db *MDBM) StoreRWithPlock(key interface{}, val interface{}, flags int, lockflags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypePlock, C.int(lockflags))
}

// StoreRWithTryLock the record specified by the key and val parameters With TryLock()
func (db *MDBM) StoreRWithTryLock(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeTryLock, lockFlagsSkip)
}

// StoreRWithTryLockSmart the record specified by the key and val parameters With TryLockSmart()
func (db *MDBM) StoreRWithTryLockSmart(key interface{}, val interface{}, flags int, lockflags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeTrySmart, C.int(lockflags))
}

// StoreRWithTryLockShared the record specified by the key and val parameters With TryLockShared()
func (db *MDBM) StoreRWithTryLockShared(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeTryShared, lockFlagsSkip)
}

// StoreRWithTryPlock the record specified by the key and val parameters With TryPlock()
func (db *MDBM) StoreRWithTryPlock(key interface{}, val interface{}, flags int, lockflags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeTryPlock, C.int(lockflags))
}

// StoreR stores the record specified by the key and val parameters.
func (db *MDBM) StoreR(key interface{}, val interface{}, flags int, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.storeRWithAnyLock(key, val, flags, iter, lockTypeNone, lockFlagsSkip)
}

// BUG: tail \00
func (db *MDBM) storeStrWithAnyLock(key interface{}, val interface{}, flags int, lockType C.int, lockFlags C.int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
		rv, err := C.set_mdbm_store_str_with_lock(db.pmdbm, k, v, C.int(flags), lockType, lockFlags)
		return int(rv), err
	})

	return rv, err
}

// StoreStrWithLock stores the record specified by the key and val parameters With Lock()
func (db *MDBM) StoreStrWithLock(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeLock, lockFlagsSkip)
}

// StoreStrWithLockSmart stores the record specified by the key and val parameters With LockSmart()
func (db *MDBM) StoreStrWithLockSmart(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeSmart, C.int(lockflags))
}

// StoreStrWithLockShared stores the record specified by the key and val parameters With LockShared()
func (db *MDBM) StoreStrWithLockShared(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeShared, lockFlagsSkip)
}

// StoreStrWithPlock stores the record specified by the key and val parameters With Plock()
func (db *MDBM) StoreStrWithPlock(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypePlock, C.int(lockflags))
}

// StoreStrWithTryLock stores the record specified by the key and val parameters With TryLock()
func (db *MDBM) StoreStrWithTryLock(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeTryLock, lockFlagsSkip)
}

// StoreStrWithTryLockSmart stores the record specified by the key and val parameters With TryLockSmart()
func (db *MDBM) StoreStrWithTryLockSmart(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeTrySmart, C.int(lockflags))
}

// StoreStrWithTryLockShared stores the record specified by the key and val parameters With TryLockShared()
func (db *MDBM) StoreStrWithTryLockShared(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeTryShared, lockFlagsSkip)
}

// StoreStrWithTryPlock stores the record specified by the key and val parameters With TryPlock()
func (db *MDBM) StoreStrWithTryPlock(key interface{}, val interface{}, flags int, lockflags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeTryPlock, C.int(lockflags))
}

// StoreStr stores the record specified by the key and val parameters
func (db *MDBM) StoreStr(key interface{}, val interface{}, flags int) (int, error) {
	return db.storeStrWithAnyLock(key, val, flags, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) fetchWithAnyLock(key interface{}, lockType C.int, lockFlags C.int) (string, error) {

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	var retval string

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return retval, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	_, _, err = db.cgoRun(func() (int, error) {
		v, err := C.get_mdbm_fetch_with_lock(db.pmdbm, k, lockType, lockFlags)

		retval = C.GoString(v)
		C.free(unsafe.Pointer(v))
		return 0, err
	})

	return retval, err
}

// FetchWithLock returns fetche the record specified by the key argument and returns a value With Lock()
func (db *MDBM) FetchWithLock(key interface{}) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeLock, lockFlagsSkip)
}

// FetchWithLockSmart returns fetche the record specified by the key argument and returns a value With LocckSmart()
func (db *MDBM) FetchWithLockSmart(key interface{}, lockflags int) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeSmart, C.int(lockflags))
}

// FetchWithLockShared returns fetche the record specified by the key argument and returns a value With LockShared()
func (db *MDBM) FetchWithLockShared(key interface{}) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeShared, lockFlagsSkip)
}

// FetchWithPlock returns fetche the record specified by the key argument and returns a value With Plock()
func (db *MDBM) FetchWithPlock(key interface{}, lockflags int) (string, error) {
	return db.fetchWithAnyLock(key, lockTypePlock, C.int(lockflags))
}

// FetchWithTryLock returns fetche the record specified by the key argument and returns a value With TryLock()
func (db *MDBM) FetchWithTryLock(key interface{}) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeTryLock, lockFlagsSkip)
}

// FetchWithTryLockSmart returns fetche the record specified by the key argument and returns a value With TryLockSmart()
func (db *MDBM) FetchWithTryLockSmart(key interface{}, lockflags int) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeTrySmart, C.int(lockflags))
}

// FetchWithTryLockShared returns fetche the record specified by the key argument and returns a value with TryShared()
func (db *MDBM) FetchWithTryLockShared(key interface{}) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeTryShared, lockFlagsSkip)
}

// FetchWithTryPlock returns fetche the record specified by the key argument and returns a value With TryPlock()
func (db *MDBM) FetchWithTryPlock(key interface{}, lockflags int) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeTryPlock, C.int(lockflags))
}

// Fetch fetchs the record specified by the key and val parameters.
func (db *MDBM) Fetch(key interface{}) (string, error) {
	return db.fetchWithAnyLock(key, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) fetchRWithAnyLock(key interface{}, iter *C.MDBM_ITER, lockType C.int, lockFlags C.int) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	err := db.checkAvailable()
	if err != nil {
		return rv, retval, db.convertIter(iter), err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.get_mdbm_fetch_r_with_lock(db.pmdbm, &k, &v, iter, lockType, lockFlags)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)

	goiter := db.convertIter(iter)

	return rv, retval, goiter, err
}

// FetchRWithLock returns fetche the record specified by the key argument and returns a value With Lock()
func (db *MDBM) FetchRWithLock(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeLock, lockFlagsSkip)
}

// FetchRWithLockSmart returns fetche the record specified by the key argument and returns a value With LocckSmart()
func (db *MDBM) FetchRWithLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeSmart, C.int(lockflags))
}

// FetchRWithLockShared returns fetche the record specified by the key argument and returns a value With LockShared()
func (db *MDBM) FetchRWithLockShared(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeShared, lockFlagsSkip)
}

// FetchRWithPlock returns fetche the record specified by the key argument and returns a value With Plock()
func (db *MDBM) FetchRWithPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypePlock, C.int(lockflags))
}

// FetchRWithTryLock returns fetche the record specified by the key argument and returns a value With TryLock()
func (db *MDBM) FetchRWithTryLock(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeTryLock, lockFlagsSkip)
}

// FetchRWithTryLockSmart returns fetche the record specified by the key argument and returns a value With TryLockSmart()
func (db *MDBM) FetchRWithTryLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeTrySmart, C.int(lockflags))
}

// FetchRWithTryLockShared returns fetche the record specified by the key argument and returns a value with TryShared()
func (db *MDBM) FetchRWithTryLockShared(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeTryShared, lockFlagsSkip)
}

// FetchRWithTryPlock returns fetche the record specified by the key argument and returns a value With TryPlock()
func (db *MDBM) FetchRWithTryPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeTryPlock, C.int(lockflags))
}

// FetchR fetchs the record specified by the key and val parameters.
func (db *MDBM) FetchR(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchRWithAnyLock(key, iter, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) fetchStrWithAnyLock(key interface{}, lockType C.int, lockFlags C.int) (string, error) {

	var retval string

	err := db.checkAvailable()
	if err != nil {
		return retval, err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return retval, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)
	defer C.free(unsafe.Pointer(k))

	_, _, err = db.cgoRun(func() (int, error) {

		val, err := C.get_mdbm_fetch_str_with_lock(db.pmdbm, k, lockType, lockFlags)
		retval = C.GoString(val)
		return 0, err
	})

	return retval, err
}

// FetchStrWithLock fetchs the record specified by the key and val parameters With Lock()
func (db *MDBM) FetchStrWithLock(key interface{}) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeLock, lockFlagsSkip)
}

// FetchStrWithLockSmart fetchs the record specified by the key and val parameters With LockSmart()
func (db *MDBM) FetchStrWithLockSmart(key interface{}, lockflags int) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeSmart, C.int(lockflags))
}

// FetchStrWithLockShared fetchs the record specified by the key and val parameters With LockShared()
func (db *MDBM) FetchStrWithLockShared(key interface{}) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeShared, lockFlagsSkip)
}

// FetchStrWithPlock fetchs the record specified by the key and val parameters With Plock()
func (db *MDBM) FetchStrWithPlock(key interface{}, lockflags int) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypePlock, C.int(lockflags))
}

// FetchStrWithTryLock fetchs the record specified by the key and val parameters With Lock()
func (db *MDBM) FetchStrWithTryLock(key interface{}) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeTryLock, lockFlagsSkip)
}

// FetchStrWithTryLockSmart fetchs the record specified by the key and val parameters With LockSmart()
func (db *MDBM) FetchStrWithTryLockSmart(key interface{}, lockflags int) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeTrySmart, C.int(lockflags))
}

// FetchStrWithTryLockShared fetchs the record specified by the key and val parameters With LockShared()
func (db *MDBM) FetchStrWithTryLockShared(key interface{}) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeTryShared, lockFlagsSkip)
}

// FetchStrWithTryPlock fetchs the record specified by the key and val parameters With Plock()
func (db *MDBM) FetchStrWithTryPlock(key interface{}, lockflags int) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeTryPlock, C.int(lockflags))
}

// FetchStr fetchs the record specified by the key and val parameters With Plock()
func (db *MDBM) FetchStr(key interface{}) (string, error) {
	return db.fetchStrWithAnyLock(key, lockTypeNone, lockFlagsSkip)
}

// FetchDupR fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
func (db *MDBM) fetchDupRWithAnyLock(key interface{}, iter *C.MDBM_ITER, lockType C.int, lockFlags C.int) (int, string, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string

	err := db.checkAvailable()
	if err != nil {
		return rv, retval, db.convertIter(iter), err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, retval, db.convertIter(iter), errors.Wrapf(err, "failured")
	}

	var k, v C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))

	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.get_mdbm_fetch_dup_r_with_lock(db.pmdbm, &k, &v, iter, lockType, lockFlags)
		return int(rv), err
	})

	retval = C.GoStringN(v.dptr, v.dsize)
	return rv, retval, db.convertIter(iter), err
}

// FetchDupRWithLock fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With Lock()
func (db *MDBM) FetchDupRWithLock(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeLock, lockFlagsSkip)
}

// FetchDupRWithLockSmart fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With LockSmart()
func (db *MDBM) FetchDupRWithLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeSmart, C.int(lockflags))
}

// FetchDupRWithLockShared fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With LockShared()
func (db *MDBM) FetchDupRWithLockShared(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeShared, lockFlagsSkip)
}

// FetchDupRWithPlock fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With Plock()
func (db *MDBM) FetchDupRWithPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypePlock, C.int(lockflags))
}

// FetchDupRWithTryLock fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With TryLock()
func (db *MDBM) FetchDupRWithTryLock(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeTryLock, lockFlagsSkip)
}

// FetchDupRWithTryLockSmart fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With TryLockSmart()
func (db *MDBM) FetchDupRWithTryLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeTrySmart, C.int(lockflags))
}

// FetchDupRWithTryLockShared fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With TryLockShared()
func (db *MDBM) FetchDupRWithTryLockShared(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeTryShared, lockFlagsSkip)
}

// FetchDupRWithTryPlock fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
// With TryPlock()
func (db *MDBM) FetchDupRWithTryPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeTryPlock, C.int(lockflags))
}

// FetchDupR fetches the next value for a key inserted via StoreR() with the mdbm.InsertDup flag set
// The order of values returned by iterating via this function is not guaranteed to be the same order as the values were inserted.
// As with any db iteration, record insertion and deletion during iteration may cause the iteration to skip and/or repeat records.
// Calling this function with an iterator initialized via iter(=mdbm.GetNewIter()) will cause this function to return the first value for the given key.
func (db *MDBM) FetchDupR(key interface{}, iter *C.MDBM_ITER) (int, string, Iter, error) {
	return db.fetchDupRWithAnyLock(key, iter, lockTypeNone, lockFlagsSkip)
}

// FetchInfo ...
func (db *MDBM) FetchInfo(key interface{}, sbuf *string, iter *C.MDBM_ITER) (int, string, FetchInfo, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	var retval string
	var info C.struct_mdbm_fetch_info

	err := db.checkAvailable()
	if err != nil {
		return rv, retval, db.convertFetchInfo(info), db.convertIter(iter), err
	}

	err = db.isVersion3Above()
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
func (db *MDBM) deleteWithAnyLock(key interface{}, lockType C.int, lockFlags C.int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	bkey, err := db.convertToArByte(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	var k C.datum
	k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
	k.dsize = C.int(len(bkey))

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_delete_with_lock(db.pmdbm, k, lockType, lockFlags)
		return int(rv), err
	})

	return rv, err
}

// DeleteWithLock deletes a specific record with Lock()
func (db *MDBM) DeleteWithLock(key interface{}) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeLock, lockFlagsSkip)
}

// DeleteWithLockSmart deletes a specific record with LockSmart()
func (db *MDBM) DeleteWithLockSmart(key interface{}, lockflags int) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeSmart, C.int(lockflags))
}

// DeleteWithLockShared deletes a specific record with LockShared()
func (db *MDBM) DeleteWithLockShared(key interface{}) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeShared, lockFlagsSkip)
}

// DeleteWithPlock deletes a specific record with Plock()
func (db *MDBM) DeleteWithPlock(key interface{}, lockflags int) (int, error) {
	return db.deleteWithAnyLock(key, lockTypePlock, C.int(lockflags))
}

// DeleteWithTryLock deletes a specific record With TryLock()
func (db *MDBM) DeleteWithTryLock(key interface{}) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeTryLock, lockFlagsSkip)
}

// DeleteWithTryLockSmart deletes a specific record With TryLockSmart()
func (db *MDBM) DeleteWithTryLockSmart(key interface{}, lockflags int) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeTrySmart, C.int(lockflags))
}

// DeleteWithTryLockShared deletes a specific record With TryLockShared()
func (db *MDBM) DeleteWithTryLockShared(key interface{}) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeTryShared, lockFlagsSkip)
}

// DeleteWithTryPlock deletes a specific record With TryPlock()
func (db *MDBM) DeleteWithTryPlock(key interface{}, lockflags int) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeTryPlock, C.int(lockflags))
}

// Delete Deletes the record specified by the key and val parameters.
func (db *MDBM) Delete(key interface{}) (int, error) {
	return db.deleteWithAnyLock(key, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) deleteRWithAnyLock(key interface{}, iter *C.MDBM_ITER, lockType C.int, lockFlags C.int) (int, Iter, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertIter(iter), err
	}

	var k C.datum

	if key != nil {
		bkey, err := db.convertToArByte(key)
		if err != nil {
			return rv, db.convertIter(iter), errors.Wrapf(err, "failured")
		}

		k.dptr = (*C.char)(unsafe.Pointer(&bkey[0]))
		k.dsize = C.int(len(bkey))

	} else {
		k.dptr = nil
		k.dsize = C.int(-1)
	}

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_delete_r_with_lock(db.pmdbm, k, iter, lockType, lockFlags)
		return int(rv), err
	})

	return rv, db.convertIter(iter), err
}

// DeleteRWithLock deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With Lock()
func (db *MDBM) DeleteRWithLock(key interface{}, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeLock, lockFlagsSkip)
}

// DeleteRWithLockSmart deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With LockSmart()
func (db *MDBM) DeleteRWithLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeSmart, C.int(lockflags))
}

// DeleteRWithLockShared deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With LockShared()
func (db *MDBM) DeleteRWithLockShared(key interface{}, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeShared, lockFlagsSkip)
}

// DeleteRWithPlock deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With Plock()
func (db *MDBM) DeleteRWithPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypePlock, C.int(lockflags))
}

// DeleteRWithTryLock deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With TryLock()
func (db *MDBM) DeleteRWithTryLock(key interface{}, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeTryLock, lockFlagsSkip)
}

// DeleteRWithTryLockSmart deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With TryLockSmart()
func (db *MDBM) DeleteRWithTryLockSmart(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeTrySmart, C.int(lockflags))
}

// DeleteRWithTryLockShared deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With TryLockSahred()
func (db *MDBM) DeleteRWithTryLockShared(key interface{}, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeTryShared, lockFlagsSkip)
}

// DeleteRWithTryPlock deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
// With TryPlock()
func (db *MDBM) DeleteRWithTryPlock(key interface{}, iter *C.MDBM_ITER, lockflags int) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeTryPlock, C.int(lockflags))
}

// DeleteR deletes the record currently addressed by the iter argument.
// After deletion, the key and/or value returned by the iterating function is no longer valid.
// Calling NextR() on the iterator will return the key/value for the entry following the entry that was deleted.
func (db *MDBM) DeleteR(key interface{}, iter *C.MDBM_ITER) (int, Iter, error) {
	return db.deleteRWithAnyLock(key, iter, lockTypeNone, lockFlagsSkip)
}

func (db *MDBM) deleteStrWithAnyLock(key interface{}, lockType C.int, lockFlags C.int) (int, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	rv := -1
	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return rv, errors.Wrapf(err, "failured")
	}

	k := C.CString(skey)

	rv, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_delete_str_with_lock(db.pmdbm, k, lockType, lockFlags)
		return int(rv), err
	})

	return rv, err
}

// DeleteStrWithLock deletes a string from the MDBM with Lock()
func (db *MDBM) DeleteStrWithLock(key interface{}) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeLock, lockFlagsSkip)
}

// DeleteStrWithLockSmart deletes a string from the MDBM with LockSmart()
func (db *MDBM) DeleteStrWithLockSmart(key interface{}, lockflags int) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeSmart, C.int(lockflags))
}

// DeleteStrWithLockShared deletes a string from the MDBM with LockShared()
func (db *MDBM) DeleteStrWithLockShared(key interface{}) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeShared, lockFlagsSkip)
}

// DeleteStrWithPlock deletes a string from the MDBM with Plock()
func (db *MDBM) DeleteStrWithPlock(key interface{}, lockflags int) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypePlock, C.int(lockflags))
}

// DeleteStrWithTryLock deletes a string from the MDBM with TryLock()
func (db *MDBM) DeleteStrWithTryLock(key interface{}) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeTryLock, lockFlagsSkip)
}

// DeleteStrWithTryLockSmart deletes a string from the MDBM with TryLockSmart()
func (db *MDBM) DeleteStrWithTryLockSmart(key interface{}, lockflags int) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeTrySmart, C.int(lockflags))
}

// DeleteStrWithTryLockShared deletes a string from the MDBM with TryLockShared()
func (db *MDBM) DeleteStrWithTryLockShared(key interface{}) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeTryShared, lockFlagsSkip)
}

// DeleteStrWithTryPlock deletes a string from the MDBM with TryPlock()
func (db *MDBM) DeleteStrWithTryPlock(key interface{}, lockflags int) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeTryPlock, C.int(lockflags))
}

// DeleteStr deletes string from the MDBM
func (db *MDBM) DeleteStr(key interface{}) (int, error) {
	return db.deleteStrWithAnyLock(key, lockTypeNone, lockFlagsSkip)
}

// First returns the first key/value pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) First() (string, string, error) {

	var kv C.kvpair
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", "", err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_first(db.pmdbm)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", errors.New("database is empty")
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

	err = db.checkAvailable()
	if err != nil {
		return "", "", err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_next(db.pmdbm)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", errors.New("database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	return key, val, nil
}

// FirstR returns the first key/value pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) FirstR(iter *C.MDBM_ITER) (string, string, Iter, error) {

	var kv C.kvpair
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", "", db.convertIter(iter), err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_first_r(db.pmdbm, iter)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", db.convertIter(iter), errors.New("database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	goiter := db.convertIter(iter)

	return key, val, goiter, nil
}

// NextR Fetches the next record in an MDBM.
// Returns the next key/value pair from the db, based on the iterator.
func (db *MDBM) NextR(iter *C.MDBM_ITER) (string, string, Iter, error) {

	var kv C.kvpair
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", "", db.convertIter(iter), err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		kv, err = C.mdbm_next_r(db.pmdbm, iter)
		return 0, err
	})

	if int(kv.key.dsize) == 0 {
		return "", "", db.convertIter(iter), errors.New("database is empty")
	}

	key := C.GoStringN(kv.key.dptr, kv.key.dsize)
	val := C.GoStringN(kv.val.dptr, kv.val.dsize)

	goiter := db.convertIter(iter)

	return key, val, goiter, nil
}

// FirstKey Returns the first key from the database.
// The order that records are returned is not specified.
func (db *MDBM) FirstKey() (string, error) {

	var k C.datum
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_firstkey(db.pmdbm)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", errors.New("database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)

	return key, err
}

// NextKey Returns the next key pair from the database.
// The order that records are returned is not specified.
func (db *MDBM) NextKey() (string, error) {

	var k C.datum
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_nextkey(db.pmdbm)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", errors.New("database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)

	return key, err
}

// FirstKeyR fetches the first key in an MDBM.
// Initializes the iterator, and returns the first key from the db.
// Subsequent calls to NextR() or NextKeyR() With this iterator will loop through the entire db.
func (db *MDBM) FirstKeyR(iter *C.MDBM_ITER) (string, Iter, error) {

	var k C.datum
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", db.convertIter(iter), err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_firstkey_r(db.pmdbm, iter)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", db.convertIter(iter), errors.New("database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)
	goiter := db.convertIter(iter)

	return key, goiter, nil
}

// NextKeyR fetches the next key in an MDBM.  Returns the next key from the db.
// Subsequent calls to NextR() or NextKeyR() With this iterator
// will loop through the entire db.
func (db *MDBM) NextKeyR(iter *C.MDBM_ITER) (string, Iter, error) {

	var k C.datum
	var err error

	err = db.checkAvailable()
	if err != nil {
		return "", db.convertIter(iter), err
	}

	_, _, err = db.cgoRun(func() (int, error) {
		k, err = C.mdbm_nextkey_r(db.pmdbm, iter)
		return 0, err
	})

	if int(k.dsize) == 0 {
		return "", db.convertIter(iter), errors.New("database is empty")
	}

	key := C.GoStringN(k.dptr, k.dsize)
	goiter := db.convertIter(iter)

	return key, goiter, nil
}

// GetCacheMode returns the current cache style of the database.
// See the cachemode parameter in SetCacheMode() for the valid values.
func (db *MDBM) GetCacheMode() (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	err = db.isVersion3Above()
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
// Tracking metadata is stored With each entry which allows MDBM to do cache eviction via LRU, LFU, and GDSF
// (greedy-dual-size-frequency). MDBM also supports clean/dirty tracking and the application can supply a callback (see SetBackingStore())
// which is called by MDBM when a dirty entry is about to be evicted allowing
// the application to sync the entry to a backing store or perform some other type of "clean" operation.
func (db *MDBM) SetCacheMode(cachemode int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	err = db.isVersion3Above()
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

	return rv, errors.Wrapf(err, "SetCachemode must be called before data is inserted.")
}

// GetCacheModeName returns the cache mode as a string. See SetCacheMode()
func (db *MDBM) GetCacheModeName(cachemode int) (string, error) {

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	err = db.isVersion3Above()
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

	err = db.checkAvailable()
	if err != nil {
		return 0, err
	}

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

	err = db.checkAvailable()
	if err != nil {
		return 0, err
	}

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

	err = db.checkAvailable()
	if err != nil {
		return 0, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_preload(db.pmdbm)
		return int(rv), err
	})

	return rv, err
}

// LockDump returns the state of lock
func (db *MDBM) LockDump() (string, error) {

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	_, out, err := db.cgoRunCapture(func() (int, error) {
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

	err = db.checkAvailable()
	if err != nil {
		return -1, err
	}

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

	err = db.checkAvailable()
	if err != nil {
		return -1, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return -1, "", err
	}

	rv, out, err := db.cgoRunCapture(func() (int, error) {
		rv, err := C.mdbm_chk_page(db.pmdbm, C.int(pagenum))
		return int(rv), err
	})

	return rv, out, err
}

// ChkError checks integrity of an entry on a page.
// NOTE: This has not been implemented.
func (db *MDBM) ChkError(pagenum int, mappedpagenum int, index int) error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	_, out, err := db.cgoRunCapture(func() (int, error) {
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

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	_, out, err := db.cgoRunCapture(func() (int, error) {
		_, err := C.mdbm_dump_page(db.pmdbm, C.int(pno))
		return 0, err
	})

	return out, err
}

// ResetStatOperations resets the stat counter and last-time performed for fetch, store, and remove operations.
func (db *MDBM) ResetStatOperations() error {

	err := db.checkAvailable()
	if err != nil {
		return err
	}

	_, _, err = db.cgoRunCapture(func() (int, error) {
		_, err := C.mdbm_reset_stat_operations(db.pmdbm)
		return 0, err
	})

	return err
}

// EnableStatOperations enables and disables gathering of stat counters and/or last-time performed for fetch, store, and remove operations.
func (db *MDBM) EnableStatOperations(flags int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return -1, 0, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	//stype between 0 (MDBM_STAT_TAG_FETCH) and 16 (MDBM_STAT_TAG_DELETE_FAILED) on mdbm 4.x
	var retval string
	_, _, err = db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stat_name(C.int(stype))
		retval = C.GoString(rv)

		return 0, err
	})

	return retval, err
}

// GetStatTime Gets the last time when an type operation was performed.
func (db *MDBM) GetStatTime(stype int) (int, uint64, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, 0, err
	}

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

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	if flags < ClockStandard || flags > ClockTSC {
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

	err := db.checkAvailable()
	if err != nil {
		return "", err
	}

	err = db.isVersion2()
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

// GetStats gets a a stats block With individual stat values.
func (db *MDBM) GetStats() (int, Stats, error) {

	var stats C.mdbm_stats_t

	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertStats(stats), err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_stats(db.pmdbm, &stats, C.sizeof_mdbm_stats_t)
		return int(rv), err
	})

	return rv, db.convertStats(stats), err
}

// GetDBInfo gets configuration information about a database
func (db *MDBM) GetDBInfo() (int, DBInfo, error) {

	var info C.mdbm_db_info_t

	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertDBInfo(info), err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_db_info(db.pmdbm, &info)
		return int(rv), err
	})

	return rv, db.convertDBInfo(info), err
}

// GetDBStats gets overall database stats.
func (db *MDBM) GetDBStats(flags int) (int, DBInfo, StatInfo, error) {

	var dbinfo C.mdbm_db_info_t
	var statinfo C.mdbm_stat_info_t

	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertDBInfo(dbinfo), db.convertStatInfo(statinfo), err
	}

	if flags != StatNolock && flags > IterateNolock {
		return -1, DBInfo{}, StatInfo{}, fmt.Errorf("not support flags=%d", flags)
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_get_db_stats(db.pmdbm, &dbinfo, &statinfo, C.int(flags))
		return int(rv), err
	})

	return rv, db.convertDBInfo(dbinfo), db.convertStatInfo(statinfo), err
}

// GetWindowStats retrieves statistics about windowing usage on the associated database.
func (db *MDBM) GetWindowStats() (int, WindowStats, error) {

	var stats C.mdbm_window_stats_t

	err := db.checkAvailable()
	if err != nil {
		return -1, db.convertWindowStat(stats), err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithPlock API
func (db *MDBM) Plock(key interface{}) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithPlock API
func (db *MDBM) Punlock(key interface{}) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithTryPlock API
func (db *MDBM) TryPlock(key interface{}) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithLockSmart API
func (db *MDBM) LockSmart(key interface{}, flags int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithLockSmart API
func (db *MDBM) UnLockSmart(key interface{}, flags int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

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
// NOTE: Unstable working on the multiple-thread(goroutine), please use the {Store|Fetch}WithTryLockSmart API
func (db *MDBM) TryLockSmart(key interface{}, flags int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	skey, err := db.convertToString(key)
	if err != nil {
		return -1, errors.Wrapf(err, "failured")
	}

	var k C.datum

	k.dptr = C.CString(skey)
	k.dsize = C.int(len(skey))
	defer C.free(unsafe.Pointer(k.dptr))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_trylock_smart(db.pmdbm, &k, C.int(flags))
		return int(rv), err
	})

	return rv, err
}

// CheckResidency checks mdbm page residency: count the number of DB pages mapped into memory.
// The counts are in units of the system-page-size (typically 4k)
func (db *MDBM) CheckResidency() (int, uint32, uint32, error) {

	var pgsin, pgsout C.mdbm_ubig_t

	err := db.checkAvailable()
	if err != nil {
		return -1, uint32(pgsin), uint32(pgsout), err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_check_residency(db.pmdbm, &pgsin, &pgsout)
		return int(rv), err
	})

	return rv, uint32(pgsin), uint32(pgsout), err
}

//EasyGetNumOfRows returns the number of rows in the MDBM file
func (db *MDBM) EasyGetNumOfRows() (uint64, error) {

	var cnt uint64

	err := db.checkAvailable()
	if err != nil {
		return 0, err
	}

	key, _, err := db.First()
	if err != nil || len(key) < 1 {
		return cnt, errors.Wrapf(err, "failed, can't run mdbm.First()")
	}

	for {

		key, _, err := db.Next()
		if err != nil || len(key) < 1 {
			break
		}
		cnt++
	}

	return cnt, nil
}

//EasyGetKeyList returns the list of key in the MDBM file
func (db *MDBM) EasyGetKeyList() ([]string, error) {

	var retval []string

	err := db.checkAvailable()
	if err != nil {
		return retval, err
	}

	key, _, err := db.First()
	if err != nil || len(key) < 1 {
		return retval, errors.Wrapf(err, "failed, can't run mdbm.First()")
	}

	retval = append(retval, key)

	for {

		key, _, err := db.Next()
		if err != nil || len(key) < 1 {
			break
		}

		retval = append(retval, key)
	}

	return retval, nil
}

// Clean does mark entries clean/re-usable in the database for the specified page. If pagenum is -1, then clean all pages.
// NOTE: V3 API
func (db *MDBM) Clean(pagenum int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	err = db.isVersion3Above()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_clean(db.pmdbm, C.int(pagenum), C.int(0)) //flags ignored
		return int(rv), err
	})

	return rv, err
}

// PreSplit forces a db to split, creating N pages.  Must be called before any data is inserted. If N is not a multiple of 2, it will be rounded up.
// N Target number of pages post split. If N is not larger than the initial size (ex., 0), a split will not be done and a success status is returned.
func (db *MDBM) PreSplit(split uint32) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_pre_split(db.pmdbm, C.mdbm_ubig_t(split))
		return int(rv), err
	})

	return rv, err
}

// Fcopy copies the contents of a database to an open file handle.
// NOTE: lock for the duration
func (db *MDBM) Fcopy(filepath string, mode int) (int, error) {

	err := db.checkAvailable()
	if err != nil {
		return -1, err
	}

	fpath := C.CString(filepath)
	defer C.free(unsafe.Pointer(fpath))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.set_mdbm_fcopy(db.pmdbm, fpath, C.int(mode))
		return int(rv), err
	})

	return rv, err
}

// SparsifyFile makes a file sparse. Read every blockssize bytes and for all-zero blocks punch a hole in the file.
// This can make a file with lots of zero-bytes use less disk space.
// NOTE: For MDBM files that may be modified, the DB should be opened, and exclusive-locked for the duration of the sparsify operation.
// NOTE: This function is linux-only.
// blocksize Minimum size to consider for hole-punching, <=0 to use the system block-size
func (db *MDBM) SparsifyFile(filepath string, blocksize int) (int, error) {

	fpath := C.CString(filepath)
	defer C.free(unsafe.Pointer(fpath))

	rv, _, err := db.cgoRun(func() (int, error) {
		rv, err := C.mdbm_sparsify_file(fpath, C.int(blocksize))
		return int(rv), err
	})

	return rv, err
}
