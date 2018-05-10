## README
debug-tool is use to help the developer to develop chaincode and test the fabric network.


### Requirements:

+ npm >=3.10.10
+ node >=6.10.0 <=6.12.0
+ docker-compose >= 1.16.1
+ docker >= 1.12.6

### Setup Config File:

Notice that config file should match the docker-compose file perfectly; be careful if you want to modify the network-config file; 

1. `./config.json` is used to setup chaincode path, app port, timeout and admin infomation.
    + if network has multiple channel, you should add all channel name into *channelsList* at `./config.json`
2. `./app/network-config.json` is used to setup fabric network
    + `./app/network-config.allorderer.json` added all orderer node, which include the unused orderer.
3. a simple fabric network at `./artifacts.old` is following with [fabric network demo](http://hyperledger-fabric.readthedocs.io/en/release/build_network.html) and the config file is `./config.old.json` and `./app/network-config.old.json`. You can use it by running `bash useold.sh` and roll back by `bash usenew.sh`.


### Setup the Environments:

Start the fabric network, this tool has been test with fabric 1.03 and 1.1-preview

1. install package from 

    ```
    npm install
    ```
2. bring up the docker-compose:

    ```
    cd fabric-docker-compose-svt
    # clean the docker environment by running bash cleanup.sh
    bash cleanup.sh 
    # will clean all docker container and chaincode image which name with dev-*
    bash bringup-wsh.sh
    # will bring up the docker-compose-e2e-couchdb
    ```

+ *if you use simple network which include 2 org and each org has 2 peers with `useold.sh`, you should:*

    ```
    bash useold.sh
    cd artifacts
    docker-compose up -d
    # if done all the test, you can use ./stopNetwork.sh and bash usenew.sh to recover the environment.
    ```

+ if you want to build your own network, use the `generateArtifacts.sh` to generate channel-artifacts and cypto-config, then all the config file must copy to `./artifacts` and replace the old file.
+ The network-config.json.old is using ca.crt, while the network-conifig for golden ticket is using msp pem file. If you want to use your own network, make sure write the right file path.

### Tool Usage

1. use `node app` to start the tool. 
    + if you want to run this tool at background, try `screen node app` or just use the log output format like `node app > stdout.log 2> stderr.log &`
2. open your broswer and enter the address
3. enroll user and select the enrolled user
4. create channel: enter the channel name and channel artifacts file
5. select channel (used as global)
6. join the peer
7. install channel: enter the chincode path 
    + Notice: chaincode root path is defined to `./artifacts/src`, you just need to enter the chaincode directory name.
    + if you want to change the root path, you can change the `CC_GOPATH` at `./config.json`
8. instantiate: this step will cost a long time and possible catch some error
9. invoke and query
10. both invoke and query support two data format: json list(defined with `[{}]`) and string list (split with `,`):

    arg type 1:

    ```
    [{
        "docType": "type1",
        "id": 1000,
        "createTime": 2017,
        "updateTime": 2018,
        "createUser": "diya",
        "updateUser": "tom",
        "gldId": "00002"
    }]

    ```

    arg type 2:

    ```
    00001,00002
    or
    00001
    ```

11. query function:
    + query by args: needs enter channel name, chaincodes, peer, args, and query function
    + query by blockID: needs channel name, blockid; returns the block timestamp and written infomation
    + query by Transaction id: needs channel name, peer； returns Transaction timestamp and written info
    + query Chaininfo: needs channel name and peer; return chain info
    + query installType: needs install Type(instantiated or installed);
    + query channel: return channel infomation.

### Troubleshooting

1. if catches `node-gyp error` install error when `npm install`, you needs to 

    ```
    npm rebuild
    ```

2. when running instantiate, if your cpu or memory not large enough, possibly catch

    ```
    client-utils.js - REQUEST TIMEOUT
    
    peer log - chincode (xxxx) is being launched.
    ```
    check log at chaincode container will get
    
    ```
    can not find local peer
    ```
    + at most of time, the tool will callback the instantiate and request again. Log at broswer may return Error, but actually the instantiate is still running. so you just wait for the chaincode container up
    + sometimes instantiate may not get success status from fabric-sdk sendChaincodeProposal at first, due to network too large; but finally the chaincode will instantiate successfully in the end until the callback function running again.
    + if you check the `docker ps -a` and find out the chaincode container Exits. Try cleanup the environment and restart the netowrk.
3. when test bigtree network, possibly will catch `enroll failure` if you are not enroll admin user, but it will still callback to enroll admin till enroll success; So don't worry it will finally return the expect result and will not affect the tool usage
    + when test other network, it may not have this error
4. When catched 

    ````
    sendPeersProposal - Promise is rejected: Error: 2 UNKNOWN: could not find chaincode with name 'xxxx' - make sure the chaincode xxxx has been successfully instantiated and try again
    
    instantiate - SERVICE UNABILABLE
    ````
    
    possibly because the fabric network cause CPU stucking.
    ```
    Message from syslogd@user ... kernel:NMI watchdog: BUG: soft lockup - CPU#1 stuck for xxxx
    ```
    Recommend to stop the network and clean the environment, till you CPU is running stable.