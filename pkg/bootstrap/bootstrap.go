package bootstrap

import (
	"log"

	"github.com/sub-usd-net/setup/pkg/client"
	"github.com/sub-usd-net/setup/pkg/genesis"
	"github.com/sub-usd-net/setup/pkg/key"
	"github.com/sub-usd-net/setup/pkg/types"
)

type Bootstrapper struct {
	km     *key.Manager
	cfg    *types.BootstrapConfig
	client *client.Client
}

func NewBootstrapper(cfg *types.BootstrapConfig) (*Bootstrapper, error) {
	km, err := key.NewKeyManager(cfg.SubnetCreatorKeyPath)
	if err != nil {
		return nil, err
	}
	c, err := client.NewClient(cfg.FujiURI, km)
	if err != nil {
		return nil, err
	}

	return &Bootstrapper{
		km:     km,
		client: c,
		cfg:    cfg,
	}, nil
}

func (b *Bootstrapper) Bootstrap() error {
	bs, err := genesis.Generate(b.cfg)
	if err != nil {
		return err
	}

	res, err := b.client.CreateSubnetAndBlockchain(b.cfg, bs)
	if err != nil {
		return err
	}
	log.Printf("== Bootstrapped ==\nSubnet: %s\nChain: %s\nVM: %s", res.Subnet, res.Chain, res.VmID)

	if err := b.client.AddSubnetToValidator(res.Subnet, b.cfg); err != nil {
		return err
	}
	log.Println("Now you can build the vm, whitelist the subnet, and restart your node")
	log.Printf("on node ./recreate/recreate.sh %s %s %s\n", res.Subnet.String(), res.VmID.String(), res.Chain.String())
	return nil
}
