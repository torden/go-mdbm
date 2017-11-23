package mdbm_test

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/torden/go-mdbm"
)

var pathList = [...]string{pathTestDBM1, pathTestDBM2, pathTestDBM3, pathTestDBMHash, pathTestDBMDup, pathTestDBMCache, pathTestDBMV2}

func init() {

	dbm := mdbm.NewMDBM()

	for _, path := range pathList {

		if _, err := os.Stat(path); err != nil {

			if os.IsExist(err) {

				err = os.Remove(path)
				if err != nil {
					log.Printf("failed remove the `%s` file", path)
				}
			}
		}

		_, err := dbm.DeleteLockFiles(path)
		if err == nil {
			log.Printf("delete lock files of %s", path)
		}
	}
	// Output:
}

func Example_mdbm_EasyOpen_EasyClose() {

	dbm := mdbm.NewMDBM()

	for _, path := range pathList {

		err := dbm.EasyOpen(path, 0666)
		if err != nil {
			log.Fatalf("failed mdbm.EasyOpen(%s), err=%v", path, err)
		}
		fmt.Println(err)

		rv, err := dbm.EnableStatOperations(mdbm.StatsTimed)
		if err != nil {
			log.Fatalf("failed dbm.EnableStatOperations(mdbm.StatsTimed), rv=%d, err=%v", rv, err)
		}

		fmt.Println(err)
		rv, err = dbm.SetStatTimeFunc(mdbm.ClockTsc)
		if err != nil {
			log.Fatalf("failed dbm.SetStatTimeFunc(mdbm.ClockTsc), rv=%d, err=%v", rv, err)
		}

		fmt.Println(err)

		dbm.EasyClose()
	}

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
}

func Example_mdbm_Open_Close() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM2, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	fmt.Println(err)

	_, err = dbm.EnableStatOperations(mdbm.StatsTimed)
	fmt.Println(err)
	_, err = dbm.SetStatTimeFunc(mdbm.ClockTsc)
	fmt.Println(err)

	dbm.Close()
	// Output:
	// <nil>
	// <nil>
	// <nil>
}

func Example_mdbm_DupHandle() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM3, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	dbm2, err2 := dbm.DupHandle()
	dbm2.Close()
	fmt.Println(err)
	fmt.Println(err2)
	// Output: <nil>
	// <nil>
}

func Example_mdbm_GetErrNo() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM1, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	fmt.Println(err)

	rv, err := dbm.GetErrNo()
	fmt.Println(rv, err)
	dbm.Close()
	// Output:  <nil>
	// 0 <nil>
}

func Example_mdbm_LogMinLevel() {

	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.LogMinLevel(mdbm.LogEmerg)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogAlert)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogCrit)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogErr)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogWarning)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogNotice)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogInfo)
	fmt.Println(err)
	err = dbm.LogMinLevel(mdbm.LogDebug)
	fmt.Println(err)

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// <nil>
}

func Example_mdbm_Sync() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println(err)
	// any
	rv, err := dbm.Sync()
	dbm.EasyClose()

	fmt.Println(rv, err)
	// Output: <nil>
	// 0 <nil>
}

func Example_mdbm_Fsync() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println(err)
	// any
	rv, err := dbm.Fsync()
	dbm.EasyClose()

	fmt.Println(rv, err)
	// Output: <nil>
	// 0 <nil>
}

func Example_mdbm_CloseFD() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println(err)
	// any
	err = dbm.CloseFD()

	fmt.Println(err)
	// Output: <nil>
	// <nil>
}

func Example_mdbm_Lock_Unlock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println("EasyOpen : ", err)

	err = dbm.Lock()
	fmt.Println("Lock : ", err)

	rv, err := dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// Lock :  <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_mdbm_IsLocked() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println("EasyOpen : ", err)

	err = dbm.TryLock()
	fmt.Println("TryLock : ", err)

	rv, err = dbm.IsLocked()
	fmt.Println("IsLocked : rv =", rv, ", err =", err)

	rv, err = dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// TryLock :  <nil>
	// IsLocked : rv = 1 , err = <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_mdbm_LockShared() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println("EasyOpen : ", err)

	rv, err = dbm.LockShared()
	fmt.Println("LockShared: rv =", rv, ", err =", err)

	rv, err = dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// LockShared: rv = 1 , err = <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_mdbm_TryLockShared() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println("EasyOpen : ", err)

	rv, err = dbm.TryLockShared()
	fmt.Println("TryLockShared: rv =", rv, ", err =", err)

	rv, err = dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// TryLockShared: rv = 1 , err = <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_mdbm_LockReset() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	fmt.Println("EasyOpen : ", err)

	err = dbm.Lock()
	fmt.Println("Lock : ", err)

	rv, err = dbm.LockReset()
	fmt.Println("LockReset : rv =", rv, ", err =", err)

	rv, err = dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// Lock :  <nil>
	// LockReset : rv = 0 , err = <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  operation not permitted
}

func Example_mdbm_ReplaceDB() {

	dbm := mdbm.NewMDBM()

	//create a dummy
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	dbm.EasyClose()

	err = dbm.ReplaceFile(pathTestDBM1, pathTestDBM3)
	fmt.Println("ReplaceFile : ", err)
	dbm.EasyClose()

	// Output:
	// ReplaceFile :  <nil>
}

func Example_mdbm_GetHash() {

	dbm := mdbm.NewMDBM()

	//create a dummy
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}

	rv, err := dbm.GetHash()
	fmt.Println("GetHash : rv =", rv, ", err = ", err)
	dbm.EasyClose()

	// mdbm.DefaultHash is 5

	// Output:
	// GetHash : rv = 5 , err =  <nil>
}

func Example_mdbm_SetHash() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMHash, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.SetHash(mdbm.HashOZ)
	if err != nil {
		log.Fatalf("failed mdbm.SetHash(mdbm.HashOZ), err=%v", err)
	}

	rv, err = dbm.GetHash()
	if rv == mdbm.HashOZ {
		fmt.Println("Setted : mdbm.HashOZ")
	} else {
		fmt.Println("Not Setted : rv = ", rv, ", err = ", err)
		fmt.Printf("Want : mdbm.HashOZ(%d)\n", mdbm.HashOZ)
	}

	// Output:
	// Setted : mdbm.HashOZ
}

func Example_mdbm_GetAlignment() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetAlignment()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv)
	// Output:
	// 0
}

func Example_mdbm_SetAlignment() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.SetAlignment(0)
	if err != nil {
		log.Println(err)
	}

	if rv == mdbm.Align8Bits {
		fmt.Println("Default :: 8-bit alignment.")
	}

	// Output:
	// Default :: 8-bit alignment.
}

func Example_mdbm_GetLimitSize() {

	var rv uint64
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetLimitSize()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv)
	// Output:
	// 0
}

func Example_mdbm_LimitDirSize() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.LimitDirSize(2)
	fmt.Println(err)
	// Output:
	// <nil>
}

func Example_mdbm_GetVersion() {

	var rv uint32
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetVersion()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 3 <nil>
}

func Example_mdbm_GetSize() {

	var rv uint64
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetSize()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 8192 <nil>
}

func Example_mdbm_GetPageSize() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetPageSize()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 4096 <nil>
}

func Example_mdbm_GetMagicNumber() {

	var rv uint32
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetMagicNumber()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 16922980 <nil>
}

func Example_mdbm_SetWindowSize() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.SetWindowSize(9216)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(err)
	// Output:
	// <nil>
}

func Example_mdbm_IsOwned() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.IsOwned()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 0 <nil>
}

func Example_mdbm_GetLockMode() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetLockMode()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 0 <nil>

	switch rv {
	case 0:
		fmt.Println("Exclusive locking")
	case mdbm.OpenNolock:
		fmt.Println("MDBM_OPEN_NOLOCK       - No locking")
	case mdbm.PartitionedLocks:
		fmt.Println("MDBM_PARTITIONED_LOCKS - Partitioned locking")
	case mdbm.RwLocks:
		fmt.Println("MDBM_RW_LOCKS          - Shared (read-write) locking")
	}

	//Exclusize Locking
}

func Example_mdbm_CompressTree() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.CompressTree()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(err)
	// Output:
	// <nil>
}

func Example_mdbm_Truncate() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.CompressTree()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(err)
	// Output:
	// <nil>
}

func Example_mdbm_Purge() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.Purge()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(err)
	// Output:
	// <nil>
}

func Example_mdbm_Check() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, _, err = dbm.Check(1, false)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 0 <nil>
}

func Example_mdbm_CheckAllPage() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.CheckAllPage()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 0 <nil>
}

func Example_mdbm_Protect() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.Protect(mdbm.ProtAccess)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)

	rv, err = dbm.Protect(mdbm.ProtNone)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
	// 0 <nil>
}

func Example_mdbm_DumpAllPage() {

	var rv string
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.DumpAllPage()
	if err != nil {
		log.Println(err)
	}

	if len(rv) > 0 {
		fmt.Println("OK")
	}

	// Output:
	// OK
}

func Example_mdbm_StoreWithLock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLock(i, i, mdbm.Replace)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

func Example_mdbm_Store() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.Store(i, i, mdbm.Replace)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

func Example_mdbm_StoreRWithLock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	for i := 0; i <= loopLimit; i++ {
		rv, _, err := dbm.StoreRWithLock(i, i, mdbm.Replace, &iter)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

func Example_mdbm_StoreR() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	for i := 0; i <= loopLimit; i++ {
		rv, _, err := dbm.StoreR(i, i, mdbm.Replace, &iter)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

// BUG: tail \00
func Example_mdbm_StoreStrWitchLock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM2, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreStrWitchLock(i, i, mdbm.Replace)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

// BUG: tail \00
func Example_mdbm_StoreStr() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM2, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreStr(i, i, mdbm.Replace)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

func Example_mdbm_Fetch() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, retval, err := dbm.Fetch(i)
		if err != nil {
			log.Fatalf("rv=%d, retval=%s, err=%v\n", rv, retval, err)
		} else {

			if retval != strconv.Itoa(i) {
				log.Fatalf("wrong return value=%s, want=%d\n", retval, i)
			}
		}
	}

	// Output:
}

func Example_mdbm_FetchR() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	for i := 0; i <= loopLimit; i++ {
		rv, retval, _, err := dbm.FetchR(i, &iter)
		if err != nil {
			log.Fatalf("rv=%d, retval=%s, err=%v\n", rv, retval, err)
		} else {

			if retval != strconv.Itoa(i) {
				log.Fatalf("wrong return value=%s, want=%d\n", retval, i)
			}
		}
	}

	// Output:
}

func Example_mdbm_FetchBuf() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {

		var retval string
		rv, copiedval, err := dbm.FetchBuf(i, &retval)

		if err != nil {
			log.Fatalf("rv=%d, retval=%s, copiedval=%s, err=%v\n", rv, retval, copiedval, err)
		} else {

			if retval != strconv.Itoa(i) {
				log.Fatalf("wrong return value=%s, want=%d\n", retval, i)
			}

			if retval != copiedval {
				log.Fatalf("not matched return value=%s, copied return valued=%s\n", retval, copiedval)
			}
		}
	}

	// Output:

}

func Example_mdbm_StoreDup() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMDup, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i < 100; i++ {
		rv, err := dbm.Store(i, i+100, mdbm.InsertDup)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	for i := 0; i < 100; i++ {
		rv, err := dbm.Store(i, i+200, mdbm.InsertDup)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	// Output:
}

/*
func Example_mdbm_FetchDupR() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMDup, 0644)
	if err != nil {
	log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i < 100; i++ {
		rv, retval, iter, err := dbm.FetchDupR(i)
		fmt.Println(rv, retval, err)
		spew.Dump(iter)
	}

	// Output:
}
*/

func Example_mdbm_FetchStr() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i < 100; i++ {
		rv, err := dbm.StoreStrWitchLock(i, i, mdbm.Replace)
		if err != nil {
			fmt.Printf("return value=%d, err=%v\n", rv, err)
		}
	}

	for i := 0; i < 100; i++ {
		retval, err := dbm.FetchStr(i)
		if err != nil {
			log.Fatalf("key=%d, retval=%s, err=%v\n", i, retval, err)
		} else {

			if retval != strconv.Itoa(i) {
				log.Fatalf("wrong return value=%s, want=%d\n", retval, i)
			}
		}
	}

	// Output:
}

/*
func Example_mdbm_FetchInfo() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
	log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i < 100; i++ {

		var retval string

		rv, copiedval, info, _, err := dbm.FetchInfo(i, &retval)
		if err != nil {
			log.Fatalf("rv=%d, retval=%s, copiedval=%s, err=%v\n", rv, retval, copiedval, err)
		} else {

			if retval != strconv.Itoa(i) {
				log.Fatalf("wrong return value=%s, want=%d\n", retval, i)
			}

			if retval != copiedval {
				log.Fatalf("not matched return value=%s, copied return valued=%s\n", retval, copiedval)
			}
		}

		pp.Println(info)
	}

	// Output:
}
*/

func Example_mdbm_Delete() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.Delete(4)
	if err != nil {
		log.Fatalf("rv=%d, err=%v\n", rv, err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
}

// func Example_mdbm_DeleteR()

func Example_mdbm_DeleteStr() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.Delete(44)
	if err != nil {
		log.Fatalf("rv=%d, err=%v\n", rv, err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
}

func Example_mdbm_First_Next_Iteration() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	key, val, err := dbm.First()
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	}

	var i int

	for {

		key, val, err := dbm.Next()
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}
		i++
	}

	fmt.Println("number of rows :", i)

	// Output:
	// number of rows : 65633
}

func Example_mdbm_FirstR_NextR_Iteration() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	var iter mdbm.Iter

	key, val, _, err := dbm.FirstR(&iter)
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	}

	var i int

	for {

		key, val, _, err := dbm.NextR(&iter)
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}
		i++
	}

	fmt.Println("number of rows :", i)

	// Output:
	// number of rows : 65633
}

func Example_mdbm_FirstKey() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	key, err := dbm.FirstKey()
	if err != nil {
		log.Fatalf("key=%s, err=%v", key, err)
	}

	// Output:
}

func Example_mdbm_NextKey() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	fkey, err := dbm.FirstKey()
	if err != nil {
		log.Fatalf("fkey=%s, err=%v", fkey, err)
	}

	nkey, err := dbm.NextKey()
	if err != nil {
		log.Fatalf("nkey=%s, err=%v", nkey, err)
	}

	// Output:
}

func Example_mdbm_FirstKeyR() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	var iter mdbm.Iter

	key, _, err := dbm.FirstKeyR(&iter)
	if err != nil {
		log.Fatalf("key=%s, err=%v", key, err)
	}

	//fmt.Printf("key=%s, iter.PageNo=%d, iter.Next=%d", key, iter.PageNo, iter.Next)
	// Output:
}

func Example_mdbm_NextKeyR() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	var iter mdbm.Iter

	fkey, _, err := dbm.FirstKeyR(&iter)
	if err != nil {
		log.Fatalf("fkey=%s, err=%v", fkey, err)
	}

	nkey, _, err := dbm.NextKeyR(&iter)
	if err != nil {
		log.Fatalf("nkey=%s, err=%v", nkey, err)
	}

	//fmt.Printf("nkey=%s, iter.PageNo=%d, iter.Next=%d", nkey, iter.PageNo, iter.Next)
	// Output:
}

func Example_mdbm_GetCacheMode() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetCacheMode()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 0 <nil>
}

func Example_mdbm_SetCacheMode() {

	var rv int
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMCache, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.SetCacheMode(mdbm.CacheModeMax)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)

	rv, err = dbm.GetCacheMode()
	if err != nil {
		log.Println(err)
	}

	if rv != mdbm.CacheModeMax {
		log.Fatalf("rv=%d, want(mdbm.CacheModeMax)=%d", rv, mdbm.CacheModeMax)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
	// 3 <nil>
}

func Example_mdbm_CountRecords() {

	var rv uint64
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.CountRecords()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv, err)
	// Output:
	// 65634 <nil>
}

func Example_mdbm_CountPages() {

	var rv uint32
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.CountPages()
	if err != nil {
		log.Println(err)
	}

	if rv < 1 {
		fmt.Printf("fail, rv=%d", rv)
	}
	// Output:
}

func Example_mdbm_GetPage() {

	var rv uint32
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err = dbm.GetPage("65500")
	if err != nil {
		log.Println(err)
	}

	if rv < 1 {
		fmt.Printf("fail, rv=%d", rv)
	}

	// Output:
}

func Example_mdbm_PreLoad() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	//for Performance
	_, err = dbm.PreLoad()
	if err != nil {
		log.Fatalf("failed mdbm.PreLoad(), err=%v", err)
	}

	key, val, err := dbm.First()
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	}

	var i int

	for {

		key, val, err := dbm.Next()
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}
		i++
	}

	fmt.Println("number of rows :", i)

	// Output:
	// number of rows : 65633
}

func Example_mdbm_LockDump() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	out, err := dbm.LockDump()
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	if len(out) > 0 {
		fmt.Println("OK")
	}

	// Output:
	// OK
}

func Example_mdbm_LockPages() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.LockPages()
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
}

func Example_mdbm_UnLockPages() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.LockPages()
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	fmt.Println(rv, err)

	//something..

	rv, err = dbm.UnLockPages()
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
	// 0 <nil>
}

func Example_mdbm_ChkPage() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, out, err := dbm.ChkPage(1)
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	fmt.Println(rv, out, err)

	// Output:
	// 0  <nil>
}

func Example_mdbm_ChkError() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	err = dbm.ChkError(1, 1, 1)
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	fmt.Println(err)

	// Output:
	// <nil>
}

func Example_mdbm_DumpPage() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	out, err := dbm.DumpPage(1)
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	if len(out) < 1 {
		log.Fatalf("empty output")
	}

	fmt.Println(err)

	// Output:
	// <nil>
}

func Example_mdbm_EnableStatOperations_ResetStatOperations() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.EnableStatOperations(mdbm.StatsBasic)
	fmt.Println(rv, err)

	err = dbm.ResetStatOperations()
	fmt.Println(err)

	rv, err = dbm.EnableStatOperations(mdbm.StatsTimed)
	fmt.Println(rv, err)

	err = dbm.ResetStatOperations()
	fmt.Println(err)

	// Output:
	// 0 <nil>
	// <nil>
	// 0 <nil>
	// <nil>
}

func Example_mdbm_GetStatCounter() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, cnt, err := dbm.GetStatCounter(mdbm.StatTypeFetch)
	fmt.Println(rv, cnt, err)

	rv, cnt, err = dbm.GetStatCounter(mdbm.StatTypeStore)
	fmt.Println(rv, cnt, err)

	rv, cnt, err = dbm.GetStatCounter(mdbm.StatTypeDelete)
	fmt.Println(rv, cnt, err)

	rv, cnt, err = dbm.GetStatCounter(mdbm.StatTypeMax)
	fmt.Println(rv, cnt, err)

	// Output:
	// 0 0 <nil>
	// 0 0 <nil>
	// 0 0 <nil>
	// 0 0 <nil>
}

func Example_mdbm_GetStatName() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	statlist := [...]int{mdbm.StatTagFetch,
		mdbm.StatTagStore,
		mdbm.StatTagDelete,
		mdbm.StatTagLock,
		mdbm.StatTagFetchUncached,
		mdbm.StatTagGetpage,
		mdbm.StatTagGetpageUncached,
		mdbm.StatTagCacheEvict,
		mdbm.StatTagCacheStore,
		mdbm.StatTagPageStore,
		mdbm.StatTagPageDelete,
		mdbm.StatTagSync,
		mdbm.StatTagFetchNotFound,
		mdbm.StatTagFetchError,
		mdbm.StatTagStoreError,
		mdbm.StatTagDeleteFailed,
		mdbm.StatTagFetchLatency,
		mdbm.StatTagStoreLatency,
		mdbm.StatTagDeleteLatency,
		mdbm.StatTagFetchTime,
		mdbm.StatTagStoreTime,
		mdbm.StatTagDeleteTime,
		mdbm.StatTagFetchUncachedLatency,
		mdbm.StatTagGetpageLatency,
		mdbm.StatTagGetpageUncachedLatency,
		mdbm.StatTagCacheEvictLatency,
		mdbm.StatTagCacheStoreLatency,
		mdbm.StatTagPageStoreValue,
		mdbm.StatTagPageDeleteValue,
		mdbm.StatTagSyncLatency,
	}

	for _, tag := range statlist {

		rv, err := dbm.GetStatName(tag)
		fmt.Println(rv, err)
	}

	// Output:
	// Fetch <nil>
	// Store <nil>
	// Delete <nil>
	// Lock <nil>
	// FetchUncached <nil>
	// GetPage <nil>
	// GetPageUncached <nil>
	// CacheEvict <nil>
	// CacheStore <nil>
	// PageStore <nil>
	// PageDelete <nil>
	// MdbmSync <nil>
	// FetchNotFound <nil>
	// FetchError <nil>
	// StoreError <nil>
	// DeleteFailed <nil>
	// FetchLatency <nil>
	// StoreLatency <nil>
	// DeleteLatency <nil>
	// FetchTime <nil>
	// StoreTime <nil>
	// DeleteTime <nil>
	// FetchUncachedLatency <nil>
	// GetPageLatency <nil>
	// GetPageUncachedLatency <nil>
	// CacheEvictLatency <nil>
	// CacheStoreLatency <nil>
	// PageStoreVal <nil>
	// PageDeleteVal <nil>
	// MdbmSyncLatency <nil>
}

func Example_mdbm_GetStatTime() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	statlist := [...]int{
		mdbm.StatTypeFetch,
		mdbm.StatTypeStore,
		mdbm.StatTypeDelete,
		mdbm.StatTypeMax,
	}

	for _, tag := range statlist {

		rv, vv, err := dbm.GetStatTime(tag)
		fmt.Println(rv, vv, err)
	}

	// Output:
	// 0 0 <nil>
	// 0 0 <nil>
	// 0 0 <nil>
	// 0 0 <nil>
}

func Example_mdbm_SetStatTimeFunc() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.SetStatTimeFunc(mdbm.ClockStandard)
	if err != nil {
		log.Fatalf("rv=%d, err=%v", rv, err)
	}

	fmt.Println(rv, err)

	rv, err = dbm.SetStatTimeFunc(mdbm.ClockTsc)
	if err != nil {
		log.Fatalf("rv=%d, err=%v", rv, err)
	}

	fmt.Println(rv, err)

	// Output:
	// 0 <nil>
	// 0 <nil>
}

func Example_mdbm_StatAllPage() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	out, err := dbm.StatAllPage()

	fmt.Println(out, err)
	// Output:
	// there is only a v2 implementation. v3 or higher version not currently supported (current version=3)
}

func Example_mdbm_GetStats() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, stats, err := dbm.GetStats()

	fmt.Println(rv, err)

	fmt.Println("stat.Size =", stats.Size)
	fmt.Println("stat.PageSize =", stats.PageSize)
	fmt.Println("stat.PageCount =", stats.PageCount)
	fmt.Println("stat.PagesUsed =", stats.PagesUsed)
	fmt.Println("stat.BytesUsed =", stats.BytesUsed)
	fmt.Println("stat.NumEntries =", stats.NumEntries)
	fmt.Println("stat.MinLevel =", stats.MinLevel)
	fmt.Println("stat.MaxLevel =", stats.MaxLevel)
	fmt.Println("stat.LargePageSize =", stats.LargePageSize)

	// Output:
	// 0 <nil>
	// stat.Size = 5193728
	// stat.PageSize = 4096
	// stat.PageCount = 1268
	// stat.PagesUsed = 510
	// stat.BytesUsed = 633718
	// stat.NumEntries = 65634
	// stat.MinLevel = 1
	// stat.MaxLevel = 1
	// stat.LargePageSize = 4096
}

func Example_mdbm_GetDBInfo() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, info, err := dbm.GetDBInfo()

	fmt.Println(rv, err)

	fmt.Println("DBInfo.PageSize =", info.PageSize)
	fmt.Println("DBInfo.NumPages =", info.NumPages)
	fmt.Println("DBInfo.MaxPages =", info.MaxPages)
	fmt.Println("DBInfo.NumDirPages =", info.NumDirPages)
	fmt.Println("DBInfo.DirWidth =", info.DirWidth)
	fmt.Println("DBInfo.MaxDirShift =", info.MaxDirShift)
	fmt.Println("DBInfo.DirMinLevel =", info.DirMinLevel)
	fmt.Println("DBInfo.DirMaxLevel =", info.DirMaxLevel)
	fmt.Println("DBInfo.DirNumNodes =", info.DirNumNodes)
	fmt.Println("DBInfo.HashFunc =", info.HashFunc)
	fmt.Println("DBInfo.HashFuncName =", info.HashFuncName)

	// Output:
	// 0 <nil>
	// DBInfo.PageSize = 4096
	// DBInfo.NumPages = 1268
	// DBInfo.MaxPages = 0
	// DBInfo.NumDirPages = 1
	// DBInfo.DirWidth = 2
	// DBInfo.MaxDirShift = 1
	// DBInfo.DirMinLevel = 1
	// DBInfo.DirMaxLevel = 1
	// DBInfo.DirNumNodes = 1
	// DBInfo.HashFunc = 5
	// DBInfo.HashFuncName = FNV
}

func Example_mdbm_GetDBStats() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	// Do not lock for overall operation
	//rv, info, stat, err := dbm.GetDBStats(mdbm.StatNolock)
	rv, _, _, err := dbm.GetDBStats(mdbm.StatNolock)
	fmt.Println(rv, err)
	//pp.Println(info)
	//pp.Println(stat)

	// Do no lock for page-based iteration
	//rv, info, stat, err = dbm.GetDBStats(mdbm.IterateNolock)
	rv, _, _, err = dbm.GetDBStats(mdbm.IterateNolock)
	fmt.Println(rv, err)
	//pp.Println(info)
	//pp.Println(stat)

	// Output:
	// 0 <nil>
	// 0 <nil>
}

func Example_mdbm_GetWindowStats() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, wstats, err := dbm.GetWindowStats()

	fmt.Println(rv, err)

	fmt.Printf("num of reused=%d, ", wstats.WnumReused)
	fmt.Printf("num of remapped=%d, ", wstats.WnumRemapped)
	fmt.Printf("windows size=%d, ", wstats.WwindowSize)
	fmt.Printf("max window used=%d\n", wstats.WmaxWindowUsed)

	// Output:
	// 0 <nil>
	// num of reused=0, num of remapped=0, windows size=0, max window used=0
}

func Example_mdbm_GetHashValue() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	hashlist := [...]int{mdbm.HashCRC32,
		mdbm.HashEJB,
		mdbm.HashPHONG,
		mdbm.HashOZ,
		mdbm.HashTOREK,
		mdbm.HashFNV,
		mdbm.HashSTL,
		mdbm.HashMD5,
		mdbm.HashSHA1,
		mdbm.HashJENKINS,
		mdbm.HashHSIEH,
		mdbm.MaxHash,
		mdbm.DefaultHash,
	}

	for _, hashtype := range hashlist {
		rv, err := dbm.GetHashValue(1, hashtype)
		fmt.Println(rv, err)
	}

	// Output:
	// 2667302803 <nil>
	// 17 <nil>
	// 2621031278 <nil>
	// 49 <nil>
	// 49 <nil>
	// 1224750888 <nil>
	// 49 <nil>
	// 943901380 <nil>
	// 723085877 <nil>
	// 2366665294 <nil>
	// 3927678806 <nil>
	// 3927678806 <nil>
	// 1224750888 <nil>
}

func Example_mdbm_Plock_Punlock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.Plock(1)
	if err != nil {
		log.Fatalf("Plock(1) : rv=%d, err=%v", rv, err)
	}

	if rv == 1 {
		fmt.Println("Plock(1) : Success, partition lock was acquired")
	} else {
		log.Fatalf("Plock(1) : rv=%d, err=%v", rv, err)
	}

	rv, err = dbm.Plock(2)
	if err != nil {
		log.Fatalf("Plock(2) : rv=%d, err=%v", rv, err)
	}

	if rv == 1 {
		fmt.Println("Plock(2) : Success, partition lock was acquired")
	} else {
		log.Fatalf("Plock(2) : rv=%d, err=%v", rv, err)
	}

	rv, err = dbm.Store(1, 777, mdbm.Replace)
	if err != nil {
		log.Fatalf("Store : rv=%d, err=%v", rv, err)
	}

	fmt.Println("Store(2) :", rv, err)

	rv, err = dbm.Punlock(1)
	if err != nil {
		log.Fatalf("Punlock(1) : rv=%d, err=%v", rv, err)
	}

	if rv == 1 {
		fmt.Println("Punlock(1) : Success, partition lock was released")
	} else {
		log.Fatalf("Punlock(1) : rv=%d, err=%v", rv, err)
	}

	// Output:
	// Plock(1) : Success, partition lock was acquired
	// Plock(2) : Success, partition lock was acquired
	// Store(2) : 0 <nil>
	// Punlock(1) : Success, partition lock was released
}

func Example_mdbm_TryPlock() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, err := dbm.Plock(1)
	if err != nil {
		log.Fatalf("Plock(1) : rv=%d, err=%v", rv, err)
	}

	if rv == 1 {
		fmt.Println("Plock(1) : Success, partition lock was acquired")
	} else {
		log.Fatalf("Plock(1) : rv=%d, err=%v", rv, err)
	}

	rv, err = dbm.TryPlock(1)
	if err != nil {
		log.Fatalf("TryPlock(1) : rv=%d, err=%v", rv, err)
	}

	rv, err = dbm.Punlock(1)
	if err != nil {
		log.Fatalf("Punlock(1) : rv=%d, err=%v", rv, err)
	}

	if rv == 1 {
		fmt.Println("Punlock(1) : Success, partition lock was released")
	} else {
		log.Fatalf("Punlock(1) : rv=%d, err=%v", rv, err)
	}

	// Output:
	// Plock(1) : Success, partition lock was acquired
	// Punlock(1) : Success, partition lock was released
}

// ISSUE : sometimes it is deadlock
// HIGH RECOMMDED : mdbm.StoreWithLockSmart
/*
func Example_mdbm_LockSamrt_Store_UnLockSmart() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
	log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {

		dbm.LockSmart(i, mdbm.Rdrw)

		rv, err := dbm.Store(i, i, mdbm.Replace)
		if err != nil {
			log.Fatalf("Store(%s,%s,mdbm.Replace) : rv=%d, err=%v", i, i, rv, err)
		}

		dbm.UnLockSmart(i, mdbm.Rdrw)
	}

	// Output:
}
*/

func Example_mdbm_StoreWithLockSamrt() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	_, err = dbm.LockReset()
	if err != nil {
		log.Fatalf("failed mdbm.LockReset(), err=%v", err)
	}

	for i := 0; i <= loopLimit; i++ {

		rv, err := dbm.StoreWithLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		if err != nil {
			log.Fatalf("Store(%d,%d,mdbm.Replace) : rv=%d, err=%v", i, i, rv, err)
		}
	}

	// Output:
}

func Example_mdbm_StoreRWithLockSamrt() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	_, err = dbm.LockReset()
	if err != nil {
		log.Fatalf("failed mdbm.LockReset(), err=%v", err)
	}

	for i := 0; i <= loopLimit; i++ {

		rv, _, err := dbm.StoreRWithLockSmart(i, i, mdbm.Replace, mdbm.Rdrw, &iter)
		if err != nil {
			log.Fatalf("Store(%d,%d,mdbm.Replace) : rv=%d, err=%v", i, i, rv, err)
		}
	}

	// Output:
}

func Example_mdbm_CheckResidency() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	rv, pgsin, pgsout, err := dbm.CheckResidency()

	fmt.Println(rv, pgsin, pgsout, err)

	// Output:
	// 0 1268 0 <nil>
}