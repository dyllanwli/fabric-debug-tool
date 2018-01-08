function sub(){
	var name=$("#username").val();
	var pwd=$("#password").val();
	if(name==""){
		$("#nameerr").text("please enter the name");
		return false;
	}
	if(pwd==""){
		$("#pwderr").text("please enter the password");
		return false;
	}
	$("#nameerr").text("");
	$("#pwderr").text("");
	return true;
}
