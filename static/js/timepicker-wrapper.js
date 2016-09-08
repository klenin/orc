requirejs.config({
    paths: {
        'jquery-ui-timepicker': 'vendor/jqueryui-timepicker-addon/dist/jquery-ui-timepicker-addon.min',
        'jquery-ui-timepicker-ru': 'vendor/jqueryui-timepicker-addon/dist/i18n/jquery-ui-timepicker-ru'
    },
    shim: {
        'jquery-ui-timepicker-ru': {
            deps: ['jquery', 'jquery-ui-timepicker']
        }
    }
});

define(['jquery', 'jquery-ui', 'jquery-ui-timepicker', 'jquery-ui-timepicker-ru'], function($) {
    $.timepicker.setDefaults({
        timeFormat: "HH:mm:ss"
    });

    return $.timepicker;
});
