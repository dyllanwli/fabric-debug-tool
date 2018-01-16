function createChannel() {
    channelName = $("#channelName").val()
    channelConfigPath = $("#channelConfigPath").val()
    $.ajax({
        // async: false,
        type: "post",
        url: "/channels",
        data: JSON.stringify({
            channelName: channelName,
            channelConfigPath: channelConfigPath
        }),
        // headers: {
        //     "Authorization": "Bearer " + sessionStorage.token,
        //     "content-type": "application/json"

        // },
        dataType: "json",
        beforeSend: function (xhr) {
            xhr.setRequestHeader("authorization", "Bearer " + sessionStorage.token);
            xhr.setRequestHeader("content-type ", "application/json");
        },
        success: function (data) {
            var response = xhr.responseText
            var ele = document.getElementById("resultArea")
            if (response == "{}") {
                alert("Check if channel are existing")
                return
            }
            ele.value += response + "\n\n"
        },
        error: function (data) {
            var response = xhr.responseText
            alert(response);
            var ele = document.getElementById("resultArea")
            ele.value += response + "\n\n "
        }
    });
}
$("#createChannel").click(function () {
    createChannel();
});