define(function() {

    function postRequest(data, callback, url) {
        $.ajax({
            method: "post",
            type: "post",
            dataType: "json",
            url: url,
            async: false,
            data: JSON.stringify(data),
            ContentType: "application/json; charset=utf-8",
            success: function(data) {
                console.log(data);
                callback(data);
            },
            error: function(ajaxRequest, ajaxOptions, thrownError) {
                console.log(thrownError);
                console.log(ajaxOptions);
                console.log(ajaxRequest);
            }
        });
    };

    function areAlive() {
        postRequest(
            {
                "action": "are-alive",
            },
            function(data) {
                if (data) {
                    $("#logout-btn, #cabinet-btn").css("visibility", "visible");
                }
            },
            "/handler"
        );
    }

    return {
        postRequest: postRequest,
        areAlive: areAlive
    };

});