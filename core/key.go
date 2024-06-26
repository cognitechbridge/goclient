package core

import (
	"bytes"
	"crypto/rand"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/curve25519"
)

var (
	ErrInvalidPublicKey = errors.New("invalid public key")
)

type PublicKey struct {
	value []byte
}

func EmptyPublicKey() PublicKey {
	return PublicKey{}
}

// NewPublicKeyFromEncoded creates a PublicKey from an encoded base58 string.
func NewPublicKeyFromEncoded(encoded string) (PublicKey, error) {
	if len(encoded) != 44 {
		return EmptyPublicKey(), ErrInvalidPublicKey
	}
	return PublicKey{
		value: base58.Decode(encoded),
	}, nil
}

// NewPublicKeyFromBytes creates a PublicKey from a byte slice.
func NewPublicKeyFromBytes(bytes []byte) PublicKey {
	return PublicKey{
		value: bytes,
	}
}

// String returns the base58 encoded string representation of the PublicKey.
func (key PublicKey) String() string {
	return key.Encode()
}

// MarshalJSON returns the JSON encoding of the PublicKey.
func (key PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + key.Encode() + `"`), nil
}

// Encode returns the base58 encoded string representation of the PublicKey.
func (key PublicKey) Encode() string {
	return base58.Encode(key.value)
}

// Bytes returns the byte slice representation of the PublicKey.
func (key PublicKey) Bytes() []byte {
	return key.value
}

// Equals returns true if the PublicKey is equal to the other PublicKey.
func (key PublicKey) Equals(other PublicKey) bool {
	return bytes.Equal(key.value, other.value)
}

// PrivateKey represents a private key used for cryptographic operations.
type PrivateKey struct {
	value []byte
}

func EmptyPrivateKey() PrivateKey {
	return PrivateKey{}
}

// NewPrivateKeyFromEncoded creates a PrivateKey from an encoded base58 string.
func NewPrivateKeyFromEncoded(encoded string) (PrivateKey, error) {
	if len(encoded) != 44 {
		return EmptyPrivateKey(), ErrInvalidPublicKey
	}
	return PrivateKey{
		value: base58.Decode(encoded),
	}, nil
}

// NewPrivateKeyFromRand creates a PrivateKey from a random byte slice.
func NewPrivateKeyFromRand() (PrivateKey, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return EmptyPrivateKey(), err
	}
	return PrivateKey{
		value: key,
	}, nil
}

// NewPrivateKeyFromBytes creates a PrivateKey from a byte slice.
func NewPrivateKeyFromBytes(bytes []byte) PrivateKey {
	return PrivateKey{
		value: bytes,
	}
}

// Bytes returns the byte slice representation of the PrivateKey.
func (key PrivateKey) Bytes() []byte {
	return key.value
}

// Unsafe returns an UnsafePrivateKey for unsafe operations.
func (key PrivateKey) Unsafe() UnsafePrivateKey {
	return UnsafePrivateKey{key}
}

// ToPublicKey returns the PublicKey representation of the PrivateKey.
func (key PrivateKey) ToPublicKey() (PublicKey, error) {
	publicKey, err := curve25519.X25519(key.Bytes(), curve25519.Basepoint)
	if err != nil {
		return EmptyPublicKey(), err
	}
	return NewPublicKeyFromBytes(publicKey), nil
}

// MarshalJSON makes sure that the PrivateKey is not marshaled to JSON. This is to prevent leaking the private key.
// To marshal the private key, use the Unsafe method.
func (key PrivateKey) MarshalJSON() ([]byte, error) {
	return []byte(``), nil
}

// String makes sure that the PrivateKey is not converted to a string. This is to prevent leaking the private key.
func (key PrivateKey) String() string {
	return ""
}

// UnsafePrivateKey is used for unsafe operations on a private key.
type UnsafePrivateKey struct {
	PrivateKey
}

// Encode returns the base58 encoded string representation of the PrivateKey (Unsafe).
func (key UnsafePrivateKey) Encode() string {
	return base58.Encode(key.Bytes())
}

// String returns the base58 encoded string representation of the PrivateKey (Unsafe).
func (key UnsafePrivateKey) String() string {
	return key.Encode()
}

// MarshalJSON returns the JSON encoding of the PrivateKey (Unsafe).
func (key UnsafePrivateKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + key.Encode() + `"`), nil
}

func GenerateKey() (*KeyInfo, error) {
	key := NewKeyFromRand()
	keyId, err := NewUid()
	if err != nil {
		return nil, err
	}
	keyInfo := NewKeyInfo(keyId, key)
	return &keyInfo, nil
}
