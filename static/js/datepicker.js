define("datepicker", function() {

    return function() {
        var currYear = (new Date).getFullYear();

        $.datepicker.setDefaults($.datepicker.regional['']);

        $('input[id$="date"]')
        .attr('type', 'text')
        .datepicker({
            showAnim: 'slideDown',
            dateFormat: 'yy-mm-dd',
            changeMonth: true,
            changeYear: true,
            yearRange: String(currYear-120 + ':' + currYear)
        })
        .attr('type', 'text');
    }

});
