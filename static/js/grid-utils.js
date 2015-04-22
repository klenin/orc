define(["utils"], function(utils) {

    function showServerPromtInGrid(gridId, prompt) {
        console.log("showServerPromtInGrid");

        var myInfo = '<div class="ui-state-highlight ui-corner-all">'
            + '<span class="ui-icon ui-icon-info" style="float: left; margin-right: .3em;"></span>'
            + '<strong>'+ prompt + '</strong><br/>' + '</div>';
        var infoTR = $("table#TblGrid_"+$("#"+gridId)[0].id+">tbody>tr.tinfo");
        var infoTD = infoTR.children("td.topinfo");

        infoTD.html(myInfo);
        infoTR.show();

        setTimeout(
            function() {
                infoTD.children("div").fadeOut(
                    "slow",
                    function() {
                        infoTR.hide();
                    }
                );
            },
            3000
        );
    }

    function showServerPromtInDialog(dialogId, prompt) {
        console.log("showServerPromtInDialog");

        var serverAns = $("<div/>", {
                class: "ui-state-highlight ui-corner-all"
            })
            .append($("<span/>", {class: "ui-icon ui-icon-info", style: "float: left; margin-right: .3em;"}))
            .append($("<strong/>", {text: prompt}))
            .append($("<br/>"));

        $("#"+dialogId).append(serverAns);
        setTimeout(
            function() {
                serverAns.fadeOut( "slow", function() {} );
            },
            3000
        );
    }

    function errTextFormat(data, gridId) {
        console.log("errTextFormat: ", data);

        var prompt, errMsg, noErr;

        switch (data.status) {
        case 200:
            noErr = true;
            prompt = "Запрос успешно выполнен";
            errMsg = "";
            break;

        case 304:
            noErr = false;
            prompt = "Используйте другое значение";
            errMsg = "Нарушение ограничения уникальности";
            break;

        case 401:
            noErr = false;
            prompt = "Войдите на сайт";
            errMsg = "Ошибка авторизации";
            break;

        case 403:
            noErr = false;
            prompt = "";
            errMsg = "У Вас не хватает прав";
            break;

        case 405:
            noErr = false;
            prompt = "";
            errMsg = data.responseText;
            break;
        }

        showServerPromtInGrid(gridId, prompt)

        return [noErr, errMsg];
    };

    function getCurrRowId(gridId) {
        var id = $("#"+gridId).jqGrid("getGridParam", "selarrrow");
        console.log("ids: ", id);

        if (id.length > 1 || id.length == 0) {
            showErrorMsg("<strong>Выберите одну запись.</strong>");
            return -1;
        }

        return id[0];
    }

    function showErrorMsg(msg) {
        $("#error").empty();
        $("#error").append(msg);
        $("#error").dialog({
            model: true,
            height: "auto",
            width: "auto",
            buttons: {
                "Закрыть": function() {
                    $(this).dialog("close");
                }
            }
        });
    }

    function ConfirmOrRejectPersonRequest(gridId, confirm) {
        console.log("ConfirmOrRejectPersonRequest");

        var id = getCurrRowId(gridId);
        if (id == -1) return false;

        var event_id = $("#"+gridId).jqGrid("getCell", id, "event_id");
        var data = { "reg_id": id, "event_id": event_id, "confirm": confirm };
        console.log("ConfirmOrRejectPersonRequest: ", data);

        utils.postRequest(
            { "reg_id": id, "event_id": event_id, "confirm": confirm},
            function(data) { showServerPromtInGrid(gridId, data["result"]); },
            "/gridhandler/confirmorrejectpersonrequest"
        );
    }

    return {
        showErrorMsg: showErrorMsg,
        getCurrRowId: getCurrRowId,
        errTextFormat: errTextFormat,
        showServerPromtInGrid: showServerPromtInGrid,
        showServerPromtInDialog: showServerPromtInDialog,
        ConfirmOrRejectPersonRequest: ConfirmOrRejectPersonRequest,
    };

});
