	var tempsavelog;
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
		//sessionStorage.token="<%=token%>";
		//sessionStorage.user="<%=user%>";
		//getAllLog("queryLogsByUser",sessionStorage.username);
        $("#page").bootstrapPaginator(options);
		getAll(1,1);
		$("#sendbtn").click(function(){
			$("#front").hide();
			$("#sendlogdiv").show();
		});
		$("#logfilediv").click(function(){
			$("#logfile").click();
		});
		$("#logfile").change(function(){
			$("#logname2").val($(this).val());
		});
		$("#sendlog").click(function(){
			var filedata=new FormData();
			var file=$("#logfile");
			var logname=$("#logname2").val();
			var logg=$("#log2").val();
			if(file.val()==""){
				if(logname!=""&&logg!=""){
					var arg=[logname,logg];
					$.ajax({
						type:"post",
						url:"/channels/logchannel/chaincodes/logcc",
						data:JSON.stringify({"args":arg,"fcn":"uploadLog"}),
						dataType:"text",
						beforeSend:function(xhr){
							xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
							xhr.setRequestHeader("content-type","application/json");
						},
						success:function(data){
							//console.log(data);
							if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
								alert("上传失败");
							}else{
								alert("上传成功");
							}
							getAll(1,1);
							$("#front").show();
							$("#sendlogdiv").hide();
							$("#logfile").val("");
							$("#logname2").val("");
							$("#log2").val("");
						}
					});
				}else{
					alert("请填写完整log信息后上传");
				}
			}else{
				filedata.append("logfile",file[0].files[0]);
				uplogfile(filedata);
			}
		});
		$("#cancel").click(function(){
			$("#front").show();
			$("#sendlogdiv").hide();
			$("#logfile").val("");
			$("#logname2").val("");
			$("#log2").val("");
		});
		$("#dellog").click(function(){
			$("#downtips").show();
			$("#downtips2").show();
		});
		$("#seltotal").click(function(){
			if($(this).is(':checked')){
				$(".tlog").prop("checked",true);
			}else{
				$(".tlog").prop("checked",false);
			}

		});
		$("#dellog3").click(function(){
			var tlogs=$(".tlog");
			for(var i=0;i<tlogs.length;i++){
				if($(tlogs[i]).is(":checked")){
					var name=$(tlogs[i]).parent().parent().children("td").eq(1).html();
					var owner=$(tlogs[i]).parent().parent().children("td").eq(4).html();
					if(owner==sessionStorage.username){
						delOneLog(name);
					}
				}
			}
			$("#downtips").hide();
			$("#downtips2").hide();
		});
		$("#canceldel").click(function(){
			$("#downtips").hide();
			$("#downtips2").hide();
		});
		$("#loglist").on("click",".delone",function(){
			$(this).parent().parent().children(":first").children().attr("checked",true);
			$("#downtips").show();
			$("#downtips2").show();
		});
		$("#loglist").on("click",".look",function(){
			var infoname=$(this).parent().parent().children("td").eq(1).html();
			var infoinfo=$(this).parent().parent().children("td").eq(2).html();
			var logtime=$(this).parent().parent().children("td").eq(3).html();
			$("#infologname").html(infoname);
			$("#infologinfo").html(infoinfo);
			$("#logtime").html("时间:"+logtime);
			checklog(infoname);
		});
		$("#closeinfo").click(function(){
			$("#oneloginfo").hide();
			$("#oneloginfoback").hide();
		});
		$("#getChainCode").click(function(){
			var lname=$("#logname").val();
			var useraccount=$("#username").val();
			var fcn="";
			var args="";
			if(lname!=null&&lname!==""){
				fcn="readLog";
				args=lname;			
			}else if(useraccount!=null&&useraccount!=""){
				fcn="queryLogsByUser";
				args=useraccount;
			}else{
				getAll(1,2);
				return;			
			}
			getAllLog(fcn,args);
		})
		$("#download").click(function(){
			DownLoad();
		});
	});
	function delOneLog(name){
		$.ajax({
				type:"post",
				url:"/channels/logchannel/chaincodes/logcc",
				data:JSON.stringify({"args":[name],"fcn":"deleteLog"}),
				dataType:"text",
				beforeSend:function(xhr){
					xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
					xhr.setRequestHeader("content-type","application/json");
				},
				success:function(data){
					//console.log(data);
					if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
						alert("删除失败");
					}else{
						
					}
					getAll(1,1);
					//getAllLog();
					//$("#downtips").hide();
				}
			})
	}
function getAllLog(fcn,args){
$("#loglist").empty();
		$.ajax({
			type:"get",
			url:"/channels/logchannel/chaincodes/logcc?peer=peer1&fcn="+fcn+"&args=%5B%22"+args+"%22%5D",
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				//console.log(data);
				var st=data.indexOf("[");
				if(data.indexOf("no record in the page")!=-1||data.indexOf("not exist")!=-1){
					var $trs=$("<tr><td colspan='3'>no record in the page</td></tr>");
					$("#loglist").append($trs);
					options.totalPages=1;
					$("#page").bootstrapPaginator(options);
					return;
				}
				if(st!=-1){
				var end=data.lastIndexOf("]");
				var newdata=data.substring(st,end+1);
				var jsdata=JSON.parse(newdata);
				tempsavelog=jsdata;
				for(var i=0;i<jsdata.length;i++){
					var record=jsdata[i].Record;
					var $trs=$("<tr><td><input type='checkbox' class='tlog' pid='"+i+"'></td><td>"+record.name+"</td><td>"+record.logContent+"</td><td>"+record.uploadTime+"</td><td>"+record.user+"</td><td><a style='margin:0 5px;' class='look'>查看</a><a style='margin:0 5px;' class='delone'>删除</a></td></tr>");
					$("#loglist").append($trs);
				}
				}else{
					st=data.indexOf("{");
					var end=data.indexOf("}");
				var newdata=data.substring(st,end+1);
				var jsdata=JSON.parse(newdata);
					var $trs=$("<tr><td><input type='checkbox' class='tlog' pid='1'></td><td>"+jsdata.name+"</td><td>"+jsdata.logContent+"</td><td>"+jsdata.uploadTime+"</td><td>"+jsdata.user+"</td><td><a style='margin:0 5px;' class='look'>查看</a><a style='margin:0 5px;' class='delone'>删除</a></td></tr>");
					$("#loglist").append($trs);
				}
			},
			error:function(data){
			console.log(data);
			}	
		});
}
function getAll(page,topic){
	$("#loglist").empty();
	$.ajax({
			type:"get",
			url:"/getallinfo/channels/logchannel/chaincodes/logcc?peer=peer1&topic="+topic+"&page="+page,
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				//console.log(data);
                $("#loglist").empty();
				var st=data.indexOf("no record in the page");
				if(st==-1){
					var jsdata=JSON.parse(data);
					tempsavelog=jsdata;
					var lh=jsdata.length
					for(var i=0;i<lh-1;i++){
						var record=jsdata[i].Record;
						var $trs=$("<tr><td><input type='checkbox' class='tlog' pid='"+i+"'></td><td>"+record.name+"</td><td>"+record.logContent+"</td><td>"+record.uploadTime+"</td><td>"+record.user+"</td><td><a style='margin:0 5px;' class='look'>查看</a><a style='margin:0 5px;' class='delone'>删除</a></td></tr>");
						$("#loglist").append($trs);
					}
					options.totalPages=jsdata[lh-1].totalpages;	
				}else{	
					var $trs=$("<tr><td colspan='3'>no record in the page</td></tr>");
					$("#loglist").append($trs);
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
function uplogfile(file){
	$.ajax({
			type:"post",
			url:"/uplogfile",
			data:file,
			cache: false,
    		contentType: false,
    		processData: false,
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
			},
			success:function(data){
				//console.log(data);
				if(data.indexOf("Fail")!=-1||data.indexOf("Error")!=-1){
					alert("上传失败");
				}else{
					alert("上传成功");
				}
				getAll(1,1);
				$("#front").show();
				$("#sendlogdiv").hide();
                $("#logfile").val("");
				$("#logname2").val("");
				$("#log2").val("");
			}
	});
}
function checklog(logname){
	$.ajax({
			type:"get",
			url:"/checklog?logname="+logname,
			dataType:"text",
			beforeSend:function(xhr){
				xhr.setRequestHeader("authorization","Bearer "+sessionStorage.token);
				xhr.setRequestHeader("content-type","application/json");
			},
			success:function(data){
				//console.log(data);
				if(data!='no log'){
					$("#logpath").text(data);
					$("#download").show();
				}else{
					$("#logpath").text("");
					$("#download").hide();
				}
				$("#oneloginfoback").show();
				$("#oneloginfo").show();				
			},
			error:function(data){
				console.log(data);
				$("#oneloginfoback").show();
				$("#oneloginfo").show();
			}	
		});
}
function DownLoad() { 
	var path=$("#logpath").text();
    var form = $("<form>");   //定义一个form表单
    form.attr('style', 'display:none');   //在form表单中添加查询参数
    form.attr('target', '');
    form.attr('method', 'post');
    form.attr('action', '/users/downloadlogfile');

    var input1 = $('<input>');
    input1.attr('type', 'hidden');
    input1.attr('name', 'logpath');
    input1.attr('value', path);
    $('body').append(form);  //将表单放置在web中 
    form.append(input1);   //将查询参数控件提交到表单上
    form.submit();
 }
