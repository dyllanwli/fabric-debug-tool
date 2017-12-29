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
					getAll(page,1);
                }
            }
$(function(){
    $("#page").bootstrapPaginator(options);
	getAll(1,4);
	$("#getitem1").click(function(){
		var name=$("#itemname").val();
		if(name!=null&&name!=""){
			$("#itemlist").empty();
			getOneItem(name);
		}else{
			getAll(1,4);
		}
	});
	$("#upitem").click(function(){
		$("#front").hide();
		$("#upitemdiv").show();
	});
	$("#upitem2").click(function(){
		var iname=$("#itemname2").val();
		var iproperty=$("#itemproperty2").val();
		var iprice=$("#itemprice2").val();
		var args=[iname,iproperty,iprice];
		upitem(args);
	});
	$("#seltotal").click(function(){
		if($(this).is(':checked')){
			$(".titem").prop("checked",true);
		}else{
			$(".titem").prop("checked",false);
		}

	});
	$("#cancelup").click(function(){
		$("#front").show();
		$("#upitemdiv").hide();
	});
	$("#downitem").click(function(){
		$("#downtips").show();
		$("#downtips2").show();
	});
	$("#downitem2").click(function(){
		var titems=$(".titem");
			for(var i=0;i<titems.length;i++){
				if($(titems[i]).is(":checked")){
					var name=$(titems[i]).parent().parent().children("td").eq(1).html();
					var property=$(titems[i]).parent().parent().children("td").eq(2).html();
					var args=[name,property];
					delOneItem(args);
				}
			}
			$("#downtips").hide();
			$("#downtips2").hide();
			//getAllItem("queryLogsByUser",sessionStorage.username);
	});
	$("#canceldown").click(function(){
		$("#downtips").hide();
		$("#downtips2").hide();
	});
	$("#itemlist").on("click",".downitem3",function(){
		$(this).parent().parent().children(":first").children().attr("checked",true);
		$("#downtips").show();
		$("#downtips2").show();
		});
})
function getOneItem(name){
	$.ajax({
		type:"get",
		url:"/channels/itemchannel/chaincodes/itemcc?peer=peer1&fcn=queryItemsByItemOwner&args=%5B%22"+name+"%22%2c%22"+sessionStorage.username+"%22%5D",
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
					var $trs=$("<tr><td><input type='checkbox' class='titem' pid='"+i+"'></td><td>"+record.name+"</td><td>"+record.property+"</td><td>"+record.price+"</td><td><a style='margin:0 5px;' class='downitem3'>下链</a></td></tr>");
					$("#itemlist").append($trs);
				}
		},
		error:function(data){
			console.log(data);
		}
	});
}
function upitem(args){
	$.ajax({
		type:"post",
		url:"/channels/itemchannel/chaincodes/itemcc",
		data:JSON.stringify({"fcn":"initItem","args":args}),
		dataType:"text",
		beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
		success:function(data){
			if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
					alert("上传失败");
				}else{
					alert("上传成功");
				}
				getAll(1,4);
				$("#front").show();
				$("#upitemdiv").hide();
		},
		error:function(data){
			alert("上传失败");
			$("#front").show();
			$("#upitemdiv").hide();
		}

	})
}
function delOneItem(args){
	$.ajax({
		type:"post",
		url:"/channels/itemchannel/chaincodes/itemcc",
		data:JSON.stringify({"fcn":"deleteItem","args":args}),
		dataType:"text",
		beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
		success:function(data){
			if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
						alert("删除失败");
					}else{
						
					}
					getAll(1,4);
		},
		error:function(data){
			console.log(data);
		}
	})
}
function getOwnerItem(){
	$.ajax({
		type:"get",
		url:"/channels/itemchannel/chaincodes/itemcc?peer=peer1&fcn=queryItemsByItemOwner&args=%5B%22%22%2c%22%22%2c%22"+sessionStorage.username+"%22%5D",
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
				//console.log(data);
				var st=data.indexOf("no record in the page");
				if(st==-1){
					var jsdata=JSON.parse(data);
					tempsavelog=jsdata;
					for(var i=0;i<jsdata.length-1;i++){
						var record=jsdata[i].Record;
						var $trs=$("<tr><td><input type='checkbox' class='titem' pid='"+i+"'></td><td>"+record.name+"</td><td>"+record.property+"</td><td>"+record.price+"</td><td><a style='margin:0 5px;' class='downitem3'>下链</a></td></tr>");
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
function clearNoNum(obj){
   obj.value = obj.value.replace(/[^\d]/g,"");  //清除“数字”以外的字符
   obj.value = obj.value.replace(/^0/g,"");  //验证第一个字符不是0
}
