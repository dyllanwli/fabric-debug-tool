function getUserInfo(args) {
	var re = {};
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
				re= JSON.parse(data)
			} else {
				re.password = "get-failure";
				re.channelsList ="get-failure"
			}
		},
		error: function (data) {
			alert("got error")
		}
	});
	return re;
}

$(function () {
	$("#username_").html(sessionStorage.username);
	$("#org_").html(sessionStorage.userorg);
	var us = getUserInfo(sessionStorage.username)
	$("#password_").html(us.password);
	var cl = getUserInfo(sessionStorage.username).channelsList
	$("#channel_list_").html(cl)
})

