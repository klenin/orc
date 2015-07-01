require(["utils"],
function(utils) {

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

    $("#cabinet-btn").click(function() {
        location.href = "/usercontroller/showcabinet/users/";
    });

    $("#home-btn").click(function() {
        location.href = "/";
    });

});
