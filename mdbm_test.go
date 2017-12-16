package mdbm_test

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/torden/go-mdbm"
	"github.com/torden/go-strutil"
)

var assert = strutils.NewAssert()

func TestMain(t *testing.T) {

	pathList := [...]string{pathTestDBM1, pathTestDBM2, pathTestDBM3, pathTestDBMLarge, pathTestDBMHash, pathTestDBMDup, pathTestDBMCache, pathTestDBMV2, pathTestDBMLock1, pathTestDBMDelete, pathTestDBMLock2}

	dbm := mdbm.NewMDBM()

	for _, path := range pathList {

		_, err := dbm.DeleteLockFiles(path)
		if err == nil {
			log.Printf("delete lock files of %s", path)
		}

		err = os.Remove(path)
		if err != nil {
			log.Printf("not exists the `%s` file", path)
		} else {
			log.Printf("remove the `%s` file", path)
		}

	}
}

func Test_mdbm_EasyOpen_EasyClose(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_Open_Close(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM2, mdbm.Create|mdbm.Rdrw, 0644, 0, 0)
	defer dbm.Close()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_OrdinaryInsertData_Store1(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.Open(pathTestDBM1, mdbm.Create|mdbm.Rdrw, 0644, 0, 0)
	defer dbm.Close()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, execute to mdbm.Sync()\nrv=%d, err=%v", rv, err)
}

func Test_mdbm_OrdinaryInsertData_Store2(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM2, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM2, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryInsertData_StoreWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM3, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLock(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryFetchData_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.Fetch(i)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_RandomFetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r1.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLock(i)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLockSmart(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLockShared(i)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithPlock(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%d, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_LockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.LockShared()
	assert.AssertEquals(t, 1, rv, "failured, Locks the database for shared access by readers, excluding access to writers., path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, Locks the database for shared access by readers, excluding access to writers., path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_TryLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.TryLockShared()
	assert.AssertEquals(t, 1, rv, "failured, locks the database for shared access by readers, excluding access to writers, path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, locks the database for shared access by readers, excluding access to writers, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_GetLockMode(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.GetLockMode()
	assert.AssertEquals(t, 1, rv, "failured, gets the mdbm's lock mode, path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, gets the mdbm's lock mode, path=%s, err=%v", pathTestDBM1, err)

}

func Test_mdbm_MutipleDataType_Store(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMLarge, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects, 0644, 0, 0)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	rv, err = dbm.StoreWithLock(true, time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock(false, time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("true", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("false", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock(byte(77), time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		rv, err = dbm.StoreWithLock(int8(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(int16(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint16(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint32(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int31(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int63(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Uint32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float64(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(int64(r.Int()), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint64(r.Int()), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)
	}
}

func Test_mdbm_DupHandle_AfterClose(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMLarge, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects, 0644, 0, 0)
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	_, err = dbm.DupHandle()
	assert.AssertNil(t, err, "failured, a pointer of the Duplicate an existing database handle, path=%s, err=%v", pathTestDBM1, err)

	dbm.EasyClose()

	_, err = dbm.DupHandle()
	assert.AssertNotNil(t, err, "failured, return of closed db hanlder, err=%v", err)

}

func Test_mdbm_LogMinLevel_WrongOption(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogMinLevel(mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "oops!. mdbm.LogMinLevel can't check a argument=%d", int(mdbm.LargeObjects))
}

func Test_mdbm_LogPlugin(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogPlugin(mdbm.LogToFile)
	assert.AssertNil(t, err, "failured, can't set logging to file, err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToSkip)
	assert.AssertNil(t, err, "failured, can't set logging to /dev/null , err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToStdErr)
	assert.AssertNil(t, err, "failured, can't set loging to /dev/stdout, err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToSysLog)
	assert.AssertNil(t, err, "failured, can't set loging to syslog, err=%v", err)

	err = dbm.LogPlugin(mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "oops!. mdbm.LogPlugin can't check a argument=%d", int(mdbm.LargeObjects))

}

func Test_mdbm_LogToAutoFile(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogPlugin(mdbm.LogToFile)
	assert.AssertNil(t, err, "failured, can't set logging to file, err=%v", err)

	rv, err := dbm.LogToAutoFile()
	assert.AssertNil(t, err, "failured, can't set logging to file, rv=%d, err=%v", rv, err)
}

func Test_mdbm_TryLock_UnLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.TryLock()
	assert.AssertNil(t, err, "failured, can't try-locking, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	rv, err = dbm.Unlock()
	assert.AssertNil(t, err, "failured, can't un-locking, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)
}

func Test_mdbm_FetchInfo(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	uint32zero := uint32(0)
	intzero := int(0)

	for i := 0; i <= 1; i++ {

		var retval string

		rv, copiedval, info, goiter, err := dbm.FetchInfo(i, &retval, &iter)
		assert.AssertNil(t, err, "failured, can't get fetch infor, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)
		assert.AssertEquals(t, string(i), copiedval, "failured, Return Value mismatch.\nExpected: %v\nActual: %v", i, copiedval)

		assert.AssertEquals(t, info.Flags, uint32zero, "failured, Return Value mismatch.\nExpected: %v\nActual: %v", info.Flags, uint32zero)
		assert.AssertEquals(t, info.CacheNumAccesses, uint32zero, "failured, Return Value mismatch.\nExpected: %v\nActual: %v", info.CacheNumAccesses, uint32zero)
		assert.AssertEquals(t, info.CacheAccessTime, uint32zero, "failured, Return Value mismatch.\nExpected: %v\nActual: %v", info.CacheAccessTime, uint32zero)

		assert.AssertNotEquals(t, goiter.PageNo, uint32zero, "failured, Return Value mismatch.\nExpected: %d\nActual: %d", goiter.PageNo, uint32zero)
		assert.AssertNotEquals(t, goiter.Next, intzero, "failured, Return Value mismatch.\nExpected: %d\nActual: %d", goiter.Next, intzero)
	}
}

func Test_mdbm_DeleteWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMDelete, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)
	}

	dbm.Sync()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.DeleteWithLock(i)
		assert.AssertNil(t, err, "failured, can't delete record, return value=%v, err=%v\n", rv, err)
	}

	dbm.Sync()
	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(i)
		assert.AssertNotNil(t, err, "failured, can't delete record, value=%v, err=%v\n", val, err)
	}
}

func Test_mdbm_EasyGetNumOfRows(t *testing.T) {
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	cnt, err := dbm.EasyGetNumOfRows()

	assert.AssertNil(t, err, "failured, can't obtain the count of number of rows, err=%v\n", err)
	assert.AssertEquals(t, cnt, uint64(772351), "failured, Return Value mismatch.\nExpected: %v\nActual: %v", 772351, cnt)

}

func Test_mdbm_EasyGetKeyList(t *testing.T) {
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	keys, err := dbm.EasyGetKeyList()

	assert.AssertNil(t, err, "failured, can't obtain the list of key, err=%v\n", err)

	assert.AssertEquals(t, len(keys), uint64(772351), "failured, Return Value mismatch.\nExpected: %v\nActual: %v", 772351, len(keys))
}

func Test_mdbm_Truncate(t *testing.T) {

	var rv int
	var val string
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMDelete, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, Return Value mismatch. value=%v, err=%v\n", rv, err)
	}

	dbm.Sync()

	err = dbm.Truncate()
	assert.AssertNil(t, err, "failured, can't truncate mdbm, err=%v\n", err)

	dbm.Sync()
	for i := 0; i <= loopLimit; i++ {

		val, err = dbm.Fetch(i)
		assert.AssertNotNil(t, err, "failured, can't delete record, value=%v, err=%v\n", val, err)
	}
}
