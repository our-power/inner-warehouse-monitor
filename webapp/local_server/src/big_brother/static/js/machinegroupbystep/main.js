$(function () {

    $(".hasDatePicker").datepicker({
        "dateFormat": "yy-mm-dd"
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
        if(href === "#machinelist") {
            return false
        }
        /*
        console.log(step);
        console.log(queryDate);
        console.log(href);
        */
        var req = $.ajax({
           "type": "get",
            "url": "/get_step_indicator_data?step="+step+"&date="+queryDate+"&indicator="+href,
            "dataType": "json"
        });

        req.done(function(resp){
           console.log(resp.data);
        });
    });
});