require(["auth", "utils"],
function(auth, utils) {

    var valid = false;

    function listEvents(data) {
        for (i in data["data"]) {
            var p = $("</p>", {});
            $(p).append($("<a/>", {
                text: data["data"][i]["name"],
                href: "/handler/getrequest/event/" + data["data"][i]["id"],
                class: "form-row",
            }));
            $(p).appendTo("div#list-events");
        }
    }

    $(document).ready(function() {
        utils.postRequest(
            {
                "table": "events",
                "fields": ["id", "name"]
            },
            listEvents,
            "/handler/geteventlist"
        );

        //by default
        $("#tab2").stop(false, false).show();

        $("#navigation li a").each(function(i) {
            $("#navigation li a:eq(" + i + ")").click(function() {
                var tab_id = i + 1;
                $("#content div").stop(false, false).hide();
                $("#tab" + tab_id).stop(false, false).show();
                return false;
            })
        });

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

});
