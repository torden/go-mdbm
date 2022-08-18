package mdbm_test

const (
	loopLimit                = 65534
	pathTestDBM1             = "/tmp/test1.mdbm"
	pathTestDBM2             = "/tmp/test2.mdbm"
	pathTestDBM3             = "/tmp/test3.mdbm"
	pathTestDBMLarge         = "/tmp/test_large.mdbm"
	pathTestDBMHash          = "/tmp/test_hash.mdbm"
	pathTestDBMDup           = "/tmp/test_dup.mdbm"
	pathTestDBMCache         = "/tmp/test_cache.mdbm"
	pathTestDBMCacheNoneData = "/tmp/test_cache_nonedata.mdbm"
	pathTestDBMV2            = "/tmp/test_v2.mdbm"
	pathTestDBMLock1         = "/tmp/test_lock1.mdbm"
	pathTestDBMDelete        = "/tmp/test_delete.mdbm"
	pathTestDBMLock2         = "/tmp/test_lock2.mdbm"
	pathTestDBMLock3         = "/tmp/test_lock3.mdbm"
	pathTestDBMAnyDataType1  = "/tmp/test_anydatatype1.mdbm"
	pathTestDBMAnyDataType2  = "/tmp/test_anydatatype2.mdbm"
	pathTestDBMStr           = "/tmp/test_str.mdbm"
	pathTestDBMR             = "/tmp/test_r.mdbm"
	pathTestDBMStrAnyLock    = "/tmp/test_str_anylock.mdbm"
	// pathTestDBMReplace1                   = "/tmp/test_replace1.mdbm"
	// pathTestDBMReplace2                   = "/tmp/test_replace2.mdbm"
	pathTestDBMReplace3                   = "/tmp/test_replace3.mdbm"
	pathTestDBMRemove                     = "/tmp/test_remove.mdbm"
	pathTestDBMEmpty                      = "/tmp/test_empty.mdbm"
	pathTestDBMFcopy                      = "/tmp/test_fcopy.mdbm"
	pathTestDBMBenchmark1                 = "/tmp/test_benchmark1.mdbm"
	pathTestDBMBenchmark2                 = "/tmp/test_benchmark2.mdbm"
	pathTestDBMBenchmark3                 = "/tmp/test_benchmark3.mdbm"
	pathTestBoltDBBenchmark1              = "/tmp/test_benchmark1.boltdb"
	pathTestBoltDBBenchmark2              = "/tmp/test_benchmark2.boltdb"
	pathTestBoltDBBenchmark3              = "/tmp/test_benchmark3.boltdb"
	pathTestDBMBenchmarkForRandomFetch    = "/tmp/test_benchmark_fetch.mdbm"
	pathTestBoltDBBenchmarkForRandomFetch = "/tmp/test_benchmark_fetch.boltdb"
)

var gPathList = [...]string{
	pathTestDBM1,
	pathTestDBM2,
	pathTestDBM3,
	pathTestDBMLarge,
	pathTestDBMHash,
	pathTestDBMDup,
	pathTestDBMCache,
	pathTestDBMV2,
	pathTestDBMLock1,
	pathTestDBMDelete,
	pathTestDBMLock2,
	pathTestDBMAnyDataType1,
	pathTestDBMAnyDataType2,
	pathTestDBMStr,
	pathTestDBMR,
	pathTestDBMStrAnyLock,
	pathTestDBMCacheNoneData,
	pathTestDBMRemove,
	pathTestDBMEmpty,
	pathTestDBMFcopy,
	pathTestDBMBenchmark1,
	pathTestDBMBenchmark2,
	pathTestDBMBenchmark3,
	pathTestBoltDBBenchmark1,
	pathTestBoltDBBenchmark2,
	pathTestBoltDBBenchmark3,
	pathTestDBMBenchmarkForRandomFetch,
	pathTestBoltDBBenchmarkForRandomFetch,
}
