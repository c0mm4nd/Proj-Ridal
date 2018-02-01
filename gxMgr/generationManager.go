package generationManager

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/maoxs2/Proj-Ridal/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

/* Do Sth Abt Generation (Local) */

type Generation struct {
	Id                uint32
	Hash              [32]byte // sha256 checksum
	PreviousHash      [32]byte
	Timestamp         int64 // time.Now().Unix()
	PreviousTimestamp int64
	Data              [1024]byte
}

var err error
var LocalHeight uint32
var LastGeneration *Generation
var CurrentGeneration *Generation
var levelDB *leveldb.DB

func init() {
	levelDB, err = leveldb.OpenFile("generationData", nil)

	if err != nil {
		log.Print("err is ", err)
	} else {
		// log.Print("DB is ", levelDB)
	}
}

func GetLocalHeight(db *leveldb.DB) uint32 {
	iter := db.NewIterator(nil, nil)
	var i uint32
	i = 0
	for iter.Next() {
		i = i + 1
	}
	iter.Release()
	return i
}

func InitLocalGeneration() *leveldb.DB {
	genesisGeneration := &Generation{
		Id:                0,
		Hash:              [32]byte{},
		PreviousHash:      [32]byte{},
		Timestamp:         time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(), // time.Now().Unix()
		PreviousTimestamp: time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(),
		Data:              [1024]byte{},
	}

	data, err := levelDB.Get([]byte(strconv.Itoa(0)), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			log.Print("LocalDB has NO data. Initializing")
			SaveNewGeneration(genesisGeneration)
		} else {
			log.Print("err:", err, " || Program closed")
		}
	} else {
		if data == nil {
			LocalHeight = 0
			SaveNewGeneration(genesisGeneration)
		} else {
			LocalHeight = GetLocalHeight(levelDB)
		}
	}

	return levelDB
}

func SaveNewGeneration(newestGeneration *Generation) {
	// Just save the total Generation Object
	// into the leveldb in the form of Json.
	// key: id (height) value: object json
	jsonNewestGeneration, err := json.Marshal(newestGeneration)
	if err == nil {
		hash := utils.CalcSha256Hash(jsonNewestGeneration)
		if LocalHeight == 0 {
			leveldbSaver(levelDB, int(LocalHeight), hash)
		} else {
			leveldbSaver(levelDB, int(LocalHeight+1), hash)
		}
	} else {
		log.Fatal("Error on saveNewGeneration")
	}
}

func leveldbSaver(db *leveldb.DB, key int, value []byte) {
	log.Print("key is", []byte(strconv.Itoa(key)))
	log.Print("value is ", value)
	err = db.Put([]byte(strconv.Itoa(key)), value, nil)
	if err != nil {
		log.Fatal("Err on db.Put", err)
	} else {
		log.Print("Saved:", string(value))
	}
}
