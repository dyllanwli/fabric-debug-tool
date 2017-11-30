var http = require('http');
var express = require('express');
var app = express();
var terminal  = require("web-terminal");
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var url = require("url");
var fs = require("fs");

//////////
// If not install the package gobally, find the package in node_modules
///////

// var getParameter =  require("./api/getParameter.js")

var userData = {}
var page = function(){
    app.use(express.static('public'));
    app.get("public", function(req,res){
        console.log("this is get public")
        res.sendFile(__dirname + "/public/" + "index.html");
    });

    app.get('/enrollAdmin', function (req, res) {
        console.log("I get enroll")
        var response = {
            "username":req.query.username,
            "orgName":req.query.orgName
        };
        logger.info(response);
        res.end(JSON.stringify(response));
    });
    app.get('/createChannel', function (req, res) {
        console.log("I get channel")
        var response = {
            "channelName":req.query.channelName,
            "channelConfigPath":req.query.channelConfigPath
        };
        logger.info(response);
        res.end(JSON.stringify(response));
        userData.channelName = response.channelName;
        userData.channelConfigPath = response.channelConfigPath;
    });

    logger.info(userData,"this is userData");

    var server = app.listen(8087, function() {
        logger.info('Listening on port %d', server.address().port);
    });
    // web terminal 
    var ter = http.createServer(function (req, res) {
        res.writeHead(200, {"Content-Type": "text/plain"});
        res.end();
    });
    terminal(ter);
    logger.info("Web-terminal accessible at http://localhost:8088/terminal");
}

exports.page = page;
