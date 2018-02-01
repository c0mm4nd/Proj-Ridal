package bliss

import (
	"github.com/LoCCS/bliss/poly"
	Crypto "github.com/maoxs2/Proj-Ridal/crypto"
	"io"
)

type DSA interface {

	// ----------------------------------------------------------------------------
	// Private keys
	//
	// NewPrivateKey instantiates a new private key for the given data
	NewPrivateKey(s1, s2, a *poly.PolyArray) Crypto.PrivateKey

	// PrivKeyFromBytes calculates the public key from serialized bytes,
	// and returns both it and the private key.
	PrivKeyFromBytes(pk []byte) (Crypto.PrivateKey, Crypto.PublicKey)

	// PrivKeyBytesLen returns the length of a serialized private key.
	PrivKeyBytesLen() int

	// ----------------------------------------------------------------------------
	// Public keys
	//
	// NewPublicKey instantiates a new public key (point) for the given data.
	NewPublicKey(a *poly.PolyArray) Crypto.PublicKey

	// ParsePubKey parses a serialized public key for the given
	// curve and returns a public key.
	ParsePubKey(pubKeyStr []byte) (Crypto.PublicKey, error)

	// PubKeyBytesLen returns the length of the default serialization
	// method for a public key.
	PubKeyBytesLen() int

	// ----------------------------------------------------------------------------
	// Signatures
	//
	// NewSignature instantiates a new signature
	NewSignature(z1, z2 *poly.PolyArray, c []uint32) Crypto.Signature

	// ParseDERSignature parses a DER encoded signature .
	// If the method doesn't support DER signatures, it
	// just parses with the default method.
	ParseDERSignature(sigStr []byte) (Crypto.Signature, error)

	// ParseSignature a default encoded signature
	ParseSignature(sigStr []byte) (Crypto.Signature, error)

	// RecoverCompact recovers a public key from an encoded signature
	// and message, then verifies the signature against the public
	// key.
	RecoverCompact(signature, hash []byte) (Crypto.PublicKey, bool, error)

	// ----------------------------------------------------------------------------
	// Bliss
	//
	// GenerateKey generates a new private and public keypair from the
	// given reader.
	GenerateKey(rand io.Reader) (Crypto.PrivateKey, Crypto.PublicKey, error)

	// Sign produces an Bliss signature using a private key and a message.
	Sign(priv Crypto.PrivateKey, hash []byte) (Crypto.Signature, error)

	// Verify verifies an Bliss signature against a given message and
	// public key.
	Verify(pub Crypto.PublicKey, hash []byte, sig Crypto.Signature) bool
}

const (
	BSTypeBliss = 4

	BlissVersion = 1

	BlissPubKeyLen = 897

	BlissPrivKeyLen = 385
)

// Secp256k1 is the secp256k1 curve and ECDSA system used in Bitcoin.
var Bliss = newBlissDSA()
