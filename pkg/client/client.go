package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/avalanchego/vms/platformvm/validator"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/sub-usd-net/setup/pkg/key"
	"github.com/sub-usd-net/setup/pkg/types"
)

type SubnetCreationResult struct {
	Subnet ids.ID
	Chain  ids.ID
	VmID   ids.ID
}

type Client struct {
	avaxAssetID ids.ID
	pvm         platformvm.Client
	wallet      primary.Wallet
	uri         string
	km          *key.Manager
}

func NewClient(uri string, km *key.Manager) (*Client, error) {
	wallet, err := primary.NewWalletFromURI(context.Background(), uri, km.SubnetKeychain)
	if err != nil {
		return nil, err
	}

	return &Client{
		avaxAssetID: wallet.P().AVAXAssetID(),
		pvm:         platformvm.NewClient(uri),
		wallet:      wallet,
		uri:         uri,
		km:          km,
	}, nil
}

func (c *Client) CreateSubnetAndBlockchain(cfg *types.BootstrapConfig, genesis []byte) (*SubnetCreationResult, error) {
	fees := c.wallet.P().CreateSubnetTxFee() + c.wallet.P().CreateBlockchainTxFee()
	b, err := c.wallet.P().Builder().GetBalance()
	if err != nil {
		return nil, fmt.Errorf("error creating subnet/chain: cannot get user balance: %s", err)
	}
	if b[c.avaxAssetID] < fees {
		return nil, fmt.Errorf("error creating subnet/chain: insufficienet balance. Need = %d Have = %d", fees, b[c.avaxAssetID])
	}

	log.Println("Creating Subnet")
	sid, err := c.wallet.P().IssueCreateSubnetTx(c.getOwner())
	if err != nil {
		return nil, fmt.Errorf("error creating subnet: %s", err)
	}
	log.Println("...created")

	c.waitFor(sid)
	vid := nameToVmID(cfg.VmName)
	log.Println("Creating Chain")
	cid, err := c.wallet.P().IssueCreateChainTx(sid, genesis, vid, nil, cfg.ChainName)
	if err != nil {
		return nil, fmt.Errorf("error creating chain: %s", err)
	}
	log.Println("...created")

	return &SubnetCreationResult{
		Subnet: sid,
		Chain:  cid,
		VmID:   vid,
	}, nil
}

func (c *Client) AddSubnetToValidator(subnet ids.ID, cfg *types.BootstrapConfig) error {
	vs, err := c.pvm.GetCurrentValidators(context.Background(), constants.PrimaryNetworkID, []ids.NodeID{cfg.ValidatorNodeID})
	if err != nil {
		return err
	}
	if len(vs) != 1 {
		return fmt.Errorf("error checking validator end time. Expected 1 validator got %d", len(vs))
	}

	log.Println("Adding SubnetValidator")
	_, err = c.wallet.P().IssueAddSubnetValidatorTx(&validator.SubnetValidator{
		Subnet: subnet,
		Validator: validator.Validator{
			NodeID: cfg.ValidatorNodeID,
			Wght:   cfg.SubnetValidatorWeight,
			Start:  uint64(time.Now().Add(15 * time.Second).Unix()),
			End:    vs[0].EndTime,
		},
	})
	if err != nil {
		return err
	}
	log.Println("...added")

	return err
}

func (c *Client) waitFor(tx ids.ID) {
	for {
		_, err := c.pvm.GetTx(context.Background(), tx)
		if err == nil {
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) getOwner() *secp256k1fx.OutputOwners {
	return &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     c.km.SubnetKeychain.Addrs.List(),
	}
}

func nameToVmID(name string) ids.ID {
	var vmID [32]byte
	prefix := []byte(fmt.Sprintf("vm-%s", name))
	for i, b := range prefix {
		if i < 32 {
			vmID[i] = b
		}
	}
	return vmID
}
