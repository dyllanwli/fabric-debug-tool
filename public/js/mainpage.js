$(function(){
	$("#user").html(sessionStorage.username);
	$("#logapi").click(function(){
		getleft("log");
		getright("orderlog");
	}).hover(over,out);
	$("#protransaction").click(function(){
		getleft("product");
		getright("orderproduct");
	}).hover(over,out);
	$("#apihelp").click(function(){
		getleft("leftapi");
		getright("apihelp");
	}).hover(over,out);
	$("#phonescams").click(function(){
		getleft("phonescams");
		getright("phonescams");
	}).hover(over,out);
	$("#leftmenu").on('click','#explorer1',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		$("#rightbody").empty();
		//$("#rightbody").html("<iframe width=800px height=800px src=></iframe>")
		$explo=$("<h4>explorer1 Under developing</h4>");
		$("#rightbody").append($explo);
	});
	$("#leftmenu").on('click','#explorer2',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		$("#rightbody").empty();
		//$("#rightbody").html("<iframe width=800px height=800px src=></iframe>")
		$explo=$("<h4>explorer2 Under developing</h4>");
		$("#rightbody").append($explo);
	});
	$("#leftmenu").on('click','#orderlog',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		getright("orderlog");
	});
	$("#leftmenu").on('click','#orderproduct',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		getright("orderproduct");
	});
	$("#leftmenu").on('click','#producttransaction',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		getright("producttransaction");
	});
	$("#leftmenu").on('click','#accountinfo',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		getright("accountinfo");
	});
	$("#leftmenu").on('click','#help',function(){
		$(".nav").children("li").removeClass("active");
		$(this).addClass("active");
		getright("help");
	});
	$("#leftmenu").on('click','#phonemark',function(){
		$(".nav").children("li").removeClass("active");
                $(this).addClass("active");
		getright('phonescams');	
	})
	$("#logout").click(function(){
		sessionStorage.clear();
		window.location.href="/users/login";
	});
	$("#user").click(function(){
		getright("accountinfo");
	});
});
function getleft(topic){
	$("#leftmenu").empty();
	$.ajax({
			type:"get",
			url:"/left/?menu="+topic,
			dataType:"html",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				$("#leftmenu").html(data);
			},
			error:function(data){
				console.log(data);
			}
	});

}
function getright(topic){
	$("#rightbody").empty();
	$.ajax({
			type:"get",
			url:"/right/?menu="+topic,
			dataType:"html",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				$("#rightbody").html(data);
			},
			error:function(data){
				console.log(data);
			}
	});

}
function over(){
	$(this).addClass("cur");
}
function out(){
	$(this).removeClass("cur");
}
