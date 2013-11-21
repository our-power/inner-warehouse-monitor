$(function() {
    function chart_pie(obj,titletext,series) {
        var chart = new Highcharts.Chart({
            chart : {
                renderTo : obj,
                plotBackgroundColor : null,
                plotBorderWidth : null,
                plotShadow : false,
                height:300
            },
            title : {
                text : titletext
            },

            plotOptions : {
                pie : {
                    allowPointSelect : true,
                    cursor : 'pointer',
                    dataLabels: {
                        enabled: true,
                        color: '#000000',
                        connectorColor: '#000000',
                        formatter: function() {
                            return '<b>'+ this.point.name +'</b>: '+ parseInt(this.percentage) +' %';
                        }
                    }
                }
            },
            series : eval('('+series+')')
        });
    }

    var groups = ["test"]
    for (var i in groups) {
        if (groups.hasOwnProperty(i)) {
            group = groups[i];
            var row = $("<div></div>", {
                class: "row-fluid"
            });
            var pie_div = $("<div></div>", {
                class: "span5",
                id: group
            });
            var table_div = $("<div></div>", {
                class: "span5"
            });
            var table_table = $("<table></table>", {
                class: "table table-bordered table-condensed"
            });
            var table_head = $("<tr></tr>").append($("<th>MAC</th>")).append($("<th>Host Name</th>")).append($("<th>状态</th>"));
            table_table.append(table_head);
            table_div.append(table_table);
            row.append(pie_div).append(table_div);
            $("#chart-div").append(row);
            $.getJSON("/api/statusoverview?role=" + group, function(data){
            console.log(data);
                formatted = '[{"type":"pie", "name":"' + group + '", "data":[["正常", ' + data[0] + '], ["关机", ' + data[1] + '], ["异常", ' + data[2] +']]}]';
                chart_pie(group, group + '概况', formatted);
                for (var m in data[3]) {
                    if(data[3].hasOwnProperty(m)){
                        var status = data[3][m]["Status"];
                        var desc;
                        if (status == 1) {
                            desc = "正常";
                        } else if (status == 0) {
                            desc = "关机";
                        } else {
                            desc = "异常";
                        }
                        var data_row = $("<tr></tr>")
                        .append($("<td>" + data[3][m]["Hardware_addr"] + "</td>"))
                        .append($("<td>" + data[3][m]["Host_name"] + "</td>"))
                        .append($("<td>" + desc + "</td>"));
                        table_table.append(data_row);
                    }
                }
            })
        }
    }

});