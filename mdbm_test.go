package mdbm_test

import (
	"testing"

	"github.com/torden/go-mdbm"
	"github.com/torden/go-strutil"
)

var assert = strutils.NewAssert()

func TestNewMDBM(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen("./test1.mdbm", 0644)
	assert.AssertNil(t, err, "Error : %v", err)
	dbm.EasyClose()
}

func TestOpen(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open("./test1.mdbm", mdbm.Create|mdbm.Rdrw, 0666, 0, 0)
	assert.AssertNil(t, err, "Error : %v", err)
	dbm.Close()
}

func TestSync(t *testing.T) {

}
