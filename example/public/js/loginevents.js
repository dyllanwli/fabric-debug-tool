function sub(){
	var name=$("#username").val();
	var pwd=$("#password").val();
	if(name==""){
		$("#nameerr").text("请输入账号");
		return false;
	}
	if(pwd==""){
		$("#pwderr").text("请输入密码");
		return false;
	}
	$("#nameerr").text("");
	$("#pwderr").text("");
	return true;
}
