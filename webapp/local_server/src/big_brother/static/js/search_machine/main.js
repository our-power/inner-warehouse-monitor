$(function () {
    $(".hasDatePicker").datepicker({
        dateFormat: "yy-mm-dd",
        minDate: -30,
        maxDate: new Date()
    });
    $("#query_date_input").datepicker("setDate", new Date());

    $("#data-tab a").on("click", function (e) {
        e.preventDefault();
        $(this).tab('show');
    });

    $("#search").on("click", function (e) {
        e.preventDefault();
        var machineList = $.trim($("#machine-list").val());
        var queryDate = $.trim($("#query_date_input").val());
        if (machineList === "" || queryDate === "") {
            alertify.alert("机器列表 和 查询日期 均不能为空！");
            return false;
        }else{
            var machineArray = machineList.split("\n"),
                machineNum = machineArray.length;

            var ipList = new Array(),
                hostNameList = new Array(),
                hardwareAddrList = new Array();

            var machine;
            var ipPattern = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,
                hardwareAddrPattern = /^([0-9a-fA-F]{2})(([/\s:-][0-9a-fA-F]{2}){5})$ /;
            for(var index=0; index < machineNum; index++){
                machine = machineArray[index];
                if(ipPattern.test(machine)){
                    ipList.push(machine);
                }else if(hardwareAddrPattern.test(machine)){
                    hardwareAddrList.push(machine);
                }else{
                    hostNameList.push(machine);
                }
            }

            var urlParam = new Array();
            if(ipList.length > 0){
                urlParam.push("iplist=" + ipList.join(","));
            }
            if(hostNameList.length > 0){
                urlParam.push("hostnamelist=" + hostNameList.join(","));
            }
            if(hardwareAddrList.length > 0 ){
                urlParam.push("hardwareaddrlist="+hardwareAddrList.join(","))
            }
        }
    });
});