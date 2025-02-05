package main

import (
	"encoding/hex"
	"net/http"

	"quantum-blockchain/crypto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Store keys (for demo purposes)
var keyPair *crypto.KeyPair

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
		keyPair, err = crypto.GenerateKeyPair()
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

		signature, err := crypto.SignTransaction(request.Message, keyPair.PrivateKey)
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

		valid := crypto.VerifySignature(request.Message, signatureBytes, publicKeyBytes)

		c.JSON(http.StatusOK, gin.H{"valid": valid})
	})

	router.Run(":8080") // Run on port 8080
}
