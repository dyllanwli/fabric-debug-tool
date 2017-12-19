window.onload = function(){   
    
    // enroll admin
    var btn1 = document.getElementById("enrollAdmin");
    btn1.onclick = function(){
        var xhr = new XMLHttpRequest();
        var vusername = document.getElementById("username").value;
        var vorgName = document.getElementById("orgName").value;
        var form = "username="+ vusername + "&orgName=" + vorgName;
        xhr.open("POST", "/users", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send(form);
    }
    
    // load enrolled user
    var number = 2
    function loadUser(){
        var element1 = document.createElement("input");
        element1.setAttribute('type','radio');
        element1.setAttribute('name','token');
        element1.setAttribute('value','no token');
        document.getElementById('div_identified').appendChild(element1);
        
        var element2 = document.createElement("label");
        element2.setAttribute('for','token');
        element2.appendChild(document.createTextNode('no token'));
        document.getElementById('div_identified').appendChild(element2);
    }
    loadUser();
    // token func
    var token;
    $(document).ready(function() {
        $('input[type=radio][name=token]').change(function() {
            if (this.checked == true) {
                token = this.value
            }
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
        xhr.send(jsonData);
    };

    // invoke chaincode 
    // TODO: support json invoke
    var xhr = new XMLHttpRequest();
    var btn6 = document.getElementById("invokeTransaction");
    btn6.onclick =function(){
        var vchannelName = document.getElementById("invoke_channelName").value;
        var vchaincodeName = document.getElementById("invoke_chaincodeName").value;
        // TODO
        var vargs = document.getElementById("invoke_args").value.split(",");
        var vfcn = document.getElementById("invoke_fcn").value;
        var jsonData = JSON.stringify({
            args: vargs,
            fcn:vfcn
        });
        window.alert(jsonData);
        xhr.open("POST", "/channels/"+vchannelName+"/chaincodes/"+vchaincodeName, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send(jsonData);
    }

    //query chaincode 
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
        xhr.responseType = "blob";
        xhr.send();
    };
    query2.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"/blocks/"+parameter.vquery_blockId+"?peer="+ parameter.vquery_peer;
        window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send();
    };
    query3.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"/transactions/"+parameter.vquery_trxnId+"?peer="+parameter.vquery_peer;
        // window.alert(url);        
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send();
    };
    query4.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/channels/"+parameter.vquery_channelName+"?peer="+parameter.vquery_peer;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send();
    };
    query5.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "/chaincodes?peer="+parameter.vquery_peer+"&type="+parameter.vquery_type;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send();
    };
    query6.onclick=function(){
        var parameter = new Object();
        getQueryParameter(parameter);
        var xhr = new XMLHttpRequest();
        url = "channels?peer="+parameter.vquery_peer;
        // window.alert(url);
        xhr.open("GET",url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('authorization', ' Bearer '+ token);
        xhr.send();
    };
}