package crypto

import (
	"crypto/ecdsa"
	"math/big"
	// "github.com/HcashOrg/hcashd/chaincfg/chainec"
)

type PublicKey interface {
	// Serialize is the default serialization method.
	Serialize() []byte

	// SerializeUncompressed serializes to the uncompressed format (if
	// available).
	SerializeUncompressed() []byte

	// SerializeCompressed serializes to the compressed format (if
	// available).
	SerializeCompressed() []byte

	// SerializeHybrid serializes to the hybrid format (if
	// available).
	SerializeHybrid() []byte

	// ToECDSA converts the public key to an ECDSA public key.
	ToECDSA() *ecdsa.PublicKey

	// GetCurve returns the current curve as an interface.
	GetCurve() interface{}

	// GetX returns the point's X value.
	GetX() *big.Int

	// GetY returns the point's Y value.
	GetY() *big.Int

	// GetType returns the ECDSA type of this key.
	GetType() int
}

type PrivateKey interface {
	// Serialize serializes the 32-byte private key scalar to a
	// byte slice.
	Serialize() []byte

	// SerializeSecret serializes the secret to the default serialization
	// format. Used for Ed25519.
	SerializeSecret() []byte

	// Public returns the (X,Y) coordinates of the point produced
	// by scalar multiplication of the scalar by the base point,
	// AKA the public key.
	Public() (*big.Int, *big.Int)

	// GetD returns the value of the private scalar.
	GetD() *big.Int

	// GetType returns the ECDSA type of this key.
	GetType() int

	PublicKey() PublicKey
}

type Signature interface {
	// Serialize serializes the signature to the default serialization
	// format.
	Serialize() []byte

	// GetR gets the R value of the signature.
	GetR() *big.Int

	// GetS gets the S value of the signature.
	GetS() *big.Int

	// GetType returns the ECDSA type of this key.
	GetType() int
}

type PublicKeyAdapter struct {
}

type PrivateKeyAdapter struct {
}

type SignatureAdapter struct {
}

//PrivateKeyAdapter
func (pa PrivateKeyAdapter) Serialize() []byte {
	return nil
}

func (pa PrivateKeyAdapter) SerializeSecret() []byte {
	return nil
}

func (pa PrivateKeyAdapter) Public() (*big.Int, *big.Int) {
	return nil, nil
}

func (pa PrivateKeyAdapter) PublicKey() PublicKey {
	return nil
}

func (pa PrivateKeyAdapter) GetD() *big.Int {
	return nil
}

func (pa PrivateKeyAdapter) GetType() int {
	return 0
}

//PublicKeyAdapter
func (pa PublicKeyAdapter) Serialize() []byte {
	return nil
}

func (pa PublicKeyAdapter) SerializeUncompressed() []byte {
	return nil
}

func (pa PublicKeyAdapter) SerializeCompressed() []byte {
	return nil
}

func (pa PublicKeyAdapter) SerializeHybrid() []byte {
	return nil
}

func (pa PublicKeyAdapter) ToECDSA() *ecdsa.PublicKey {
	return nil
}

func (pa PublicKeyAdapter) GetCurve() interface{} {
	return nil
}

func (pa PublicKeyAdapter) GetX() *big.Int {
	return nil
}

func (pa PublicKeyAdapter) GetY() *big.Int {
	return nil
}

func (pa PublicKeyAdapter) GetType() int {
	return 0
}

//SignatureAdapter
func (s SignatureAdapter) Serialize() []byte {
	return nil
}

func (s SignatureAdapter) GetR() *big.Int {
	return nil
}

func (s SignatureAdapter) GetS() *big.Int {
	return nil
}

func GetType() int {
	return 0
}
