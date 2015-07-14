define(["utils", "grid_lib", "datepicker/datepicker", "kladr/kladr"],
function(utils, gridLib, datepicker, kladr) {

    function drawParam(data, for_saving, admin) {
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

        } else if (data["type"] === "textarea") {
            block = $("<textarea/>", {});

        } else if (data["type"] === "date") {
            block = $("<input/>", {type: "date"});
            datepicker.initDatePicker(block);
        }

        block.attr("id", data["param_id"]);
        block.attr("for-saving", for_saving);
        block.attr("name", data["param_name"]);

        if (data["value"]) {
            block.val(data["value"]);
            block.attr("param_val_id", data["param_val_id"]);
        }

        if (data["required"]) {
            block.attr("required", true);
        }

        var lable = $("<label/>", {
            text: data["param_name"],
        });

        if (!data["editable"] && !admin) {
            block.attr("readonly", true);
        }

        return $("<p/>").append(lable).append(block);
    }

    function showParam(data, forSaving, admin) {
        console.log("showParam");

        var block = $("<div/>", {style: "border: 1px solid #4c9ac3;"});

        block.attr("id", data["param_id"]);
        block.attr("for-saving", forSaving);
        block.attr("name", data["param_name"]);

        block.text(data["value"]);
        block.attr("param_val_id", data["param_val_id"]);

        var lable = $("<label/>", {
            text: data["param_name"],
        });

        block.attr("readonly", true);

        return $("<p/>").append(lable).append(block);
    }

    function getFormData(name) {
        console.log("getFormData");

        var values = [];
        var empty = false;
        var pattern = /^[ \t\v\r\n\f]{0,}$/;
        var data = $("#"+name+" [for-saving=true]");
        console.log(data);

        for (var i = 0; i < data.length; ++i) {
            var elem = data[i];
            if (pattern.test($(elem).val()) && $(elem).attr("required")) {
                alert("Поле '" + $(elem).attr("name") + "' обязательное к заполнению.");
                empty = true;
                break;
            }

            values.push({
                "value": $(elem).val(),
                "param_val_id": $(elem).attr("param_val_id"),
                "id": $(elem).attr("id"),
            });
        };

        return empty ? false : values;
    }

    function getListHistoryEvents(historyDiv, formIds) {
        console.log("getListHistoryEvents: formIds: ", formIds);

        utils.postRequest(
            { "form_ids": formIds },
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
            "/blankcontroller/getlisthistoryevents"
        );
    }

    function ShowPersonBlankFromGroup(groupRegId, faceId, dialogId, formType) {
        console.log("ShowPersonBlankFromGroup");

        if (!groupRegId || !faceId) {
            return false;
        }

        var data = {
            "group_reg_id": groupRegId,
            "face_id": faceId,
            "personal": formType,
        };
        console.log("ShowPersonBlankFromGroup: ", data);

        $("#"+dialogId).empty();

        utils.postRequest(
            data,
            function(data) {
                ShowBlank(data["data"], dialogId, data["role"], data["regId"].toString(), formType, drawParam);
                $("#"+dialogId+" #history").hide();
            },
            "/blankcontroller/getpersonblankfromgroup"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить изменения": function() {
                    var values = getFormData(dialogId+" div");
                    if (!values) {
                        console.log("Не все поля заполнены.");
                        return false;
                    }

                    utils.postRequest(
                        { "data": values },
                        function(data) { gridLib.showServerPromtInDialog($("#"+dialogId), data["result"]); },
                        "/blankcontroller/editparams"
                    );
                },
                "Отмена": function() {
                    $(this).empty();
                    $(this).dialog("close");
                },
            }
        });

        return true;
    }

    function ShowBlank(d, dialogId, role, regId, formType, drawFunc) {
        console.log("ShowBlank data: ", d);
        console.log("ShowBlank role: ", role);
        console.log("ShowBlank formType: ", formType);

        if (d.length == 0) {
            return false;
        }

        var history = $("<p/>", {id: "history"})
            .append($("<b/>", {text: "Ранее заполненные анкеты"})).append("<br/>")
            .append($("<select/>", {}))
            .append($("<input/>", {type: "button", value: "выбрать", id: "send-btn", name: "submit"}));

        $("#"+dialogId).append(history);

        if (regId) {
            $("#"+dialogId)
                .append($("<input/>", {type: "checkbox", id: "edit-history-box", width: "auto"}))
                .append($("<label/>", {id: "edit-history", style: "display:inline;"}).text("Информация о редактировании полей"));

            $("#"+dialogId+" #edit-history-box").change(function() {
                console.log("ShowBlank: ", { "reg_id": regId });
                utils.postRequest(
                    { "reg_id": regId, "personal": formType },
                    function(response) {
                        if (response["result"] !== "ok") {
                            ShowServerAns(-1, response, "now #server-answer");
                            return false;
                        }

                        if ($("#"+dialogId+" #edit-history-box").is(":checked")) {
                            SetEditHistoryData(response["data"], dialogId);
                        } else {
                            ClearEditHistoryData(response["data"], dialogId);
                        }
                    },
                    "/blankcontroller/getedithistorydata"
                );
            });
        }

        $("#"+dialogId).append($("<h1/>")).append($("<div/>"));

        var formIds = [];

        $("#" + dialogId + " h1").text(d[0]["event_name"]);

        var div_forms = $("<div/>", {id: "event-" + d[0]["event_id"]});
        var ul_forms = $("<ul/>", {});

        $(div_forms).append(ul_forms);
        $(div_forms).appendTo("#" + dialogId + " div");

        for (i = 0; i < d.length; ++i) {
            if ($("#" + dialogId +" div#form-" + d[i]["form_id"]).attr("id") == undefined) {
                var li_form = $("<li/>", {});
                var a_form = $("<a/>", {href: "#" + "form-" + d[i]["form_id"]}).text(d[i]["form_name"]);

                $(li_form).append(a_form);
                $(ul_forms).append(li_form);

                var div_tab_form = $("<div/>", {id: "form-" + d[i]["form_id"]});
                $(div_forms).append(div_tab_form);

                formIds.push(parseInt(d[i]["form_id"]));

                var div_params = $("<div/>", {id: "params-" + d[i]["form_id"]});
                $(div_tab_form).append(div_params);

                var table = $("<table/>");
                div_params.append(table);

            }

            var tr = $("<tr/>").appendTo($("#" + dialogId +" div#form-" + d[i]["form_id"] + " table"));
            var td_1 = $("<td/>").appendTo(tr);
            $(td_1).append(drawFunc(d[i], true, role));
            tr.append($("<td/>", {id: "export-param-"+d[i]["param_id"]}));
            tr.append($("<td/>", {id: "export-val-"+d[i]["param_id"]}));
            tr.append($("<td/>", {id: "export-edit-history-"+d[i]["param_id"]}));
        }

        $("#" + dialogId + " #" + "event-" + d[0]["event_id"]).tabs();

        console.log("formIds: ", formIds);

        $("#"+dialogId+" #history #send-btn").click(function() {
            utils.postRequest(
                {
                    "event_id": $("#"+dialogId+" #history select").find(":selected").attr("value")
                },
                function(response) {
                    if (response["result"] !== "ok") {
                        ShowServerAns(-1, response, "now #server-answer");
                        return false;
                    }
                    ExportDataLoad(response["data"], dialogId);
                },
                "/blankcontroller/gethistoryrequest"
            );
        });

        kladr.kladr();

        return formIds;
    }

    function ShowServerAns(event_id, data, responseId) {
        console.log("ShowServerAns");

        if (data.result === "ok") {
            var msg = "Запрос успешно выполнен. ";
            if (event_id != 1) {
                msg += "Ваша заявка на участие будет рассмотрена.";
            } else {
                msg += "На вашу электронную почту было отправлено письмо, содержащее ссылку для подтверждения регистрации. "
                + "Воспользуйтесь этой ссылкой для продолжения работы.";
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

        } else if (data.result === "Unauthorized") {
            $("#"+responseId).text("Пользователь не авторизован.").css("color", "red");

        } else if (data.result === "authorized") {
            $("#"+responseId).text("Пользователь уже авторизован.").css("color", "red");

        } else if (data.result === "badEmail") {
            $("#"+responseId).text("Проверьте правильность введенного Вами email.").css("color", "red");

        } else {
            $("#"+responseId).text(data.result).css("color", "red");
        }
    }

    function ShowPersonBlank(dialogId, regId) {
        console.log("ShowPersonBlank: reg_id = ", regId);
        $("#"+dialogId).empty();

        utils.postRequest(
            { "reg_id": regId },
            function(data) {
                var formIds = ShowBlank(data["data"], dialogId, data["role"], regId, "true", drawParam);
                if (!formIds) {
                    return false;
                }
                getListHistoryEvents(dialogId+" #history", formIds);
            },
            "/blankcontroller/getblankbyregid"
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
                        console.log("Не все поля заполнены.");
                        alert("Не все поля заполнены.");
                        return false;
                    }

                    utils.postRequest(
                        { "data": values },
                        function(data) { gridLib.showServerPromtInDialog($("#"+dialogId), data["result"]); },
                        "/blankcontroller/editparams"
                    );

                },
                "Отмена": function() {
                    $(this).empty();
                    $(this).dialog("close");
                },
            }
        });

        return true;
    }

    function ShowGroupBlank(groupRegId, dialogId) {
        if (!groupRegId) {
            return false;
        }

        var data = { "group_reg_id": groupRegId };
        console.log("ShowGroupBlank: ", data);

        $("#"+dialogId).empty();

        utils.postRequest(
            data,
            function(data) {
                ShowBlank(data["data"], dialogId, false, false, false, showParam);
                $("#"+dialogId+" #history").hide();
            },
            "/blankcontroller/getgroupblank"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Отмена": function() {
                    $(this).empty();
                    $(this).dialog("close");
                },
            }
        });

        return true;
    }

    function ExportDataLoad(data, dialogId) {
        console.log("ExportDataLoad: ", data);

        for (var i = 0; i < data.length; ++i) {
            var f_id = data[i]["form_id"];
            var p_id = data[i]["param_id"];
            var p_v = data[i]["value"];

            $("#"+dialogId+" #params-"+f_id +" table #export-param-"+p_id+" input").remove();
            $("#"+dialogId+" #params-"+f_id +" table #export-val-"+p_id+" p").remove();
            $("#"+dialogId+" #params-"+f_id +" table #export-param-"+p_id+" br").remove();
            if (data[i]["value"] != "") {
                $("#"+dialogId+" #params-"+f_id +" table #export-val-"+p_id).append(drawParam(data[i], 0, false));
                $("<br/>").appendTo("#"+dialogId+" #params-"+f_id+" table #export-param-"+p_id);
                $("<input/>", {
                    "id": "export-btn-"+f_id+"-"+p_id,
                    "type": "button",
                    "value": "←",
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

    function SetEditHistoryData(data, dialogId) {
        console.log("SetEditHistoryData: ", data);

        for (var i = 0; i < data.length; ++i) {
            var f_id = data[i]["form_id"];
            var p_id = data[i]["param_id"];
            var p_v = data[i]["edit_date"] ? data[i]["edit_date"].replace(/[T,Z]/g, " ")+" - "+data[i]["login"] : data[i]["login"];

            $("#"+dialogId+" #params-"+f_id +" table #export-edit-history-"+p_id+" div").remove();
            $("#"+dialogId+" #params-"+f_id +" table #export-edit-history-"+p_id).append($("<div/>"));
            $("#"+dialogId+" #params-"+f_id +" table #export-edit-history-"+p_id+" div").append($("<br/>"));
            $("#"+dialogId+" #params-"+f_id +" table #export-edit-history-"+p_id+" div").append($("<div/>", {text: p_v}));
        }
    }

    function ClearEditHistoryData(data, dialogId) {
        console.log("ClearEditHistoryData: ", data);

        for (var i = 0; i < data.length; ++i) {
            var f_id = data[i]["form_id"];
            var p_id = data[i]["param_id"];

            $("#"+dialogId+" #params-"+f_id +" table #export-edit-history-"+p_id+" div").remove();
        }
    }

    return {
        drawParam: drawParam,
        getFormData: getFormData,
        getListHistoryEvents: getListHistoryEvents,
        ShowBlank: ShowBlank,
        ShowPersonBlank: ShowPersonBlank,
        ShowGroupBlank: ShowGroupBlank,
        ShowPersonBlankFromGroup: ShowPersonBlankFromGroup,
        ShowServerAns: ShowServerAns,
    };

});
