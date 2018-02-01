package bliss

import (
	"github.com/LoCCS/bliss"
	Crypto "github.com/maoxs2/Proj-Ridal/crypto"
)

type PublicKey struct {
	Crypto.PublicKeyAdapter
	bliss.PublicKey
}

func (p PublicKey) GetType() int {
	return pqcTypeBliss
}

func (p PublicKey) Serialize() []byte {
	return p.PublicKey.Serialize()
}

func (p PublicKey) SerializeCompressed() []byte {
	return p.Serialize()
}

func (p PublicKey) SerializeUnCompressed() []byte {
	return p.Serialize()
}
