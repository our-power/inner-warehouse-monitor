$(function () {

    $(".hasDatePicker").datepicker({
        dateFormat: "yy-mm-dd",
        minDate: -30,
        maxDate: new Date()
    });
    $("#query_date_input").datepicker("setDate", new Date());
    $("#compare_date_input").datepicker("setDate", +1);

    $("#data-tab a").on("click", function (e) {
        e.preventDefault();
        $(this).tab('show');
    });

    $("#data-tab li").on("click", function (e) {
        e.preventDefault();
        var step = $(".navbar .active>a").attr("href").split("=")[1];
        var queryDate = $("#query_date_input").val();
        var href = $(this).find("a").attr("href");
        var indicator = href.substr(1)
        if (indicator === "machine_list") {
            return false
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
        switch (indicator) {
            case "cpu_view":
                title = queryDate + " , CPU使用率(%)";
                xAxisTitle = "时间";
                yAxisTitle = "使用率(%)";
                seriesName = "CPU使用率(%)";
                break;
            case "memory_view":
                title = queryDate + " , 内存使用量(MB)";
                xAxisTitle = "时间";
                yAxisTitle = "使用量(MB)";
                seriesName = "使用量(MB)";
                break;
            case "netflow_view":
                netflow_packets_title = queryDate + " , 网络包量(个)";
                netflow_bytes_title = queryDate + " , 网络流量(byte)";
                netflow_xAxisTitle = "时间";
                netflow_packets_yAxisTitle = "包量(个)";
                netflow_bytes_yAxisTitle = "流量(byte)";
                break
        }

        var req = $.ajax({
            "type": "get",
            "url": "/api/get_step_indicator_data?step=" + step + "&date=" + queryDate + "&indicator=" + indicator,
            "dataType": "json"
        });

        var wait_tip_div = $('<div></div>', {
            'class': 'wait-tips lead'
        });
        wait_tip_div.append($('<i></i>', {
            'class': 'fa fa-spinner fa-spin fa-3x'
        }));
        wait_tip_div.append(' 请稍等，我正努力地加载数据！');
        $(href).empty().append(wait_tip_div);

        req.done(function (resp) {
                $(href).empty();
                if (resp != null && resp.length > 0) {
                    var machineCount = resp.length;
                    if (indicator === "netflow_view") {
                        for (var x = 0; x < machineCount; x++) {
                            var packetsElementId = indicator + "_nc_packets_" + x;
                            var bytesElementId = indicator + "_nc_bytes_" + x;
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
                                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                                    data: resp[x].Data.In_packets
                                },
                                {
                                    type: "area",
                                    name: "出包量(个)",
                                    pointInterval: 30 * 1000,
                                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                                    data: resp[x].Data.Out_packets
                                }
                            ];
                            var data = {
                                container: packetsElementId,
                                chartTitle: resp[x].Host_name + " , " + netflow_packets_title,
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
                                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                                    data: resp[x].Data.In_bytes
                                },
                                {
                                    type: "area",
                                    name: "出流量(byte)",
                                    pointInterval: 30 * 1000,
                                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                                    data: resp[x].Data.Out_bytes
                                }
                            ];
                            data = {
                                container: bytesElementId,
                                chartTitle: resp[x].Host_name + " , " + netflow_bytes_title,
                                xAxisTitle: netflow_xAxisTitle,
                                yAxisTitle: netflow_bytes_yAxisTitle,
                                series: bytesSeries
                            };
                            plotHighCharts(data);
                        }
                    } else {
                        for (var index = 0; index < machineCount; index++) {
                            var elementId = indicator + "_" + index;
                            $(href).append($("<div></div>", {
                                "id": elementId
                            }));
                            var series = [
                                {
                                    type: 'area',
                                    name: seriesName,
                                    pointInterval: 30 * 1000,
                                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                                    data: resp[index].Data
                                }
                            ];
                            var data = {
                                container: elementId,
                                chartTitle: resp[index].Host_name + " , " + title,
                                xAxisTitle: xAxisTitle,
                                yAxisTitle: yAxisTitle,
                                series: series
                            }
                            plotHighCharts(data);
                        }
                    }
                }
            });
    });

    function plotHighCharts(data) {
        new Highcharts.Chart({
            chart: {
                renderTo: data.container,
                zoomType: 'x',
                spacingRight: 20
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
                    text: data.xAxisTitle
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
});