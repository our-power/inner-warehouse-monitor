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
                title = queryDate + " , 网络流量";
                break
        }

        var req = $.ajax({
            "type": "get",
            "url": "/get_step_indicator_data?step=" + step + "&date=" + queryDate + "&indicator=" + indicator,
            "dataType": "json"
        });

        req.done(function (resp) {
            $(href).empty()
            if (resp != null && resp.length > 0) {
                var machineCount = resp.length;
                for (var index = 0; index < machineCount; index++) {
                    $(href).append($("<div></div>", {
                        "id": indicator + "_" + index
                    }));
                    var container = indicator + "_" + index
                    plotHighCharts(container, resp[index].Host_name + " , " + title, xAxisTitle, yAxisTitle, seriesName, resp[index].Data);
                }
            }
        });
    });

    function plotHighCharts(container, chartTitle, xAxisTitle, yAxisTitle, seriesName, seriesData) {
        new Highcharts.Chart({
            chart: {
                renderTo: container,
                zoomType: 'x',
                spacingRight: 20
            },
            credits: {
                enabled: false
            },
            title: {
                text: chartTitle
            },
            xAxis: {
                type: 'datetime',
                maxZoom: 5 * 60 * 1000, // five minutes
                title: {
                    text: xAxisTitle
                }
            },
            yAxis: {
                title: {
                    text: yAxisTitle
                },
                min: -5,
                startOnTick: false
            },
            tooltip: {
                shared: true
            },
            legend: {
                enabled: false
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
            series: [
                {
                    type: 'area',
                    name: seriesName,
                    pointInterval: 30 * 1000,
                    pointStart: Date.UTC(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()),
                    data: seriesData
                }
            ]
        });
    }
});