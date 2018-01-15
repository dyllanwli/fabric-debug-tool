function getUserInfo(args) {
	var password;
	var re;
	$.ajax({
		async :false,
		type: "get",
		url: "/users/password?username="+args,
		dataType: "text",
		beforeSend: function (xhr) {
			xhr.setRequestHeader("authorization", "Bearer " + sessionStorage.token);
			xhr.setRequestHeader("content-type", "application/json");
		},
		success: function (data) {
			if (data.indexOf("Error") == -1) {
				password = data;
			} else {
				password = "get-fail";
			}
		},
		error: function (data) {
			alert("got error")
		}
	});
	re["password"] = password
	return re;
}

$(function () {
	$("#username_").html(sessionStorage.username);
	$("#org_").html(sessionStorage.userorg);
	var pd = getUserInfo(sessionStorage.username).password
	$("#password_").html(pd);
	var cl = getUserInfo(sessionStorage.username).password
	$("#channel_list_").html(cl)
})

