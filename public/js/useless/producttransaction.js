var options={
        bootstrapMajorVersion:3,    //版本
        currentPage:1,    //当前页数
        numberOfPages:5,    //最多显示Page页
        totalPages:1,    //所有数据可以显示的页数
		itemTexts: function (type, page, current) {
		    switch (type) {
		        case "first":
		            return "首页";
		        case "prev":
		            return "上一页";
		        case "next":
		            return "下一页";
		        case "last":
		            return "末页";
		        case "page":
		            return page;
           	 }
        },
        onPageClicked:function(e,originalEvent,type,page){
			getAll(page,3);

        }
    }
$(function(){
            $("#page").bootstrapPaginator(options);
	getAll(1,3);
	$("#getitem2").click(function(){
		var name=$("#itemname").val();
		var business=$("#business").val();
		var args=[name,business];
		$("#itemlist").empty();
		if(name!=null&&name!=""||(business!=''&&business!=null)){			
			getItem(args);
		}else{
			getAll(1,3);
		}
	});
	$("#itemlist").on("click",".buy",function(){
		var name=$(this).parent().parent().children("td").eq(0).html();
		var price=$(this).parent().parent().children("td").eq(1).html();
		var property=$(this).parent().parent().children("td").eq(2).html();
		var owner=$(this).parent().parent().children("td").eq(3).html();
		var args=[name,property,price,owner];
		buyItem(args);
		});
})
function getItem(args){
	$.ajax({
		type:"get",
		url:"/channels/itemchannel/chaincodes/itemcc?peer=peer1&fcn=queryItemsByItemOwner&args=%5B%22"+args[0]+"%22%2c%22"+args[1]+"%22%5D",
		dataType:"text",
		beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
		success:function(data){
			console.log(data);
				var jsdata=JSON.parse(data);
				tempsavelog=jsdata;
				for(var i=0;i<jsdata.length;i++){
					var record=jsdata[i].Record;
					var $trs=$("<tr><td>"+record.name+"</td><td>"+record.price+"</td><td>"+record.property+"</td><td>"+record.owner+"</td><td><a style='margin:0 5px;' class='buy'>购买</a></td></tr>");
					$("#itemlist").append($trs);
				}
		},
		error:function(data){
			console.log(data);
		}
	});
}
function getOwnerItem(business){
	$.ajax({
		type:"get",
		url:"/channels/itemchannel/chaincodes/itemcc?peer=peer1&fcn=queryItemsByItemOwner&args=%5B%22%22%2c%22"+business+"%22%5D",
		dataType:"text",
		beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
		success:function(data){

		},

	})
}
function getAll(page,topic){
	$("#itemlist").empty();
	$.ajax({
			type:"get",
			url:"/getallinfo/channels/itemchannel/chaincodes/itemcc?peer=peer1&topic="+topic+"&page="+page,
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				console.log(data);
				var st=data.indexOf("no record in the page");
				if(st==-1){
					var jsdata=JSON.parse(data);
					tempsavelog=jsdata;
					for(var i=0;i<jsdata.length-1;i++){
						var record=jsdata[i].Record;
						var $trs=$("<tr><td>"+record.name+"</td><td>"+record.price+"</td><td>"+record.property+"</td><td>"+record.owner+"</td><td><a style='margin:0 5px;' class='buy'>购买</a></td></tr>");
						$("#itemlist").append($trs);
					}
					options.totalPages=jsdata[jsdata.length-1].totalpages;	
				}else{	
					var $trs=$("<tr><td colspan='3'>no record in the page</td></tr>");
					$("#itemlist").append($trs);
					options.totalPages=1;
				}
				if(page==1){
				$("#page").bootstrapPaginator(options);}
			},
			error:function(data){
			console.log(data);
			}	
		});
}
function buyItem(args){
	$.ajax({
		type:"post",
		url:"/channels/itemchannel/chaincodes/itemcc",
		data:JSON.stringify({"fcn":"transferItem","args":args}),
		dataType:"text",
		beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
		},
		success:function(data){
			if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
				alert("购买失败");
			}else{
				alert("购买成功");
			}
			getAll(1,3);
		},
		error:function(data){
			console.log(data);
		}
	})
}
