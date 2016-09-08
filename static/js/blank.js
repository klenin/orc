define(["jquery", "utils", "grid_lib", "datepicker", "timepicker", "kladr"],
function($, utils, gridLib) {

    function drawParam(data, admin) {
        console.log("drawParam");

        var block;

        switch (data["type"]) {
            case "textarea":
                block = $("<textarea/>", {});
                break;
            case "address":
                block = $("<input/>", {type: "text"}).kladr({
                    parentInput: null,
                    select: null,
                    oneString: true
                });
                break;
            case "region":
            case "district":
            case "city":
            case "street":
            case "building":
                block = $("<input/>", {type: "text"}).kladr({
                    type: $.kladr.type[data["type"]]
                });
                break;
            case "date":
            case "time":
            case "datetime":
                block = $("<input/>", {type: "text"});
                block.one("DOMNodeInserted", block[data["type"] + "picker"].bind(block));
                break;
            default:
                block = $("<input/>", {type: data["type"]});
        }

        block.attr("id", data["param_id"]);
        block.attr("for-saving", true);
        block.attr("name", data["param_name"]);

        if (data["value"]) {
            block.val(data["value"]);
            block.attr("param_val_id", data["param_val_id"]);
        }

        if (data["required"]) {
            block.attr("required", true);
        }

        var label = $("<label/>", {
            text: data["param_name"],
        });

        if (!data["editable"] && !admin) {
            block.attr("readonly", true);
        }

        block.change(function() {
            block.attr("wasChanged", true);
        });

        return $("<p/>").append(label).append(block);
    }

    function showParam(data, admin) {
        console.log("showParam");

        var block = $("<div/>");

        block.attr("id", data["param_id"]);
        block.attr("name", data["param_name"]);

        block.text(data["value"]);
        block.attr("param_val_id", data["param_val_id"]);

        var label = $("<label/>", {
            text: data["param_name"],
            style: "color: #CC6600; font-weight: bold;",
        });

        block.attr("readonly", true);

        return $("<p/>").append(label).append(block);
    }

    function getFormData(id) {
        console.log("getFormData");

        var values = [];
        var empty = false;
        var pattern = /^[ \t\v\r\n\f]{0,}$/;
        var data = $("#"+id+" [for-saving=true][wasChanged=true]");
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

    function showBlank(d, dialogId, admin, regId, formType, drawFunc) {
        console.log("showBlank data: ", d);
        console.log("showBlank admin: ", admin);
        console.log("showBlank formType: ", formType);

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
                console.log("showBlank: ", { "reg_id": regId });
                utils.postRequest(
                    { "reg_id": regId, "personal": formType },
                    function(response) {
                        if (response["result"] !== "ok") {
                            showServerAns(-1, response, "now #server-answer");
                            return false;
                        }

                        if ($("#"+dialogId+" #edit-history-box").is(":checked")) {
                            setEditHistoryData(response["data"], dialogId);
                        } else {
                            clearEditHistoryData(response["data"], dialogId);
                        }
                    },
                    "/blankcontroller/getedithistorydata"
                );
            });
        }

        $("#"+dialogId).append($("<h1/>")).append($("<div/>"));

        var formIds = [];

        $("#"+dialogId+" h1").text(d[0]["event_name"]);

        var divForms = $("<div/>", {id: "event-"+d[0]["event_id"]});
        var ulForms = $("<ul/>", {});

        $(divForms).append(ulForms);
        $(divForms).appendTo("#"+dialogId+" div");

        for (i = 0; i < d.length; ++i) {
            if ($("#"+dialogId +" div#form-"+d[i]["form_id"]).attr("id") == undefined) {
                var liForm = $("<li/>", {});
                var aForm = $("<a/>", {href: "#"+"form-"+d[i]["form_id"]}).text(d[i]["form_name"]);

                $(liForm).append(aForm);
                $(ulForms).append(liForm);

                var divTabForm = $("<div/>", {id: "form-"+d[i]["form_id"]});
                $(divForms).append(divTabForm);

                formIds.push(parseInt(d[i]["form_id"]));

                var divParams = $("<div/>", {id: "params-"+d[i]["form_id"]});
                $(divTabForm).append(divParams);

                var table = $("<table/>");
                divParams.append(table);

            }

            $.kladr.setDefault("parentInput", "#" + dialogId + " div#form-" + d[i]["form_id"]);

            var tr = $("<tr/>").appendTo($("#"+dialogId +" div#form-"+d[i]["form_id"]+" table"));
            var td = $("<td/>").appendTo(tr);
            $(td).append(drawFunc(d[i], admin));
            tr.append($("<td/>", {id: "export-param-"+d[i]["param_id"]}));
            tr.append($("<td/>", {id: "export-val-"+d[i]["param_id"]}));
            tr.append($("<td/>", {id: "export-edit-history-"+d[i]["param_id"]}));
        }

        $("#"+dialogId+" #"+"event-"+d[0]["event_id"]).tabs();

        console.log("formIds: ", formIds);

        $("#"+dialogId+" #history #send-btn").click(function() {
            utils.postRequest(
                {
                    "event_id": $("#"+dialogId+" #history select").find(":selected").attr("value")
                },
                function(response) {
                    if (response["result"] !== "ok") {
                        showServerAns(-1, response, "now #server-answer");
                        return false;
                    }
                    exportDataLoad(response["data"], dialogId);
                },
                "/blankcontroller/gethistoryrequest"
            );
        });

        return formIds;
    }

//-----------------------------------------------------------------------------
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

    function showPersonBlankFromGroup(groupRegId, faceId, dialogId, formType) {
        console.log("showPersonBlankFromGroup");

        if (!groupRegId || !faceId) {
            return false;
        }

        var data = {
            "group_reg_id": groupRegId,
            "face_id": faceId,
            "personal": formType,
        };
        console.log("showPersonBlankFromGroup: ", data);

        $("#"+dialogId).empty();

        utils.postRequest(
            data,
            function(data) {
                showBlank(data["data"], dialogId, data["role"], data["regId"].toString(), formType, drawParam);
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

    function showPersonBlank(dialogId, regId) {
        console.log("showPersonBlank: reg_id = ", regId);
        $("#"+dialogId).empty();

        utils.postRequest(
            { "reg_id": regId },
            function(data) {
                var formIds = showBlank(data["data"], dialogId, data["role"], regId, "true", drawParam);
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

    function showGroupBlank(groupRegId, dialogId) {
        if (!groupRegId) {
            return false;
        }

        var data = { "group_reg_id": groupRegId };
        console.log("showGroupBlank: ", data);

        $("#"+dialogId).empty();

        utils.postRequest(
            data,
            function(data) {
                showBlank(data["data"], dialogId, false, false, false, showParam);
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

//-----------------------------------------------------------------------------
    function exportDataLoad(data, dialogId) {
        console.log("exportDataLoad: ", data);

        for (var i = 0; i < data.length; ++i) {
            var formId = data[i]["form_id"];
            var paramId = data[i]["param_id"];
            var value = data[i]["value"];

            $("#"+dialogId+" #params-"+formId+" table #export-param-"+paramId+" input").remove();
            $("#"+dialogId+" #params-"+formId+" table #export-val-"+paramId+" p").remove();
            $("#"+dialogId+" #params-"+formId+" table #export-param-"+paramId+" br").remove();

            if (data[i]["value"] != "") {
                $("#"+dialogId+" #params-"+formId+" table #export-val-"+paramId).append(drawParam(data[i], 0, false));
                $("<br/>").appendTo("#"+dialogId+" #params-"+formId+" table #export-param-"+paramId);

                $("<input/>", {
                    "id": "export-btn-"+formId+"-"+paramId,
                    "type": "button",
                    "value": "←",
                    "data-event-type-id": formId,
                    "data-param-id": paramId,
                    "data-param-val": value,
                }).appendTo("#"+dialogId+" #params-"+formId+" table #export-param-"+paramId);

                $("#export-btn-"+formId+"-"+paramId).click(function() {
                    var formId = $(this).attr("data-event-type-id");
                    var paramId = $(this).attr("data-param-id");
                    var value = $(this).attr("data-param-val");
                    $("#"+dialogId+" #params-"+formId+" table #"+paramId).val(value);
                });
            }
        }
    }

//-----------------------------------------------------------------------------
    function setEditHistoryData(data, dialogId) {
        console.log("setEditHistoryData: ", data);

        for (var i = 0; i < data.length; ++i) {
            var formId = data[i]["form_id"];
            var paramId = data[i]["param_id"];
            var value = data[i]["edit_date"] ? data[i]["edit_date"].replace(/[T,Z]/g, " ")+" - "+data[i]["login"] : data[i]["login"];

            $("#"+dialogId+" #params-"+formId+" table #export-edit-history-"+paramId+" div").remove();
            $("#"+dialogId+" #params-"+formId+" table #export-edit-history-"+paramId).append($("<div/>"));
            $("#"+dialogId+" #params-"+formId+" table #export-edit-history-"+paramId+" div").append($("<br/>"));
            $("#"+dialogId+" #params-"+formId+" table #export-edit-history-"+paramId+" div").append($("<div/>", {text: value}));
        }
    }

    function clearEditHistoryData(data, dialogId) {
        console.log("clearEditHistoryData: ", data);

        for (var i = 0; i < data.length; ++i) {
            var formId = data[i]["form_id"];
            var paramId = data[i]["param_id"];

            $("#"+dialogId+" #params-"+formId+" table #export-edit-history-"+paramId+" div").remove();
        }
    }

//-----------------------------------------------------------------------------
    function showServerAns(event_id, data, responseId) {
        console.log("showServerAns");

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

    return {
        drawParam: drawParam,
        getFormData: getFormData,
        getListHistoryEvents: getListHistoryEvents,
        showBlank: showBlank,
        showPersonBlank: showPersonBlank,
        showGroupBlank: showGroupBlank,
        showPersonBlankFromGroup: showPersonBlankFromGroup,
        showServerAns: showServerAns,
    };

});
