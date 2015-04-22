define(["utils", "grid-utils", "datepicker/datepicker"], function(utils, gridUtils, datepicker) {

    function drawParam(data, for_saving) {
        var block;

        if (data["type"] === "region"
            || data["type"] === "district"
            || data["type"] === "city"
            || data["type"] === "street"
            || data["type"] === "building"
            || data["type"] === "text"
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
        console.log("F_ids: ", F_ids);

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
        var person_id = gridUtils.getCurrRowId(gridId);
        if (person_id == -1) return false;

        var data = { "group_reg_id": group_reg_id, "person_id": person_id };
        console.log("ShowPersonBlankFromGroup: ", data);

        $("#"+dialogId).empty();
        $("#"+dialogId).append($("<h1/>")).append($("<div/>"));

        utils.postRequest(
            data,
            function(data) {
                ShowBlank(data["data"], dialogId);
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

        var F_ids = {};

        $("#" + dialogId + " h1").text(d[0]["event_name"]);

        var div_forms = $("<div/>", {id: "event-" + d[0]["event_id"]});
        var ul_forms = $("<ul/>", {});

        $(div_forms).append(ul_forms);
        $(div_forms).appendTo("#" + dialogId + " div");

        F_ids["form_id"] = [];

        for (i = 0; i < d.length; ++i) {
            if ($("#" + dialogId +" div#form-" + d[i]["form_id"]).attr("id")==undefined) {
                console.log("создать")
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

        return F_ids;
    }

    function ShowServerAns(event_id, data, responseId) {
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

        } else {
            $("#"+responseId).text(data.result).css("color", "red");
        }
    }

    function ShowPersonBlank(dialogId, gridId) {
        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("reg_id", id)

        $("#"+dialogId).empty();

        $("#"+dialogId).empty();
        $("#"+dialogId).append($("<h1/>")).append($("<div/>"));

        utils.postRequest(
            { "reg_id": id },
            function(data) { ShowBlank(data["data"], dialogId); },
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

    return {
        drawParam: drawParam,
        getFormData: getFormData,
        getListHistoryEvents: getListHistoryEvents,
        ShowBlank: ShowBlank,

        ShowPersonBlankFromGroup: ShowPersonBlankFromGroup,
        ShowServerAns: ShowServerAns,
        ShowPersonBlank: ShowPersonBlank,
    };

});
