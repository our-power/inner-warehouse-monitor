$(function(){
    var fileTypeAllowed = new Array("zip");

    $("button:submit").on("click", function(e){
        var filePath = $("input[name='uploadfile']").val(),
            pos = filePath.lastIndexOf('.'),
            fileType = filePath.substr(pos + 1),
            fileType = fileType.toLowerCase();
        if(fileTypeAllowed.indexOf(fileType) == -1){
            alert("请上传zip格式文件");
            e.preventDefault();
            return false;
        }
    });
});
