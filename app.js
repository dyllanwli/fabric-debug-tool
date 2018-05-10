var fs = require("fs");
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var express = require('express');
var session = require('express-session');
var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');
var http = require('http');
var util = require('util');
var app = express();
var expressJWT = require('express-jwt');
var jwt = require('jsonwebtoken');
var bearerToken = require('express-bearer-token');
var cors = require('cors');

var config_json = require('./config.json')
var config = require('./config.js');
var hfc = require('fabric-client');
var network_config = require('./app/network-config.json')

var helper = require('./app/helper.js');
var channels = require('./app/create-channel.js');
var join = require('./app/join-channel.js');
var install = require('./app/install-chaincode.js');
var instantiate = require('./app/instantiate-chaincode.js');
var invoke = require('./app/invoke-transaction.js');
var query = require('./app/query.js');
var host = process.env.HOST || hfc.getConfigSetting('host');
var port = process.env.PORT || hfc.getConfigSetting('port');
///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// SET configuration ////////////////////////////
///////////////////////////////////////////////////////////////////////////////
app.use(express.static('public'));
app.options('*', cors());
app.use(cors());
//parsing of application/json type post data
app.use(bodyParser.json());
//parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
	extended: false
}));

// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
	secret: 'thisismysecret'
}).unless({
	path: ['/getchannelslist', '/users', '/', ]
}));
app.use(bearerToken());
// use bearer token
app.use(function (req, res, next) {
	// the code below will go to servers start
	// logger.info('jwt verify REQUEST Time:', Date.now());
	if (req.originalUrl.indexOf('/users') >= 0) {
		return next();
	}
	if (req.originalUrl.indexOf('/getchannelslist') >= 0) {
		return next();
	}
	if (req.originalUrl.length <= 2) {
		return next();
	}

	var token = req.token;
	jwt.verify(token, app.get('secret'), function (err, decoded) {
		if (err) {
			res.send({
				success: false,
				message: 'Failed to authenticate token. Make sure to include the ' +
					'token returned from /users call in the authorization header ' +
					' as a Bearer token'
			});
			return;
		} else {
			// add the decoded user name and org name to the request object
			// for the downstream code to use
			req.username = decoded.username;
			req.orgname = decoded.orgName;
			logger.info("successed")
			logger.info(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.username, decoded.orgName));
			return next();
		}
	});
});



///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// START server /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

var server = http.createServer(app).listen(port, function () {
	logger.info('****************** SERVER STARTED ************************');
});
logger.info('************** http://' + host + ':' + port + '******************');
server.timeout = 240000;


/////////////////////////////above is main page////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////





//////////////////////////////////////////////////////////////////////////////////
///////////////////////////////below is error debug///////////////////////////////

function getErrorMessage(field) {
	var response = {
		success: false,
		message: field + ' field is missing or Invalid in the request'
	};
	return response;
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////// rest end point start here ///////////////////////////
///////////////////////////////////////////////////////////////////////////////
// Register and enroll user
app.post('/users', function (req, res) {
	req.headers['content-type'] = "application/x-www-form-urlencoded";
	var username = req.body.username;
	var orgName = req.body.orgName;
	logger.debug('End point : /users');
	logger.debug('User name : ' + username);
	logger.debug('Org name  : ' + orgName);
	if (!username) {
		res.json(getErrorMessage('\'username\''));
		return;
	}
	if (!orgName) {
		res.json(getErrorMessage('\'orgName\''));
		return;
	}
	var token = jwt.sign({
		exp: Math.floor(Date.now() / 1000) + parseInt(hfc.getConfigSetting('jwt_expiretime')),
		username: username,
		orgName: orgName
	}, app.get('secret'));
	helper.getRegisteredUsers(username, orgName, true).then(function (response) {
		if (response && typeof response !== 'string') {
			response.token = token;
			res.json(response);
		} else {
			res.json({
				success: false,
				message: response
			});
		}
	});
});

app.get('/getchannelslist', function (req, res) {
	res.json({
		channelsList: config_json.channelsList,
		channelsConfigPath: config_json.channelsConfigPath,
	})
});


// Create Channel
app.post('/channels', function (req, res) {

	// get data
	logger.info('====================CREATE CHANNEL====================');
	logger.debug('End point : /channels');
	var channelName = req.body.channelName;
	var channelConfigPath = config_json.channelsConfigPath[channelName];

	logger.debug('Channel name : ' + channelName);
	logger.debug('channelConfigPath : ' + channelConfigPath);
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!channelConfigPath) {
		res.json(getErrorMessage('\'channelConfigPath\''));
		return;
	}
	channels.createChannel(channelName, channelConfigPath, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		}, (err) => {
			res.send(err)
		});
});

// Join Channel 
app.post('/channels/:channelName/peers', function (req, res) {
	logger.info('==================== JOIN CHANNEL====================');
	var channelName = req.params.channelName;
	var peers = req.body.peers;
	logger.debug('channelName : ' + channelName);
	logger.debug('peers : ' + peers);
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
		return;
	}

	if (peers = "joinall") {
		peers = []
		for (peer in network_config['network-config'][req.orgname]["peers"]) {
			peers.push(peer)
		}
	}

	join.joinChannel(channelName, peers, req.username, req.orgname)
		.then(function (message) {
			response_message = {
				message: message,
				peers: {
					org: req.orgname,
					peers: peers
				}
			}
			res.send(JSON.stringify(response_message));
		});
});
// 
// Install chaincode on target peers 
app.post('/chaincodes', function (req, res) {

	logger.debug('==================== INSTALL CHAINCODE ==================');
	var peers = req.body.peers;
	var channelName = req.body.channelName;
	var chaincodeName = req.body.chaincodeName;
	var chaincodePath = req.body.chaincodePath;
	var chaincodeVersion = req.body.chaincodeVersion;
	logger.debug('channelName : ' + channelName);
	logger.debug('peers : ' + peers);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('chaincodePath  : ' + chaincodePath);
	logger.debug('chaincodeVersion  : ' + chaincodeVersion);
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!chaincodePath) {
		res.json(getErrorMessage('\'chaincodePath\''));
		return;
	}
	if (!chaincodeVersion) {
		res.json(getErrorMessage('\'chaincodeVersion\''));
		return;
	}
	var func_list = []
	var path = "./artifacts/src/" + chaincodePath

	function getFiles(dir, files_) {
		files_ = files_ || [];
		var files = fs.readdirSync(dir);
		for (var i in files) {
			if (files[i].indexOf("main") > -1) {
				var name = dir + '/' + files[i];
				var lines = fs.readFileSync(name).toString().split("\n")
				for (var line of lines) {
					if (line.indexOf('case') > -1) {
						fcn = line.slice(7, -2)
						files_.push(fcn)
					}
				}
			}
		}
		return files_;
	}
	// func_list = getFiles(path)

	install.installChaincode(peers, channelName, chaincodeName, chaincodePath, chaincodeVersion, req.username, req.orgname)
		.then(function (message) {
			var mes = {}
			mes.message = message
			// mes.func_list = func_list
			res.send(mes.message);
		});
});
// 
// 
// Instantiate chaincode on target peers
app.post('/channels/:channelName/chaincodes', function (req, res) {

	logger.debug('==================== INSTANTIATE CHAINCODE==================');
	var chaincodeName = req.body.chaincodeName;
	var chaincodeVersion = req.body.chaincodeVersion;
	var channelName = req.params.channelName;
	var fcn = req.body.fcn;
	var args = req.body.args;
	logger.debug('channelName  : ' + channelName);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('chaincodeVersion  : ' + chaincodeVersion);
	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!chaincodeVersion) {
		res.json(getErrorMessage('\'chaincodeVersion\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}
	instantiate.instantiateChaincode(channelName, chaincodeName, chaincodeVersion, fcn, args, req.username, req.orgname)
		.then(function (message) {
			if (message.indexOf("Failed") > -1) {
				message += "\ninstantiate proposal not all good\nsend Promise Twice...\n"
				logger.debug(message)
				res.send(message)
			} else {
				res.send(message);
			}
		});
});
// 
// 
// Invoke transaction on chaincode on target peers
app.post('/channels/:channelName/chaincodes/:chaincodeName', function (req, res) {

	logger.debug('==================== INVOKE ON CHAINCODE ==================');
	var peers = req.body.peers;
	var chaincodeName = req.params.chaincodeName;
	var channelName = req.params.channelName;
	var fcn = req.body.fcn;
	var args = req.body.args;
	logger.debug('channelName  : ' + channelName);
	logger.debug('peers  : ' + peers);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);
	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!fcn) {
		res.json(getErrorMessage('\'fcn\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}

	invoke.invokeChaincode(peers, channelName, chaincodeName, fcn, args, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		});
});
///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////
//////////////////////////below is query///////////////////////////////////////////////
// Got 
// Query to fetch channels
// Query on chaincode on target peers
// Query Get Block by BlockNumber
// Query Get Transaction by Transaction ID
// Query Get Block by Hash
// Query for Channel Information
// Query to fetch all Installed/instantiated chaincodes
// 
// 
// 
// Query on chaincode on target peers
app.get('/channels/:channelName/chaincodes/:chaincodeName', function (req, res) {
	logger.debug('==================== QUERY BY CHAINCODE ==================');
	var channelName = req.params.channelName;
	var chaincodeName = req.params.chaincodeName;
	let args = req.query.args;
	let fcn = req.query.fcn;
	let peer = req.query.peer;

	logger.debug('channelName : ' + channelName);
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('fcn : ' + fcn);
	logger.debug('args : ' + args);

	if (!chaincodeName) {
		res.json(getErrorMessage('\'chaincodeName\''));
		return;
	}
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!fcn) {
		res.json(getErrorMessage('\'fcn\''));
		return;
	}
	if (!args) {
		res.json(getErrorMessage('\'args\''));
		return;
	}
	try {
		args = JSON.parse(args)
		for (let i = 0; i < args.length; i++) {
			args[i] = JSON.stringify(args[i])
		}
	} catch (e) {
		if (args.indexOf(",") > -1) {
			args = args.split(",")
		} else {
			args = [args]
		}
	}
	query.queryChaincode(peer, channelName, chaincodeName, fcn, args, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		});
});
// 
// 
// 
// 
//  Query Get Block by BlockNumber
app.get('/channels/:channelName/blocks/:blockId', function (req, res) {
	logger.debug('==================== GET BLOCK BY NUMBER ==================');
	let blockId = req.params.blockId;
	let peer = req.query.peer;
	let channelName = req.params.channelName;
	logger.debug('channelName : ' + req.params.channelName);
	logger.debug('BlockID : ' + blockId);
	logger.debug('Peer : ' + peer);
	if (!blockId) {
		res.json(getErrorMessage('\'blockId\''));
		return;
	}

	query.getBlockByNumber(peer, channelName, blockId, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		});
});
// Query Get Transaction by Transaction ID
app.get('/channels/:channelName/transactions/:trxnId', function (req, res) {
	logger.debug(
		'================ GET TRANSACTION BY TRANSACTION_ID ======================'
	);
	logger.debug('channelName : ' + req.params.channelName);
	let trxnId = req.params.trxnId;
	let channelName = req.params.channelName;
	let peer = req.query.peer;
	if (!trxnId) {
		res.json(getErrorMessage('\'trxnId\''));
		return;
	}

	query.getTransactionByID(peer, channelName, trxnId, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		});
});


// 
// Query Get Block by Hash
app.get('/channels/:channelName/blocks', function (req, res) {
	logger.debug('================ GET BLOCK BY HASH ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let hash = req.query.hash;
	let channelName = req.params.channelName;
	let peer = req.query.peer;
	if (!hash) {
		res.json(getErrorMessage('\'hash\''));
		return;
	}

	query.getBlockByHash(peer, channelName, hash, req.username, req.orgname).then(
		function (message) {
			res.send(message);
		});
});

// 
//Query for Channel Information
app.get('/channels/:channelName', function (req, res) {
	logger.debug(
		'================ GET CHANNEL INFORMATION ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let peer = req.query.peer;
	let channelName = req.params.channelName;

	query.getChainInfo(peer, channelName, req.username, req.orgname).then(
		function (message) {
			res.send(message);
		});
});

// Query to fetch all Installed/instantiated chaincodes
app.get('/chaincodes', function (req, res) {
	var peer = req.query.peer;
	var installType = req.query.type;
	var channelName = req.query.channelName;
	// able to add constants
	if (installType === 'installed') {
		logger.debug(
			'================ GET INSTALLED CHAINCODES ======================');
	} else if (installType === "instantiated") {
		logger.debug(
			'================ GET INSTANTIATED CHAINCODES ======================');
	} else {
		logger.debug("Add Constant Type")
	}

	query.getInstalledChaincodes(peer, channelName, installType, req.username, req.orgname)
		.then(function (message) {
			res.send(message);
		});
});
// Query to fetch channels
app.get('/channels', function (req, res) {
	logger.debug('================ GET CHANNELS ======================');
	logger.debug('peer: ' + req.query.peer);
	var peer = req.query.peer;
	if (!peer) {
		res.json(getErrorMessage('\'peer\''));
		return;
	}

	query.getChannels(peer, req.username, req.orgname)
		.then(function (
			message) {
			res.send(message);
		});
});