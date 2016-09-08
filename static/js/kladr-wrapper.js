define(['jquery', 'vendor/jquery.kladr/jquery.kladr.min'], function($, kladr) {
    $.kladr.setDefault({
        token: '51dfe5d42fb2b43e3300006e',
        key: '86a2c2a06f1b2451a87d05512cc2c3edfdf41969',
        select: function (obj) {
            $(this).parent().find('label').text(obj.type);
        }
    });

    return kladr;
});
