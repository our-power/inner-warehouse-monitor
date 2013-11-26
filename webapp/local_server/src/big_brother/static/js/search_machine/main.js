$(function () {

    function plotHighCharts(data) {
        new Highcharts.Chart({
            chart: {
                renderTo: data.container,
                zoomType: 'x',
                spacingRight: 20,
                width: $.browser.msie ? 900 : null,
                height: $.browser.msie ? 400 : null
            },
            credits: {
                enabled: false
            },
            title: {
                text: data.chartTitle
            },
            xAxis: {
                type: 'datetime',
                maxZoom: 5 * 60 * 1000, // five minutes
                title: {
                    text: null
                }
            },
            yAxis: {
                title: {
                    text: data.yAxisTitle
                },
                min: -5,
                startOnTick: false
            },
            tooltip: {
                shared: true
            },
            legend: {
                enabled: true
            },
            plotOptions: {
                area: {
                    fillColor: {
                        linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1},
                        stops: [
                            [0, Highcharts.getOptions().colors[0]],
                            [1, Highcharts.Color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                        ]
                    },
                    lineWidth: 1,
                    marker: {
                        enabled: false
                    },
                    shadow: false,
                    states: {
                        hover: {
                            lineWidth: 1
                        }
                    },
                    threshold: null
                }
            },
            series: data.series
        });
    }

    function newIEDate(dateStr){
        // dateStr的格式：xxxx-xx-xx
        var dateParts = dateStr.split("-");
        return new Date(parseInt(dateParts[0]), parseInt(dateParts[1]), parseInt(dateParts[2]));
    }
    function getMachineIndicatorData(jqObj) {
        var queryDate = $("#query_date_input").val();
        var objDate;
        if($.browser.msie){
            objDate = newIEDate(queryDate);
        }else{
            objDate = new Date(queryDate);
        }

        var href = jqObj.find("a").attr("href");
        var indicator = href.substr(1)
        if (indicator === "machine_list") {
            $("#form-machine-list").show();
            $("#form-date").hide();
            return false
        } else {
            $("#form-machine-list").hide();
            $("#form-date").show();
        }
        var title = "";
        var xAxisTitle = "";
        var yAxisTitle = "";
        var seriesName = "";

        var netflow_packets_title = "";
        var netflow_bytes_title = "";
        var netflow_xAxisTitle = "";
        var netflow_packets_yAxisTitle = "";
        var netflow_bytes_yAxisTitle = "";
        var trueIndicator = "";
        switch (indicator) {
            case "cpu_view":
                title = queryDate + " , CPU使用率(%)";
                xAxisTitle = "时间";
                yAxisTitle = "使用率(%)";
                seriesName = "CPU使用率(%)";
                trueIndicator = "cpu_usage";
                break;
            case "memory_view":
                title = queryDate + " , 内存使用量(MB)";
                xAxisTitle = "时间";
                yAxisTitle = "使用量(MB)";
                seriesName = "使用量(MB)";
                trueIndicator = "mem_usage";
                break;
            case "netflow_view":
                netflow_packets_title = queryDate + " , 网络包量(个)";
                netflow_bytes_title = queryDate + " , 网络流量(byte)";
                netflow_xAxisTitle = "时间";
                netflow_packets_yAxisTitle = "包量(个)";
                netflow_bytes_yAxisTitle = "流量(byte)";
                trueIndicator = "net_flow";
                break
        }
        $(href).empty();
        var hardwareAddrList = $("#machine_list .hardware-addr");
        var machineNum = hardwareAddrList.length;
        for (var index = 0; index < machineNum; index++) {
            var hardwareAddr = $.trim($(hardwareAddrList[index]).text());
            if (hardwareAddr != "") {
                var machineInChartTitle = $(hardwareAddrList[index]).siblings(".target-machine").text();
                var machineStatus = $(hardwareAddrList[index]).siblings(".machine-status").text();
                var req = $.ajax({
                    "async": false,
                    "type": "get",
                    "url": "/api/get_machine_indicator_data?hardware_addr=" + hardwareAddr + "&date=" + queryDate + "&indicator=" + trueIndicator,
                    "dataType": "json"
                });
                req.done(function (resp) {
                        var date, timeIndexNum;
                        // 对于非“正常运行中”的机器，补充最后收到数据的时间点之后时间点的数据
                        date = new Date();
                        var sameYear = objDate.getFullYear() === date.getFullYear();
                        var sameMonth;
                        if($.browser.msie){
                            sameMonth = objDate.getMonth() === date.getMonth()+1;
                        }else{
                            sameMonth = objDate.getMonth() === date.getMonth();
                        }
                        var sameDate = objDate.getDate() === date.getDate();
                        if (sameYear && sameMonth && sameDate) {
                            if (machineStatus != "正常运行中") {
                                timeIndexNum = date.getHours() * 60 * 2 + date.getMinutes() * 2;
                            }
                        } else {
                            timeIndexNum = 23 * 60 * 2 + 59 * 2 + 1;
                        }
                       var seriesPointStart = Date.UTC(objDate.getFullYear(), objDate.getMonth() + 1, objDate.getDate());
                        if (indicator === "netflow_view") {
                            if (resp === null) {
                                resp = {
                                    In_packets: [],
                                    Out_packets: [],
                                    In_bytes: [],
                                    Out_bytes: []
                                };
                            }
                            var numToReplenish = timeIndexNum - resp.In_packets.length
                            if (numToReplenish > 0) {
                                while (numToReplenish) {
                                    resp.In_packets.push(-1);
                                    resp.Out_packets.push(-1);
                                    resp.In_bytes.push(-1);
                                    resp.Out_bytes.push(-1);
                                    numToReplenish--;
                                }
                            }
                            var packetsElementId = indicator + "_nc_packets_" + index
                            var bytesElementId = indicator + "_nc_bytes_" + index
                            $(href).append($("<div></div>", {
                                    "id": packetsElementId
                                })).append($("<div></div>", {
                                    "id": bytesElementId
                                }));
                            // 先绘制包量数据图
                            var packetsSeries = [
                                {
                                    type: "area",
                                    name: "入包量(个)",
                                    pointInterval: 30 * 1000,
                                    pointStart: seriesPointStart,
                                    data: resp.In_packets
                                },
                                {
                                    type: "area",
                                    name: "出包量(个)",
                                    pointInterval: 30 * 1000,
                                    pointStart: seriesPointStart,
                                    data: resp.Out_packets
                                }
                            ];
                            var data = {
                                container: packetsElementId,
                                chartTitle: machineInChartTitle + " , " + netflow_packets_title,
                                xAxisTitle: netflow_xAxisTitle,
                                yAxisTitle: netflow_packets_yAxisTitle,
                                series: packetsSeries
                            };
                            plotHighCharts(data);

                            // 再绘制流量数据图
                            var bytesSeries = [
                                {
                                    type: "area",
                                    name: "入流量(byte)",
                                    pointInterval: 30 * 1000,
                                    pointStart: seriesPointStart,
                                    data: resp.In_bytes
                                },
                                {
                                    type: "area",
                                    name: "出流量(byte)",
                                    pointInterval: 30 * 1000,
                                    pointStart: seriesPointStart,
                                    data: resp.Out_bytes
                                }
                            ];
                            data = {
                                container: bytesElementId,
                                chartTitle: machineInChartTitle + " , " + netflow_bytes_title,
                                xAxisTitle: netflow_xAxisTitle,
                                yAxisTitle: netflow_bytes_yAxisTitle,
                                series: bytesSeries
                            };
                            plotHighCharts(data);
                        } else {
                            if (resp === null) {
                                resp = []
                            }
                            var numToReplenish = timeIndexNum - resp.length
                            if (numToReplenish > 0) {
                                while (numToReplenish) {
                                    resp.push(-1);
                                    numToReplenish--;
                                }
                            }
                            var elementId = indicator + "_" + index;
                            $(href).append($("<div></div>", {
                                "id": elementId
                            }));
                            var series = [
                                {
                                    type: 'area',
                                    name: seriesName,
                                    pointInterval: 30 * 1000,
                                    pointStart: seriesPointStart,
                                    data: resp
                                }
                            ];
                            var data = {
                                container: elementId,
                                chartTitle: machineInChartTitle + " , " + title,
                                xAxisTitle: xAxisTitle,
                                yAxisTitle: yAxisTitle,
                                series: series
                            }
                            plotHighCharts(data);
                        }
                    }
                )
                ;
            }
        }
    }

    $("#form-date").hide();

    $("textarea#machine-list").on("focusin",function (e) {
        e.preventDefault();
        $(this).attr("rows", "5");
    }).on("focusout", function (e) {
            e.preventDefault();
            $(this).attr("rows", "1");
        });
    // 为“查询日期”输入框绑定datepicker
    $(".hasDatePicker").datepicker({
        dateFormat: "yy-mm-dd",
        minDate: -30,
        maxDate: new Date()
    });
    $("#query_date_input").datepicker("setDate", new Date());

    $("#data-tab a").on("click", function (e) {
        e.preventDefault();
        if ($(this).attr("href") != "#machine-list" && $.trim($("#machine_list>table>tbody").html()) === "") {
            alertify.alert("请先查询机器！");
            return false;
        }
        $(this).tab('show');
    });

    // 搜索过滤要查询的机器，支持域名/ip/MAC地址，并将结果展示在“机器列表”中
    $("#search").on("click", function (e) {
        e.preventDefault();
        if ($(".tab-content>.active").attr("id") != "machine_list") {
            var that = $("#data-tab>.active");
            getMachineIndicatorData(that);
        } else {
            var machineList = $.trim($("#machine-list").val());
            var queryDate = $.trim($("#query_date_input").val());
            if (machineList === "" || queryDate === "") {
                alertify.alert("机器列表 和 查询日期 均不能为空！");
                return false;
            } else {
                var machineArray = machineList.split("\n"),
                    machineNum = machineArray.length;

                var ipList = new Array(),
                    hostNameList = new Array(),
                    hardwareAddrList = new Array();

                var machine;
                var ipPattern = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,
                    hardwareAddrPattern = /^([0-9a-fA-F]{2})(([/\s:-][0-9a-fA-F]{2}){5})$/;
                for (var index = 0; index < machineNum; index++) {
                    machine = machineArray[index];
                    if (ipPattern.test(machine)) {
                        ipList.push(machine);
                    } else if (hardwareAddrPattern.test(machine)) {
                        hardwareAddrList.push(machine);
                    } else {
                        if (machine != "") {
                            hostNameList.push(machine);
                        }
                    }
                }
                var urlParam = new Array();
                if (ipList.length > 0) {
                    urlParam.push("iplist=" + ipList.join(","));
                }
                if (hostNameList.length > 0) {
                    urlParam.push("hostnamelist=" + hostNameList.join(","));
                }
                if (hardwareAddrList.length > 0) {
                    urlParam.push("hardwareaddrlist=" + hardwareAddrList.join(","))
                }

                var req = $.ajax({
                    "type": "get",
                    "url": "/filter_machine_list?" + urlParam.join("&"),
                    "dataType": "json"
                });

                var wait_tip_div = $('<div></div>', {
                    'class': 'wait-tips lead'
                });
                wait_tip_div.append($('<i></i>', {
                    'class': 'fa fa-spinner fa-spin fa-3x'
                }));
                wait_tip_div.append(' 请稍等，我正努力地加载数据！');
                $("#machine_list tbody").append(wait_tip_div);

                req.done(function (resp) {
                    $("#machine_list tbody").empty();
                    var machineExistList = new Array();
                    var machineExistNum = resp.length;
                    if (resp.length > 0) {
                        for (var index = 0; index < machineExistNum; index++) {
                            if (resp[index].IsExisted) {
                                machineExistList.push(resp[index].SearchItem)
                            }
                            $("#machine_list tbody").append("<tr><td>" + (index + 1) + "</td><td class='target-machine'>" + resp[index].SearchItem + "</td><td class='host-name'>" + resp[index].Host_name + "</td><td class='hardware-addr'>" + resp[index].Hardware_addr + "</td><td>" + resp[index].Ip + "</td><td>" + resp[index].Machine_role + "</td><td class='machine-status'>" + resp[index].Status + "</td></tr>");
                        }
                    }
                    $("#machine-list").val(machineExistList.join("\n"))
                });
            }
        }
    });

    $("#data-tab>li ").on("click", function (e) {
        e.preventDefault();
        var that = $(this);
        getMachineIndicatorData(that);
    });
})
;