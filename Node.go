package main

import (
	"bytes"
	sha256 "crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	// "os"
	"strconv"
	"time"

	// glue "github.com/desertbit/glue"
	rpc "github.com/gorilla/rpc"
	jsonRpcCodec "github.com/gorilla/rpc/json"
	// websocket "github.com/gorilla/websocket"
	wsServer "github.com/googollee/go-socket.io"
	wsClient "github.com/zhouhui8915/go-socket.io-client"
)

var http_port = 33333
var p2p_port = 44444

type HttpService struct{}

type Block struct {
	Index     int64  `json:"index,string"`
	PrevHash  []byte `json:"prevhash,string"`
	Timestamp int64  `json:"timestamp,string"`
	Data      []byte `json:"data,string"`
	Hash      []byte `json:"hash,string"`
}

func getGenesisBlock() *Block {
	b := &Block{
		Index:     int64(0),
		PrevHash:  []byte("0"),
		Timestamp: int64(1465154705),
		Data:      []byte("Genesis"),
		Hash:      []byte("1"),
	}
	return b
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

var BlockChain = [](*Block){}

type MineBlockArgs struct {
	Data string
}

type MineBlockReply struct {
	Res string
}

func (h *HttpService) MineBlocks(r *http.Request, args *MineBlockArgs, reply *MineBlockReply) error {
	data := args.Data
	newBlock := generateNextBlock([]byte(data))
	if newBlock.Timestamp-getLatestBlock().Timestamp < 1000 {
		return fmt.Errorf("Too fast")
	}
	addBlock(newBlock)
	var emptySocket wsServer.Socket
	broadcast(emptySocket, responseLatestMsg())
	jo, _ := json.Marshal(newBlock)
	log.Println("block added: ", jo)
	return nil
}

type LocalBlocksArgs struct{}
type LocalBlocksReply struct {
	BlockList [](*Block) `json:"blocks"`
}

func (h *HttpService) LocalBlocks(r *http.Request, args *LocalBlocksArgs, reply *LocalBlocksReply) error {
	reply.BlockList = BlockChain
	return nil
}

type LocalPeersArgs struct{}
type LocalPeersReply struct {
	PeerList [](*Peer)
}

func (h *HttpService) LocalPeers(r *http.Request, args *LocalPeersArgs, reply *LocalPeersReply) error {
	// TODO
	return nil
}

type AddPeersArgs struct {
	PeerIpAddr string
	PeerPort   string
}
type AddPeersReply struct {
	Res string
}

func (h *HttpService) AddPeers(r *http.Request, args *string, reply *string) error {
	// TODO
	return nil
}

func initHttp() {
	s := rpc.NewServer()
	s.RegisterCodec(jsonRpcCodec.NewCodec(), "application/json")
	s.RegisterService(new(HttpService), "")
	h := &http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	log.Fatal(h.ListenAndServe())
}

func initP2P() {
	// Create a new glue server.
	// server := glue.NewServer(glue.Options{
	// 	HTTPListenAddress: ":8081",
	// })

	// // Release the glue server on defer.
	// // This will block new incoming connections
	// // and close all current active sockets.
	// defer server.Release()

	// // Set the glue event function to handle new incoming socket connections.
	// server.OnNewSocket(onNewSocket)

	// // Run the glue server.
	// err := server.Run()
	// if err != nil {
	// 	log.Fatalln("P2P Server Runs on: %v", err)
	// }

	server, err := wsServer.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(ws wsServer.Socket) {
		// initConnection(ws)
		// ServerSocketsList = append(ServerSocketsList, ws)

		ws.Join("ConnectedClients") // instead of append

		ws.On("message", func(ws wsServer.Socket, message string) {
			var msg *MsgFmt
			json.Unmarshal([]byte(message), msg)
			log.Println("Received message", message)
			switch msg.Type {
			case MSGTYPE_QUERY_LATEST:
				ws.Emit("message", responseLatestMsg())
				break
			case MSGTYPE_QUERY_ALL:
				ws.Emit("message", responseChainMsg())
				break
			case MSGTYPE_RESPONSE_BLOCKCHAIN:
				handleBlockchainResponse(ws, msg.Data)
				break
			}
		})

		ws.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(ws wsServer.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/", server)
	log.Println("Serving P2P at localhost:8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

type MsgFmt struct {
	Type int
	Data string // json fmt
}

// func onNewSocket(s *glue.Socket) {
// 	// act as initConn
// 	// ServerSocketsList = append(ServerSocketsList, s)

// 	// // Set a function which is triggered as soon as the socket is closed.
// 	// s.OnClose(func() {
// 	// 	log.Printf("socket closed with remote address: %s", s.RemoteAddr())
// 	// })

// 	// // Set a function which is triggered during each received message.
// 	// s.OnRead(func(message string) {
// 	// 	// Echo the received data back to the client.
// 	// 	var msg *MsgFmt
// 	// 	json.Unmarshal([]byte(message), &msg)

// 	// 	log.Println("Received message", message)
// 	// 	switch msg.Type {
// 	// 	case MSGTYPE_QUERY_LATEST:
// 	// 		s.Write(responseLatestMsg())
// 	// 		break
// 	// 	case MSGTYPE_QUERY_ALL:
// 	// 		s.Write(responseChainMsg())
// 	// 		break
// 	// 	case MSGTYPE_RESPONSE_BLOCKCHAIN:
// 	// 		handleBlockchainResponse(s, msg.Data)
// 	// 		break
// 	// 	}

// 	// })

// 	// // Send a welcome string to the client.
// 	// s.Write(queryChainLengthMsg())

// }

func addBlock(newBlock *Block) {
	if isValidNewBlock(newBlock, getLatestBlock()) == true {
		BlockChain = append(BlockChain, newBlock)
	}
}

func isValidNewBlock(newBlock *Block, PrevBlock *Block) bool {
	if PrevBlock.Index+1 != newBlock.Index {
		log.Println("invalid Index")
		return false
	} else if bytes.Equal(PrevBlock.Hash, newBlock.PrevHash) == false {
		log.Println("invalid previoushash")
		return false
	} else if bytes.Equal(calculateHashForBlock(newBlock), newBlock.Hash) == false {
		// log.Println(typeof (newBlock.hash) + " " + typeof calculateHashForBlock(newBlock));
		log.Println("invalid hash: " + string(calculateHashForBlock(newBlock)[:]) + " " + string(newBlock.Hash[:]))
		return false
	}
	return true
}

type Peer struct {
	IpAddr net.IP
	Port   int
	Conn   *net.Conn
}

const (
	MSGTYPE_QUERY_LATEST = iota
	MSGTYPE_QUERY_ALL
	MSGTYPE_RESPONSE_BLOCKCHAIN
)

var initialPeer = &Peer{net.ParseIP("119.28.176.178"), 8081, nil}

var initialPeers = [](*Peer){initialPeer}

func connectToPeers(newPeers [](*Peer)) {
	// act as Client
	for i := range newPeers {
		var peer = newPeers[i]
		opts := &wsClient.Options{
			Transport: "websocket",
			Query:     make(map[string]string),
		}

		uri := "http://" + peer.IpAddr.String() + ":" + strconv.Itoa(peer.Port)

		c, err := wsClient.NewClient(uri, opts)
		if err != nil {
			log.Println("Failed on NewClient:", uri)
		}
		c.On("open", func() { log.Println("Open") })
		c.On("error", func() { log.Println("2P Conn Failed ") })
	}
}

func initConnection(ws wsServer.Socket) {
	// Just for Client
	// About Sever side, Turn to onNewSocket

	// ServerSocketsList = append(ServerSocketsList, ws)

	ws.On("message", func(data string) {
		var msg *Msg
		json.Unmarshal([]byte(data), msg)
		log.Println("Received message", data)
		switch msg.Type {
		case MSGTYPE_QUERY_LATEST:
			ws.Emit("message", responseLatestMsg())
			break
		case MSGTYPE_QUERY_ALL:
			ws.Emit("message", responseChainMsg())
			break
		case MSGTYPE_RESPONSE_BLOCKCHAIN:
			// handleBlockchainResponse(ws, msg.Data)
			break
		}
	})

	ws.Emit("message", queryChainLengthMsg())
	// sockets.push(ws);
	// initMessageHandler(ws);
	// initErrorHandler(ws);
	// write(ws, queryChainLengthMsg());

}

type MsgDataFmt struct{}

func handleBlockchainResponse(ws wsServer.Socket, data string) {
	var blocksData [](*Block)
	json.Unmarshal([]byte(data), blocksData)
	receivedBlocks := blocksData
	var latestBlockReceived = receivedBlocks[len(receivedBlocks)-1]
	var latestBlockHeld = getLatestBlock()
	if latestBlockReceived.Index > latestBlockHeld.Index {
		log.Println("blockchain possibly behind. We got: ", latestBlockHeld.Index, " Peer got: ", latestBlockReceived.Index)
		if bytes.Equal(latestBlockHeld.Hash, latestBlockReceived.PrevHash) {
			log.Println("We can append the received block to our chain")
			BlockChain = append(BlockChain, latestBlockReceived)
			// broadcast(responseLatestMsg())
			ws.BroadcastTo("ConnectedClients", "message", responseLatestMsg(), func(ws wsServer.Socket, data string) {
				log.Println("Broadcasted ")
			})
		} else if len(receivedBlocks) == 1 {
			log.Println("We have to query the chain from our peer")
			// broadcast(queryAllMsg())
			ws.BroadcastTo("ConnectedClients", "message", queryAllMsg(), func(ws wsServer.Socket, data string) {
				log.Println("Broadcasted ")
			})
		} else {
			log.Println("Received blockchain is longer than current blockchain")
			replaceChain(receivedBlocks)
		}
	} else {
		log.Println("received blockchain is not longer than current blockchain. Do nothing")
	}
}

func generateNextBlock(blockData []byte) *Block {
	var prevBlock = getLatestBlock()
	var nextIndex = prevBlock.Index + 1
	var nextTimestamp = time.Now().Unix() // / 1000
	var nextHash = calculateHash(nextIndex, prevBlock.Hash, nextTimestamp, blockData)

	return &Block{nextIndex, prevBlock.Hash, nextTimestamp, blockData, nextHash}
}

func calculateHash(Index int64, PrevHash []byte, Timestamp int64, Data []byte) []byte {
	h := sha256.New()
	var combination = [][]byte{Int64ToBytes(Index), PrevHash, Int64ToBytes(Timestamp), Data}
	h.Write(bytes.Join(combination, []byte("")))
	return h.Sum(nil)
}

func calculateHashForBlock(b *Block) []byte {
	return calculateHash(b.Index, b.PrevHash, b.Timestamp, b.Data)
}

func replaceChain(newBlocks [](*Block)) {
	if isValidChain(newBlocks) && len(newBlocks) > len(BlockChain) {
		log.Println("Received blockchain is valid. Replacing current blockchain with received blockchain")
		BlockChain = newBlocks
		var emptySocket wsServer.Socket
		broadcast(emptySocket, responseLatestMsg())
	} else {
		log.Println("Received blockchain invalid")
	}

}

func isValidChain(blockchainToValidate [](*Block)) bool {
	UnknownGenesisBlock := blockchainToValidate[0]
	if UnknownGenesisBlock != getGenesisBlock() {
		return false
	}
	var tempBlocks = [](*Block){blockchainToValidate[0]}
	for i := range blockchainToValidate {
		if isValidNewBlock(blockchainToValidate[i], tempBlocks[i-1]) {
			tempBlocks = append(tempBlocks, blockchainToValidate[i])
		} else {
			return false
		}
	}
	return true
}

var ServerSocketsList [](wsServer.Socket)

var PeerSocketList [](*wsClient.Client)

func broadcast(ws wsServer.Socket, msg string) {
	// for i := range SocketsList {
	// 	SocketsList[i].Write(msg)
	// }
	ws.BroadcastTo("ConnectedClients", "message", msg, func(ws wsServer.Socket, data string) {
		log.Println("Broadcasted ", msg)
	})
}

func getLatestBlock() *Block {
	return BlockChain[len(BlockChain)-1]
}

// Msg
type Msg struct {
	Type int
	Data []byte
}

func queryChainLengthMsg() string {
	r := new(Msg)
	r.Type = MSGTYPE_QUERY_LATEST
	rtn, _ := json.Marshal(r)
	return string(rtn)
}

func queryAllMsg() string {
	r := new(Msg)
	r.Type = MSGTYPE_QUERY_ALL
	rtn, _ := json.Marshal(r)
	return string(rtn)
}

func responseChainMsg() string {
	r := new(Msg)
	r.Type = MSGTYPE_RESPONSE_BLOCKCHAIN
	r.Data, _ = json.Marshal(BlockChain)
	rtn, _ := json.Marshal(r)
	return string(rtn)
}

func responseLatestMsg() string {
	r := new(Msg)
	r.Type = MSGTYPE_RESPONSE_BLOCKCHAIN
	r.Data, _ = json.Marshal(getLatestBlock())
	rtn, _ := json.Marshal(r)
	return string(rtn)
}

// addition utils
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func main() {
	genesisBlock := getGenesisBlock()
	BlockChain = append(BlockChain, genesisBlock)
	// jo, _ := json.Marshal(genesisBlock)
	// os.Stdout.Write(jo)
	// debug.Dump(getGenesisBlock())	)
	go connectToPeers(initialPeers)
	go initHttp()
	go initP2P()
	for {
		time.Sleep(1000)
	}
}

// 约定俗称的Block内容大致不变
// func: addBlock()

// isVaildNewBlock() 检测新Block

// replaceChain

// isVaildChain

// P2P网络工作内容：
//     建立socket通道，~~发送~~沟通节点状态（最新？、完整？、处理反馈）

// http端（jsonrpc）内容：
//     blocks()返回本地区块
//     mineBlocks()先生成新区块，再addBlock增加该区块，然后广播并取回回复
//     peers链接节点列表
//     addPeer增加连接节点

// 必须的内容补充：
//     mineBlocks的PoW实现
//     Blocks的存储

//     rpc和ws（P2P）内容可以融合
