package bliss

import (
	"github.com/LoCCS/bliss"
	Crypto "github.com/maoxs2/Proj-Ridal/crypto"
)

type Signature struct {
	Crypto.SignatureAdapter
	bliss.Signature
}

func (s Signature) GetType() int {
	return pqcTypeBliss
}

func (s Signature) Serialize() []byte {
	return s.Signature.Serialize()
}
