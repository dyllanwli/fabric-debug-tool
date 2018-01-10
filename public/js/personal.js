function getUserInfo(args) {
	var password;
	$.ajax({
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
			alert(password)
		},
		error: function (data) {
		}
	});
	// password = "get-fail"
	return password;
}

$(function () {
	$("#username_").html(sessionStorage.username);
	$("#org_").html(sessionStorage.userorg);
	var pd = getUserInfo(sessionStorage.username)
	$("#password_").html(pd);
})

