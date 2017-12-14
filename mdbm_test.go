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

	pathList := [...]string{pathTestDBM1, pathTestDBM2, pathTestDBM3, pathTestDBMHash, pathTestDBMDup, pathTestDBMCache, pathTestDBMV2}

	dbm := mdbm.NewMDBM()

	for _, path := range pathList {

		err := os.Remove(path)
		if err != nil {
			log.Printf("not exists the `%s` file", path)
		} else {
			log.Printf("remove the `%s` file", path)
		}

		_, err = dbm.DeleteLockFiles(path)
		if err == nil {
			log.Printf("delete lock files of %s", path)
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
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryInsertData_StoreWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM3, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLock(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryFetchData_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, val, err := dbm.Fetch(i)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.Fetch(r1.Intn(loopLimit))
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, val, err := dbm.FetchWithLock(i)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, val, err := dbm.FetchWithLockSmart(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, val, err := dbm.FetchWithLockShared(i)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "Return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, val, err := dbm.FetchWithPlock(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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

		rv, val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "return value=%d, err=%v\n", rv, err)
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
	assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock(false, time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("true", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("false", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		rv, err = dbm.StoreWithLock(int8(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(int16(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint16(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint32(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int31(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int63(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Uint32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Uint64(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float64(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "return value=%v, err=%v\n", rv, err)
	}
}
