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
        var id = $("#grid-table").jqGrid("getGridParam", "selarrrow");
        if (id.length > 1 || id.length == 0) {
            $("#error").empty();
            $("#error").append("<strong>Выберите одну запись.</strong>");
            $("#error").dialog({
                model: true,
                buttons: {
                    "Закрыть": function() {
                        $(this).dialog("close");
                    }
                }
            });
            return false;
        }

        $("#password-1, #password-2").val("");

        $("#dialog-confirm").dialog({
            modal: true,
            toTop: "150",
            buttons: {
                "Сохранить": function() {
                    if (valid) {
                        utils.postRequest(
                            {
                                "pass": $("#password-1").val(),
                                "id": id[0]
                            },
                            function() {},
                            "/gridhandler/resetpassword"
                        );
                        $(this).dialog("close");
                    } else {
                        $("#error").empty();
                        $("#error").append("Неверные значения паролей.\n"
                            + "Пароль должен иметь длину от 6 до 36 символов.");
                        $("#error").dialog({
                            model: true,
                            buttons: {
                                "Закрыть": function() {
                                    $(this).dialog("close");
                                }
                            }
                        });
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
            {
                "table": tableName,
                "id": row_id,
                "index": index
            },
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
                //viewPagerButtons:   false,
                closeAfterEdit:     true,
                afterSubmit:        function(response, postdata) {
                                        $('#grid-table').jqGrid().trigger('reloadGrid');
                                        return [true, "", response.responseText];
                                    }
            },
            {   //add
                width: "100%",
                recreateForm: true,
                //viewPagerButtons:   false,
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

/*          $("#" + subgrid_table_id).jqGrid (
            "navButtonAdd",
            "#" + pager_id,
            {
                caption: "", buttonicon: "ui-icon-calculator", title: "Выбрать колонки",
                onClickButton: function() {
                    $("#"+subgrid_table_id).jqGrid("columnChooser", {
                        done: function(perm) {
                            if (perm) {
                                $("#" + subgrid_table_id).jqGrid("remapColumns", perm, true);
                            }
                        }
                    });
                }
            }
        );
*/
        $("#" + subgrid_table_id).jqGrid(
            "filterToolbar",
            {
                stringResult:  true,
                searchOnEnter: false,
                defaultSearch: "cn"
            }
        );

}

    return {
        GetColModelItem: GetColModelItem,
        ResetPassword: ResetPassword,
        AddSubTable: AddSubTable
    };

});
