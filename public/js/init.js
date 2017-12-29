window.onload = function () {
    var host = "http://39.106.141.206:4000/"

    // clear log
    clear()

    // result area
    function loadResult(xhr) {
        xhr.onload = (response) => {
            var ele = document.getElementById("resultArea")
            ele.appendChild(document.createTextNode(response))
        }
    }

    // enroll admin
    var btn1 = document.getElementById("enrollAdmin")
    btn1.onclick = function () {
        var xhr = new XMLHttpRequest()
        var vusername = document.getElementById("username").value
        var vorgName = document.getElementById("orgName").value
        var form = "username=" + vusername + "&orgName=" + vorgName
        xhr.open("POST", "/users", true)
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded")
        // callback function
        xhr.onreadystatechange = function () { //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    try {
                        response = JSON.parse(response)
                    } catch (e) {
                        alert("enroll failure: no response")
                        return
                    }
                    tk = response.token
                    // regular the response
                    delete response.token
                    delete response.secret
                    response = JSON.stringify(response)
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    ele.value += response + "\n\n"
                    loadUser(vusername + '_' + vorgName, tk)
                }
            }
        }
        // call backend
        xhr.send(form)
    }

    // token func
    var token
    $(document).ready(function () {
        $(document).on('change', "input[name='token']", function () {
            token = $(this).val()
        })
    })

    // channel func
    $(document).ready(function () {
        $(document).on('change', "input[name='channel']", function () {
            channelName = $(this).val()
        })
    })

    // create channels
    var xhr = new XMLHttpRequest()
    var btn2 = document.getElementById("createChannel")
    btn2.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var vchannelName = document.getElementById("channelName").value
        if (vchannelName == channelName) {
            alert("Channel exits, no need to create channel")
            return
        }
        loadChannel(vchannelName)
        var vchannelConfigPath = document.getElementById("channelConfigPath").value
        var jsonData = JSON.stringify({
            channelName: vchannelName,
            channelConfigPath: vchannelConfigPath
        })
        xhr.open("POST", "/channels", true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () { 
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    if (response == "{}") {
                        alert("Check if channel are existing")
                        return
                    }
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        // call backend
        xhr.send(jsonData)
    }

    // join channels
    var xhr = new XMLHttpRequest()
    var btn3 = document.getElementById("joinChannel")
    btn3.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var vchannelName = channelName
        var vpeers = document.getElementById("join_peers").value.split(",")
        var temp = document.getElementsByName("peers")
        if (temp.length == 0) {
            
        } else {
            for (i = 0; i < temp.length; i++) {
                if (vpeers.includes(temp[i].id)) {
                    alert(temp[i].id + " has joinned.")
                    return
                }
            }
        }
        loadPeers(vpeers)
        var jsonData = JSON.stringify({
            peers: vpeers
        })
        xhr.open("POST", "/channels/" + vchannelName + "/peers", true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () {
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        // call backend
        xhr.send(jsonData)
    }
    // join all peers
    var xhr = new XMLHttpRequest()
    var btn3_all = document.getElementById("join_all")
    btn3_all.onclick = function(){
        if(token == null){
            alert("token is undefined")
            return
        }
        var vchannelName = channelName
        var vpeers = "peer1,peer2".split(",")
        var temp = document.getElementsByName("peers")
        if (temp.length == 0) {

        } else {
            for (i = 0; i < temp.length; i++) {
                if (vpeers.includes(temp[i].id)) {
                    alert(temp[i].id + " has joinned.")
                    return
                }
            }
        }
        loadPeers(vpeers)
        var jsonData = JSON.stringify({
            peers: vpeers
        })
        xhr.open("POST", "/channels/" + vchannelName + "/peers", true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () {
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        // call backend
        xhr.send(jsonData)
    }



    // install chaincode
    var xhr = new XMLHttpRequest()
    var btn4 = document.getElementById("installChaincode")
    btn4.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var vchaincodeName = document.getElementById("install_chaincodeName").value
        var vpeers = document.getElementById("install_peers").value.split(",")
        var vchaincodeVersion = document.getElementById("install_chaincodeVersion").value
        var vchaincodePath = document.getElementById("install_chaincodePath").value
        var jsonData = JSON.stringify({
            peers: vpeers,
            chaincodeName: vchaincodeName,
            chaincodePath: vchaincodePath,
            chaincodeVersion: vchaincodeVersion
        })
        xhr.open("POST", "/chaincodes", true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () { //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    if (response.includes("Failed") == true) {
                        return
                    }
                    var ele = document.getElementById("resultArea")
                    ele.value += response + "\n\n"
                    loadChaincode(vpeers,vchaincodeName,vchaincodePath,vchaincodeVersion)
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        // call backend
        xhr.send(jsonData)
    }

    // instantiate 
    var xhr = new XMLHttpRequest()
    var btn5 = document.getElementById("instantiateChaincode")
    btn5.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var vargs = document.getElementById("instan_args").value.split(",")
        var vchannelName = channelName
        var vchaincodeName = document.getElementById("instan_chaincodeName").value
        var vchaincodeVersion = document.getElementById("instan_chaincodeVersion").value
        var jsonData = JSON.stringify({
            args: vargs,
            chaincodeName: vchaincodeName,
            chaincodeVersion: vchaincodeVersion
        })
        xhr.open("POST", "/channels/" + vchannelName + "/chaincodes", true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () { 
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        // call backend
        xhr.send(jsonData)
    }

    // install and instantiate
    var xhr = new XMLHttpRequest()
    var btn5_install_instantiate = document.getElementById("install_instantiate")
    btn5_install_instantiate.onclick = function(){
        if(token == null){
            alert("token is undefined")
            return
        }
        btn4.onclick()
        btn5.onclick()
    }

    // invoke chaincode 
    var xhr = new XMLHttpRequest()
    var btn6 = document.getElementById("invokeTransaction")
    btn6.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var vchannelName = document.getElementById("invoke_channelName").value
        var vchaincodeName = document.getElementById("invoke_chaincodeName").value
        var vargs = document.getElementById("invoke_args").value.replace(/\{|\}/gi, "")
        vargs = vargs.replace(/\:/gi, ",").split(',')
        var vfcn = document.getElementById("invoke_fcn").value
        var jsonData = JSON.stringify({
            args: vargs,
            fcn: vfcn
        })
        xhr.open("POST", "/channels/" + vchannelName + "/chaincodes/" + vchaincodeName, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        // callback function
        xhr.onreadystatechange = function () { 
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    if (response.includes("Failed") == true) {
                        ele.value += "Invoke Failed:\n" + response + "\n\n"
                    } else {
                        ele.value += "Invoke Successful. The transaction ID is:\n" + response + "\n\n"
                    }
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }

        }
        // call backend
        xhr.send(jsonData)
    }

    //query by args 
    var query1 = document.getElementById("query1")
    var query2 = document.getElementById("query2")
    var query3 = document.getElementById("query3")
    var query4 = document.getElementById("query4")
    var query5 = document.getElementById("query5")
    var query6 = document.getElementById("query6")
    // click incident
    query1.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        parameter.vquery_args = "[\"" + parameter.vquery_args + "\"]"
        parameter.vquery_args = escape(parameter.vquery_args)
        url = "/channels/" + parameter.vquery_channelName + "/chaincodes/" + parameter.vquery_chaincodeName + "?peer=" + parameter.vquery_peer + "&fcn=" + parameter.vquery_fcn + "&args=" + parameter.vquery_args
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)

        xhr.onreadystatechange = function () { 
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }

        }
        xhr.send()
    }
    // query by blockId
    query2.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        url = "/channels/" + parameter.vquery_channelName + "/blocks/" + parameter.vquery_blockId + "?peer=" + parameter.vquery_peer
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)

        xhr.onreadystatechange = function () { 
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    try {
                        response = JSON.parse(response)
                    } catch (e) {
                        alert("Query failure: no response, please check parameter and block status.")
                        return
                    }
                    writeData = response.data.data[0].payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1]
                    namespace = JSON.stringify(writeData.namespace)
                    writes = JSON.stringify(writeData.rwset.writes)

                    channel_header = response.data.data[0].payload.header.channel_header
                    time = JSON.stringify(channel_header.timestamp)
                    tx_id = JSON.stringify(channel_header.tx_id)
                    channnel_id = JSON.stringify(channel_header.channel_id)
                    version = JSON.stringify(channel_header.version)

                    w = "The Blockid query:\n======\nTime: " + time + "\nTransaction ID: " + tx_id + "\nChannel ID: " + channnel_id + "\nVersion: " + version + "\nChaincode: " + namespace + "\nWrites: " + writes + "\n=====\n"
                    ele.value += w + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }

        }
        xhr.send()
    }

    // qeury by transaction id
    query3.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        url = "/channels/" + parameter.vquery_channelName + "/transactions/" + parameter.vquery_trxnId + "?peer=" + parameter.vquery_peer
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)
        xhr.onreadystatechange = function () {
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    try {
                        response = JSON.parse(response)
                    } catch (e) {
                        alert("Query failure: no response, please check parameter and block status.")
                        return
                    }
                    channel_header = response.transactionEnvelope.payload.header.channel_header
                    time = JSON.stringify(channel_header.timestamp)
                    tx_id = JSON.stringify(channel_header.tx_id)
                    channnel_id = JSON.stringify(channel_header.channel_id)
                    version = JSON.stringify(channel_header.version)

                    writeData = response.transactionEnvelope.payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1]
                    namespace = JSON.stringify(writeData.namespace)
                    write = JSON.stringify(writeData.rwset.writes)
                    w = "The TransactionID query:\n=====\nTime: " + time + "\nTransaction ID: " + tx_id + "\nChannel ID: " + channnel_id + "\nVersion: " + version + "\nChaincode: " + namespace + "\nWrites: " + write + "\n=====\n"
                    ele.value += w + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }
        }
        xhr.send()
    }

    // query chaininfo
    query4.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        url = "/channels/" + parameter.vquery_channelName + "?peer=" + parameter.vquery_peer
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)

        xhr.onreadystatechange = function () {
            //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    try {
                        response = JSON.parse(response)
                    } catch (e) {
                        alert("Query failure: no response, please check parameter")
                    }
                    height = response.height
                    height = JSON.stringify(height)
                    ele.value += "The chaininfo.height: " + height + "\n" + "currentBlockHash and previousBlockHash are hidden\n\n"

                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }

        }
        xhr.send()
    }

    // query installType
    query5.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        url = "/chaincodes?peer=" + parameter.vquery_peer + "&type=" + parameter.vquery_type
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)

        xhr.onreadystatechange = function () { //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    ele.value += response + "\n\n"
                }
            }

        }
        xhr.send()
    }

    // query Channel
    query6.onclick = function () {
        if(token == null){
            alert("token is undefined")
            return
        }
        var parameter = new Object()
        getQueryParameter(parameter)
        var xhr = new XMLHttpRequest()
        url = "channels?peer=" + parameter.vquery_peer
        xhr.open("GET", url, true)
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.setRequestHeader('authorization', ' Bearer ' + token)

        xhr.onreadystatechange = function () { //Call a function when the state changes.
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    var response = xhr.responseText
                    var ele = document.getElementById("resultArea")
                    // ele.appendChild(document.createTextNode(response+"\n\n"))
                    ele.value += response + "\n\n"
                } else if (xhr.status == 401) {
                    alert("Response: 401, check if you have selected the user.")
                }
            }

        }
        xhr.send()
    }
}