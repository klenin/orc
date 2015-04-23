define(["utils", "grid-utils", "datepicker/datepicker", "kladr/kladr"], function(utils, gridUtils, datepicker, kladr) {

    function drawParam(data, for_saving) {
        console.log("drawParam");

        var block;

        if (data["type"] === "region"
            || data["type"] === "district"
            || data["type"] === "city"
            || data["type"] === "street"
            || data["type"] === "building"
            || data["type"] === "text"
            || data["type"] === "email"
            || data["type"] === "password"
            || data["type"] === "phon") {
            block = $("<input/>", {type: data["type"]});
            block.attr("id", data["param_id"]);
            if (data["value"]) { block.val(data["value"]); block.attr("param_val_id", data["param_val_id"]); }
            block.attr("for-saving", for_saving);

        } else if (data["type"] === "textarea") {
            block = $("<textarea/>", {});
            block.attr("id", data["param_id"]);
            if (data["value"]) { block.val(data["value"]); block.attr("param_val_id", data["param_val_id"]); }
            block.attr("for-saving", for_saving);

        } else if (data["type"] === "date") {
            block = $("<input/>", {type: "date"});
            block.attr("id", data["param_id"]);
            block.attr("for-saving", for_saving);
            if (data["value"]) { block.val(data["value"]); block.attr("param_val_id", data["param_val_id"]); }
            datepicker.initDatePicker(block);
        }

        var lable = $("<label/>", {
            text: data["param_name"]
        });

        return $("<p/>").append(lable).append(block);
    }

    function getFormData(name) {
        console.log("getFormData");

        var values = [];
        var empty = false;

        $("#"+name+" [for-saving=true]").each(function() {
            if ($(this).val() == "") {
                empty = true;
            }

            values.push({
                "value": $(this).val(),
                "param_val_id": $(this).attr("param_val_id"),
                "id": $(this).attr("id"),
                "event_type_id": $(this).attr("event_type_id"),
            });
        });

        return empty ? false : values;
    }

    function getListHistoryEvents(historyDiv, F_ids) {
        console.log("getListHistoryEvents: F_ids: ", F_ids);

        utils.postRequest(
            { "form_ids": F_ids, },
            function(response) {
                console.log("getListHistoryEvents: ", response);

                if (response["data"]) {
                    var data = response["data"];
                    for (var i = 0; i < data.length; ++i) {
                        $("#"+historyDiv+" select").append($("<option></option>")
                            .attr("value", data[i]["id"])
                            .text(data[i]["name"])
                        );
                    }
                    $("#"+historyDiv).show();
                }
            },
            "/handler/getlisthistoryevents"
        );
    }

    function ShowPersonBlankFromGroup(group_reg_id, dialogId, gridId) {
        console.log("ShowPersonBlankFromGroup");

        var person_id = gridUtils.getCurrRowId(gridId);
        if (person_id == -1) return false;

        var data = { "group_reg_id": group_reg_id, "person_id": person_id };
        console.log("ShowPersonBlankFromGroup: ", data);

        $("#"+dialogId).empty();

        utils.postRequest(
            data,
            function(data) {
                ShowBlank(data["data"], dialogId);
                $("#"+dialogId+" #history").hide();
            },
            "/gridhandler/getrequest"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить изменения": function() {
                    var values = getFormData(dialogId+" div");
                    if (values == false) {
                        alert("Не все поля заполнены.");
                        return;
                    }

                    utils.postRequest(
                        { "data": values },
                        function(data) { gridUtils.showServerPromtInDialog(dialogId, data["result"]); },
                        "/gridhandler/editparams"
                    );
                },
                "Отмена": function() {
                    $(this).empty();
                    $(this).dialog("close");
                },
            }
        });
    }

    function ShowBlank(d, dialogId) {
        console.log("ShowBlank data: ", d);

        if (d.length == 0) {
            alert("Нет данных.");
            return;
        }

        var history = $("<p/>", {id: "history"})
            .append($("<h5/>", {text: "Ранее заполненные анкеты"}))
            .append($("<select/>", {}))
            .append($("<input/>", {type: "button", value: "выбрать", id: "send-btn", name: "submit"}));

        $("#"+dialogId).append(history).append($("<h1/>")).append($("<div/>"));

        var F_ids = {};

        $("#" + dialogId + " h1").text(d[0]["event_name"]);

        var div_forms = $("<div/>", {id: "event-" + d[0]["event_id"]});
        var ul_forms = $("<ul/>", {});

        $(div_forms).append(ul_forms);
        $(div_forms).appendTo("#" + dialogId + " div");

        F_ids["form_id"] = [];

        for (i = 0; i < d.length; ++i) {
            if ($("#" + dialogId +" div#form-" + d[i]["form_id"]).attr("id") == undefined) {
                var li_form = $("<li/>", {});
                var a_form = $("<a/>", {href: "#" + "form-" + d[i]["form_id"]}).text(d[i]["form_name"]);

                $(li_form).append(a_form);
                $(ul_forms).append(li_form);

                var div_tab_form = $("<div/>", {id: "form-" + d[i]["form_id"]});
                $(div_forms).append(div_tab_form);

                F_ids["form_id"].push(parseInt(d[i]["form_id"]));

                var div_params = $("<div/>", {id: "params-" + d[i]["form_id"]});
                $(div_tab_form).append(div_params);

                var table = $("<table/>");
                div_params.append(table);

            }

            var tr = $("<tr/>").appendTo($("#" + dialogId +" div#form-" + d[i]["form_id"] + " table"));
            var td_1 = $("<td/>").appendTo(tr);
            $(td_1).append(drawParam(d[i], true));
            tr.append($("<td/>", {id: "export-param-"+d[i]["param_id"]}));
            tr.append($("<td/>", {id: "export-val-"+d[i]["param_id"]}));
        }

        $("#" + dialogId + " #" + "event-" + d[0]["event_id"]).tabs();

        console.log("F_ids: ", F_ids);

        $("#"+dialogId+" #history #send-btn").click(function() {
            utils.postRequest(
                {
                    "event_id": $("#"+dialogId+" #history select").find(":selected").attr("value")
                },
                function(response) {
                    if (response["result"] !== "ok") {
                        ShowServerAns(-1, response, "now #server-answer");
                        return;
                    }
                    ExportDataLoad(response["data"], dialogId);
                },
                "/handler/gethistoryrequest"
            );
        });

        kladr.kladr();

        return F_ids;
    }

    function ShowServerAns(event_id, data, responseId) {
        console.log("ShowServerAns");

        if (data.result === "ok") {
            var msg = "Запрос успешно выполнен.";
            if (event_id != 1) {
                msg += " Ваша заявка на участие будет рассмотрена.";
            } else {
                msg += " Проверьте, пожалуйста, свою почту.";
            }
            $("#"+responseId).text(msg).css("color", "green");

        } else if (data.result === "loginExists") {
            $("#"+responseId).text("Такой логин уже существует.").css("color", "red");

        } else if (data.result === "badLogin") {
            $("#"+responseId).text("Логин может содержать латинские буквы и/или "
                + "цифры и иметь длину от 2 до 36 символов.").css("color", "red");

        } else if (data.result === "badPassword") {
            $("#"+responseId).text("Пароль должен иметь длину от 6 "
                + "до 36 символов.").css("color", "red");

        } else if (data.result === "notAuthorized") {
            $("#"+responseId).text("Пользователь не авторизован.").css("color", "red");

        } else if (data.result === "authorized") {
            $("#"+responseId).text("Пользователь уже авторизован.").css("color", "red");

        } else if (data.result === "badEmail") {
            $("#"+responseId).text("Проверьте правильность введенного Вами email.").css("color", "red");

        } else {
            $("#"+responseId).text(data.result).css("color", "red");
        }
    }

    function ShowPersonBlank(dialogId, gridId) {
        console.log("ShowPersonBlank");

        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("reg_id", id)

        $("#"+dialogId).empty();

        utils.postRequest(
            { "reg_id": id },
            function(data) {
                var f_ids = ShowBlank(data["data"], dialogId);
                getListHistoryEvents(dialogId+" #history", f_ids);
            },
            "/gridhandler/getpersonrequest"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить изменения": function() {
                    var values = getFormData(dialogId+" div");
                    if (values == false) {
                        alert("Не все поля заполнены.");
                        return;
                    }

                    utils.postRequest(
                        { "data": values },
                        function(data) { gridUtils.showServerPromtInDialog(dialogId, data["result"]); },
                        "/gridhandler/editparams"
                    );

                },
                "Отмена": function() {
                    $(this).empty();
                    $(this).dialog("close");
                },
            }
        });
    }

    function ExportDataLoad(data, dialogId) {
        console.log("ExportDataLoad: ", data);

        for (var i = 0; i < data.length; ++i) {
            var f_id = data[i]["form_id"];
            var p_id = data[i]["param_id"];
            var p_v = data[i]["value"];

            $("#"+dialogId+" #params-"+f_id +" table #export-param-"+p_id+" input").remove();
            $("#"+dialogId+" #params-"+f_id +" table #export-val-"+p_id+" p").remove();
            if (data[i]["value"] != "") {
                $("#"+dialogId+" #params-"+f_id +" table #export-val-"+p_id).append(drawParam(data[i], 0, false));

                $("<input/>", {
                    "id": "export-btn-"+f_id+"-"+p_id,
                    "type": "button",
                    "value": "экспорт",
                    "data-event-type-id": f_id,
                    "data-param-id": p_id,
                    "data-param-val": p_v,
                }).appendTo("#"+dialogId+" #params-"+f_id+" table #export-param-"+p_id);

                $("#export-btn-"+f_id+"-"+p_id).click(function() {
                    var f_id = $(this).attr("data-event-type-id");
                    var p_id = $(this).attr("data-param-id");
                    var p_v = $(this).attr("data-param-val");
                    $("#"+dialogId+" #params-"+f_id+" table #"+p_id).val(p_v);
                });
            }
        }
    }

    return {
        drawParam: drawParam,
        getFormData: getFormData,
        getListHistoryEvents: getListHistoryEvents,
        ShowBlank: ShowBlank,
        ExportDataLoad: ExportDataLoad,

        ShowPersonBlankFromGroup: ShowPersonBlankFromGroup,
        ShowServerAns: ShowServerAns,
        ShowPersonBlank: ShowPersonBlank,
    };

});
