package mdbm

import (
	"testing"

	"github.com/torden/go-strutil"
)

var assert = strutils.NewAssert()

var intzero = int(0)
var uint32zero = uint32(0)
var uint64zero = uint64(0)

func Test_mdbm_unexport_convertIterToC(t *testing.T) {

	var iter Iter

	dbm := NewMDBM()
	cIter := dbm.convertIterToC(iter)
	goIter := dbm.convertIter(&cIter)

	assert.AssertEquals(t, goIter.PageNo, uint32zero, "failured, converted value mismatch, goIter.Pageno=%+v", goIter.PageNo)
	assert.AssertEquals(t, goIter.Next, intzero, "failured, converted value mismatch, goIter.Next=%+v", goIter.Next)
}

func Test_mdbm_unexport_convertFetchInfo(t *testing.T) {

	//golang not supported testing include the cgo routine

}

func Test_mdbm_unexport_convertStatsToC(t *testing.T) {

	var stats Stats

	dbm := NewMDBM()

	cStats := dbm.convertStatsToC(stats)
	goStats := dbm.convertStats(cStats)

	assert.AssertEquals(t, goStats.Size, uint32zero, "failured, converted value mismatch, goStats.Size=%+v", goStats.Size)
	assert.AssertEquals(t, goStats.PageSize, uint32zero, "failured, converted value mismatch, goStats.PageSize=%+v", goStats.PageSize)
	assert.AssertEquals(t, goStats.PageCount, uint32zero, "failured, converted value mismatch, goStats.PageCount=%+v", goStats.PageCount)
	assert.AssertEquals(t, goStats.PagesUsed, uint32zero, "failured, converted value mismatch, goStats.PagesUsed=%+v", goStats.PagesUsed)
	assert.AssertEquals(t, goStats.BytesUsed, uint32zero, "failured, converted value mismatch, goStats.BytesUsed=%+v", goStats.BytesUsed)
	assert.AssertEquals(t, goStats.NumEntries, uint32zero, "failured, converted value mismatch, goStats.NumEntries=%+v", goStats.NumEntries)
	assert.AssertEquals(t, goStats.MinLevel, uint32zero, "failured, converted value mismatch, goStats.MinLevel=%+v", goStats.MinLevel)
	assert.AssertEquals(t, goStats.MaxLevel, uint32zero, "failured, converted value mismatch, goStats.MaxLevel=%+v", goStats.MaxLevel)
	assert.AssertEquals(t, goStats.LargePageSize, uint32zero, "failured, converted value mismatch, goStats.LargePageSize=%+v", goStats.LargePageSize)
	assert.AssertEquals(t, goStats.LargePageCount, uint32zero, "failured, converted value mismatch, goStats.LargePageCount=%+v", goStats.LargePageCount)
	assert.AssertEquals(t, goStats.LargeThreshold, uint32zero, "failured, converted value mismatch, goStats.LargeThreshold=%+v", goStats.LargeThreshold)
	assert.AssertEquals(t, goStats.LargePagesUsed, uint32zero, "failured, converted value mismatch, goStats.LargePagesUsed=%+v", goStats.LargePagesUsed)
	assert.AssertEquals(t, goStats.LargeNumFreeEntries, uint32zero, "failured, converted value mismatch, goStats.LargeNumFreeEntries=%+v", goStats.LargeNumFreeEntries)
	assert.AssertEquals(t, goStats.LargeMaxFree, uint32zero, "failured, converted value mismatch, goStats.LargeMaxFree=%+v", goStats.LargeMaxFree)
	assert.AssertEquals(t, goStats.LargeNumEntries, uint32zero, "failured, converted value mismatch, goStats.LargeNumEntries=%+v", goStats.LargeNumEntries)
	assert.AssertEquals(t, goStats.LargeBytesUsed, uint32zero, "failured, converted value mismatch, goStats.LargeBytesUsed=%+v", goStats.LargeBytesUsed)
	assert.AssertEquals(t, goStats.LargeMinSize, uint32zero, "failured, converted value mismatch, goStats.LargeMinSize=%+v", goStats.LargeMinSize)
	assert.AssertEquals(t, goStats.LargeMaxSize, uint32zero, "failured, converted value mismatch, goStats.LargeMaxSize=%+v", goStats.LargeMaxSize)
	assert.AssertEquals(t, goStats.CacheMode, uint32zero, "failured, converted value mismatch, goStats.CacheMode=%+v", goStats.CacheMode)

}

func Test_mdbm_unexport_convertDBInfo(t *testing.T) {

	var dbinfo DBInfo

	dbm := NewMDBM()

	cDBInfo := dbm.convertDBInfoToC(dbinfo)
	goDBInfo := dbm.convertDBInfo(cDBInfo)

	assert.AssertEquals(t, goDBInfo.PageSize, uint32zero, "failured, converted value mismatch, goDBInfo.PageSize=%+v", goDBInfo.PageSize)
	assert.AssertEquals(t, goDBInfo.NumPages, uint32zero, "failured, converted value mismatch, goDBInfo.NumPages=%+v", goDBInfo.NumPages)
	assert.AssertEquals(t, goDBInfo.MaxPages, uint32zero, "failured, converted value mismatch, goDBInfo.MaxPages=%+v", goDBInfo.MaxPages)
	assert.AssertEquals(t, goDBInfo.NumDirPages, uint32zero, "failured, converted value mismatch, goDBInfo.NumDirPages=%+v", goDBInfo.NumDirPages)
	assert.AssertEquals(t, goDBInfo.DirWidth, uint32zero, "failured, converted value mismatch, goDBInfo.DirWidth=%+v", goDBInfo.DirWidth)
	assert.AssertEquals(t, goDBInfo.MaxDirShift, uint32zero, "failured, converted value mismatch, goDBInfo.MaxDirShift=%+v", goDBInfo.MaxDirShift)
	assert.AssertEquals(t, goDBInfo.DirMinLevel, uint32zero, "failured, converted value mismatch, goDBInfo.DirMinLevel=%+v", goDBInfo.DirMinLevel)
	assert.AssertEquals(t, goDBInfo.DirMaxLevel, uint32zero, "failured, converted value mismatch, goDBInfo.DirMaxLevel=%+v", goDBInfo.DirMaxLevel)
	assert.AssertEquals(t, goDBInfo.DirNumNodes, uint32zero, "failured, converted value mismatch, goDBInfo.DirNumNodes=%+v", goDBInfo.DirNumNodes)
	assert.AssertEquals(t, goDBInfo.HashFunc, uint32zero, "failured, converted value mismatch, goDBInfo.HashFunc=%+v", goDBInfo.HashFunc)
	assert.AssertEquals(t, goDBInfo.HashFuncName, "", "failured, converted value mismatch, goDBInfo.HashFuncName=%+v", goDBInfo.HashFuncName)
	assert.AssertEquals(t, goDBInfo.SpillSize, uint32zero, "failured, converted value mismatch, goDBInfo.SpillSize=%+v", goDBInfo.SpillSize)
	assert.AssertEquals(t, goDBInfo.CacheMode, uint32zero, "failured, converted value mismatch, goDBInfo.CacheMode=%+v", goDBInfo.CacheMode)
}

func Test_mdbm_unexport_convertBucketStat(t *testing.T) {

	var bucketstat BucketStat

	dbm := NewMDBM()

	cBucketStat := dbm.convertBucketStatToC(bucketstat)
	goBucketStat := dbm.convertBucketStat(cBucketStat)

	assert.AssertEquals(t, goBucketStat.NumPages, uint32zero, "failured, converted value mismatch, goBucketStat.NumPages=%+v", goBucketStat.NumPages)
	assert.AssertEquals(t, goBucketStat.MinBytes, uint32zero, "failured, converted value mismatch, goBucketStat.MinBytes=%+v", goBucketStat.MinBytes)
	assert.AssertEquals(t, goBucketStat.MaxBytes, uint32zero, "failured, converted value mismatch, goBucketStat.MaxBytes=%+v", goBucketStat.MaxBytes)
	assert.AssertEquals(t, goBucketStat.MinFreeBytes, uint32zero, "failured, converted value mismatch, goBucketStat.MinFreeBytes=%+v", goBucketStat.MinFreeBytes)
	assert.AssertEquals(t, goBucketStat.MaxFreeBytes, uint32zero, "failured, converted value mismatch, goBucketStat.MaxFreeBytes=%+v", goBucketStat.MaxFreeBytes)
	assert.AssertEquals(t, goBucketStat.SumEntries, uint64zero, "failured, converted value mismatch, goBucketStat.SumEntries=%+v", goBucketStat.SumEntries)
	assert.AssertEquals(t, goBucketStat.SumBytes, uint64zero, "failured, converted value mismatch, goBucketStat.SumBytes=%+v", goBucketStat.SumBytes)
	assert.AssertEquals(t, goBucketStat.SumFreeBytes, uint64zero, "failured, converted value mismatch, goBucketStat.SumFreeBytes=%+v", goBucketStat.SumFreeBytes)
}

func Test_mdbm_unexport_convertStatInfo(t *testing.T) {

	var statinfo StatInfo

	dbm := NewMDBM()
	cStatInfo := dbm.convertStatInfoToC(statinfo)
	goStatInfo := dbm.convertStatInfo(cStatInfo)

	assert.AssertEquals(t, goStatInfo.Flags, intzero, "failured, converd value mismatch, goStatInfo.Flags=%+v", goStatInfo.Flags)
	assert.AssertEquals(t, goStatInfo.NumActiveEntries, uint64zero, "failured, converd value mismatch, goStatInfo.NumActiveEntries=%+v", goStatInfo.NumActiveEntries)
	assert.AssertEquals(t, goStatInfo.NumActiveLobEntries, uint64zero, "failured, converd value mismatch, goStatInfo.NumActiveLobEntries=%+v", goStatInfo.NumActiveLobEntries)
	assert.AssertEquals(t, goStatInfo.SumKeyBytes, uint64zero, "failured, converd value mismatch, goStatInfo.SumKeyBytes=%+v", goStatInfo.SumKeyBytes)
	assert.AssertEquals(t, goStatInfo.SumLobValBytes, uint64zero, "failured, converd value mismatch, goStatInfo.SumLobValBytes=%+v", goStatInfo.SumLobValBytes)
	assert.AssertEquals(t, goStatInfo.SumNormalValBytes, uint64zero, "failured, converd value mismatch, goStatInfo.SumNormalValBytes=%+v", goStatInfo.SumNormalValBytes)
	assert.AssertEquals(t, goStatInfo.SumOverheadBytes, uint64zero, "failured, converd value mismatch, goStatInfo.SumOverheadBytes=%+v", goStatInfo.SumOverheadBytes)
	assert.AssertEquals(t, goStatInfo.MinEntryBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MinEntryBytes=%+v", goStatInfo.MinEntryBytes)
	assert.AssertEquals(t, goStatInfo.MaxEntryBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MaxEntryBytes=%+v", goStatInfo.MaxEntryBytes)
	assert.AssertEquals(t, goStatInfo.MinKeyBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MinKeyBytes=%+v", goStatInfo.MinKeyBytes)
	assert.AssertEquals(t, goStatInfo.MaxKeyBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MaxKeyBytes=%+v", goStatInfo.MaxKeyBytes)
	assert.AssertEquals(t, goStatInfo.MinValBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MinValBytes=%+v", goStatInfo.MinValBytes)
	assert.AssertEquals(t, goStatInfo.MaxValBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MaxValBytes=%+v", goStatInfo.MaxValBytes)
	assert.AssertEquals(t, goStatInfo.MinLobBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MinLobBytes=%+v", goStatInfo.MinLobBytes)
	assert.AssertEquals(t, goStatInfo.MaxLobBytes, uint32zero, "failured, converd value mismatch, goStatInfo.MaxLobBytes=%+v", goStatInfo.MaxLobBytes)
	assert.AssertEquals(t, goStatInfo.MaxPageUsedSpace, uint32zero, "failured, converd value mismatch, goStatInfo.MaxPageUsedSpace=%+v", goStatInfo.MaxPageUsedSpace)
	assert.AssertEquals(t, goStatInfo.MaxDataPages, uint32zero, "failured, converd value mismatch, goStatInfo.MaxDataPages=%+v", goStatInfo.MaxDataPages)
	assert.AssertEquals(t, goStatInfo.NumFreePages, uint32zero, "failured, converd value mismatch, goStatInfo.NumFreePages=%+v", goStatInfo.NumFreePages)
	assert.AssertEquals(t, goStatInfo.NumActivePages, uint32zero, "failured, converd value mismatch, goStatInfo.NumActivePages=%+v", goStatInfo.NumActivePages)
	assert.AssertEquals(t, goStatInfo.NumNormalPages, uint32zero, "failured, converd value mismatch, goStatInfo.NumNormalPages=%+v", goStatInfo.NumNormalPages)
	assert.AssertEquals(t, goStatInfo.NumOversizedPages, uint32zero, "failured, converd value mismatch, goStatInfo.NumOversizedPages=%+v", goStatInfo.NumOversizedPages)
	assert.AssertEquals(t, goStatInfo.NumLobPages, uint32zero, "failured, converd value mismatch, goStatInfo.NumLobPages=%+v", goStatInfo.NumLobPages)
	assert.AssertEquals(t, goStatInfo.MaxPageEntries, uint32zero, "failured, converd value mismatch, goStatInfo.MaxPageEntries=%+v", goStatInfo.MaxPageEntries)
	assert.AssertEquals(t, goStatInfo.MinPageEntries, uint32zero, "failured, converd value mismatch, goStatInfo.MinPageEntries=%+v", goStatInfo.MinPageEntries)

	for key, item := range goStatInfo.Buckets {

		assert.AssertEquals(t, item.NumPages, uint32zero, "failured, convertd value mismatch, goStatInfo[%d].NumPages=%+v", key, item.NumPages)
		assert.AssertEquals(t, item.MinBytes, uint32zero, "failured, convertd value mismatch, goStatInfo[%d].MinBytes=%+v", key, item.MinBytes)
		assert.AssertEquals(t, item.MaxBytes, uint32zero, "failured, convertd value mismatch, goStatInfo[%d].MaxBytes=%+v", key, item.MaxBytes)
		assert.AssertEquals(t, item.MinFreeBytes, uint32zero, "failured, convertd value mismatch, goStatInfo[%d].MinFreeBytes=%+v", key, item.MinFreeBytes)
		assert.AssertEquals(t, item.MaxFreeBytes, uint32zero, "failured, convertd value mismatch, goStatInfo[%d].MaxFreeBytes=%+v", key, item.MaxFreeBytes)
		assert.AssertEquals(t, item.SumEntries, uint64zero, "failured, convertd value mismatch, goStatInfo[%d].SumEntries=%+v", key, item.SumEntries)
		assert.AssertEquals(t, item.SumBytes, uint64zero, "failured, convertd value mismatch, goStatInfo[%d].SumBytes=%+v", key, item.SumBytes)
		assert.AssertEquals(t, item.SumFreeBytes, uint64zero, "failured, convertd value mismatch, goStatInfo[%d].SumFreeBytes=%+v", key, item.SumFreeBytes)

	}
}

func Test_mdbm_unexport_convertWindowStat(t *testing.T) {

	var wstats WindowStats
	dbm := NewMDBM()

	cWstats := dbm.convertWindowStatToC(wstats)
	goWStats := dbm.convertWindowStat(cWstats)

	assert.AssertEquals(t, goWStats.WnumReused, uint64zero, "failured, converted value mismatch, goWStats.WnumReused=%+v", goWStats.WnumReused)
	assert.AssertEquals(t, goWStats.WnumRemapped, uint64zero, "failured, converted value mismatch, goWStats.WnumRemapped=%+v", goWStats.WnumRemapped)
	assert.AssertEquals(t, goWStats.WwindowSize, uint32zero, "failured, converted value mismatch, goWStats.WwindowSize=%+v", goWStats.WwindowSize)
	assert.AssertEquals(t, goWStats.WmaxWindowUsed, uint32zero, "failured, converted value mismatch, goWStats.WmaxWindowUsed=%+v", goWStats.WmaxWindowUsed)
}
