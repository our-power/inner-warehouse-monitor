/**
 * Created by yongfengxia on 14-1-17.
 */

$(function () {

    // 点击“添加新用户”按钮
    $("button#add-new-user").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-user>.modal-header>h3").text("添加新用户");

        // 清除可能上次修改操作使用modal留下的数据
        $("#modal-for-user form>#user-id").remove();
        $("#new-passwd-controlgroup").remove();
        $("#re-new-passwd-controlgroup").remove();

        $("#modal-for-user input[name='user-name']").val("");
        $("#modal-for-user input[name='email']").val("");
        $("#modal-for-user input[name='which-role']").val("");

        var newPasswdControlGroup = $("<div></div>", {
            "class": "control-group",
            "id": "new-passwd-controlgroup"
        });
        var newPasswdControls = $("<div></div>", {
            "class": "controls"
        });
        newPasswdControls.append('<input type="password" id="new-passwd" name="new-passwd" />');
        newPasswdControlGroup.append('<label class="control-label" for="new-passwd">密码</label>').append(newPasswdControls);
        $("#modal-for-user form").append(newPasswdControlGroup);

        var reNewPasswdControlGroup = $("<div></div>", {
            "class": "control-group",
            "id": "re-new-passwd-controlgroup"
        });
        var reNewPasswdControls = $("<div></div>", {
            "class": "controls"
        });
        reNewPasswdControls.append('<input type="password" id="re-new-passwd" name="re-new-passwd" />');
        reNewPasswdControlGroup.append('<label class="control-label" for="re-new-passwd">密码确认</label>').append(reNewPasswdControls);
        $("#modal-for-user form").append(reNewPasswdControlGroup);

        $("#modal-for-user").modal("show");
    });

    // 点击“添加新角色”按钮
    $("button#add-new-role").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-role>.modal-header>h3").text("添加新角色");

        // 清除可能上次修改操作使用modal留下的数据
        $("#modal-for-role form>#role-id").remove();
        $("label.checkbox").find("input").removeAttr("checked");
        $("label.checkbox").removeClass("checked");
        $("#modal-for-role input[name='role-name']").val("");

        $("#modal-for-role").modal("show");
    });

    // 点击用户“修改”按钮
    $("button#modify-this-user").on("click", function (e) {
        e.preventDefault();
        $("#modal-for-user>.modal-header>h3").text("修改用户信息");

        var parent = $(this).parent(),
            id = $.trim(parent.siblings('.id').text()),
            userName = $.trim(parent.siblings('.user-name').text()),
            email = $.trim(parent.siblings('.email').text()),
            roleId = $.trim(parent.siblings('.role-id').text());
        $("#modal-for-user form>#user-id").remove();
        $("#new-passwd-controlgroup").remove();
        $("#re-new-passwd-controlgroup").remove();

        $("#modal-for-user form").prepend('<input type="hidden" id="user-id" name="user-id" value="' + id + '"');
        $("#modal-for-user input[name='user-name']").val(userName);
        $("#modal-for-user input[name='email']").val(email);
        $("#modal-for-user input[name='which-role']").val(roleId);

        $("#modal-for-user").modal("show");
    });

    // 点击角色“修改”按钮
    $("button#modify-this-role").on("click", function (e) {
        e.preventDefault();

        $("#modal-for-role>.modal-header>h3").text("修改角色信息");

        var parent = $(this).parent(),
            role_id = $.trim(parent.siblings('.id').text()),
            role_type = $.trim(parent.siblings('.role-type').text()),
            permissions = $.trim(parent.siblings('.permission').text());
        $("#modal-for-role form>#role-id").remove();
        $("#modal-for-role form").prepend('<input type="hidden" id="role-id" name="role-id" value="' + role_id + '"');
        $("#modal-for-role input[name='role-name']").val(role_type);

        $("label.checkbox").removeClass("checked");
        $("label.checkbox").find("input").removeAttr("checked");
        var permission_list = permissions.split("|"),
            permission_num = permission_list.length;
        for (var index = 0; index < permission_num; index++) {
            var target = $("#" + permission_list[index] + "-permission");
            target.attr("checked", "true");
            target.parent("label.checkbox").addClass("checked");
        }

        $("#modal-for-role").modal("show");
    });

    // 删除某用户
    $("button#del-this-user").on("click", function (e) {
        e.preventDefault();

        var userId = $.trim($(this).parent().siblings(".id").text());

        alertify.confirm('确定删除该用户吗？', function (e) {
            if (e) {
                var req = $.ajax({
                    "type": "post",
                    "url": "/admin/api/del_user",
                    "data": {
                        "user_id": userId
                    },
                    "dataType": "json"
                });
                req.done(function (resp) {
                    if (resp.Status === "failure") {
                        alertify.log(resp.Msg, "error", 5000);
                        return false;
                    }
                    setTimeout("window.location.href='/admin'", 1500)
                });
            }
        })
    });

    // 删除某角色
    $("button#del-this-role").on("click", function (e) {
        e.preventDefault();

        var roleId = $.trim($(this).parent().siblings(".id").text());
        alertify.confirm('确定删除该角色吗？', function (e) {
            if (e) {
                var req = $.ajax({
                    "type": "post",
                    "url": "/admin/del_role",
                    "data": {
                        "role_id": roleId
                    },
                    "dataType": "json"
                });
                req.done(function (resp) {
                    if (resp.Status === "failure") {
                        alertify.log(resp.Msg, "error", 5000);
                        return false;
                    }
                    setTimeout("window.location.href='/admin'", 1500);
                });
            }
        });

    });

    // 提交修改或新增用户信息
    $("button#save-user").on("click", function (e) {
        e.preventDefault();
        if ($("#new-passwd-controlgroup").length == 0) {
            var userId = $.trim($("#user-id").val()),
                userName = $.trim($("#user-name").val()),
                email = $.trim($("#email").val()),
                roleType = $.trim($("#which-role").val());
            var req = $.ajax({
                "type": "post",
                "url": "/admin/api/modify_user",
                "data": {
                    "user_id": userId,
                    "user_name": userName,
                    "email": email,
                    "role_type": roleType
                },
                "dataType": "json"
            });
            $("#modal-for-user").modal("hide");
            req.done(function (resp) {
                if (resp.Status === "failure") {
                    alertify.log(resp.Msg, "error", 5000);
                    return false;
                }
                setTimeout("window.location.href='/admin'", 1500);
            });
        } else {
            var userName = $.trim($("#userName").val()),
                newPasswd = $.trim($("#new-passwd").val()),
                reNewPasswd = $.trim($("#re-new-passwd").val()),
                userName = $.trim($("#user-name").val()),
                email = $.trim($("#email").val()),
                roleType = $.trim($("#which-role").val());
            if (newPasswd == "" || newPasswd != reNewPasswd) {
                $("#new-passwd-error").remove();
                $("#re-new-passwd").after('<span class="label label-important" id="new-passwd-error">密码不一致！</span>');
            } else {
                var req = $.ajax({
                    "type": "post",
                    "url": "/admin/api/add_user",
                    "data": {
                        "user_name": userName,
                        "passwd": newPasswd,
                        "email": email,
                        "role_type": roleType
                    },
                    "dataType": "json"
                });
                $("#modal-for-user").modal("hide");
                req.done(function (resp) {
                    if (resp.Status === "failure") {
                        alertify.log(resp.Msg, "error", 5000);
                        return false;
                    }
                    setTimeout("window.location.href='/admin'", 1500);
                });
            }
        }

    });

    // 提交修改或新增角色信息
    $("button#save-role").on("click", function (e) {
        e.preventDefault();

        $("#modal-for-role").modal("hide");

        var roleId = $.trim($("#role-id").val());
        var roleName = $.trim($("#role-name").val());
        var checkboxLabels = $("label.checkbox"),
            checkboxLabel_num = checkboxLabels.length;
        var checkedPermissionList = new Array();
        for (var index = 0; index < checkboxLabel_num; index++) {
            var permissionElement = $(checkboxLabels[index]).find("input");
            if (permissionElement.prop("checked")) {
                checkedPermissionList.push(permissionElement.val());
            }
        }
        var permissions = checkedPermissionList.join("|");
        if (roleId === "") {
            var req = $.ajax({
                "type": "post",
                "url": "/admin/api/add_role",
                "data": {
                    "role_name": roleName,
                    "permissions": permissions
                },
                "dataType": "json"
            });
            req.done(function (resp) {
                if (resp.Status === "failure") {
                    alertify.log(resp.Msg, "error", 5000);
                    return false;
                }
                setTimeout("window.location.href='/admin'", 1500);
            });
        } else {
            var req = $.ajax({
                "type": "post",
                "url": "/admin/api/modify_role",
                "data": {
                    "role_id": roleId,
                    "role_name": roleName,
                    "permissions": permissions
                },
                "dataType": "json"
            });
            req.done(function (resp) {
                if (resp.Status === "failure") {
                    alertify.log(resp.Msg, "error", 5000);
                    return false;
                }
                setTimeout("window.location.href='/admin'", 1500);
            });
        }

    });

    $("label.checkbox").on("click", function (e) {

        // 这里不能加e.preventDefault();

        if ($(this).hasClass("checked")) {
            $(this).find("input").removeAttr("checked");
        } else {
            $(this).find("input").attr("checked", "true");
        }
        $(this).toggleClass("checked");
    });

    // 修改用户密码
    $("button#change-passwd").on("click", function (e) {
        e.preventDefault();

        var parent = $(this).parent(),
            userId = $.trim(parent.siblings('.id').text()),
            userName = $.trim(parent.siblings('.user-name').text());

        $("#userid-change-passwd").remove();
        $("#modal-for-change-passwd form").prepend('<input type="hidden" id="userid-change-passwd" name="userid-change-passwd" val="' + userId + '"');
        $("#username-to-change-passwd").text(userName);

        $("#modal-for-change-passwd").modal("show");
    });

    $("#save-new-passwd").on("click", function (e) {
        e.preventDefault();

        var passwd = $.trim($("#passwd").val()),
            rePasswd = $.trim($("#re-passwd").val());
        if (passwd != rePasswd) {
            $("#passwd-error").remove();
            $("#re-passwd").after('<span class="label label-important" id="passwd-error">密码不一致！</span>');
            return false;
        }
        var userId = $.trim($("#userId-change-passwd").val());
        var req = $.ajax({
            "type": "post",
            "url": "/admin/api/change_passwd",
            "data": {
                "user_id": userId,
                "new_passwd": passwd
            },
            "dataType": "json"
        });
        req.done(function (resp) {
            if (resp.Status === "failure") {
                alertify.log(resp.Msg, "error", 5000);
                $("#passwd").val("");
                $("#re-passwd").val("");
                return false
            }
            $("#modal-for-change-passwd").modal("hide");
            alertify.log("密码修改成功！", "success", 2000);
        });
    });
});
