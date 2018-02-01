package bliss

import (
	"github.com/LoCCS/bliss"
	Crypto "github.com/maoxs2/Proj-Ridal/crypto"
)

type PrivateKey struct {
	Crypto.PrivateKeyAdapter
	bliss.PrivateKey
}

// Public returns the PublicKey corresponding to this private key.
func (p PrivateKey) PublicKey() Crypto.PublicKey {
	blissPkp := p.PrivateKey.PublicKey()
	pk := &PublicKey{
		PublicKey: *blissPkp,
	}
	return pk
}

// GetType satisfies the bliss PrivateKey interface.
func (p PrivateKey) GetType() int {
	return pqcTypeBliss
}

func (p PrivateKey) Serialize() []byte {
	return p.PrivateKey.Serialize()
}
