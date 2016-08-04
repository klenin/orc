define(['jquery-ui'], function() {

    function initDatePicker(root) {

        var currYear = (new Date).getFullYear();

        $.datepicker.setDefaults($.datepicker.regional['ru']);

        var datepickerOptions = {
            showAnim: 'slideDown',
            dateFormat: 'yy-mm-dd',
            changeMonth: true,
            changeYear: true,
            yearRange: String(currYear) + ':' + String(currYear+10)
        };

        $(root).datepicker(datepickerOptions);
    }

    return {
        initDatePicker: initDatePicker,
    };

});
