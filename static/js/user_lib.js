define(["utils", "grid_lib"], function(utils, gridLib) {

    function GetServerMsg(msg) {
        console.log("GetServerMsg");

        if (msg === "ok") {
            return "Пароль изменен";
        } else if (msg === "badPassword") {
            return "Неверные значения паролей. Пароль должен иметь длину от 6 до 36 символов";
        } else if (msg === "differentPasswords") {
            return "Пароли не совпадают";
        }
    }

    function CheckPass(passId1, passId2) {
        console.log("CheckPass");

        var pattern = /^.{6,36}$/;

        if (pattern.test($("#"+passId1).val())
            && pattern.test($("#"+passId2).val())
            && $("#"+passId1).val() === $("#"+passId2).val()) {
            return { "result": true, "msg": GetServerMsg("ok") };

        } else if ($("#"+passId1).val() !== $("#"+passId2).val()) {
            return { "result": false, "msg": GetServerMsg("differentPasswords") };

        } else if (!pattern.test($("#"+passId1).val()) || !pattern.test($("#"+passId2).val())) {
            return { "result": false, "msg": GetServerMsg("badPassword") };
        }
    }

    function ResetPassword(dialogId, gridId, passId1, passId2) {
        console.log("ResetPassword");

        var id = gridLib.getCurrRowId(gridId);
        if (id == -1) return false;

        $("#"+passId1+", #"+passId2).val("");

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить": function() {
                    var result = CheckPass(passId1, passId2);
                    if (!result.result) {
                        gridLib.showErrorMsg(result.msg);
                        return false;
                    }
                    utils.postRequest(
                        { "pass": $("#"+passId1).val(), "id": id },
                        function(data) { gridLib.showServerPromtInDialog($("#"+dialogId), data["result"]); },
                        "/usercontroller/resetpassword"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    function CheckSession() {
        utils.postRequest(
            null,
            function(data) {
                if (data["result"] === "ok") {
                    $("#logout-btn, #cabinet-btn").css("visibility", "visible");
                    $("#wrap #content").css("visibility", "hidden");
                } else {
                    $("#wrap #content").css("visibility", "visible");
                    $("#logout-btn, #cabinet-btn").css("visibility", "hidden");
                }
            },
            "/usercontroller/checksession"
        );
    };

    function Login(gridId) {
        var user_id = gridLib.getCurrRowId(gridId);
        if (user_id == -1) return false;

        location.href = "/usercontroller/login/"+user_id;
    }

    function SendEmailWellcomeToProfile(gridId, dialog) {
        var user_id = gridLib.getCurrRowId(gridId);
        if (user_id == -1) return false;

        utils.postRequest(
            { "user_id": user_id },
            function(data) { gridLib.showServerPromtInDialog(dialog, data["result"]); },
            "/usercontroller/sendemailwellcometoprofile"
        );
    }

    function ConfirmOrRejectPersonRequest(gridId, confirm) {
        console.log("ConfirmOrRejectPersonRequest");

        var id = gridLib.getCurrRowId(gridId);
        if (id == -1) return false;

        var event_id = $("#"+gridId).jqGrid("getCell", id, "event_id");
        var data = { "reg_id": id, "event_id": event_id, "confirm": confirm };
        console.log("ConfirmOrRejectPersonRequest: ", data);

        utils.postRequest(
            { "reg_id": id, "event_id": event_id, "confirm": confirm},
            function(data) { gridLib.showServerPromtInGrid(gridId, data["result"]); },
            "/usercontroller/confirmorrejectpersonrequest"
        );
    }

    return {
        ResetPassword: ResetPassword,
        CheckPass: CheckPass,
        CheckSession: CheckSession,
        Login: Login,
        SendEmailWellcomeToProfile: SendEmailWellcomeToProfile,
        ConfirmOrRejectPersonRequest: ConfirmOrRejectPersonRequest,
    };

});
