require(["auth", "utils"],
function(auth, utils) {

    $(document).ready(function() {

        //by default
        $("#tab2").stop(false, false).show();

        $("#navigation li a").each(function(i) {
            $("#navigation li a:eq(" + i + ")").click(function() {
                var tab_id = i + 1;
                $("#content div").stop(false, false).hide();
                $("#tab" + tab_id).stop(false, false).show();
                return false;
            })
        })

        utils.postRequest(
            {
                "action": "select",
                "table": "events",
                "fields": ["id", "name"],
                "count": "2"
            },
            listEvents,
            "/handler"
        );

    });

    $("#register-btn").click(function() {
        auth.jsonHandle("register", auth.registerCallback);
    });

    $("#login-btn").click(function() {
        auth.jsonHandle("login", auth.loginCallback);
    });

    $("#logout-btn").click(function() {
        auth.jsonHandle("logout", auth.logoutCallback);
    });

    $("#cabinet-btn").click(function() {
        location.href = "/handler/showcabinet/users/";
    });

    $("#home-btn").click(function() {
        location.href = "/";
    });

    function listEvents(data) {
        var d = data["data"];
        for (i in d) {
            var p = $("</p>", {});
            $(p).append($("<a/>", {
                text: d[i]["name"],
                href: "/handler/show/event/" + d[i]["id"],
                class: "form-row",
            }))
            $(p).appendTo("div#list-events");
        }
    }

});