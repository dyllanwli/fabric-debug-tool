window.onload = function(){   
    var host = "http://39.106.141.206:4000/";

    // result area
    function loadResult(xhr){
        xhr.onload = (response) => {
            var ele = document.getElementById("resultArea");
            ele.appendChild(document.createTextNode(response))
        }
    }

    // load enrolled user
    function loadUser(name,resToken){
        var element1 = document.createElement("input");
        element1.setAttribute('type','radio');
        element1.setAttribute('name','token');
        element1.setAttribute('value',resToken);
        document.getElementById('div_identified').appendChild(element1);
        
        var element2 = document.createElement("label");
        element2.setAttribute('for','token');
        element2.appendChild(document.createTextNode(name));
        document.getElementById('div_identified').appendChild(element2);
    }

    // enroll admin
    var btn1 = document.getElementById("enrollAdmin");
    btn1.onclick = function(){
        var xhr = new XMLHttpRequest();
        var vusername = document.getElementById("username").value;
        var vorgName = document.getElementById("orgName").value;
        var form = "username="+ vusername + "&orgName=" + vorgName;
        xhr.open("POST", "/users", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                response = JSON.parse(response);
                tk = response.token;
                // regular the response
                delete response.token;
                delete response.secret;
                response = JSON.stringify(response)
                ele.appendChild(document.createTextNode(response+"\n\n"));
                loadUser(vusername+'_'+vorgName,tk);
            }
        }
        // call backend
        xhr.send(form);
    }
    
    // token func
    var token;
    $(document).ready(function() {
        $(document).on('change',"input[name='token']",function(){
            token = $(this).val();
        });
    });
    
    // create channels
    var xhr = new XMLHttpRequest();
    var btn2 = document.getElementById("createChannel");
    btn2.onclick = function(){
        var vchannelName = document.getElementById("channelName").value;
        var vchannelConfigPath = document.getElementById("channelConfigPath").value;
        var jsonData = JSON.stringify({
            channelName: vchannelName,
            channelConfigPath: vchannelConfigPath
        })
        xhr.open("POST", "/channels", true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
            }
        }
        // call backend
        xhr.send(jsonData);
    };
    
    // join channels
    var xhr = new XMLHttpRequest();
    var btn3 = document.getElementById("joinChannel");
    btn3.onclick = function(){
        var vchannelName = document.getElementById("join_channelName").value;
        var vpeers = document.getElementById("join_peers").value.split(",");
        var jsonData = JSON.stringify({
            peers: vpeers
        });
        // window.alert(jsonData);
        xhr.open("POST", "/channels/"+vchannelName+"/peers", true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
                
            }
        }
        // call backend
        xhr.send(jsonData);
    };
    

    // install chaincode
    var xhr = new XMLHttpRequest();
    var btn4 = document.getElementById("installChaincode");
    btn4.onclick = function(){
        var vchaincodeName = document.getElementById("install_chaincodeName").value;
        var vpeers = document.getElementById("install_peers").value.split(",");
        var vchaincodeVersion = document.getElementById("install_chaincodeVersion").value;
        var vchaincodePath = document.getElementById("install_chaincodePath").value;
        var jsonData = JSON.stringify({
            peers: vpeers,
            chaincodeName: vchaincodeName,
            chaincodePath: vchaincodePath,
            chaincodeVersion: vchaincodeVersion
        });
        // window.alert(jsonData);
        xhr.open("POST", "/chaincodes", true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
                
            }
        }
        // call backend
        xhr.send(jsonData);
    };

    // instantiate 
    var xhr = new XMLHttpRequest();
    var btn5 = document.getElementById("instantiateChaincode");
    btn5.onclick = function(){
        var vargs = document.getElementById("instan_args").value.split(",");
        var vchannelName = document.getElementById("instan_channelName").value;
        var vchaincodeName = document.getElementById("instan_chaincodeName").value;
        var vchaincodeVersion = document.getElementById("instan_chaincodeVersion").value;
        // var vfcn = document.getElementById("instan").values
        var jsonData = JSON.stringify({
            args: vargs,
            chaincodeName: vchaincodeName,
            chaincodeVersion: vchaincodeVersion
        });
        // window.alert(jsonData);
        xhr.open("POST", "/channels/"+vchannelName+"/chaincodes", true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
            }
        }
        // call backend
        xhr.send(jsonData);
    };

    // invoke chaincode 
    // TODO: support json invoke
    var xhr = new XMLHttpRequest();
    var btn6 = document.getElementById("invokeTransaction");
    btn6.onclick =function(){
        var vchannelName = document.getElementById("invoke_channelName").value;
        var vchaincodeName = document.getElementById("invoke_chaincodeName").value;
        var vargs = document.getElementById("invoke_args").value.replace(/\{|\}/gi,"");
        vargs = vargs.replace(/\:/gi,",").split(',');
        var vfcn = document.getElementById("invoke_fcn").value;
        var jsonData = JSON.stringify({
            args: vargs,
            fcn:vfcn
        });
        // window.alert(jsonData);
        xhr.open("POST", "/channels/"+vchannelName+"/chaincodes/"+vchaincodeName, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        // callback function
        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode("The transaction ID is: "+response+"\n\n"));
            }
        }
        // call backend
        xhr.send(jsonData);
    }

    //query by args 
    var query1 = document.getElementById("query1");
    var query2 = document.getElementById("query2");
    var query3 = document.getElementById("query3");
    var query4 = document.getElementById("query4");
    var query5 = document.getElementById("query5");
    var query6 = document.getElementById("query6");
    // click incident
    query1.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        parameter.vquery_args = "[\""+parameter.vquery_args+"\"]";
        parameter.vquery_args = escape(parameter.vquery_args);
        url = "/channels/"+parameter.vquery_channelName+"/chaincodes/"+parameter.vquery_chaincodeName+"?peer="+ parameter.vquery_peer+"&fcn="+parameter.vquery_fcn+"&args="+parameter.vquery_args;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);

        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
            }
        }
        xhr.send();
    };
    // query by blockId
    query2.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"/blocks/"+parameter.vquery_blockId+"?peer="+ parameter.vquery_peer;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);

        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                response = JSON.parse(response);
                num = response.header.number;
                channel_header = response.data.data[0].payload.header.channel_header;
                delete channel_header.extension
                num = JSON.stringify(num);
                channel_header = JSON.stringify(num);
                ele.appendChild(document.createTextNode("The Block.header.number: "+ num+"\n"+"The channel_header: "+channel_header+"\n\n"));
            }
        }
        xhr.send();
    };

    // qeury by transaction id
    query3.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"/transactions/"+parameter.vquery_trxnId+"?peer="+parameter.vquery_peer;
        // window.alert(url);        
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.onreadystatechange = function() {
            //Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                response = JSON.parse(response);
                channel_header = response.transactionEnvelope.payload.header.channel_header;
                delete channel_header.extension
                channel_header = JSON.stringify(channel_header);
                ele.appendChild(document.createTextNode("The channel_header: "+channel_header+"\n\n"));
            }
        }
        xhr.send();
    };
    
    // query chaininfo
    query4.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"?peer="+parameter.vquery_peer;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);

        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                response = JSON.parse(response);
                height = response.height
                height = JSON.stringify(height);
                ele.appendChild(document.createTextNode("The chaininfo.height: "+height+"\n"+"currentBlockHash and previousBlockHash are hidden\n\n"));
                
            }
        }
        xhr.send();
    };

    // query installType
    query5.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/chaincodes?peer="+parameter.vquery_peer+"&type="+parameter.vquery_type;
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);

        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
            }
        }

        xhr.send();
    };

    // query Channel
    query6.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "channels?peer="+parameter.vquery_peer;
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);

        xhr.onreadystatechange = function() {//Call a function when the state changes.
            if(xhr.readyState == XMLHttpRequest.DONE && xhr.status == 200) {
                var response = xhr.responseText;
                var ele = document.getElementById("resultArea");
                ele.appendChild(document.createTextNode(response+"\n\n"));
            }
        }
        xhr.send();
    };
}