var path = require('path');
var fs = require('fs');
var util = require('util');
var hfc = require('fabric-client');
var Peer = require('fabric-client/lib/Peer.js');
var EventHub = require('fabric-client/lib/EventHub.js');
var helper = require('./helper.js');
var logger = helper.getLogger('instantiate-chaincode');
hfc.addConfigFile(path.join(__dirname, 'network-config.json'));
var ORGS = hfc.getConfigSetting('network-config');
var tx_id = null;
var eh = null;

var instantiateChaincode = function (channelName, chaincodeName, chaincodeVersion, functionName, args, username, org) {
    logger.debug('\n============ Instantiate chaincode on organization ' + org +
        ' ============\n');

    var channel = helper.getChannelForOrg(org, channelName);
    var client = helper.getClientForOrg(org);

    return helper.getOrgAdmin(org).then((user) => {
        // read the config block from the orderer for the channel
        // and initialize the verify MSPs based on the participating
        // organizations
        return channel.initialize();
    }, (err) => {
        logger.error('Failed to enroll user \'' + username + '\'. ' + err);
        throw new Error('Failed to enroll user \'' + username + '\'. ' + err);
    }).then((success) => {
        tx_id = client.newTransactionID();
        // ep = {
        //     identities: [],
        //     policy: []
        // }
        // endorsementPolicy = helper.getEndorsementpolicy(ep, org)
        // 取消设置背书策略

        // send proposal to endorser
        var request = {
            chaincodeId: chaincodeName,
            chaincodeVersion: chaincodeVersion,
            args: args,
            txId: tx_id
        };
        // if (endorsementPolicy)
        //     request['endorsement-policy'] = endorsementPolicy
        // 设置背书策略

        if (functionName)
            request.fcn = functionName;
        return channel.sendInstantiateProposal(request, 60000 * 60);
    }, (err) => {
        logger.error('Failed to initialize the channel');
        throw new Error('Failed to initialize the channel');
    }).then((results) => {
        // function sleep(ms, results) {
        //     logger.debug('Sleeping ' + ms + ' to wait for launching the chaincode.')
        //     return new Promise(resolve => setTimeout(resolve, ms)).then((results) => {
        //         logger.debug("end sleep")

        //     });
        //     // out sleep function
        // }
        // sleep(60 * 1000, results)
        // 设置sleeptime 用来等待chaincode launch 回导致其他时间更长
        function mySetTimeout(ms) {
            var currentTime = new Date().getTime();
            while (new Date().getTime() < currentTime + ms);
        }
        // mySetTimeout(60*1000);        
        var proposalResponses = results[0];
        var proposal = results[1];
        var header = results[2];
        var all_good = true;
        for (var i in proposalResponses) {
            let one_good = false;
            if (proposalResponses && proposalResponses[i].response &&
                proposalResponses[i].response.status === 200) {
                one_good = true;
                logger.info('instantiate proposal was good');
            } else {
                logger.error('instantiate proposal was bad');
            }
            all_good = all_good & one_good;
        }
        if (all_good) {
            logger.info(util.format(
                'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s',
                proposalResponses[0].response.status, proposalResponses[0].response.message,
                proposalResponses[0].response.payload, proposalResponses[0].endorsement
                .signature));
            var request = {
                proposalResponses: proposalResponses,
                proposal: proposal,
                header: header
            };
            // set the transaction listener and set a timeout of 30sec
            // if the transaction did not get committed within the timeout period,
            // fail the test
            var deployId = tx_id.getTransactionID();

            eh = client.newEventHub();
            let data = fs.readFileSync(path.join(__dirname, ORGS[org].peers['peer0'][
                'tls_cacerts'
            ]));
            eh.setPeerAddr(ORGS[org].peers['peer0']['events'], {
                pem: Buffer.from(data).toString(),
                'ssl-target-name-override': ORGS[org].peers['peer0']['server-hostname']
            });
            eh.connect();
            // 设置的默认peer0

            let txPromise = new Promise((resolve, reject) => {
                let handle = setTimeout(() => {
                    eh.disconnect();
                    reject();
                }, 60000 * 60);

                eh.registerTxEvent(deployId, (tx, code) => {
                    logger.info(
                        'The chaincode instantiate transaction has been committed on peer ' +
                        eh._ep._endpoint.addr);
                    clearTimeout(handle);
                    eh.unregisterTxEvent(deployId);
                    eh.disconnect();

                    if (code !== 'VALID') {
                        logger.error('The chaincode instantiate transaction was invalid, code = ' + code);
                        reject();
                    } else {
                        logger.info('The chaincode instantiate transaction was valid.');
                        resolve();
                    }
                });
            });

            var sendPromise = channel.sendTransaction(request);
            return Promise.all([sendPromise].concat([txPromise])).then((results) => {
                logger.debug('Event promise all complete and testing complete');
                return results[0]; // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call
            }).catch((err) => {
                logger.error(
                    util.format('Failed to send instantiate transaction and get notifications within the timeout period. %s', err)
                );
                return 'Failed to send instantiate transaction and get notifications within the timeout period.';
            });
        } else {
            logger.error(
                'Failed to send instantiate Proposal or receive valid response. Response null or status is not 200. exiting...'
            );
            return 'Failed to send instantiate Proposal or receive valid response. Response null or status is not 200. exiting...';
        }
    }, (err) => {
        logger.error('Failed to send instantiate proposal due to error: ' + err.stack ?
            err.stack : err);
        return 'Failed to send instantiate proposal due to error: ' + err.stack ?
            err.stack : err;
    }).then((response) => {
        if (response.status === 'SUCCESS' || response.status === '200' ) {
            logger.info('Successfully sent transaction to the orderer.');
            logger.info('Chaincode Instantiation is SUCCESS.');
            return 'Chaincode Instantiation is SUCCESS';
        } else {
            logger.error('Failed to order the transaction. Error code: ' + err.stack ?
            err.stack : err);
            return 'Failed to order the transaction. Error code: ' + err.stack ?
            err.stack : err;
        }
    }, (err) => {
        logger.error('Failed to send instantiate due to error: ' + err.stack ? err
            .stack : err);
        return 'Failed to send instantiate due to error: ' + err.stack ? err.stack :
            err;
    });
};
exports.instantiateChaincode = instantiateChaincode;