package mdbm

/*
#include <mdbm.h>
#include <mdbm_log.h>
*/
import "C"

const (
	// Rdonly is Read-only access
	Rdonly = C.MDBM_O_RDONLY

	// Wronly is Write-only access (deprecated in V3)
	Wronly = C.MDBM_O_WRONLY

	// Rdrw is Read and write access
	Rdrw = C.MDBM_O_RDWR

	// Create is Create file if it does not exist
	Create = C.MDBM_O_CREAT

	// Trunc is Truncate file
	Trunc = C.MDBM_O_TRUNC

	// Fsync is Sync file on close
	Fsync = C.MDBM_O_FSYNC

	// Async is Perform asynchronous writes
	Async = C.MDBM_O_ASYNC

	// Direct is Perform direction I/O
	Direct = C.MDBM_O_DIRECT

	// NoDirty is Do not not track clean/dirty status
	NoDirty = C.MDBM_NO_DIRTY

	// SingleArch is User *promises* not to mix 32/64-bit access
	SingleArch = C.MDBM_SINGLE_ARCH

	// OpenWindowed is Use windowing to access db
	OpenWindowed = C.MDBM_OPEN_WINDOWED

	// Protect is Protect database except when locked
	Protect = C.MDBM_PROTECT

	// DBSizeMB is Dbsize is specific in MB
	DBSizeMB = C.MDBM_DBSIZE_MB

	// StatOperations is collect stats for fetch, store, delete
	StatOperations = C.MDBM_STAT_OPERATIONS

	// LargeObjects is Support large objects - obsolete
	LargeObjects = C.MDBM_LARGE_OBJECTS

	// PartitionedLocks is Partitioned locks
	PartitionedLocks = C.MDBM_PARTITIONED_LOCKS

	// RwLocks is Read-write locks
	RwLocks = C.MDBM_RW_LOCKS

	// AnyLocks is Open, even if existing locks don't match flags
	AnyLocks = C.MDBM_ANY_LOCKS

	// OpenNolock is Don't lock during open
	OpenNolock = C.MDBM_OPEN_NOLOCK

	// LogEmerg is system is unusable
	LogEmerg = C.LOG_EMERG

	// LogAlert is action must be taken immediately
	LogAlert = C.LOG_ALERT

	// LogCrit is critical conditions
	LogCrit = C.LOG_CRIT

	// LogErr is error conditions
	LogErr = C.LOG_ERR

	// LogWarning is warning conditions
	LogWarning = C.LOG_WARNING

	// LogNotice is normal but significant condition
	LogNotice = C.LOG_NOTICE

	// LogInfo is informational
	LogInfo = C.LOG_INFO

	// LogDebug is debug-level messages
	LogDebug = C.LOG_DEBUG

	// HashCRC32 is table based 32bit crc
	HashCRC32 = C.MDBM_HASH_CRC32

	// HashEJB is from hsearch
	HashEJB = C.MDBM_HASH_EJB

	// HashPHONG is congruential hash
	HashPHONG = C.MDBM_HASH_PHONG

	// HashOZ is from sdbm
	HashOZ = C.MDBM_HASH_OZ

	// HashTOREK is from Berkeley db
	HashTOREK = C.MDBM_HASH_TOREK

	// HashFNV is Fowler/Vo/Noll hash
	HashFNV = C.MDBM_HASH_FNV

	// HashSTL is STL string hash
	HashSTL = C.MDBM_HASH_STL

	// HashMD5 is MD5
	HashMD5 = C.MDBM_HASH_MD5

	// HashSHA1 is SHA_1
	HashSHA1 = C.MDBM_HASH_SHA_1

	// HashJENKINS is JENKINS
	HashJENKINS = C.MDBM_HASH_JENKINS

	// HashHSIEH is HSIEH SuperFastHash
	HashHSIEH = C.MDBM_HASH_HSIEH

	// MaxHash is bump up if adding more
	MaxHash = C.MDBM_MAX_HASH

	// DefaultHash is MDBM_HASH_FNV is best
	DefaultHash = C.MDBM_DEFAULT_HASH

	// Align8Bits is 1-Byte data alignment
	Align8Bits = C.MDBM_ALIGN_8_BITS

	// Align16Bits is 2-Byte data alignment
	Align16Bits = C.MDBM_ALIGN_16_BITS

	// Align32Bits is 4-Byte data alignment
	Align32Bits = C.MDBM_ALIGN_32_BITS

	// Align64Bits is 8-Byte data alignment
	Align64Bits = C.MDBM_ALIGN_64_BITS

	// Magic is V2 file identifier
	Magic = C._MDBM_MAGIC

	// MagicNew is V2 file identifier with large objects
	MagicNew = C._MDBM_MAGIC_NEW

	// MagicNew2 is V3 file identifier
	MagicNew2 = C._MDBM_MAGIC_NEW2

	// ProtNone is Page no access
	ProtNone = C.MDBM_PROT_NONE

	// ProtRead is Page read access
	ProtRead = C.MDBM_PROT_READ

	// ProtWrite is Page write access
	ProtWrite = C.MDBM_PROT_WRITE

	// ProtNoaccess is Page no access
	ProtNoaccess = C.MDBM_PROT_NOACCESS

	// ProtAccess is Page protection mask
	ProtAccess = C.MDBM_PROT_ACCESS

	// Insert is Insert if key does not exist; fail if exists
	Insert = C.MDBM_INSERT

	// Replace is Update if key exists; insert if does not exist
	Replace = C.MDBM_REPLACE

	// InsertDup is Insert new record (creates duplicate if key exists)
	InsertDup = C.MDBM_INSERT_DUP

	// Modify is Update if key exists; fail if does not exist
	Modify = C.MDBM_MODIFY

	// StoreMask is Mask for all store options
	StoreMask = C.MDBM_STORE_MASK

	// CacheModeNone is No caching behavior
	CacheModeNone = C.MDBM_CACHEMODE_NONE

	// CacheModeLfu is Entry with smallest number of accesses is evicted
	CacheModeLfu = C.MDBM_CACHEMODE_LFU

	// CacheModeLru is Entry with oldest access time is evicted
	CacheModeLru = C.MDBM_CACHEMODE_LRU

	// CacheModeGdsf is Greedy dual-size frequency (size and frequency) eviction
	CacheModeGdsf = C.MDBM_CACHEMODE_GDSF

	// CacheModeMax is Maximum cache mode value
	CacheModeMax = C.MDBM_CACHEMODE_MAX

	// StatsBasic is enables gathering only the stats counters.
	StatsBasic = C.MDBM_STATS_BASIC

	// StatsTimed is enables gathering only the stats timestamps.
	StatsTimed = C.MDBM_STATS_TIMED

	// StatTypeFetch is fetch* operations
	StatTypeFetch = C.MDBM_STAT_TYPE_FETCH

	// StatTypeStore is store* operations
	StatTypeStore = C.MDBM_STAT_TYPE_STORE

	// StatTypeDelete is delete* operations
	StatTypeDelete = C.MDBM_STAT_TYPE_DELETE

	// StatTypeMax is C.MDBM_STAT_TYPE_DELETE
	StatTypeMax = C.MDBM_STAT_TYPE_MAX

	// ClockTsc is Enables use of TSC
	ClockTsc = C.MDBM_CLOCK_TSC

	// ClockStandard is Disables use of TSC
	ClockStandard = C.MDBM_CLOCK_STANDARD

	// StatTagFetch is Successful fetch stats-callback counter
	StatTagFetch = C.MDBM_STAT_TAG_FETCH

	// StatTagStore is Successful store stats-callback counter
	StatTagStore = C.MDBM_STAT_TAG_STORE

	// StatTagDelete is Successful delete stats-callback counter
	StatTagDelete = C.MDBM_STAT_TAG_DELETE

	// StatTagLock is lock stats-callback counter (not implemented)
	StatTagLock = C.MDBM_STAT_TAG_LOCK

	// StatTagFetchUncached is Cache-miss with cache+backingstore
	StatTagFetchUncached = C.MDBM_STAT_TAG_FETCH_UNCACHED

	// StatTagGetpage is Generic access counter in windowed mode
	StatTagGetpage = C.MDBM_STAT_TAG_GETPAGE

	// StatTagGetpageUncached is Windowed-mode "miss" (load new page into window)
	StatTagGetpageUncached = C.MDBM_STAT_TAG_GETPAGE_UNCACHED

	// StatTagCacheEvict is Cache evict stats-callback counter
	StatTagCacheEvict = C.MDBM_STAT_TAG_CACHE_EVICT

	// StatTagCacheStore is Successful cache store counter (BS only)
	StatTagCacheStore = C.MDBM_STAT_TAG_CACHE_STORE

	// StatTagPageStore is Successful page-level store indicator
	StatTagPageStore = C.MDBM_STAT_TAG_PAGE_STORE

	// StatTagPageDelete is Successful page-level delete indicator
	StatTagPageDelete = C.MDBM_STAT_TAG_PAGE_DELETE

	// StatTagSync is Counter of mdbm_syncs and fsyncs
	StatTagSync = C.MDBM_STAT_TAG_SYNC

	// StatTagFetchNotFound is Fetch cannot find a key in MDBM
	StatTagFetchNotFound = C.MDBM_STAT_TAG_FETCH_NOT_FOUND

	// StatTagFetchError is Error occurred during fetch
	StatTagFetchError = C.MDBM_STAT_TAG_FETCH_ERROR

	// StatTagStoreError is Error occurred during store (e.g. MODIFY failed)
	StatTagStoreError = C.MDBM_STAT_TAG_STORE_ERROR

	// StatTagDeleteFailed is Delete failed: cannot find a key in MDBM
	StatTagDeleteFailed = C.MDBM_STAT_TAG_DELETE_FAILED

	// StatTagFetchLatency is Fetch latency (expensive to collect)
	StatTagFetchLatency = C.MDBM_STAT_TAG_FETCH_LATENCY

	// StatTagStoreLatency is Store latency (expensive to collect)
	StatTagStoreLatency = C.MDBM_STAT_TAG_STORE_LATENCY

	// StatTagDeleteLatency is Delete latency (expensive to collect)
	StatTagDeleteLatency = C.MDBM_STAT_TAG_DELETE_LATENCY

	// StatTagFetchTime is timestamp of last fetch (not yet implemented)
	StatTagFetchTime = C.MDBM_STAT_TAG_FETCH_TIME

	// StatTagStoreTime is timestamp of last store (not yet implemented)
	StatTagStoreTime = C.MDBM_STAT_TAG_STORE_TIME

	// StatTagDeleteTime is timestamp of last delete (not yet implemented)
	StatTagDeleteTime = C.MDBM_STAT_TAG_DELETE_TIME

	// StatTagFetchUncachedLatency is Cache miss latency for cache+Backingstore only (expensive to collect)
	StatTagFetchUncachedLatency = C.MDBM_STAT_TAG_FETCH_UNCACHED_LATENCY

	// StatTagGetpageLatency is access latency in windowed mode (expensive to collect)
	StatTagGetpageLatency = C.MDBM_STAT_TAG_GETPAGE_LATENCY

	// StatTagGetpageUncachedLatency is windowed-mode miss latency (expensive to collect)
	StatTagGetpageUncachedLatency = C.MDBM_STAT_TAG_GETPAGE_UNCACHED_LATENCY

	// StatTagCacheEvictLatency is cache evict latency (expensive to collect)
	StatTagCacheEvictLatency = C.MDBM_STAT_TAG_CACHE_EVICT_LATENCY

	// StatTagCacheStoreLatency is Cache store latency in Cache+backingstore mode only (expensive to collect)
	StatTagCacheStoreLatency = C.MDBM_STAT_TAG_CACHE_STORE_LATENCY

	// StatTagPageStoreValue is Indicates a delete occurred on a particular page.
	StatTagPageStoreValue = C.MDBM_STAT_TAG_PAGE_STORE_VALUE

	// StatTagPageDeleteValue is Indicates a delete occurred on a particular page.
	StatTagPageDeleteValue = C.MDBM_STAT_TAG_PAGE_DELETE_VALUE

	// StatTagSyncLatency is mdbm_sync/fsync latency (expensive to collect)
	StatTagSyncLatency = C.MDBM_STAT_TAG_SYNC_LATENCY

	// IterateEntries is Iterate over page entries
	IterateEntries = C.MDBM_ITERATE_ENTRIES

	// IterateNolock is Iterate without locking
	IterateNolock = C.MDBM_ITERATE_NOLOCK

	// StatNolock is Do not lock for stat operation
	StatNolock = C.MDBM_STAT_NOLOCK
)
