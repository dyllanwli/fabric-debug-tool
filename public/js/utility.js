// use to get query parameter
function getQueryParameter(obj) {
    obj.vquery_peer = document.getElementById("query_peer").value;
    obj.vquery_chaincodeName = document.getElementById("query_chaincodeName").value;
    obj.vquery_channelName = document.getElementById("query_channelName").value;
    // obj.vquery_fcn = document.getElementById("query_fcn").value;
    obj.vquery_fcn = "query";
    // TODO: suupport json query
    obj.vquery_args = document.getElementById("query_args").value;
    obj.vquery_blockId = document.getElementById("query_blockId").value;
    obj.vquery_trxnId = document.getElementById("query_trxnId").value;
    obj.vquery_hash = document.getElementById("query_hash").value;
    obj.vquery_type = document.getElementById("query_type").value;
}

// clear result area
function clear() {
    clear_btn = document.getElementById("clear")
    clear_btn.onclick = function () {
        tx = document.getElementById("resultArea").value = ""
        tempalert("Log cleared",500)
    }
}
function tempalert(msg, duration) {
    var el = document.createElement("div");
    el.setAttribute("style", "position:absolute;top:5%;left:70%;background-color:white;text-align:center;");
    el.innerHTML = msg;
    setTimeout(function () {
        el.parentNode.removeChild(el);
    }, duration);
    document.body.appendChild(el);
}


// load enrolled user
function loadUser(name,resToken){
    // input
    var element1 = document.createElement("input")
    element1.setAttribute('type','radio')
    element1.setAttribute('name','token')
    element1.setAttribute('value',resToken)
    document.getElementById('div_select_token').appendChild(element1)
    // label
    var element2 = document.createElement("label")
    element2.setAttribute('for','token')
    element2.appendChild(document.createTextNode(name))
    document.getElementById('div_select_token').appendChild(element2)
}

// load channel
function loadChannel(name){
    var element1 = document.createElement("input")
    element1.setAttribute('type','radio')
    element1.setAttribute('name','channel')
    element1.setAttribute('value',name)
    document.getElementById('div_channels').appendChild(element1)
    // label
    var element2 = document.createElement("label")
    element2.setAttribute('for','channel')
    element2.appendChild(document.createTextNode(name))
    document.getElementById('div_channels').appendChild(element2)
}

function loadPeers(name){
    for(i = 0;i<name.length;i++){
        // label
        var element2 = document.createElement("label")
        element2.setAttribute('name','peers')
        element2.setAttribute('id',name[i])
        element2.appendChild(document.createTextNode(name[i]+", "))
        document.getElementById('div_join_form').appendChild(element2)
    }
}