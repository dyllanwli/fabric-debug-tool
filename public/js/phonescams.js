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
	



})

function getOnePhone(phonenumber){

}
function markPhone(phonenumber){
}
function cancelMark(phonenumber){
}
function delOnePhone(phonenumber){
}
function getAllPhone(){
}
