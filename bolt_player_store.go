package main

import (
	"cmp"
	"encoding/binary"
	"log"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

var playerScoreBucket = []byte("PlayerScores")

type BoltPlayerStore struct {
	mu sync.Mutex
	db *bolt.DB
}

func NewBoltPlayerStore(db_name string) (*BoltPlayerStore, error) {
	db_path := "./databases/" + db_name
	db, err := bolt.Open(db_path, 0600, &bolt.Options{Timeout: 1 * time.Second})
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

func (b *BoltPlayerStore) GetLeague() []Player {
	var players []Player
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(playerScoreBucket)

		bucket.ForEach(func(name, wins []byte) error {
			player := Player{string(name), btoi(wins)}

			players = append(players, player)

			return nil
		})

		return nil
	})

	slices.SortFunc(players, func(a Player, b Player) int {
		return cmp.Compare(b.Wins, a.Wins)
	})

	sort.Slice(players, func(i, j int) bool {
		return i < j
	})

	return players
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
