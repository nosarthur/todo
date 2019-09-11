package db

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// FIXME: We should not know the bucket name
var bucketName = []byte("todos")
var db *bolt.DB

// Add a task
func Add(task string) error {
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		t := time.Now()
		return b.Put([]byte(t.Format("2006-01-02T22:08:41")), []byte(task))
	})
}

// List all tasks
func List() error {
	defer db.Close()
	return db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(bucketName).Cursor()
		i := 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%d. %s\n", i, v)
			i++
		}
		return nil
	})
}

// Rm deletes tasks
func Rm(nums map[int]struct{}) error {
	defer db.Close()
	var keysToDelete [][]byte
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		i := 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if _, ok := nums[i]; ok {
				keysToDelete = append(keysToDelete, k)
				// FIXME: we should print out after delete
				fmt.Printf("'%s' is deleted\n", v)
			}
			i++
		}

		for _, k := range keysToDelete {
			err := b.Delete(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// MustInit starts the db
func MustInit(dbPath string) {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("todos"))
		return err
	})
	if err != nil {
		panic(err)
	}

}
