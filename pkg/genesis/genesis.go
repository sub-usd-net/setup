package genesis

import (
	"encoding/json"
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sub-usd-net/setup/pkg/types"
)

var (
	microToWeiFactor = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)
)

func Generate(c *types.BootstrapConfig) ([]byte, error) {
	defaultAllowList := precompile.AllowListConfig{
		BlockTimestamp:  big.NewInt(0),
		AllowListAdmins: []common.Address{c.GenesisAddress},
	}

	cfg := params.SubnetEVMDefaultChainConfig
	cfg.ChainID = big.NewInt(int64(c.ChainId))

	cfg.ContractDeployerAllowListConfig = precompile.ContractDeployerAllowListConfig{AllowListConfig: defaultAllowList}
	cfg.ContractNativeMinterConfig = precompile.ContractNativeMinterConfig{AllowListConfig: defaultAllowList}
	cfg.FeeManagerConfig = precompile.FeeConfigManagerConfig{AllowListConfig: defaultAllowList}

	targetCostInWeiUSD := microToWei(c.TargetTransferCostInMicroUSD)
	minBaseFee := new(big.Int).Div(targetCostInWeiUSD, big.NewInt(21000))

	cfg.AllowFeeRecipients = true
	cfg.DefaultFeeRecipient = c.DefaultFeeReceiverAddress
	cfg.FeeConfig = commontype.FeeConfig{
		GasLimit:                 big.NewInt(20000000),
		TargetBlockRate:          2,
		TargetGas:                big.NewInt(100000000),
		MinBaseFee:               minBaseFee,
		BaseFeeChangeDenominator: big.NewInt(48),
		MinBlockGasCost:          big.NewInt(0),
		MaxBlockGasCost:          big.NewInt(4500000000),
		BlockGasCostStep:         big.NewInt(5000000),
	}

	alloc := make(map[common.Address]core.GenesisAccount)
	alloc[c.GenesisAddress] = core.GenesisAccount{
		Balance: microToWei(c.GenesisAllocationInMicroUSD),
	}

	genesis := core.Genesis{
		Config:     cfg,
		Alloc:      alloc,
		Difficulty: big.NewInt(0),
		GasLimit:   20000000,
	}

	return json.Marshal(genesis)
}

func microToWei(v uint64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(v)), microToWeiFactor)
}
