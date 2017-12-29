$(function(){
	$("#username").html(sessionStorage.username);
	$("#username2").html(sessionStorage.username);
	$("#org2").html(sessionStorage.userorg);
	var bal=getUserInfo(sessionStorage.username);
console.log(bal);
	$("#balance2").html(bal);
})
function getUserInfo(args){
var tmpBalance;
		$.ajax({
			async:false,
			type:"get",
			url:"/channels/itemchannel/chaincodes/itemcc?peer=peer1&fcn=queryBalance&args=%5B%22"+args+"%22%5D&startKey=''&endKey=''",
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				//console.log(data);
				//var start=data.indexOf("[");
				//var end=data.indexOf("]");
				//var balance=data.substring(start+1,end);
				if(data.indexOf("Error")==-1){
				tmpBalance=data;}else{tmpBalance="get fail";}
				//console.log(tmpBalance);
			},
				error:function(data){
				console.log(data);
			}	
		});
		return tmpBalance;
}
