$(function(){
	$("#user").html(sessionStorage.username);
	$("#chaincode_").click(function(){
		getright("chaincode");
	}).hover(over,out);
	$("#invoke_").click(function(){
		getright("invoke");
	}).hover(over,out);
	$("#channel_").click(function(){
		getright("channel_","channel.js");
	}).hover(over,out);
	$("#query_").click(function(){
		getright("query");
	}).hover(over,out);
	// 
	// debug
	// $("#leftresult").on('click','#explorer1',function(){
	// 	$(".nav").children("li").removeClass("active");
	// 	$(this).addClass("active");
	// 	$("#rightbody").empty();
	// 	$explo=$("<h4>explorer1 Under developing</h4>");
	// 	$("#rightbody").append($explo);
	// });
	$("#logout").click(function(){
		sessionStorage.clear();
		window.location.href="/users/login";
	});
	$("#user").click(function(){
		getright("accountinfo");
	});
});
function getright(topic,script){
	$("#rightbody").empty();
	$.ajax({
			type:"get",
			url:"/right/?menu="+topic,
			dataType:"html",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+ sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				$("#rightbody").html(data);
			},
			error:function(data){
				console.log(data);
			},
			complete: function(){
				if(script){
					$.getScript("/js/"+script)
					alert("script loaded")
				}
			}
	});

}
function over(){
	$(this).addClass("cur");
}
function out(){
	$(this).removeClass("cur");
}
