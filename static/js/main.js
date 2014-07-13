require(["auth", "utils"],
function(auth, utils) {

    var valid = false;

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
                "fields": ["id", "name"]
            },
            listEvents,
            "/handler"
        );

        utils.areAlive();

    });

    $("#register-btn").click(function() {
        if (valid == true) {
            auth.jsonHandle("register", auth.registerCallback);
        } else {
            $("#server-answer").text("Не все поля заполнены.").css("color", "red");
        }
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

    $("#fname, #lname, #pname").blur(function() {
        //if ($("#fname").val() != "" && $("#lname").val() != "" && $("#pname").val() != ""){
        if ($(this).val() != ""){
            valid = true;
            $(this).css({"border": "2px solid green"});
        } else {
            valid = false;
            $(this).css({"border": "2px solid red"});
        }
    });

    function listEvents(data) {
    console.log(data);
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