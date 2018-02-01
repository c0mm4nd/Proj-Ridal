package main

import (
	"crypto/rand"
	"github.com/maoxs2/Proj-Ridal/crypto/bliss"
	"log"
)

func main() {
	// b := []byte{}
	// a :=
	pri, pub, err := bliss.Bliss.GenerateKey(rand.Reader)
	if err != nil {
		log.Println(err)
	}
	log.Println("pri is", string(pri.Serialize()))
	log.Println("pub is", pub.SerializeCompressed())
}
