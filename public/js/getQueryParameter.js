// use to get query parameter
function getQueryParameter(obj) {
    obj.token = document.getElementById("token").value;
    obj.vquery_peer = document.getElementById("query_peer").value;
    obj.vquery_chaincodeName = document.getElementById("query_chaincodeName").value;
    obj.vquery_channelName = document.getElementById("query_channelName").value;
    obj.vquery_fcn = document.getElementById("query_fcn").value;
    // TODO: suupport json query
    obj.vquery_args = document.getElementById("query_args").value;
    obj.vquery_blockId = document.getElementById("query_blockId").value;
    obj.vquery_trxnId = document.getElementById("query_trxnId").value;
    obj.vquery_hash = document.getElementById("query_hash").value;
    obj.vquery_type = document.getElementById("query_type").value;
    // var parameter = JSON.stringify({
    //     token:vtoken,
    //     peer:vquery_peer,
    //     chaincodeName:vquery_chaincodeName,
    //     channelName:vquery_channelName,
    //     fcn:vquery_fcn,
    //     args:vquery_args,
    //     trxnId:vquery_trxnId,
    //     hash:vquery_hash,
    //     type:vquery_hash
    // });
}