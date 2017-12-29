/**
 * Copyright 2017 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an 'AS IS' BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
'use strict';
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var express = require('express');
var session = require('express-session');
//var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');
var path = require('path');
var http = require('http');
var util = require('util');
var crypto = require('crypto');
//var async=require('async');
var app = express();
var expressJWT = require('express-jwt');
var jwt = require('jsonwebtoken');
var bearerToken = require('express-bearer-token');
var cors = require('cors');
var multer=require('multer');
var fs=require('fs');
var upload=multer({dest:'./uploads/'});
var unless = require('express-unless');
var mysql=require('mysql');
var pool = mysql.createPool({     
  host     : 'localhost',       
  user     : 'root',              
  password : '1234',       
  port: '3306',                   
  database: 'fabricexplorer'
}); 

require('./config.js');
var hfc = require('fabric-client');

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
//////////////////////////////// SET CONFIGURATONS ////////////////////////////
///////////////////////////////////////////////////////////////////////////////
app.options('*', cors());
app.use(cors());
//support parsing of application/json type post data
app.use(bodyParser.json());
//support parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
	extended: false
}));
app.use(express.static(path.join(__dirname, 'public')));
app.set('views', path.join(__dirname, 'views'));
//app.set('view engine', 'ejs');
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
	secret: 'thisismysecret'
}).unless({
	path: ['/users/login','/users/register','/users/regvali','/','/users','/favicon.ico','/users/forgetpwd','/users/downloadlogfile']
}));
app.use(bearerToken());
app.use(function(req, res, next) {
    var url=req.originalUrl;
	if (url.indexOf('/users') >= 0) {
		return next();
	}
	if(!req.token){
		res.redirect('/users/login');
		return;
	}
	var token = req.token;
	jwt.verify(token, app.get('secret'), function(err, decoded) {
		if (err) {
			res.send({
				success: false,
				message: 'Failed to authenticate token. Make sure to include the ' +
					'token returned from /users call in the authorization header ' +
					' as a Bearer token'
			});
            res.redirect('/users/login');
			return;
		} else {
			// add the decoded user name and org name to the request object
			// for the downstream code to use
			req.username = decoded.username;
			req.orgname = decoded.orgName;
			logger.debug(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.username, decoded.orgName));
			return next();
		}
	});
});

///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// START SERVER /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
var server = http.createServer(app).listen(port, function() {});
logger.info('****************** SERVER STARTED ************************');
logger.info('**************  http://' + host + ':' + port +
	'  ******************');
server.timeout = 240000;

function getErrorMessage(field) {
	var response = {
		success: false,
		message: field + ' field is missing or Invalid in the request'
	};
	return response;
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////// REST ENDPOINTS START HERE ///////////////////////////
///////////////////////////////////////////////////////////////////////////////
app.post('/users', function(req, res) {
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
	helper.getRegisteredUsers(username, orgName, true).then(function(response) {
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
// 登录
var orgMap={'org1':'A省公司','org2':'B省公司','org3':'C省公司'};
app.route('/users/login')
 .get(function(req, res) {
     res.render('login', { title: '用户登录' });
 })
.post(function(req, res) {
    var sql='select * from fabricusers where username=?';
    var username=req.body.username;
    var password=req.body.password;
    pool.query(sql,[username],function (err, result) {
        if(err){
          console.log('[SELECT ERROR] - ',err.message);
          return;
        }
        if(result==null||result==''){
            res.render('login',{loginerr:'nameerr'});
        }else{
            if(result[0].userpassword==password){
			var uorg=result[0].org;
                var user={
                  username:result[0].username,
                  password:result[0].userpassword,
                  orgName:result[0].org,
                  orgRname:orgMap[uorg]
                }
                var token = jwt.sign({
                exp: Math.floor(Date.now() / 1000) + parseInt(hfc.getConfigSetting('jwt_expiretime')),
                username: user.username,
                orgName: user.orgName
                }, app.get('secret'));
                res.render('mainpage',{token:token,username:user.username,userorg:user.orgRname});
            }else{
                res.render('login',{loginerr:'pwderr'});
            }
        }
    });
});
//注册验证账户唯一
app.post('/users/regvali',function(req,res){
    var username=req.body.username;
    var sql='select userid from fabricusers where username=?';
    pool.query(sql,[username],function(err,result){
        if(err){
          console.log('[SELECT ERROR] - ',err.message);
          return;
        }
        //console.log(result);
        if(result!=null&&result!=''){
            res.json({err:'账户名已存在'});
        }else{
            res.json({err:''});
        }
    })
})
//用户注册
app.route('/users/register')
 .get(function(req, res) {
     res.render('register', { title: '用户注册' });
 })
 .post(function(req, res) {
     var user={
         username: req.body.username,
         password: req.body.password,
         org:req.body.org,
         phonenumber:req.body.phonenumber
     }
     helper.getRegisteredUsers(user.username, user.org, true).then(function(response) {
         if (response && typeof response !== 'string') {
             var peers;
             logger.info(user.username);
             invoke.invokeChaincode(peers, "itemchannel", "itemcc", "initUser", [], user.username, user.org).then(function(message){
				if(message.indexOf("Failed")==-1){
					initusertosql(user,res);
				}else{
                    invoke.invokeChaincode(peers, "itemchannel", "itemcc", "initUser", [], user.username, user.org).then(function(data){
                        if(data.indexOf("Failed")==-1) {
                            initusertosql(user,res);
                        }else{
                            res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
                            res.write('<script>alert("注册失败");window.location.href="/users/register"</script>');
                            res.end();
						}
					});
				}
			 });
         } else {
             res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
             res.write('<script>alert("注册失败");window.location.href="/users/register"</script>');
             res.end();
         }
     });
 });
function initusertosql(user,res){
    var sql='insert into fabricusers values(0,?,?,?,?)';
    pool.query(sql,[user.username,user.password,user.phonenumber,user.org],function(err,result){
        if(err){
            console.log('[SELECT ERROR] - ',err.message);
            res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
            res.write('<script>alert("注册失败");window.location.href="/users/register"</script>');
            res.end();
            return;
        }
        if(result.affectedRows=='1'){
            res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
            res.write('<script>alert("注册成功");window.location.href="/users/login"</script>');
            res.end();
        }else{
            res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
            res.write('<script>alert("注册失败");window.location.href="/users/register"</script>');
            res.end();
        }
    });
}
 
 //忘记密码
 app.route('/users/forgetpwd')
 .get(function(req,res){
 	if(req.query.name){
 		var name=req.query.name;
 		var phone=req.query.phone;
 		var sql1="select phonenumber from fabricusers where username=?";
 		pool.query(sql1,[name],function(err,result){
	 		if(err){
	          console.log('[SELECT ERROR] - ',err.message);
	          return;
	        }
	        if(result[0].phonenumber==phone){
	        	res.json({err:""});
	        }else{
	        	res.json({err:"账号与手机号不匹配"});
	        }
 		});		
 	}else{
 		res.render('forgetpwd',{title:"忘记密码"});
 	}
 })
 .post(function(req,res){
 	var name=req.body.username;
 	var phone=req.body.phonenumber;
 	var pwd=req.body.password;
 	var sql1="select phonenumber from fabricusers where username=?";
 		pool.query(sql1,[name],function(err,result){
	 		if(err){
	          console.log('[SELECT ERROR] - ',err.message);
	          return;
	        }
	        if(result[0].phonenumber==phone){
	        	var sql2="update fabricusers set userpassword=? where username=?";
	        	pool.query(sql2,[pwd,name],function(err1,result1){
	        		if(err1){
			          console.log('[SELECT ERROR] - ',err.message);
			          return;
	        		}
	        		if(result1.affectedRows=='1'){
	        			res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
 						res.write('<script>alert("密码修改成功");window.location.href="/users/login"</script>');
 						res.end();
	        		}else{
	        			res.writeHead(200, {'Content-type' : 'text/html;charset=utf-8'});
 						res.write('<script>alert("密码修改失败");window.location.href="/users/forgetpwd"</script>');
 						res.end();
	        		}
	        	});
	        }
 		});		
 });
 
 
/* app.get('/logout', function(req, res) {
     req.session.user = null;
     res.redirect('/');
 });*/
  //请求菜单页面
 app.get('/left',function(req,res){
 	var topic=req.query.menu;
 	if(topic=="log"){
 		res.render('leftlog');
 	}else if(topic=="product"){
 		res.render('leftproduct');
 	}else if(topic=="leftapi"){
 		res.render('leftapi');
 	}
 });
 //请求主体界面
 app.get('/right',function(req,res){
 	var topic=req.query.menu;
 	if(topic=="orderlog"){
 		res.render('rightlog');
 	}else if(topic=="explorer"){
	var opt = {  
         host:'172.20.29.20',  
         port:'8080',  
         method:'GET',  
         path:'/'  
    	};
	var sreq=http.request(opt,function(sres){
	sres.pipe(res);
	});
	req.pipe(sreq);
 		//res.render('rightexplorer');
 	}else if(topic=="orderproduct"){
 		res.render('rightorderproduct');
 	}else if(topic=="producttransaction"){
 		res.render('rightproducttransaction');
 	}else if(topic=="accountinfo"){
 		res.render('rightaccountinfo');
 	}else if(topic=="help"){
 		res.render('help')
 	}else if(topic=="apihelp"){
 		res.render('apihelp')
 	}
 });
 
// Create Channel
app.post('/channels', function(req, res) {
	logger.info('<<<<<<<<<<<<<<<<< C R E A T E  C H A N N E L >>>>>>>>>>>>>>>>>');
	logger.debug('End point : /channels');
	var channelName = req.body.channelName;
	var channelConfigPath = req.body.channelConfigPath;
	logger.debug('Channel name : ' + channelName);
	logger.debug('channelConfigPath : ' + channelConfigPath); //../artifacts/channel/mychannel.tx
	if (!channelName) {
		res.json(getErrorMessage('\'channelName\''));
		return;
	}
	if (!channelConfigPath) {
		res.json(getErrorMessage('\'channelConfigPath\''));
		return;
	}

	channels.createChannel(channelName, channelConfigPath, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Join Channel
app.post('/channels/:channelName/peers', function(req, res) {
	logger.info('<<<<<<<<<<<<<<<<< J O I N  C H A N N E L >>>>>>>>>>>>>>>>>');
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

	join.joinChannel(channelName, peers, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Install chaincode on target peers
app.post('/chaincodes', function(req, res) {
	logger.debug('==================== INSTALL CHAINCODE ==================');
	var peers = req.body.peers;
	var chaincodeName = req.body.chaincodeName;
	var chaincodePath = req.body.chaincodePath;
	var chaincodeVersion = req.body.chaincodeVersion;
	logger.debug('peers : ' + peers); // target peers list
	logger.debug('chaincodeName : ' + chaincodeName);
	logger.debug('chaincodePath  : ' + chaincodePath);
	logger.debug('chaincodeVersion  : ' + chaincodeVersion);
	if (!peers || peers.length == 0) {
		res.json(getErrorMessage('\'peers\''));
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

	install.installChaincode(peers, chaincodeName, chaincodePath, chaincodeVersion, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Instantiate chaincode on target peers
app.post('/channels/:channelName/chaincodes', function(req, res) {
	logger.debug('==================== INSTANTIATE CHAINCODE ==================');
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
	.then(function(message) {
		res.send(message);
	});
});
// Invoke transaction on chaincode on target peers
app.post('/channels/:channelName/chaincodes/:chaincodeName', function(req, res) {
	logger.debug('==================== INVOKE ON CHAINCODE ==================');
	var peers = req.body.peers;
	var chaincodeName = req.params.chaincodeName;
	var channelName = req.params.channelName;
	var fcn = req.body.fcn;
	var args = req.body.args;
	logger.debug('channelName  : ' + channelName);
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
	.then(function(message) {
		res.send(message);
	});
});
//post log file and uplog
app.post('/uplogfile',upload.single('logfile'),function(req,res){
	logger.debug('===========upload file=========');
	var file=req.file;
	var peers=req.body.peers;
	var oldname=file.originalname;
	var index=oldname.lastIndexOf('.');
	var newname=oldname.substring(0,index)+"_"+Date.now();
	var logname=newname;
	var newpath='uploads/'+newname+oldname.substring(index);
	fs.renameSync(file.path,newpath);
	fs.readFile(newpath,function(err,data){
		var logbody="";
		if(err){
			logbody="log text ..."
		}
		var hasher=crypto.createHash('md5');
		hasher.update(data);
		logbody=hasher.digest('hex');
		var args=[logname,logbody];
		invoke.invokeChaincode(peers, 'logchannel', 'logcc', 'uploadLog', args, req.username, req.orgname)
		.then(function(message) {
			res.send(message);
		});
	});
	var sql='insert into logsinfo values(0,?,?,?)';
	pool.query(sql,[logname,newpath,1],function(err,result){
		if(err){
          console.log('[SELECT ERROR] - ',err.message);
          return;
        }
        if(result.affectedRows==1){
        	logger.info("mysql save ok!");
        }else{
        	logger.info("mysql save false!");
        }
	});
});
//检查是否存在log文件
app.get('/checklog',function(req,res){
	logger.debug('==================== check if has logfile ==================');
	var logname=req.query.logname;
	var sql='select * from logsinfo where logname=?';
	pool.query(sql,[logname],function(err,result){
		if(err){
          console.log('[SELECT ERROR] - ',err.message);
          res.send("sql ERROR");
          return;
        }
        if(result==''||result==null){
        	res.send("no log");
        }else{
        	res.send(result[0].logpath);
        }
	});
});
//下载log文件
app.post('/users/downloadlogfile',function(req,res){
	var logpath=req.body.logpath;
	res.download(logpath);
});

// Query on chaincode on target peers
app.get('/channels/:channelName/chaincodes/:chaincodeName', function(req, res) {
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

	args = args.replace(/'/g, '"');

	args = JSON.parse(args);
	logger.debug(args);

	query.queryChaincode(peer, channelName, chaincodeName, args, fcn, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
//  Query Get Block by BlockNumber
app.get('/channels/:channelName/blocks/:blockId', function(req, res) {
	logger.debug('==================== GET BLOCK BY NUMBER ==================');
	let blockId = req.params.blockId;
	let peer = req.query.peer;
	logger.debug('channelName : ' + req.params.channelName);
	logger.debug('BlockID : ' + blockId);
	logger.debug('Peer : ' + peer);
	if (!blockId) {
		res.json(getErrorMessage('\'blockId\''));
		return;
	}

	query.getBlockByNumber(peer, blockId, req.username, req.orgname)
		.then(function(message) {
			res.send(message);
		});
});
// Query Get Transaction by Transaction ID
app.get('/channels/:channelName/transactions/:trxnId', function(req, res) {
	logger.debug(
		'================ GET TRANSACTION BY TRANSACTION_ID ======================'
	);
	logger.debug('channelName : ' + req.params.channelName);
	let trxnId = req.params.trxnId;
	let peer = req.query.peer;
	if (!trxnId) {
		res.json(getErrorMessage('\'trxnId\''));
		return;
	}

	query.getTransactionByID(peer, trxnId, req.username, req.orgname)
		.then(function(message) {
			res.send(message);
		});
});
// Query Get Block by Hash
app.get('/channels/:channelName/blocks', function(req, res) {
	logger.debug('================ GET BLOCK BY HASH ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let hash = req.query.hash;
	let peer = req.query.peer;
	if (!hash) {
		res.json(getErrorMessage('\'hash\''));
		return;
	}

	query.getBlockByHash(peer, hash, req.username, req.orgname).then(
		function(message) {
			res.send(message);
		});
});
//Query for Channel Information
app.get('/channels/:channelName', function(req, res) {
	logger.debug(
		'================ GET CHANNEL INFORMATION ======================');
	logger.debug('channelName : ' + req.params.channelName);
	let peer = req.query.peer;

	query.getChainInfo(peer, req.username, req.orgname).then(
		function(message) {
			res.send(message);
		});
});
// Query to fetch all Installed/instantiated chaincodes
app.get('/chaincodes', function(req, res) {
	var peer = req.query.peer;
	var installType = req.query.type;
	//TODO: add Constnats
	if (installType === 'installed') {
		logger.debug(
			'================ GET INSTALLED CHAINCODES ======================');
	} else {
		logger.debug(
			'================ GET INSTANTIATED CHAINCODES ======================');
	}

	query.getInstalledChaincodes(peer, installType, req.username, req.orgname)
	.then(function(message) {
		res.send(message);
	});
});
// Query to fetch channels
app.get('/channels', function(req, res) {
	logger.debug('================ GET CHANNELS ======================');
	logger.debug('peer: ' + req.query.peer);
	var peer = req.query.peer;
	if (!peer) {
		res.json(getErrorMessage('\'peer\''));
		return;
	}

	query.getChannels(peer, req.username, req.orgname)
	.then(function(
		message) {
		res.send(message);
	});
});

//分页返回请求信息
app.get('/getallinfo/channels/:channelName/chaincodes/:chaincodeName',function(req,res){
	logger.debug('==============get allinfo===============');
	var page=req.query.page;//页码
	var topic=req.query.topic;//获取信息主题
	var fcn;//chaincode函数
	var results=new Array();//返回结果集，一个json串
	var chaincodeName = req.params.chaincodeName;
	var channelName = req.params.channelName;
	var peer = req.query.peer;
	var users=new Array();

 	//根据主题判断调用函数
	switch(topic){ 
		//log all
		case '1': 
			fcn="queryLogsByUser";
			break; 
		case '2': 
			fcn="queryLogsByUser";
			break; 
			//item not own
		case '3': 
			fcn="queryItemsByItemOwner";
			break;
		case '4': 
			fcn="queryItemsByItemOwner";
			break; 
		default: 
		res.json(getErrorMessage('\'topic\''));
		return;
	} 
	var records = new Array();
	//获取用户列表
	if (topic == 2 || topic == 4 ) {
		users.push(req.username);
		if (!users) {
		res.json(getErrorMessage('\'No users info\''));
		return;
		}
		var args=[];
		if(topic==4){
			args=['',users[0]];
		}else{
			args=[users[0]];
		}
		query.queryChaincode(peer, channelName, chaincodeName,args, fcn, req.username, req.orgname).then(function(message) {
		records =records.concat(JSON.parse(message));
		if (records.length==0) {
			res.json(getErrorMessage('\'no record in the page\''));
			return;
		}
		var results=new Array();
		for (var index = page * 10-10; (index < records.length) && (index < (page * 10)); index++) {
			var element = records[index];
			results.push(element);
		}
		var totalpage=Math.ceil(records.length/10.0);
		results.push({"totalpages":totalpage});
		res.send(results);
		return;
		});
		
	}
	else{
		
	var sql="select username from fabricusers";
	pool.query(sql,function(err,result){
        if(err){
          console.log('[SELECT ERROR] - ',err.message);
          return;
        }
	for(var j=0;j<result.length;j++){
		if(topic==3&&result[j].username==req.username){
			continue;
		}
		users.push(result[j].username);
	}	
	//logger.debug(result[0].username);
	if (!users) {
		res.json(getErrorMessage('\'No users info\''));
		return;
	}
		var promisearray=[];
		for(var i=0;i<users.length;i++){
			var args=[];
			if(topic==1){
				args=[users[i]];
			}else{
				args=['',users[i]];
			}
			logger.debug(users[i]);
			promisearray.push(query.queryChaincode(peer, channelName, chaincodeName, args, fcn, req.username, req.orgname).then(function(data){//logger.info(data);
return data;}));
		} 
		Promise.all(promisearray).then(function(data){
			//logger.info(data);
			for(var i=0;i<data.length;i++){
				records=records.concat(JSON.parse(data[i]));
			}
		if (records.length==0) {
			res.json(getErrorMessage('\'no record in the page\''));
			return;
		}
		var results=new Array();
		for (var index = page * 10-10; (index < records.length) && (index < (page * 10)); index++) {
			var element = records[index];
			results.push(element);
		}
		var totalpage=Math.ceil(records.length/10.0);
		results.push({"totalpages":totalpage});
		res.send(results);
		return;
		});
	});
	}	
});

