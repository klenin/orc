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

    function checkSession() {
        postRequest({
                "action": "checkSession",
            },
            function(data) {
                if (data["result"] === "ok") {
                    $("#logout-btn, #cabinet-btn").css("visibility", "visible");
                } else {
                    $("#logout-btn, #cabinet-btn").css("visibility", "hidden");
                }
            },
            "/handler"
        );
    }

    return {
        postRequest: postRequest,
        checkSession: checkSession,
    };

});
