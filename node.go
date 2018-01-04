package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	// "strconv"
	gxMgr "github.com/maoxs2/Proj_Ridal/generationManager"
	"time"
)

/* ==== Define Type Start ===== */

type privateKey []byte

type seedPhrase []byte // some words -> decide the privateKey

type publickey []byte // gene from private using bliss crypto X sha for test ...

type accountId uint32

type Account struct {
	accountId  uint32
	publickey  []byte
	balance    uint32
	operations uint32
	name       []byte
}

type Tx struct {
	from accountId
	to   accountId
	fee  uint32
}

type HostJob int

// type ConfirmGeneration struct {
// 	confirmRound     int
// 	confirmServerNum int
// 	generation       *Generation
// }

/* ==== Define Type End ===== */

/* ==== Claim Global Var Start ===== */
// all of guys will be used for many times, too lazy to set as arg

var err error
var nodePool []string

/* ==== Claim Global Var End ===== */

/* ===== Host Job list Start ======== */

// TODO: Maybe I need to write a wiki to normalize req & res ?

// type commandHandler func(*rpcServer, interface{}, <-chan struct{}) (interface{}, error)

// var rpcHandlersBeforeInit = map[string]commandHandler{
// 	"addnode":              handleAddNode,
// 	"createrawsstx":        handleCreateRawSStx,
// 	"createrawssgentx":     handleCreateRawSSGenTx,
// 	"createrawssrtx":       handleCreateRawSSRtx,
// 	"createrawtransaction": handleCreateRawTransaction,
// 	// "debuglevel":            handleDebugLevel,
// 	"decoderawtransaction":  handleDecodeRawTransaction,
// 	"decodescript":          handleDecodeScript,
// 	"estimatefee":           handleEstimateFee,
// 	"estimatestakediff":     handleEstimateStakeDiff,
// 	"existsaddress":         handleExistsAddress,
// 	"existsaddresses":       handleExistsAddresses,
// 	"existsmissedtickets":   handleExistsMissedTickets,
// 	"existsexpiredtickets":  handleExistsExpiredTickets,
// 	"existsliveticket":      handleExistsLiveTicket,
// 	"existslivetickets":     handleExistsLiveTickets,
// 	"existsmempooltxs":      handleExistsMempoolTxs,
// 	"generate":              handleGenerate,
// 	"getaddednodeinfo":      handleGetAddedNodeInfo,
// 	"getbestblock":          handleGetBestBlock,
// 	"getbestblockhash":      handleGetBestBlockHash,
// 	"getblock":              handleGetBlock,
// 	"getblockcount":         handleGetBlockCount,
// 	"getblockhash":          handleGetBlockHash,
// 	"getkeyblockhash":       handleGetKeyBlockHash,
// 	"getblockheader":        handleGetBlockHeader,
// 	"getblockkeyheight":     handleGetBlockKeyHeight,
// 	"getblocksubsidy":       handleGetBlockSubsidy,
// 	"getcoinsupply":         handleGetCoinSupply,
// 	"getconnectioncount":    handleGetConnectionCount,
// 	"getcurrentnet":         handleGetCurrentNet,
// 	"getdifficulty":         handleGetDifficulty,
// 	"getgenerate":           handleGetGenerate,
// 	"gethashespersec":       handleGetHashesPerSec,
// 	"getheaders":            handleGetHeaders,
// 	"getinfo":               handleGetInfo,
// 	"getmempoolinfo":        handleGetMempoolInfo,
// 	"getmininginfo":         handleGetMiningInfo,
// 	"getnettotals":          handleGetNetTotals,
// 	"getnetworkhashps":      handleGetNetworkHashPS,
// 	"getpeerinfo":           handleGetPeerInfo,
// 	"getrawmempool":         handleGetRawMempool,
// 	"getrawtransaction":     handleGetRawTransaction,
// 	"getstakedifficulty":    handleGetStakeDifficulty,
// 	"getstakeversioninfo":   handleGetStakeVersionInfo,
// 	"getstakeversions":      handleGetStakeVersions,
// 	"getticketpoolvalue":    handleGetTicketPoolValue,
// 	"getvoteinfo":           handleGetVoteInfo,
// 	"gettxout":              handleGetTxOut,
// 	"getwork":               handleGetWork,
// 	"help":                  handleHelp,
// 	"livetickets":           handleLiveTickets,
// 	"missedtickets":         handleMissedTickets,
// 	"node":                  handleNode,
// 	"ping":                  handlePing,
// 	"searchrawtransactions": handleSearchRawTransactions,
// 	"rebroadcastmissed":     handleRebroadcastMissed,
// 	"rebroadcastwinners":    handleRebroadcastWinners,
// 	"sendrawtransaction":    handleSendRawTransaction,
// 	"setgenerate":           handleSetGenerate,
// 	"stop":                  handleStop,
// 	"submitblock":           handleSubmitBlock,
// 	"ticketfeeinfo":         handleTicketFeeInfo,
// 	"ticketsforaddress":     handleTicketsForAddress,
// 	"ticketvwap":            handleTicketVWAP,
// 	"txfeeinfo":             handleTxFeeInfo,
// 	"validateaddress":       handleValidateAddress,
// 	"verifychain":           handleVerifyChain,
// 	"verifymessage":         handleVerifyMessage,
// 	"version":               handleVersion,
// }

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

func (t *HostJob) TestSaveNewGeneration(args string, reply *int) error {
	genesisGeneration := &gxMgr.Generation{
		Id:                0,
		Hash:              [32]byte{},
		PreviousHash:      [32]byte{},
		Timestamp:         time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(), // time.Now().Unix()
		PreviousTimestamp: time.Date(2017, time.December, 25, 23, 0, 0, 0, time.UTC).Unix(),
		Data:              [1024]byte{},
	}
	gxMgr.SaveNewGeneration(genesisGeneration)
	*reply = 1
	return nil
}

func (t *HostJob) ConfirmGeneration(args *gxMgr.Generation, reply *int) error {
	if args.Hash == gxMgr.CurrentGeneration.Hash { // compare with yourself
		*reply = 1
	} else {
		*reply = 0
	}
	return nil
}

func (t *HostJob) MiningCompetition(args string, reply *int) error {
	log.Print("Init the Mining Dealer.")
	// Temprorily Simplify the Mining Period
	if time.Now().Unix()-gxMgr.LastGeneration.Timestamp == 10 {
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
	if err != nil {
		log.Fatal(err)
	}
	// log.Print(tcpAddr)
	conn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	go rpcServer.Accept(conn)
}

/* ========= Host Server End ======*/

/* ========= Local Client Start ======*/

func LocalMain(rpcClient *rpc.Client) {
	log.Print("Local Starting")
	// Define: a rpc client with p2p connection persisted

	rpcClient.Call("HostJob.GetNodePool", 0, &nodePool)

	for {
		rpcClient.Call("HostJob.MiningCompetition", 0, &nodePool)
		time.Sleep(10 * time.Second)
	}
}

/* ========= Local Client End ======*/

func userFrontEnd(rpcClient *rpc.Client) {
	var userInput string
	// if nodePool.len < 1 {
	// 	return
	// }
	fmt.Print("$ ")
	fmt.Scanln(&userInput)

	callFunc := "HostJob."
	callFunc += userInput
	err = rpcClient.Call(callFunc, 0, &nodePool)
	if err != nil {
		log.Print("Err on rpcClient.Call:", callFunc)
	}
	// time.Sleep(10 * time.Second)
}

func broadcastGeneration() {
	// need 20 round Confirm
	// if broadcastSuccess {
	// 	saveNewGeneration() // save the Newest Generation to DB
	// }
}

func testMining() {}

/* ============================================================================*/
/* ============================================================================*/
/* =================================I'm the Owner==============================*/
/* ===================================the King=================================*/
/* ============================================================================*/
/* ============================================================================*/

func main() {

	gxMgr.InitLocalGeneration()

	service := "119.28.176.178:10101" // Initial Node, will able to be costomed
	// first step is sync
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal(err)
	}

	rpcClient := rpc.NewClient(conn)

	nodePool = []string{"119.28.176.178"}

	hostSwitch := flag.Bool("h", false, "Start Host")
	localSwitch := flag.Bool("l", false, "Start Client")
	flag.Parse()

	if *hostSwitch == true {
		log.Print("Host Switch is ", *hostSwitch)
		go HostMain()
	}

	if *localSwitch == true {
		log.Print("Local Switch is ", *localSwitch)
		go LocalMain(rpcClient)
	}

	if *localSwitch == false && *hostSwitch == false {
		log.Fatal("Neither Started")
	}

	for {
		if *localSwitch == true {
			userFrontEnd(rpcClient)
		} else {
			time.Sleep(10)
		}
	}
}

/* ============================================================================*/
/* ============================================================================*/
/* ======================================All===================================*/
/* ===================================Falls Down===============================*/
/* ============================================================================*/
/* ============================================================================*/
