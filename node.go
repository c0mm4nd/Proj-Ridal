package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
	//"crypto/sha256"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)

/* ==== Define Type Start ===== */
type Generation struct {
	id                int64
	hash              [32]byte // sha256 checksum
	previousHash      [32]byte
	timestamp         int64 // time.Now().Unix()
	previousTimestamp int64
	data              [1024]byte
}

type HostJob int

type ConfirmGeneration struct {
	confirmRound     int
	confirmServerNum int
	generation       Generation
}

/* ==== Define Type End ===== */

/* ==== Claim Global Var Start ===== */
// all of guys will be used for many times, too lazy to set as arg
var levelDB *leveldb.DB
var err error
var nodePool []string
var lastGeneration Generation
var currentGeneration Generation

/* ==== Claim Global Var End ===== */

/* ===== Host Job list Start ======== */

// TODO: Maybe I need to write a wiki to normalize req & res ?

func (t *HostJob) GetNodePool(args int, reply *[]string) error {
	log.Print("run GetNodePool")
	*reply = nodePool
	return nil
}

func (t *HostJob) AddToNodePool(args string, reply *int) error {
	// if err rtn 0
	if args != "" {
		nodePool = append(nodePool, args)
		// if err != nil {
		//     *reply = 0
		// }else{
		*reply = 1
		// }
	}
	return nil
}

func (t *HostJob) GatherTransactions(args string, reply *int) error {
	// currentGeneration.data = append(currentGeneration.data, args)
	*reply = 1
	return nil
}

func (t *HostJob) confirmGeneration(args Generation, reply *int) error {
	if args.hash == currentGeneration.hash { // compare with yourself
		*reply = 1
	} else {
		*reply = 0
	}
	return nil
}

func (t *HostJob) miningCompetition(args string, reply *int) error {
	log.Print("Init the Mining Dealer.")
	// Temprorily Simplify the Mining Period
	if time.Now().Unix()-lastGeneration.timestamp == 10 {
		// 10 seconds to Update and Auto Boardcast
		// Mined
		broadcastGeneration()
	} else {
		log.Print("Not reach the Time.")
	}
	*reply = 1
	return nil
}

/* ===== Host Job list End ======== */

/* ========= Host Server Start ======*/
// What req will be sent to the server?
// getNodePool, getGenerationChains, boardcastTransactions and so on.

func HostMain() {
	log.Print("Host Starting")
	rpcServer := rpc.NewServer()
	HostJob := new(HostJob)
	rpcServer.Register(HostJob)
	service := ":10101"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	log.Print(tcpAddr)
	conn, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	go rpcServer.Accept(conn)
}

/* ========= Host Server End ======*/

/* ========= Local Client Start ======*/

func LocalMain() {
	log.Print("Local Starting")
	// Define: a rpc client with p2p connection persisted

	// if len(os.Args) != 2 {
	//     fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
	//     os.Exit(1)
	// }
	service := "119.28.176.178:10101" // Initial Node, will able to be costomed
	// first step is sync
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	rpcClient := rpc.NewClient(conn)
	for {

		rpcClient.Call("HostJob.GetNodePool", 0, &nodePool)
		log.Print(nodePool)
		time.Sleep(10 * time.Second)
	}
	// nodePool :=
	// for {
	//     // period to persist the connection
	//     _, err = conn.Write([]byte("anything"))
	//     fmt.Println("Sending")
	//     checkError(err)
	//     var buf [512]byte
	//     fmt.Println("Ready for Read")
	//     n, err := conn.Read(buf[0:])
	//     checkError(err)
	//     fmt.Println(string(buf[0:n]))
	//     // os.Exit(0)
	//     time.Sleep(300)
	// }
}

/* ========= Local Client Start ======*/

func broadcastGeneration() {
	// need 20 round Confirm
	// if broadcastSuccess {
	// 	saveNewGeneration() // save the Newest Generation to DB
	// }
}

func saveNewGeneration(newestGeneration Generation) {
	// Just save the total Generation Object
	// into the leveldb in the form of Json.
	// key: id (height) value: object json
	json.Marshal(newestGeneration)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err)
		log.Fatal(err)
		os.Exit(1)
	}
}

func leveldbSaver() {

}

func startInitLocalGeneration() {

}

/* ============================================================================*/
/* ============================================================================*/
/* =================================I'm the Owner==============================*/
/* ===================================the King=================================*/
/* ============================================================================*/
/* ============================================================================*/

func main() {
	levelDB, err = leveldb.OpenFile("level.db", nil)
	if err != nil {
		log.Print("err is ", err)
	} else {
		log.Print("db is ", levelDB)
	}
	data, err := levelDB.Get([]byte("key"), nil)
	if err == nil {
		if err == leveldb.ErrNotFound {
			log.Print("LocalDB has NO data. Initializing")
			startInitLocalGeneration()
		} else {
			log.Fatal("err:", err, " || Program closed")
		}
	} else {
		log.Print(data) // TODO: check the height
	}

	genesisGeneration := &Generation{
		id:                0,
		hash:              [32]byte{},
		previousHash:      [32]byte{},
		timestamp:         time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(), // time.Now().Unix()
		previousTimestamp: time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(),
		data:              [1024]byte{},
	}
	log.Print(genesisGeneration) // for reserving the gensis, will be del

	nodePool = []string{"119.28.176.178"}
	go HostMain()
	go LocalMain()
	for {
		time.Sleep(10)
	}
}

/* ============================================================================*/
/* ============================================================================*/
/* ======================================All===================================*/
/* ===================================Falls Down===============================*/
/* ============================================================================*/
/* ============================================================================*/
