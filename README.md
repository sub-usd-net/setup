
### Overview

This repo exposes a tool that wraps `subnet-cli` to make the UX surrounding deployment of a prototype stablecoin settlement subnet for Eden
Assumes that you have access to a node on the Avalanche Fuji test network that is a validator

1. Creates a subnet
2. Creates a blockchain with genesis data appropriate to an experimental stablecoin settlement subnet
    - The genesis address, specified in the config can mint new native coins
    - Only the genesis address can deploy new contracts
    - The genesis address can update the fee configuration without a hard fork
3. Adds the subnet to the subnet-validator set

### Walkthrough:

#### (This repo) Create subnet, genesis data, chain, and add subnet-validator

Please copy and edit the sample configuration `config.sample.yaml` to `config.yaml` and add in the details for your deployment

1. Create the subnet, blockchain, and add the subnet validator

```shell
$ go run main/main.go bootstrap
```

#### (On your node) Build, join, and expose

- Build, put in avalanchego/build/plugins dir, whitelist and restart node 

```shell
$ git clone https://github.com/sub-usd-net/vm ## assumes avalanchego should be a sibling directory to this one
$ ./scripts/build.sh <YOUR_VM_ID>
$ cp <YOUR_VM_ID> ../avalanchego/build/plugins/
```

- update /home/ubuntu/.avalanchego/configs/node.json

```shell
{"whitelisted-subnets": <YOUR_SUBNET_ID>}
```

- Expose RPC. You can use the sample nginx config in `samples/` (TBD)

#### Deploy contracts for users

- Deploy gnosis safe factory (TBD)

- Allow users to create gnosis safe's by whitelisting the factory (TBD)

- Deploy bridge (TBD)
