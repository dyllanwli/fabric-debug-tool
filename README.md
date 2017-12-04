### README

Requirements:

+ fabric-network: see ./artifacts
+ `apt install jq` or `brew install jq`
+ `npm install`
+ config.json used to set port and admin info
+ `bash run.sh` TO START THE NETWORK
+ `./cleanMaterial.sh` use to clean the user info without restart the network.
+ `./stopNetwork.sh` use to pull down the docker network and cleanMaterial
+ api tests via curl is `bash ./test2.sh`
+ make sure port 4000 and 8088 is free
+ network config: see ./app/network-config.json

test api func:

1. enroll user
2. create channel
3. join channel
4. install
5. instantiate
6. invoke
7. query by differ functions.