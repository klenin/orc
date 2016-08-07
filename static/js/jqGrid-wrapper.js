requirejs.config({
    paths: {
        'jquery-jqGrid': 'vendor/jqGrid/js/jquery.jqGrid.min'
    },
    shim: {
        'jquery-jqGrid': {
            deps: ['vendor/jqGrid/js/minified/i18n/grid.locale-ru']
        }
    }
});

define(['jquery', 'jquery-jqGrid'], function($, jqGrid) {
    $.extend($.jgrid.defaults, {
        loadError: function (jqXHR, textStatus, errorThrown) {
            alert('HTTP status code: '+jqXHR.status+'\n'
                +'textStatus: '+textStatus+'\n'
                +'errorThrown: '+errorThrown);
            alert('HTTP message body: '+jqXHR.responseText);
        }
    });

    return jqGrid;
});
