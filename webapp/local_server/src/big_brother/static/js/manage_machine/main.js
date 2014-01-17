$(function() {

    alertify.set({
        labels: {
            ok: "是",
            cancel: "否"
        }
    });

    $('.del-machine').on('click', function(e) {

        e.preventDefault();
        var id = $.trim($(this).parent().siblings('.id').text());
        var mac = $.trim($(this).parent().siblings(".hardware-addr").text());
        alertify.confirm('确定删除该机器吗？', function(e) {
            if (e) {
                var req = $.ajax({
                    "type": "post",
                    "url": "/manage/del_machine",
                    "data": {
                        "id": id,
                        "mac": mac
                    },
                    "dataType": "json"
                });
                req.done(function(resp) {
                    if (resp.Status === "failure") {
                        alertify.log(resp.Msg, "error", 5000);
                        return false;
                    }
                    alertify.log("已成功删除！", "success", 1000);
                    setTimeout("window.location.href='/manage/list_machine'", 1500)
                });
            } else {
                return false;
            }
        })
    });
});