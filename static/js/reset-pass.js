define(["utils", "grid-utils"], function(utils, gridUtils) {

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

        var id = gridUtils.getCurrRowId(gridId);
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
                        gridUtils.showErrorMsg(result.msg);
                        return false;
                    }
                    utils.postRequest(
                        { "pass1": $("#"+passId1).val(), "pass2": $("#"+passId2).val(), "id": id },
                        function(data) { gridUtils.showServerPromtInDialog($(this), data["result"]); },
                        "/gridhandler/resetpassword"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        ResetPassword: ResetPassword,
    };

});