/**
 * Created by yongfengxia on 14-1-17.
 */
$(function(){
    $("button#add-new-user").on("click", function(e){
        e.preventDefault();

        $("#modal-for-role").modal("show");
    });

    $("button#add-new-role").on("click", function(e){
        e.preventDefault();

        $("#modal-for-role").modal("show");
    });

    $("button#modify-this-user").on("click", function(e){
        e.preventDefault();

    });

    $("button#del-this-user").on("click", function(e){
        e.preventDefault();

    });

    $("button#modify-this-role").on("click", function(e){
        e.preventDefault();

    });

    $("button#del-this-role").on("click", function(e){
        e.preventDefault();

    });

    $("button#save-user").on("click", function(e){
        e.preventDefault();
    });

    $("button#save-role").on("click", function(e){
        e.preventDefault();
    });
});
