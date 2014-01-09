$(function(){

    alertify.set({
        labels: {
            ok     : "是",
            cancel : "否"
        }
    });

    $('.del-machine').on('click', function(e){

        var sure_to_del = false;

        alertify.confirm('确定删除该机器吗？', function(e){
            if (e) {
                sure_to_del = true;
            }
        })
        if(! sure_to_del) {
            return false;
        }
        var id = $.trim($(this).parent().siblings('.id').text()),
        var req = $.ajax({
            "type": "post",
            "url": "/manage/del_machine",
            "data": {
                "id": id
            },
            "dataType": "json"
        });
        req.done(function(resp){
            if(resp.status === "failure") {
                alertify.log(resp.msg, "error", 5000);
                return false;
            }
            alertify.log("已成功删除！", "success", 1000);
            setTimeout("window.location.href='/manage/list_machine'", 1500)
        });
    });
});