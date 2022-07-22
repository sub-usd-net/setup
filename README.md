
### Overview

This repo exposes a tool that wraps `avalanchego/wallet` to make the UX surrounding deployment of a prototype stablecoin settlement subnet for Eden
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

- Expose an RPC endpoint. Here is a sample nginx config:

```conf
server {
  listen 80;
  server_name <YOUR_DOMAIN>;

  location /rpc/ {
    proxy_pass http://localhost:9650/ext/bc/<YOUR_CHAIN_ID/rpc;
  }
  
  location = /ws/ {
    proxy_send_timeout 5m;
    proxy_read_timeout 5m;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_pass http://localhost:9650/ext/bc/<YOUR_CHAIN_ID>/ws;
  }
}
```


#### Deploy contracts for users

- Deploy gnosis safe factory. You can follow the instructions in the [safe repo](https://github.com/sub-usd-net/safe)

- Allow users to create gnosis safe's by whitelisting the factory. Obtain the factory address from the previous step.

```shell
# allow user's to create their own gnosis safe proxies
cast send --rpc-url ${RPC_URL} --private-key ${GENESIS_KEY} 0x0200000000000000000000000000000000000000 'setEnabled(address)' "$FACTORY"
# confirm
cast call --rpc-url ${RPC_URL} 0x0200000000000000000000000000000000000000 "readAllowList(address)(uint256)" "$FACTORY"
```

- Deploy C-Chain <--> Subnet bridge (TBD)
