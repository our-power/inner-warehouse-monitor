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

    function newIEDate(dateStr) {
        // dateStr的格式：xxxx-xx-xx
        var dateParts = dateStr.split("-");
        return new Date(parseInt(dateParts[0]), parseInt(dateParts[1]), parseInt(dateParts[2]));
    }

    function ajaxCallbackClosure(params) {
        params.req.done(function (resp) {
                $(params.wait_tips_id).remove();
                var date, timeIndexNum, numToReplenish, data;
                // 对于非“正常运行中”的机器，补充最后收到数据的时间点之后时间点的数据
                date = new Date();
                var sameYear = params.objDate.getFullYear() === date.getFullYear();
                var sameMonth;
                if ($.browser.msie) {
                    sameMonth = params.objDate.getMonth() === date.getMonth() + 1;
                } else {
                    sameMonth = params.objDate.getMonth() === date.getMonth();
                }
                var sameDate = params.objDate.getDate() === date.getDate();
                if (sameYear && sameMonth && sameDate) {
                    if (params.machine_status != "正常运行中") {
                        timeIndexNum = date.getHours() * 60 * 2 + date.getMinutes() * 2;
                    }
                } else {
                    timeIndexNum = 23 * 60 * 2 + 59 * 2 + 1;
                }
                var seriesPointStart = Date.UTC(params.objDate.getFullYear(), params.objDate.getMonth() + 1, params.objDate.getDate());
                if (params.indicator === "netflow_view") {
                    if (resp === null) {
                        resp = {
                            In_packets: [],
                            Out_packets: [],
                            In_bytes: [],
                            Out_bytes: []
                        };
                    }
                    numToReplenish = timeIndexNum - resp.In_packets.length;
                    if (numToReplenish > 0) {
                        while (numToReplenish) {
                            resp.In_packets.push(-1);
                            resp.Out_packets.push(-1);
                            resp.In_bytes.push(-1);
                            resp.Out_bytes.push(-1);
                            numToReplenish--;
                        }
                    }
                    var packetsElementId = params.indicator + "_nc_packets_" + params.hardware_addr;
                    var bytesElementId = params.indicator + "_nc_bytes_" + params.hardware_addr;
                    $(params.href).append($("<div></div>", {
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
                    data = {
                        container: packetsElementId,
                        chartTitle: params.search_target + " , " + params.netflow_packets_title,
                        xAxisTitle: params.netflow_xAxisTitle,
                        yAxisTitle: params.netflow_packets_yAxisTitle,
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
                        chartTitle: params.search_target + " , " + params.netflow_bytes_title,
                        xAxisTitle: params.netflow_xAxisTitle,
                        yAxisTitle: params.netflow_bytes_yAxisTitle,
                        series: bytesSeries
                    };
                    plotHighCharts(data);
                } else {
                    if (resp === null) {
                        resp = []
                    }
                    numToReplenish = timeIndexNum - resp.length;
                    if (numToReplenish > 0) {
                        while (numToReplenish) {
                            resp.push(-1);
                            numToReplenish--;
                        }
                    }
                    var elementId = params.indicator + "_" + params.hardware_addr;
                    $(params.href).append($("<div></div>", {
                        "id": elementId
                    }));
                    var series = [
                        {
                            type: 'area',
                            name: params.seriesName,
                            pointInterval: 30 * 1000,
                            pointStart: seriesPointStart,
                            data: resp
                        }
                    ];
                    data = {
                        container: elementId,
                        chartTitle: params.search_target + " , " + params.title,
                        xAxisTitle: params.xAxisTitle,
                        yAxisTitle: params.yAxisTitle,
                        series: series
                    };
                    plotHighCharts(data);
                }
            }
        )
        ;
    }

    function loadWaitTips(parentID, tipsID) {
        // 加载等待提示
        var wait_tip_div = $('<div></div>', {
            'class': 'wait-tips lead',
            'id': tipsID
        });
        wait_tip_div.append($('<i></i>', {
            'class': 'fa fa-spinner fa-spin fa-3x'
        }));
        wait_tip_div.append(' 请稍等，我正努力地加载数据！');
        $(parentID).append(wait_tip_div);
    }

    function getMachineIndicatorData(jqObj) {
        var queryDate = $("#query_date_input").val();
        var objDate;
        if ($.browser.msie) {
            objDate = newIEDate(queryDate);
        } else {
            objDate = new Date(queryDate);
        }

        var href = jqObj.find("a").attr("href");
        var indicator = href.substr(1);

        $("button#search").prop("disabled", false);
        if (indicator === "machine_list") {
            $("#form-machine-list").show();
            $("#form-date").hide();
            return false
        } else {
            $("#form-machine-list").hide();
            $("#form-date").show();
            $("input#query_date_input").prop('disabled', false);
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

                // 加载等待提示
                loadWaitTips(href, "wait-tips-" + index);

                var machineInChartTitle = $(hardwareAddrList[index]).siblings(".target-machine").text();
                var machineStatus = $(hardwareAddrList[index]).siblings(".machine-status").text();
                var req = $.ajax({
                    "type": "get",
                    "url": "/api/get_machine_indicator_data?hardware_addr=" + hardwareAddr + "&date=" + queryDate + "&indicator=" + trueIndicator,
                    "dataType": "json"
                });
                var params = {
                    req: req,
                    indicator: indicator,
                    href: href,
                    objDate: objDate,
                    hardware_addr: hardwareAddr,
                    machine_status: machineStatus,
                    search_target: machineInChartTitle,
                    wait_tips_id: "#wait-tips-" + index
                };
                if (indicator === "netflow_view") {
                    params["netflow_xAxisTitle"] = netflow_xAxisTitle;
                    params["netflow_packets_yAxisTitle"] = netflow_packets_yAxisTitle;
                    params["netflow_bytes_yAxisTitle"] = netflow_bytes_yAxisTitle;
                    params["netflow_packets_title"] = netflow_packets_title;
                    params["netflow_bytes_title"] = netflow_bytes_title;
                } else {
                    params["xAxisTitle"] = xAxisTitle;
                    params["yAxisTitle"] = yAxisTitle;
                    params["seriesName"] = seriesName;
                    params["title"] = title;
                }
                ajaxCallbackClosure(params);
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

                // 加载等待提示
                loadWaitTips("#machine_list tbody", "");

                req.done(function (resp) {
                    $("#machine_list tbody").empty();
                    $("#machine_list tr").remove();
                    var machineExistList = new Array();
                    var machineExistNum = resp.length;
                    if (resp.length > 0) {
                        for (var index = 0; index < machineExistNum; index++) {
                            if (resp[index].IsExisted) {
                                machineExistList.push(resp[index].SearchItem)
                            }
                            var trStatusClass;
                            switch (resp[index].Status) {
                                case "正常运行中":
                                    trStatusClass = "normal";
                                    break;
                                case "运行异常":
                                    trStatusClass = "abnormal";
                                    break;
                                case "已正常关机":
                                    trStatusClass = "closed";
                                    break;
                                case "不再使用":
                                    trStatusClass = "not_use"
                                    break;
                                default:
                                    trStatusClass = "";
                            }
                            $("#machine_list tbody").append("<tr class='" + trStatusClass + "'><td>" + (index + 1) + "</td><td class='target-machine'>" + resp[index].SearchItem + "</td><td class='host-name'>" + resp[index].Host_name + "</td><td class='hardware-addr'>" + resp[index].Hardware_addr + "</td><td>" + resp[index].Ip + "</td><td>" + resp[index].Machine_role + "</td><td class='machine-status'>" + resp[index].Status + "</td></tr>");
                        }
                    }
                    //$("#machine-list").val(machineExistList.join("\n"))
                });
            }
        }
    });

    function raphaelDrawAccessibility(HTMLElement, pingData, telnetData) {
        var pingDataLength = pingData.length;
        var angleRange,
            beginAngle = 270;
        if (pingDataLength <= 10) {
            angleRange = 180;
        } else {
            angleRange = 360;
        }
        var unitAngle = angleRange / pingDataLength;

        var svgWidth,
            svgHeight = 500,
            radius,
            pingCircleCenterX,
            pingCircleCenterY,
            telnetCircleCenterX = 740,
            telnetCircleCenterY = 240,
            notOk;

        if ($.browser.msie) {
            svgWidth = 820;
            radius = 150;
            pingCircleCenterX = 180;
            pingCircleCenterY = 180;
            telnetCircleCenterX = 570;
            telnetCircleCenterY = 180;
            notOk = "fail to connect"
        } else {
            svgWidth = 1000;
            radius = 200;
            pingCircleCenterX = 240;
            pingCircleCenterY = 240;
            telnetCircleCenterX = 740;
            telnetCircleCenterY = 240;
            notOk = "不通"
        }

        var paper = Raphael(HTMLElement, svgWidth, svgHeight);
        paper.circle(pingCircleCenterX, pingCircleCenterY, 8).attr({fill: "yellow"});
        paper.text(pingCircleCenterX - 25, pingCircleCenterY, "ping").attr({"font-weight": "bold", "font-size": 15});

        var x, y;
        var xOffset, yOffset;
        for (var pingIndex = 0; pingIndex < pingDataLength; pingIndex++) {
            x = radius * Math.cos(Raphael.rad(beginAngle + pingIndex * unitAngle));
            y = radius * Math.sin(Raphael.rad(beginAngle + pingIndex * unitAngle));
            paper.circle(pingCircleCenterX + x, pingCircleCenterY + y, 8).attr({fill: "yellow"});
            paper.path(
                    ['M', pingCircleCenterX, pingCircleCenterY,
                        'l', x, y
                    ]).attr({
                    stroke: ((pingData[pingIndex].Response_time == -1) ? "red" : "green"),
                    "stroke-dasharray": ((pingData[pingIndex].Response_time == -1) ? "." : ""),
                    "stroke-width": 3,
                    'arrow-end': 'block-midium-long',
                    'arrow-start': 'none',
                    title: ((pingData[pingIndex].Response_time == -1) ? notOk : pingData[pingIndex].Response_time + "ms")
                });
            if (x < 0) {
                xOffset = -15;
            } else {
                xOffset = 15;
            }
            if (y < 0) {
                yOffset = -15;
            } else {
                yOffset = 15;
            }
            paper.text(pingCircleCenterX + x / 2, pingCircleCenterY + y / 2, ((pingData[pingIndex].Response_time == -1) ? notOk : pingData[pingIndex].Response_time + "ms")).attr({"font-weight": "bold"});
            paper.text(pingCircleCenterX + x + xOffset, pingCircleCenterY + y + yOffset, pingData[pingIndex].Target_ip).attr({"font-size": 14, "font-weight": "bold"});
        }

        var telnetDataLength = telnetData.length;
        paper.circle(telnetCircleCenterX, telnetCircleCenterY, 8).attr({fill: "yellow"});
        paper.text(telnetCircleCenterX - 25, telnetCircleCenterY, "telnet").attr({"font-weight": "bold", "font-size": 15});

        for (var telnetIndex = 0; telnetIndex < telnetDataLength; telnetIndex++) {
            x = radius * Math.cos(Raphael.rad(beginAngle + telnetIndex * unitAngle));
            y = radius * Math.sin(Raphael.rad(beginAngle + telnetIndex * unitAngle));
            paper.circle(telnetCircleCenterX + x, telnetCircleCenterY + y, 8).attr({fill: "yellow"});
            paper.path(
                    ['M', telnetCircleCenterX, telnetCircleCenterY,
                        'l', x, y
                    ]).attr({
                    stroke: ((telnetData[telnetIndex].Status === "OK") ? "green" : "red"),
                    "stroke-dasharray": ((telnetData[telnetIndex].Status === "OK") ? "" : "."),
                    "stroke-width": 3,
                    'arrow-end': 'block-midium-long',
                    'arrow-start': 'none',
                    title: (telnetData[telnetIndex].Status)
                });
            if (x < 0) {
                xOffset = -15;
            } else {
                xOffset = 15;
            }
            if (y < 0) {
                yOffset = -15;
            } else {
                yOffset = 15;
            }
            paper.text(telnetCircleCenterX + x / 2, telnetCircleCenterY + y / 2, ((telnetData[telnetIndex].Status === "NotOK") ? notOk : "OK")).attr({"font-weight": "bold"});
            paper.text(telnetCircleCenterX + x + xOffset, telnetCircleCenterY + y + yOffset, telnetData[telnetIndex].Target_url).attr({"font-size": 14, "font-weight": "bold"});
        }
    }

    function formatTime(time_index) {
        var hours = Math.floor(time_index / 2 / 60);
        var minutes = Math.floor(time_index / 2) % 60;
        var seconds = time_index % 2 * 30;
        if (hours < 10) {
            hours = "0" + hours;
        }
        if (minutes < 10) {
            minutes = "0" + minutes;
        }
        if (seconds === 0) {
            seconds = "00";
        }
        return hours + ":" + minutes + ":" + seconds;
    }

    $("#data-tab>li").on("click", function (e) {
        e.preventDefault();
        var targetTab = $(this).find("a").attr("href");
        if (targetTab != "#machine-list" && $.trim($("#machine_list>table>tbody").html()) === "") {
            alertify.alert("请先查询机器！");
            return false;
        }

        $(this).find("a").tab('show');

        if (targetTab === "#accessibility") {
            $("#form-machine-list").hide();
            $("#form-date").show();
            $("input#query_date_input").prop('disabled', true);
            $("button#search").prop("disabled", true);
            $("#accessibility").empty();
            var hardwareAddrList = $("#machine_list .hardware-addr");
            var machineNum = hardwareAddrList.length;
            for (var index = 0; index < machineNum; index++) {
                var hardwareAddr = $.trim($(hardwareAddrList[index]).text());
                if (hardwareAddr != "") {
                    var machineStatus = $(hardwareAddrList[index]).siblings(".machine-status").text();
                    if (machineStatus === "正常运行中") {

                        // 加载等待提示
                        loadWaitTips("#accessibility", "wait-tips-" + hardwareAddr.replace(/:/g, ""));

                        var req = $.ajax({
                            "type": "get",
                            "url": "/api/get_machine_accessibility_data?hardware_addr=" + hardwareAddr,
                            "dataType": "json"
                        });
                        req.done(function (resp) {
                            var pingTimeIndex = formatTime(resp.Ping_time_index);
                            var telnetTimeIndex = formatTime(resp.Telnet_time_index);
                            $("#accessibility").append($("<p></p>", {
                                    html: "<strong>机器</strong>：" + resp.Hardware_addr + "，<strong>日期</strong>：" + resp.Date + "，<strong>ping 时间</strong>：" + pingTimeIndex + "，<strong>telnet 时间</strong>：" + telnetTimeIndex
                                })).append("<hr>").append($("<div></div>", {
                                    "id": "accessibility_" + resp.Hardware_addr
                                }));

                            // 删除等待提示信息
                            var waitTipsID = "#wait-tips-" + resp.Hardware_addr.replace(/:/g, "");
                            $(waitTipsID).remove();

                            raphaelDrawAccessibility("accessibility_" + resp.Hardware_addr, resp.Ping_results, resp.Telnet_results)
                        });
                    }
                }
            }
        } else {
            var that = $(this);
            getMachineIndicatorData(that);
        }
    });

    var machine_list = $.trim($("#machine-list").val());
    if(machine_list != ''){
        $('#search').click();
    }
})
;