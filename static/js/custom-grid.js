define(["utils", "datepicker/datepicker"],
function(utils, datepicker) {

//check: change password-------------------------------------------------------
    var valid = false;

    $("#password-1").blur(function() {
        var pattern = /^.{6,36}$/;
        if (pattern.test($(this).val())) {
            $(this).css({"border": "2px solid green"});
        } else {
            valid = false;
            $(this).css({"border": "2px solid red"});
        }
    });

    $("#password-2").blur(function() {
        var pattern = /^.{6,36}$/;
        if (pattern.test($(this).val()) && $(this).val() === $("#password-1").val()) {
            valid = true;
            $(this).css({"border": "2px solid green"});
        } else {
            valid = false;
            $(this).css({"border": "2px solid red"});
        }
    });

//-----------------------------------------------------------------------------
    function GetColModelItem(refData, refFields, field) {
        function timePicker(e) {
            $(e).timepicker({"timeFormat": "HH:mm"});
        }

        function timeFormat(cellvalue, options, rowObject) {
            return cellvalue.slice(11, 19);
        }

        var data = {};
        data["name"] = field;
        data["index"] = field;
        data["editable"] = true;
        data["editrules"] = {required: true};

        // if ((field.indexOf("id") > -1)) {
        if (field == "id") {
            data["editable"] = false;

        } else if (field.indexOf("date") > -1) {
            data["formatter"] = "date";
            data["editrules"].date = true;
            data["formatoptions"] = {srcformat: 'Y-m-d', newformat: 'Y-m-d'};
            data["editoptions"] = {dataInit: datepicker.initDatePicker};

        } else if (field == "time") {
            data["formatter"] = timeFormat;
            data["editrules"].time = true;
            data["editoptions"] = {dataInit: timePicker};

        } else if (field == "topicality") {
            data["formatter"] = "checkbox";
            data["edittype"] = "checkbox";
            data["editoptions"] = {value: "true:false"};
            data["formatoptions"] = {disabled: true};

        } else if (field == "url") {
            data["formatter"] = "link";
            data["editrules"].required = false;

        } else if (field == "avatar") {
            data["manual"] = true;
            data["edittype"] = "file";
            data["sortable"] = false;
            data["search"] = false;
            data["editoptions"] = {enctype: "multipart/form-data"};
        }

        if (refData[field] != null) {
            data["formatter"] = "select";
            data["edittype"] = "select";
            data["stype"] = "select";
            data["search"] = true;
            var str = "", f;
            for (var i = 0; i < refData[field].length; ++i) {
                for (var k = 0; k < refFields.length; ++k) {
                    if (refFields[k] in refData[field][i]) {
                        f = refFields[k];
                        console.log("GetColModelItem-f: ", f)
                        break;
                    }
                }
                str += refData[field][i]["id"]+":"+refData[field][i][f]+";";
            }
            data["editoptions"] = {value: str.slice(0, -1)};
            data["searchoptions"] = {value: ":Все;"+str.slice(0, -1)};
        }
        return data;
    }

    function ResetPassword() {
        var id = getRowId();
        if (id == -1) return false;

        $("#password-1, #password-2").val("");

        $("#dialog-confirm").dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить": function() {
                    if (valid) {
                        utils.postRequest(
                            {
                                "pass": $("#password-1").val(),
                                "id": id
                            },
                            function() {},
                            "/gridhandler/resetpassword"
                        );
                        $(this).dialog("close");
                    } else {
                        showErrorMsg("Неверные значения паролей.\n"
                            + "Пароль должен иметь длину от 6 до 36 символов.");
                    }
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    function AddSubTable(subgrid_id, row_id, index, subgrid_table_id, pager_id, tableName) {

        var subTableCaption = "";
        var subTableName    = "";
        var subColNames     = [];
        var subColModel     = [];
        var subData         = [];
        var subColumns      = [];
        var subRefData      = [];

        $("#" + subgrid_id).append(
            "<table id='" + subgrid_table_id + "' class='scroll'></table>"
            + "<div id='" + pager_id + "' class='scroll'></div></br>"
        );

        function collbackSUB(data) {
            subTableCaption = data["caption"]
            subTableName    = data["name"];
            subColNames     = data["colnames"];
            subColumns      = data["columns"];
            subData         = data["data"];
            subRefData      = data["refdata"];
            subRefFields    = data["reffields"]
        }

        utils.postRequest(
            { "table": tableName, "id": row_id, "index": index },
            collbackSUB,
            "/gridhandler/getsubtable"
        );

        for (var i = 0; i < subColumns.length; ++i) {
            subColModel.push(GetColModelItem(subRefData, subRefFields, subColumns[i]));
        }

        $("#" + subgrid_table_id).jqGrid({
            datatype:   "local",
            data:        subData,
            colNames:    subColNames,
            colModel:    subColModel,
            rowNum:      5,
            rowList:     [5, 10, 20, 50],
            pager:       pager_id,
            caption:     subTableCaption,
            sortname:    "num",
            sortorder:   "asc",
            height:      '100%',
            width:         $("#grid-table").width()-65,
            multiselect: true,
            editurl:     "/gridhandler/edit/" + subTableName,
        });

        $("#" + subgrid_table_id).navGrid(
            "#" + pager_id,
            {
                edit:    true,    //edittext:   "Редактировать",
                add:     true,    //addtext:    "Создать",
                del:     true,    //deltext:    "Удалить",
                refresh: false,
                view:    false,
                search:  false
            },
            {   //edit
                width: "100%",
                recreateForm: true,
                closeAfterEdit:     true,
                afterSubmit:        function(response, postdata) {
                                        $('#grid-table').jqGrid().trigger('reloadGrid');
                                        return [true, "", response.responseText];
                                    }
            },
            {   //add
                width: "100%",
                recreateForm: true,
                clearAfterAdd:      true,
                closeAfterAdd:      true,
                addedrow:           "last",
                afterSubmit:        function(response, postdata) {
                                        $('#grid-table').jqGrid().trigger('reloadGrid');
                                        return [true, "", response.responseText];
                                    }
            },
            {   //del
                closeAfterAdd:      true,
                viewPagerButtons:   false
            }
        );

    }

    function listEventTypes(data) {
        if (data["result"] !== "ok") {
            showErrorMsg(data["result"]);
            return;
        }

        for (i in data["data"]) {
            $("#dialog-confirm-import select").append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }
    }

    function ImportForms() {
        var id = getRowId();
        if (id == -1) return false;

        $("#dialog-confirm-import select").empty();

        utils.postRequest(
            { "event_id": id },
            listEventTypes,
            "/gridhandler/geteventtypesbyeventid"
        );

        $("#dialog-confirm-import").dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Импорт": function() {
                    var ids = [];
                    $("#dialog-confirm-import select option:selected").each(function(i, selected) {
                       ids[i] = $(selected).val();
                    });
                    utils.postRequest(
                        { "event_id": id, "event_types_ids": ids },
                        function() {},
                        "/gridhandler/importforms"
                    );
                    $(this).dialog("close");

                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    function listPersons(data) {
        console.log("listPersons: ", data)

        if (data["result"] !== "ok") {
            showErrorMsg(data["result"]);
            return;
        }

        var result = $("<div/>");

        for (i in data["data"]) {
            result.append($("<div/>", {
                id: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }

        w = window.open();
        w.document.title = "Участники";
        $(w.document.body).html(result);
    }

    function listParams(data) {
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
        $("#dialog-persons").append(select);
    }

    function GetPersons() {
        var id = getRowId();
        if (id == -1) return false;

        $("#dialog-persons").empty();

        utils.postRequest(
            { "event_id": id },
            listParams,
            "/gridhandler/getparamsbyeventid"
        );

        $("#dialog-persons").dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Получить список участников": function() {
                    var ids = [];
                    $("#dialog-persons select option:selected").each(function(i, selected) {
                       ids[i] = $(selected).val();
                    });
                    utils.postRequest(
                        { "event_id": id, "params_ids": ids },
                        listPersons,
                        "/gridhandler/getpersonsbyeventid"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
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

    function showServerPromt(prompt) {
        var myInfo = '<div class="ui-state-highlight ui-corner-all">'
            + '<span class="ui-icon ui-icon-info" style="float: left; margin-right: .3em;"></span>'
            + '<strong>'+ prompt + '</strong><br/>' + '</div>';
        var infoTR = $("table#TblGrid_"+$("#grid-table")[0].id+">tbody>tr.tinfo");
        var infoTD = infoTR.children("td.topinfo");

        infoTD.html(myInfo);
        infoTR.show();

        setTimeout(function() {
            infoTD.children("div").fadeOut('slow',
                function() {
                    infoTR.hide();
            });
        }, 3000);
    }

    function getRowId() {
        var id = $("#grid-table").jqGrid("getGridParam", "selarrrow");

        if (id.length > 1 || id.length == 0) {
            showErrorMsg("<strong>Выберите одну запись.</strong>");
            return -1;
        }

        return id[0];
    }

    function DrowPersonRequest(data) {
        if (data["result"] !== "ok") {
            showErrorMsg(data["result"]);
            return;
        }

        for (i in data["data"]) {
            var row = $("<div/>");

            row.append($("<div/>", {
                text: data["data"][i]["name"]+": "+data["data"][i]["value"],
            }));

            $("#dialog-persons-request").append(row);
        }
    }

    function ShowPersonsRequest() {
        var id = getRowId();
        if (id == -1) return false;

        var reg_id = $("#grid-table").jqGrid("getCell", id, "reg_id");
        var event_id = $("#grid-table").jqGrid("getCell", id, "event_id");

        $("#dialog-persons-request").empty();

        utils.postRequest(
            { "reg_id": reg_id, "event_id": event_id },
            DrowPersonRequest,
            "/gridhandler/getpersonrequest"
        );

        $("#dialog-persons-request").dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
        });
    }

    function ConfirmOrRejectPersonRequest(confirm) {
        var id = getRowId();
        if (id == -1) return false;

        $("#grid-table").jqGrid("setCell", id, "status", false);

        var reg_id = $("#grid-table").jqGrid("getCell", id, "reg_id");
        var event_id = $("#grid-table").jqGrid("getCell", id, "event_id");

        console.log("ConfirmOrRejectPersonRequest");
        console.log({ "reg_id": reg_id, "event_id": event_id, "confirm": confirm});

        utils.postRequest(
            { "reg_id": reg_id, "event_id": event_id, "confirm": confirm},
            function(data) {
                showServerPromt(data["result"]);
            },
            "/gridhandler/confirmorrejectpersonrequest"
        );
    }

    return {
        GetColModelItem: GetColModelItem,
        ResetPassword: ResetPassword,
        AddSubTable: AddSubTable,
        ImportForms: ImportForms,
        GetPersons: GetPersons,
        ShowPersonsRequest: ShowPersonsRequest,
        ConfirmOrRejectPersonRequest: ConfirmOrRejectPersonRequest,
    };

});
