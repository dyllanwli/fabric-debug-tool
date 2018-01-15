$(function(){
    $("#username").focus();
    $("#username").blur(function(){
      var name=$(this).val();
      var rule=/^\w{3,15}$/;
      if(!rule.test(name)||name==""){
        $("#nameerr").text("请正确输入用户名");
        return;
      }else{
        $("#nameerr").text("");
      }
      $.ajax(
      {
        type:"post",
        url:"/users/regvali",
        data:{username:name},
        dataType:"json",
        success:function(data){
          $("#nameerr").text(data.err);
        }
      }
      );
    });
    $("#phonenumber").blur(function(){
      var phone=$(this).val();
      var rule=/^1\d{10}$/;
      if(!rule.test(phone)||phone==""){
        $("#phoneerr").text("请正确输入手机号");
        return;
      }else{
        $("#phoneerr").text("");
      }
    });
    $("#password").blur(function(){
      var pwd=$(this).val();
      var rule=/^[a-z0-9]{6,15}$/;
      if(!rule.test(pwd)||pwd==""){
        $("#pwderr").text("请正确输入密码");
        return;
      }else{
        $("#pwderr").text("");
        return;
      }
    });
    $("#dpassword").blur(function(){
      var dpwd=$(this).val();
      var pwd=$("#password").val();
      if(dpwd!=pwd){
        $("#dpwderr").text("两次密码输入不一致");
        return;
      }else{
        $("#dpwderr").text("");
        return;
      }
    });
    $("#confirm").click(function(){
      var err=$("#nameerr").text()+$("#phoneerr").text()+$("pwderr").text()+$("#dpwderr").text();
      if(err==""){
        return true;
      }else{
        alert("请完整填写正确信息后提交！");
        return false;
      }
    });
});
