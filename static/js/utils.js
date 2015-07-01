define(["grid-utils"], function(gridUtils) {

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
        postRequest(
            { "action": "checkSession" },
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
    };

    function login(gridId) {
        var user_id = gridUtils.getCurrRowId(gridId);
        if (user_id == -1) return false;

        location.href = "/usercontroller/login/"+user_id;
    }

    function sendEmailWellcomeToProfile(gridId, dialog) {
        var user_id = gridUtils.getCurrRowId(gridId);
        if (user_id == -1) return false;

        postRequest(
            { "user_id": user_id },
            function(data) { gridUtils.showServerPromtInDialog(dialog, data["result"]); },
            "/usercontroller/sendemailwellcometoprofile"
        );
    }

    return {
        postRequest: postRequest,
        checkSession: checkSession,
        login: login,
        sendEmailWellcomeToProfile: sendEmailWellcomeToProfile,
    };

});
