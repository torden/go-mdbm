package mdbm_test

import (
	"fmt"
	"log"

	"github.com/torden/go-mdbm"
)

const pathTestDBM1 = "./tmp/test1.mdbm"
const pathTestDBM2 = "./tmp/test2.mdbm"
const pathTestDBM3 = "./tmp/test3.mdbm"
const pathTestDBMHash = "./tmp/test_hash.mdbm"

func Example_EasyOpen_EasyClose() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	dbm.EasyClose()
	fmt.Println(err)
	// Output: <nil>
}

func Example_Open_Close() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM2, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	dbm.Close()
	fmt.Println(err)
	// Output: <nil>
}

func Example_DupHandle() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM3, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)

	dbm2, err2 := dbm.DupHandle()
	dbm2.Close()
	fmt.Println(err)
	fmt.Println(err2)
	// Output: <nil>
	// <nil>
}

func Example_GetErrNo() {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM1, mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	fmt.Println(err)

	rv, err := dbm.GetErrNo()
	fmt.Println(rv, err)
	dbm.Close()
	// Output:  <nil>
	// 0 <nil>
}

func Example_LogMinLevel() {

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

func Example_Sync() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	fmt.Println(err)
	// any
	rv, err := dbm.Sync()
	dbm.EasyClose()

	fmt.Println(rv, err)
	// Output: <nil>
	// 0 <nil>
}

func Example_Fsync() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	fmt.Println(err)
	// any
	rv, err := dbm.Fsync()
	dbm.EasyClose()

	fmt.Println(rv, err)
	// Output: <nil>
	// 0 <nil>
}

func Example_CloseFD() {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	fmt.Println(err)
	// any
	err = dbm.CloseFD()

	fmt.Println(err)
	// Output: <nil>
	// <nil>
}

func Example_Lock_Unlock() {

	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)
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

func Example_IsLocked() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)
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

func Example_LockShared() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)
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

func Example_TryLockShared() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)
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

func Example_LockReset() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM3, 0644)
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
	// LockReset : rv = -1 , err = no such file or directory
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_DeleteLockFiles() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)
	fmt.Println("EasyOpen : ", err)

	err = dbm.Lock()
	fmt.Println("Lock : ", err)

	rv, err = dbm.DeleteLockFiles()
	fmt.Println("DeleteLockFiles : rv =", rv, ", err =", err)

	rv, err = dbm.Store("iamKey", "iamValue", mdbm.Replace)
	fmt.Println("Store : rv =", rv, ", err =", err)

	err = dbm.Unlock()
	fmt.Println("Unlock : ", err)

	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// Lock :  <nil>
	// DeleteLockFiles : rv = 0 , err = <nil>
	// Store : rv = 0 , err = <nil>
	// Unlock :  <nil>
}

func Example_ReplaceDB() {

	var err error
	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBM1, 0644)

	fmt.Println("EasyOpen : ", err)
	dbm.EasyClose()

	err = dbm.EasyOpen(pathTestDBM1, 0644)
	fmt.Println("EasyOpen : ", err)

	err = dbm.ReplaceDB(pathTestDBM2)
	fmt.Println("ReplaceDB : ", err)
	dbm.EasyClose()

	// Output:
	// EasyOpen :  <nil>
	// EasyOpen :  <nil>
	// ReplaceDB :  <nil>
}

func Example_ReplaceFile() {

	var err error
	dbm := mdbm.NewMDBM()

	//create a dummy
	dbm.EasyOpen(pathTestDBM1, 0644)
	dbm.EasyClose()

	err = dbm.ReplaceFile(pathTestDBM1, pathTestDBM3)
	fmt.Println("ReplaceFile : ", err)
	dbm.EasyClose()

	// Output:
	// ReplaceFile :  <nil>
}

func Example_GetHash() {

	var err error
	dbm := mdbm.NewMDBM()

	//create a dummy
	dbm.EasyOpen(pathTestDBM1, 0644)

	rv, err := dbm.GetHash()
	fmt.Println("GetHash : rv =", rv, ", err = ", err)
	dbm.EasyClose()

	// mdbm.DefaultHash is 5

	// Output:
	// GetHash : rv = 5 , err =  <nil>
}

func Example_SetHash() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	dbm.EasyOpen(pathTestDBMHash, 0644)
	defer dbm.EasyClose()

	dbm.SetHash(mdbm.HashOZ)
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

func Example_GetAlignment() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	rv, err = dbm.GetAlignment()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv)
	// Output:
	// 0
}

func Example_SetAlignment() {

	var rv int
	var err error
	dbm := mdbm.NewMDBM()
	dbm.EasyOpen(pathTestDBM3, 0644)
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

func Example_GetLimitSize() {

	var rv uint64
	var err error
	dbm := mdbm.NewMDBM()
	dbm.EasyOpen(pathTestDBM3, 0644)
	defer dbm.EasyClose()

	rv, err = dbm.GetLimitSize()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rv)
	// Output:
	// 0
}
