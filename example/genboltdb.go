package main

import (
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

const dbPath = "/tmp/test_benchmark_large.boltdb"

func main() {
	db, err := bolt.Open(dbPath, 0644, nil)
	if err != nil {
		log.Fatalf("failured, can't open the boltdb, path=%s, err=%v", dbPath, err)
	}
	defer db.Close()

	bucketName := []byte("MyBucket")

	/* ISSUE : deadlock
	   // futex(0x550530, FUTEX_WAIT
	   	err = dlog.Update(func(tx *bolt.Tx) error {
	   		bucket, err := tx.CreateBucketIfNotExists(bucketName)
	   		if err != nil {
	   			log.Fatalf("failed, can't create or get bucket from boltDB(=%s), err=%v", dbPath, err)
	   			return err
	   		}

	   		for i := int(0); i <= 100000000; i++ {
	   			key := []byte(strconv.Itoa(i))
	   			value := []byte(strconv.Itoa(i))

	   			err = bucket.Put(key, value)
	   			if err != nil {
	   				log.Fatalf("failed, can't data(kv=%d) add to the boltdb(=%s), err=%v", i, dbPath, err)
	   				return err
	   			}
	   		}
	   		return nil
	   	})

	if err != nil {
		log.Fatalf("exception : boltdb(=%s), err=%v", dbPath, err)
	}
	*/
	for i := int(0); i <= 100000000; i++ {
		key := []byte(strconv.Itoa(i))
		value := []byte(strconv.Itoa(i))

		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				log.Fatalf("failed, can't data(kv=%d) add to the boltdb(=%s), err=%v", i, dbPath, err)
				return err
			}

			err = bucket.Put(key, value)
			if err != nil {
				log.Fatalf("failed, can't data(kv=%d) add to the boltdb(=%s), err=%v", i, dbPath, err)
				return err
			}
			return nil
		})

		if err != nil {
			log.Fatalf("exception : boltdb(=%s), err=%v", dbPath, err)
		}
	}

	log.Println("complete")
}
