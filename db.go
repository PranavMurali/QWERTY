package main

import (
	"bytes"
	"encoding/json"
	crand "crypto/rand"
	rand "math/rand"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os/user"
	"encoding/binary"
	"time"
)

// Config type
type Config struct {
	UserId   int   `json:"userid"`
	UserName string    `json:"username"`
}

// Entry type
type Entry struct {
	UserId   int   `json:"userid"`
	UserName string    `json:"username"`
	Password string    `json:"password"`
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	user, err := user.Current()
	tmp:= user.Username

	var src cryptoSource
	rnd := rand.New(src)
	ndate := rnd.Intn(100000)

	conf := Config{UserId: ndate, UserName: tmp}
	err = setConfig(db, conf)
	if err != nil {
		log.Fatal(err)
	}
	err = addWeight(db, "password", time.Now())
	if err != nil {
		log.Fatal(err)
	}
	err = addEntry(db, 21332,"red","password",time.Now())
	if err != nil {
		log.Fatal(err)
	}

	err = addEntry(db,123,"asd", "pasdwrpd",time.Now().AddDate(0, 0, -2))
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		conf := tx.Bucket([]byte("DB")).Get([]byte("CONFIG"))
		fmt.Printf("Config: %s\n", conf)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("PASSWORD"))
		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(k), string(v))
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("DB")).Bucket([]byte("ENTRIES")).Cursor()
		min := []byte(time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
		max := []byte(time.Now().AddDate(0, 0, 0).Format(time.RFC3339))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Println(string(k), string(v))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("profiles.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("PASSWORD"))
		if err != nil {
			return fmt.Errorf("could not create weight bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("ENTRIES"))
		if err != nil {
			return fmt.Errorf("could not create days bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}

func setConfig(db *bolt.DB, config Config) error {
	confBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Put([]byte("CONFIG"), confBytes)
		if err != nil {
			return fmt.Errorf("could not set config: %v", err)
		}
		return nil
	})
	fmt.Println("Set Config")
	return err
}

func addWeight(db *bolt.DB, password string, date time.Time) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("PASSWORD")).Put([]byte(date.Format(time.RFC3339)), []byte(password))
		if err != nil {
			return fmt.Errorf("could not insert weight: %v", err)
		}
		return nil
	})
	fmt.Println("Added Weight")
	return err
}

func addEntry(db *bolt.DB, userid int,username string,password string ,date time.Time) error {
	entry := Entry{UserId:userid, UserName:username,Password:password}
	entryBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("could not marshal entry json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("ENTRIES")).Put([]byte(date.Format(time.RFC3339)), entryBytes)
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}

		return nil
	})
	fmt.Println("Added Entry")
	return err
}
