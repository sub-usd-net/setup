package types

import (
	"fmt"
	"io"
	"os"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v2"
)

type NodeID ids.NodeID

func (n *NodeID) UnmarshalText(b []byte) error {
	id, err := ids.ShortFromPrefixedString(string(b), ids.NodeIDPrefix)
	if err != nil {
		return fmt.Errorf("error unmarshaling node ID: %s", err)
	}
	*n = NodeID(id)
	return nil
}

type BootstrapConfig struct {
	FujiURI               string `yaml:"fujiURI"`
	SubnetCreatorKeyPath  string `yaml:"subnetCreatorKeyPath"`
	VmName                string `yaml:"vmName"`
	ChainName             string `yaml:"chainName"`
	SubnetValidatorWeight uint64 `yaml:"subnetValidatorWeight"`

	GenesisAddress               common.Address `yaml:"genesisAddress"`
	DefaultFeeReceiverAddress    common.Address `yaml:"defaultFeeReceiverAddress"`
	GenesisAllocationInMicroUSD  uint64         `yaml:"genesisAllocationInMicroUSD"`
	TargetTransferCostInMicroUSD uint64         `yaml:"targetTransferCostInMicroUSD"`
	ChainId                      uint64         `yaml:"chainId"`

	ValidatorNodeID_ NodeID `yaml:"validatorNodeID"`
	ValidatorNodeID  ids.NodeID
}

func NewBootstrapConfigFromPath(path string) (*BootstrapConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var cfg BootstrapConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	cfg.ValidatorNodeID = ids.NodeID(cfg.ValidatorNodeID_)

	if err := ensureNotEmpty(cfg.VmName, "vmName"); err != nil {
		return nil, err
	}
	if err := ensureNotEmpty(cfg.ChainName, "chainName"); err != nil {
		return nil, err
	}

	if err := ensurePath(cfg.SubnetCreatorKeyPath, "subnetCreatorKeyPath"); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ensureNotEmpty(s, errPrefix string) error {
	if s == "" {
		return fmt.Errorf("%s cannot be empty", errPrefix)
	}
	return nil
}

func ensurePath(s, errPrefix string) error {
	if s == "" {
		return fmt.Errorf("%s: must be specified", errPrefix)
	}
	if _, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("%s = (%s) does not exist", errPrefix, s)
	}
	return nil
}
