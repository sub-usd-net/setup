package key

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ethereum/go-ethereum/common"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
)

type Manager struct {
	SubnetKeychain *secp256k1fx.Keychain
}

func NewKeyManager(subnetKeyPath string) (*Manager, error) {
	sk, err := decodeKeychainFromFile(subnetKeyPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		SubnetKeychain: sk,
	}, nil
}

func decodeKeychainFromFile(path string) (*secp256k1fx.Keychain, error) {
	key, err := decodeFromFile(path)
	if err != nil {
		return nil, err
	}

	factory := crypto.FactorySECP256K1R{}
	pkey, err := factory.ToPrivateKey(crypto2.FromECDSA(key))
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %s", err)
	}
	return secp256k1fx.NewKeychain(pkey.(*crypto.PrivateKeySECP256K1R)), nil
}

func decodeFromFile(path string) (*ecdsa.PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := strings.TrimSpace(string(content))
	return crypto2.ToECDSA(common.FromHex(s))
}
