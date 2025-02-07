package main

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	errInvalidPrivateKeySize = errors.New("invalid private key size")
	errInvalidPublicKeySize  = errors.New("invalid public key size")
)

// Store keys (for demo purposes)
var keyPair *KeyPair

type KeyPair struct {
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"-"`
}

func main() {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Generate Quantum-Safe Key
	router.GET("/generate-key", func(c *gin.Context) {
		var err error
		keyPair, err = GenerateKeyPair()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Key generation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"publicKey": hex.EncodeToString(keyPair.PublicKey),
		})
	})

	// Sign a Transaction
	router.POST("/sign", func(c *gin.Context) {
		var request struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if keyPair == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Generate keys first"})
			return
		}

		signature, err := SignTransaction(request.Message, keyPair.PrivateKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Signing failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"signature": hex.EncodeToString(signature),
		})
	})

	// Verify a Signature
	router.POST("/verify", func(c *gin.Context) {
		var request struct {
			Message   string `json:"message"`
			Signature string `json:"signature"`
			PublicKey string `json:"publicKey"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		signatureBytes, err := hex.DecodeString(request.Signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature format"})
			return
		}

		publicKeyBytes, err := hex.DecodeString(request.PublicKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid public key format"})
			return
		}

		valid, err := VerifySignature(request.Message, signatureBytes, publicKeyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"valid": valid})
	})

	router.Run(":8080") // Run on port 8080
}

// Generate a new Quantum-Resistant Key Pair using CRYSTALS-Dilithium
func GenerateKeyPair() (*KeyPair, error) {
	pk, sk, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	pkBytes := pk.Bytes()
	skBytes := sk.Bytes()

	return &KeyPair{
		PublicKey:  pkBytes,
		PrivateKey: skBytes,
	}, nil
}

// Sign a Blockchain Transaction using CRYSTALS-Dilithium

func SignTransaction(message string, privateKeyBytes []byte) ([]byte, error) {
	if len(privateKeyBytes) != mode3.PrivateKeySize {
		return nil, errInvalidPrivateKeySize
	}

	var sk mode3.PrivateKey
	var privateKeyArray [mode3.PrivateKeySize]byte
	copy(privateKeyArray[:], privateKeyBytes)

	// Corrected: Unpack without assignment
	sk.Unpack(&privateKeyArray)

	// Use crypto.Hash(0) as the third argument
	sig, err := sk.Sign(rand.Reader, []byte(message), crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// Verify a Quantum-Safe Signature using CRYSTALS-Dilithium
func VerifySignature(message string, signature []byte, publicKeyBytes []byte) (bool, error) {
	if len(publicKeyBytes) != mode3.PublicKeySize {
		return false, errInvalidPublicKeySize
	}

	// Convert public key bytes into a PublicKey struct
	var publicKey mode3.PublicKey
	var publicKeyArray [mode3.PublicKeySize]byte
	copy(publicKeyArray[:], publicKeyBytes)

	// Unpack the public key (No return value)
	publicKey.Unpack(&publicKeyArray)

	// Verify the signature
	return publicKey.Verify(signature, []byte(message)), nil
}
