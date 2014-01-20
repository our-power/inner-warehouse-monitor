/**
 * Created by yongfengxia on 14-1-17.
 */

$(function () {

    function getPasswd(userId) {
        var req = $.ajax({
            "type": "get",
            "url": "/admin/api/get_passwd?id=" + userId,
            "dataType": "json"
        })
        req.done(function (resp) {
            if (resp.Status === "failure") {
                alertify.log(resp.Msg, "error", 5000);
                return false;
            }
            $("#modal-for-user input[name='passwd']").val(resp.Passwd);
            $("#modal-for-user input[name='re-passwd']").val(resp.Passwd);
        })
    }

    $("button#add-new-user").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-user>.modal-header>h3").text("添加新用户");

        // 清除可能上次修改操作使用modal留下的数据
        $("#modal-for-user form>#user-id").remove();
        $("#modal-for-user input[name='user-name']").val("");
        $("#modal-for-user input[name='email']").val("");
        $("#modal-for-user input[name='which-role']").val("");
        $("#modal-for-user input[name='passwd']").val("");
        $("#modal-for-user input[name='re-passwd']").val("");

        $("#modal-for-user").modal("show");
    });

    $("button#add-new-role").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-user>.modal-header>h3").text("添加新角色");

        // 清除可能上次修改操作使用modal留下的数据
        $("#modal-for-role form>#role-id").remove();
        $("label.checkbox").children("input").removeClass("checked");
        $("#modal-for-role input[name='role-name']").val("");

        $("#modal-for-role").modal("show");
    });

    $("button#modify-this-user").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-user>.modal-header>h3").text("修改用户信息");

        var parent = $(this).parent(),
            id = $.trim(parent.siblings('.id').text()),
            userName = $.trim(parent.siblings('.user-name').text()),
            email = $.trim(parent.siblings('.email').text()),
            roleId = $.trim(parent.siblings('.role-id').text());
        $("#modal-for-user form>#user-id").remove();
        $("#modal-for-user > .modal-body > form").prepend('<input type="hidden" id="user-id" name="user-id" value="' + id + '"');
        $("#modal-for-user input[name='user-name']").val(userName);
        $("#modal-for-user input[name='email']").val(email);
        $("#modal-for-user input[name='which-role']").val(roleId);

        getPasswd(id);

        $("#modal-for-user").modal("show");
    });

    $("button#modify-this-role").on("click", function (e) {
        e.preventDefault();

        $("#modal-for-user>.modal-header>h3").text("修改角色信息");

        var parent = $(this).parent(),
            role_id = $.trim(parent.siblings('.id').text()),
            role_type = $.trim(parent.siblings('.role-type').text()),
            permissions = $.trim(parent.siblings('.permission').text());
        $("#modal-for-role form>#role-id").remove();
        $("#modal-for-role form").prepend('<input type="hidden" id="role-id" name="role-id" value="' + role_id + '"');
        $("#modal-for-role input[name='role-name']").val(role_type);

        $("label.checkbox").children("input").removeClass("checked");
        var permission_list = permissions.split("|"),
            permission_num = permission_list.length;
        for(var index=0; index < permission_num; index++){
            $("#"+permission_list[index]+"-permission").addClass("checked");
        }

        $("#modal-for-role").modal("show");
    });

    $("button#del-this-user").on("click", function (e) {
        e.preventDefault();

    });


    $("button#del-this-role").on("click", function (e) {
        e.preventDefault();

    });

    $("button#save-user").on("click", function (e) {
        e.preventDefault();
    });

    $("button#save-role").on("click", function (e) {
        e.preventDefault();
    });

    $("label.checkbox").on("click", function(e){
        e.preventDefault();

        $(this).children("input").toggleClass("checked");
    });
});
