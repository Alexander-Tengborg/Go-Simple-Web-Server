package stores

import (
	"encoding/binary"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

var playerScoreBucket = []byte("PlayerScores")

func NewBoltPlayerStore(db_name string) (*BoltPlayerStore, error) {
	db, err := bolt.Open(db_name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	store := &BoltPlayerStore{db: db}

	err = store.createBucket()
	if err != nil {
		return nil, err
	}

	return store, nil
}

type BoltPlayerStore struct {
	mu sync.Mutex
	db *bolt.DB
}

func (b *BoltPlayerStore) GetPlayerScore(name string) int {
	var score int
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(playerScoreBucket) //Check if bucket is not nil?

		bScore := bucket.Get([]byte(name))
		if bScore == nil {
			return nil // Return player not found error
		}

		score = btoi(bScore)

		return nil
	})

	return score
}

func (b *BoltPlayerStore) RecordWin(name string) {
	b.mu.Lock()
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(playerScoreBucket) //Check if bucket is not nil?

		var score int
		bScore := bucket.Get([]byte(name))
		if bScore != nil {
			score = int(binary.BigEndian.Uint64(bScore))
		}
		score++

		bScore = itob(score)

		return bucket.Put([]byte(name), bScore)
	})
	b.mu.Unlock()

	if err != nil {
		log.Println(err)
	}
}

func (b *BoltPlayerStore) Close() error {
	return b.db.Close()
}

func (b *BoltPlayerStore) ResetBucket() error {
	b.deleteBucket()
	err := b.createBucket()

	return err
}

func (b *BoltPlayerStore) deleteBucket() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(playerScoreBucket)
	})

	return err
}

func (b *BoltPlayerStore) createBucket() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(playerScoreBucket)
		return err
	})

	return err
}

// int to byte[]
func itob(number int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(number))
	return b
}

// byte[] to int
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
