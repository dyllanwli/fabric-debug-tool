function sub(){
	var name=$("#username").val();
	var pwd=$("#password").val();
	var org=$("#orgname").val();
	if(name==""){
		$("#nameerr").text("please enter the name");
		return false;
	}
	if(pwd==""){
		$("#pwderr").text("please enter the password");
		return false;
	}
	if(org==""){
		$("#orgerr").text("please enter the org");
		return false;
	}
	$("#nameerr").text("");
	$("#pwderr").text("");
	$("#orgerr").text("");
	return true;
}
