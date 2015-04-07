require(["auth", "utils"],
function(auth, utils) {

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
