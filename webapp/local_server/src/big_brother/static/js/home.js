$(document).ready(function() {
    function chart_pie(obj,titletext,series) {
        var chart = new Highcharts.Chart({
            chart : {
                renderTo : obj,
                plotBackgroundColor : null,
                plotBorderWidth : null,
                plotShadow : false,
                height:300
            },
            credits: {
                enabled : false
            },
            title : {
                text : titletext
            },
            colors: [
                '#228B22',
                '#c0c0c0',
                '#FF0000'
            ],
            plotOptions : {
                pie : {
                    allowPointSelect : true,
                    cursor : 'pointer',
                    dataLabels: {
                        enabled: false
                    },
                    showInLegend: true
                }
            },
            series : eval('('+series+')')
        });
    }

    var groups = ["test@财务开票"]
    for (var i in groups) {
        if (groups.hasOwnProperty(i)) {
            var sp = new Array();
            sp = groups[i].split("@");
            group = sp[0];
            label = sp[1];
            var row = $("<div></div>", {
                class: "row-fluid"
            });
            var pie_div = $("<div></div>", {
                class: "span5",
                id: group
            });
            var table_div = $("<div></div>", {
                class: "span7"
            });
            var table_table = $("<table></table>", {
                class: "table table-bordered table-condensed"
            });
            var table_head = $("<tr></tr>")
                .append($("<th>MAC</th>"))
                .append($("<th>Host Name</th>"))
                .append($("<th>状态</th>"));
            table_table.append($("<caption>" + label + "组机器情况</caption>"))
            table_table.append(table_head);
            table_div.append(table_table);
            row.append(pie_div).append(table_div);
            $("#chart-div").append(row).append($("<hr />"));
            $.getJSON("/api/status_overview?role=" + group, function(data){
            console.log(data);
                formatted = '[{"type":"pie", "name":"' + group + '", "data":[["正常", ' + data[0] + '], ["关机", ' + data[1] + '], ["异常", ' + data[2] +']]}]';
                chart_pie(group, label + ' 概况', formatted);
                for (var m in data[3]) {
                    if(data[3].hasOwnProperty(m)){
                        var status = data[3][m]["Status"];
                        var desc;
                        var bgc;
                        if (status == 1) {
                            desc = "正常运行中";
                            bgc = "#e6fcf0";
                        } else if (status == 0) {
                            desc = "已正常关机";
                            bgc = "#c0c0c0";
                        } else {
                            desc = "运行异常";
                            bgc = "#FF6347";
                        }
                        var data_row = $("<tr></tr>")
                            .append($("<td>" + data[3][m]["Hardware_addr"] + "</td>"))
                            .append($("<td>" + data[3][m]["Host_name"] + "</td>"))
                            .append($("<td>" + desc + "</td>"));
                        data_row.attr("style", "BACKGROUND-COLOR: " + bgc)
                        table_table.append(data_row);
                    }
                }
            })
        }
    }

});