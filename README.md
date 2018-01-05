## README
Hyperledger debug tool; fabric node sdk; e2e-test

## Log:

+ got problem on example/explorer/; try add xmlhttprequest 
+ set host on 0.0.0.0, look into this [issue](https://github.com/webpack/webpack-dev-server/issues/151)
+ 

### TODO：

+ enable mysql
+ add fabric explorer


#### Requirements:

+ app are based on hyperledger-node-sdk/test
+ fabric-network: see ./artifacts
+ `apt install jq` or `brew install jq` (`testAppbyCurl.sh` needs)
+ needs to `npm rebuild` if run script at first time or catch node package error on `run.sh`
+ config.json used to set port and admin info and keyValueStore( default on "/tmp/fabric-client-kvs")
+ 
+ `bash run.sh` TO START THE NETWORK
+ `./cleanMaterial.sh` use to clean the user info without restart the network.
+ `./stopNetwork.sh` use to pull down the docker network and cleanMaterial
+ api tests via curl is `bash ./testAppbyCurl.sh`
+ network config: see ./app/network-config.json
+ use `docker exec -it cli bash` to enter the peer0 org1


#### Test api func:

1. enroll user
2. create channel
3. join channel
4. install
5. instantiate
6. invoke
7. query by differ functions.

#### Add new Oragnization dynamically process (Ideal):
1. prepare new crypto-config.yaml (not sure whether it needs existing config). 
2. use cryptogen to generate channel-crypto stuff for new org and anchor peer.
3. start the anchor peer docker container solely based on crypto stuff.
4. join the existing docker network
5. start the `configtxlator`. I found hyperledger gives an [official method](https://github.com/hyperledger/fabric/tree/master/examples/configtxupdate) to adding an organization
6. possibly can write a node script as a api 
7. node api shoud be equipped with some function:
    - decode exisiting channel to json
    - modified the global config file by adding the new org config info
    - updata the channel 
    - join peer and instantiate, invoke ... on new org