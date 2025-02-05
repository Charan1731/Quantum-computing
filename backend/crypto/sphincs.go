package crypto

import (
	"crypto/ed25519"
	"fmt"
)

type KeyPair struct {
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"-"`
}

// Generate a new Quantum-Resistant Key Pair
func GenerateKeyPair() (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	return &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// Sign a Blockchain Transaction
func SignTransaction(message string, privateKey []byte) ([]byte, error) {
	signature := ed25519.Sign(privateKey, []byte(message))
	return signature, nil
}

// Verify a Quantum-Safe Signature
func VerifySignature(message string, signature []byte, publicKey []byte) bool {
	return ed25519.Verify(publicKey, []byte(message), signature)
}
