define(["jquery", "utils", "datepicker"], function($, utils) {

    function resizeSelectWidth(form) {
        var maxWidth = 0, i,
            elems = form.find('tr.FormData > td.DataTD > .FormElement');
        for (i = 0; i < elems.length; i++) {
            $(elems[i]).attr("size", $(elems[i]).val().length+10);
            maxWidth = Math.max(maxWidth, $(elems[i]).width());
        }
        for (i = 0; i < elems.length; i++) {
            $(elems[i]).width(maxWidth+"px");
        }
    };

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

    function showServerPromtInDialog(obj, prompt) {
        console.log("showServerPromtInDialog");

        var serverAns = $("<div/>", {
                class: "ui-state-highlight ui-corner-all"
            })
            .append($("<span/>", {class: "ui-icon ui-icon-info", style: "float: left; margin-right: .3em;"}))
            .append($("<strong/>", {text: prompt}))
            .append($("<br/>"));

       obj.append(serverAns);
        setTimeout(
            function() {
                serverAns.fadeOut( "slow", function() { serverAns.remove(); } );
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
        case 400:
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

//-----------------------------------------------------------------------------
    function timeFormat(cellvalue, options, rowObject) {
        return cellvalue != undefined ? cellvalue.slice(11, 19) : "";
    }

    function timeValidator(e) {
        var pattern = /^[0-2][0-9]:[0-6][0-9]:[0-6][0-9]$/;
        if (!(pattern.test(e))) {
            return [false, "Неверный формат времени. (HH:mm:ss)"];
        }

        return [true, ""];
    }

    function dateFormat(cellvalue, options, rowObject) {
        return cellvalue != undefined ? cellvalue.slice(0, 10) : "";
    }

    function timeStampFormat(cellvalue, options, rowObject) {
        console.log(cellvalue)
        return cellvalue != undefined ? cellvalue.slice(0, 10)+" "+cellvalue.slice(11, 19) : "";
    }

    function SetPrimitive(colModel) {
        colModel.forEach(function(model) {
            switch (model.type) {
                case "date":
                    model.formatter = dateFormat;
                    break;
                case "time":
                    model.editrules.custom_func = timeValidator;
                    model.formatter = timeFormat;
                    break;
                case "datetime":
                    model.formatter = timeStampFormat;
            }
            if (["date", "time", "datetime"].indexOf(model.type) >= 0)
                model.searchoptions.dataInit = model.editoptions.dataInit =
                    function(elem) { $(elem)[model.type + "picker"](); };
        });
        return colModel;
    }

//-----------------------------------------------------------------------------
    function listEventTypes(dialogId, data) {
        console.log("listEventTypes: ", data);

        if (data["result"] !== "ok") {
            showErrorMsg(data["result"]);
            return;
        }

        for (i in data["data"]) {
            $("#"+dialogId+" select").append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }
    }

    function ImportForms(dialogId, gridId) {
        console.log("ImportForms");

        var id = getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("event_id", id)

        $("#"+dialogId+" select").empty();

        utils.postRequest(
            { "event_id": id },
            function(data) { listEventTypes(dialogId, data); },
            "/gridcontroller/geteventtypesbyeventid"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Импорт": function() {
                    var ids = [];
                    $("#"+dialogId+" select option:selected").each(function(i, selected) {
                       ids[i] = $(selected).val();
                    });
                    utils.postRequest(
                        { "event_id": id, "event_types_ids": ids },
                        function(data) { showServerPromtInDialog($("#"+dialogId), data["result"]); },
                        "/gridcontroller/importforms"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

//-----------------------------------------------------------------------------
    function listParams(dialogId, data) {
        console.log("listParams: ", data)

        if (data["result"] !== "ok") {
            showErrorMsg(data["result"]);
            return;
        }

        var select = $("<select/>", { multiple: "multiple" });
        for (i in data["data"]) {
            select.append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }
        $("#"+dialogId).append(select);
    }

    function GetPersons(dialogId, gridId) {
        var id = getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("event_id", id)

        $("#"+dialogId).empty();

        utils.postRequest(
            { "event_id": id },
            function(data) { listParams(dialogId, data); },
            "/gridcontroller/getparamsbyeventid"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Получить список участников": function() {
                    var url = "/gridcontroller/getpersonsbyeventid?event="+id+"&params=";

                    $("#"+dialogId+" select option:selected").each(function(i, selected) {
                       url += $(selected).val() + ",";
                    });

                    url = url.slice(0, url.length-1);

                    location.href = url;
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        showErrorMsg: showErrorMsg,
        getCurrRowId: getCurrRowId,
        errTextFormat: errTextFormat,
        showServerPromtInGrid: showServerPromtInGrid,
        showServerPromtInDialog: showServerPromtInDialog,
        resizeSelectWidth: resizeSelectWidth,

        SetPrimitive: SetPrimitive,

        ImportForms: ImportForms,
        GetPersons: GetPersons,
    };

});
