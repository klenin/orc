requirejs.config({
    baseUrl: "/",
    paths: {
        jquery: "vendor/jquery/dist/jquery.min",
        jqGrid: "vendor/jqGrid/js/jquery.jqGrid.min",
        "jquery-ui": "vendor/jquery-ui/jquery-ui.min",
        kladr: "js/kladr/kladr",
        datepicker: "js/datepicker/datepicker",
        utils: "js/utils",
        grid_lib: "js/grid_lib",
        subgrid_lib: "js/subgrid_lib",
        group_lib: "js/group_lib",
        blank: "js/blank",
        user_lib: "js/user_lib"
    },
    shim: {
        jqGrid: {
            deps: ["vendor/jqGrid/js/minified/i18n/grid.locale-ru"]
        }
    }
});

require(["jquery", "user_lib"], function($, userLib) {
    $(document).ready(function() {
        userLib.CheckSession();
    });
});

require(["jquery", "utils"], function($, utils) {
    function loginCallback(data) {
        if (data.result === "ok") {
            $("#events #server-answer").empty();
            $("#server-answer").text("Авторизация прошла успешно.").css("color", "green");
            $("#logout-btn, #cabinet-btn").css("visibility", "visible");
            $("#password, #username").val("");

        } else if (data.result === "invalidCredentials") {
            $("#server-answer").text("Неверный логин.").css("color", "red");

        } else if (data.result === "badPassword") {
            $("#server-answer").text("Неверный пароль.").css("color", "red");
        } else if (data.result === "notEnabled") {
            $("#server-answer").text("Ваш аккаунт заблокирован.").css("color", "red");
        }
    };

    function logoutCallback(data) {
        if (data.result === "ok") {
            $("#server-answer").text("Вы вышли.").css("color", "green").css("visibility", "visible");
            console.log("Вы вышли.");
            location.href = "/"

        } else if (data.result === "badSid") {
            $("#server-answer").text("Invalid session ID.").css("color", "red");
        }
    };

    $("#login-btn").click(function() {
        var data = {};
        data["login"] = $("#content #username").val();
        data["password"] = $("#content #password").val();
        utils.postRequest(data, loginCallback, "/registrationcontroller/login");
    });

    $("#logout-btn").click(function() {
        utils.postRequest(null, logoutCallback, "/registrationcontroller/logout");
    });

});
