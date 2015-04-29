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
                alert(ajaxRequest["responseText"]);
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
                    $("#wrap #content").css("visibility", "hidden");
                } else {
                    $("#wrap #content").css("visibility", "visible");
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
