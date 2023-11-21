package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

type Wallet struct {
	privateKey     *ecdsa.PrivateKey
	publicKey      *ecdsa.PublicKey
	blockchainAddr string
}

func NewWallet() *Wallet {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	w := new(Wallet)
	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey

	// 2. Perform SHA-256 hashing on the public key (32 bytes)
	h2 := sha256.New()
	fmt.Println("h2", h2)
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	fmt.Println("h2", h2)
	fmt.Println("digest2", digest2)

	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes)
	h3 := ripemd160.New()
	fmt.Println("h3", h3)
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	fmt.Println("h3", h3)
	fmt.Println("digest3", digest3)

	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for main network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	fmt.Println("vd4", vd4)
	copy(vd4[1:], digest3[:])
	fmt.Println("vd4", vd4)

	// 5. Perform SHA-256 hash on the result of extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	fmt.Println("digest5", digest5)

	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	fmt.Println("digest6", digest6)

	// 7. Take first for bytes of the second SHA-256 hash for checksum
	chsum := digest6[:4]

	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	// 9. convert the result from a byte string to base58
	address := base58.Encode(dc8)
	w.blockchainAddr = address
	fmt.Println("block address", address)

	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddr
}

type Transaction struct {
	SenderPrivateKey   *ecdsa.PrivateKey `json:"senderPrivateKey"`
	SenderPublicKey    *ecdsa.PublicKey  `json:"senderPublicKey"`
	SenderBlockAddr    string            `json:"senderBlockAddr"`
	RecipientBlockAddr string            `json:"recipientBlockAddr"`
	Value              float32           `json:"value"`
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (t *Transaction) GenerateSignature() *Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.SenderPrivateKey, h[:])
	return &Signature{r, s}
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
