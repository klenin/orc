define(['jquery', 'jquery-ui'], function($) {
    $.ajax('/vendor/jquery-ui/ui/minified/i18n/datepicker-ru.js', {
        complete: function(response) {
            eval('(function() { var define; ' + response.responseText + ' })()');
            $.datepicker.setDefaults({ dateFormat: 'yy-mm-dd' });
        },
        dataType: 'text'
    });

    $.datepicker.setDefaults({
        changeMonth: true,
        changeYear: true
    });

    return $.datepicker;
});
