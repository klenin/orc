define(['jquery', 'vendor/jquery.kladr/jquery.kladr.min'], function($) {

function kladr() {
    var token = '51dfe5d42fb2b43e3300006e';
    var key   = '86a2c2a06f1b2451a87d05512cc2c3edfdf41969';

    var region   = $('[type="region"]');
    var district = $('[type="district"]');
    var city     = $('[type="city"]');
    var street   = $('[type="street"]');
    var building = $('[type="building"]');

    region.kladr({
        token: token,
        key: key,
        type: $.kladr.type.region,
        select: function(obj) {
            region.parent().find('label').text(obj.type);
            district.kladr('parentType', $.kladr.type.region);
            district.kladr('parentId', obj.id);
            city.kladr('parentType', $.kladr.type.region);
            city.kladr('parentId', obj.id);
        }
    });

    district.kladr({
        token: token,
        key: key,
        type: $.kladr.type.district,
        select: function(obj) {
            district.parent().find('label').text(obj.type);
            city.kladr('parentType', $.kladr.type.district);
            city.kladr('parentId', obj.id);
        }
    });

    city.kladr({
        token: token,
        key: key,
        type: $.kladr.type.city,
        select: function(obj) {
            city.parent().find('label').text(obj.type+' ');
            street.kladr('parentType', $.kladr.type.city);
            street.kladr('parentId', obj.id);
            building.kladr('parentType', $.kladr.type.city);
            building.kladr('parentId', obj.id);
        }
    });

    street.kladr({
        token: token,
        key: key,
        type: $.kladr.type.street,
        select: function(obj) {
            street.parent().find('label').text(obj.type+' ');
            building.kladr('parentType', $.kladr.type.street);
            building.kladr('parentId', obj.id);
        }
    });

    building.kladr({
        token: token,
        key: key,
        type: $.kladr.type.building,
        select: function(obj) {
            building.parent().find('label').text('Дом ');
        }
    });
}

return {
    kladr: kladr,
};

});
