define(["utils", "datepicker/datepicker", "person-request"],
function(utils, datepicker, personRequest) {

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

        } else if (field == "topicality" || field == "status" || field == "enabled") {
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
                        // break;
                    }
                }
                str += refData[field][i]["id"]+":"+refData[field][i][f]+";";
            }
            data["editoptions"] = {value: str.slice(0, -1)};
            data["searchoptions"] = {value: ":Все;"+str.slice(0, -1)};
        }
        return data;
    }

    function AddSubTable(subgrid_id, row_id, index, subgrid_table_id, pager_id, tableName, gridId) {
        console.log("AddSubTable")
        console.log(subgrid_id, row_id, index, subgrid_table_id, pager_id, tableName, gridId)

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
            console.log("collbackSUB")
            console.log(data)
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
            height:      "100%",
            width:       $("#grid-table").width()-65,
            multiselect: true,
            editurl:     "/gridhandler/edit/" + subTableName,
        });

        $("#" + subgrid_table_id).navGrid(
            "#" + pager_id,
            {
                edit:    true,
                add:     true,
                del:     true,
                refresh: false,
                view:    false,
                search:  false
            },
            {   //edit
                width: "100%",
                recreateForm: true,
                closeAfterEdit:     true,
                afterSubmit:        function(response, postdata) {
                                        $("#"+gridId).jqGrid().trigger('reloadGrid');
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
                                        $("#"+gridId).jqGrid().trigger('reloadGrid');
                                        return [true, "", response.responseText];
                                    }
            },
            {   //del
                closeAfterAdd:      true,
                viewPagerButtons:   false
            }
        );

        if (subTableName == "persons") {
            $("#" + subgrid_table_id).jqGrid(
                "navButtonAdd",
                "#" + pager_id,
                {
                    caption: "", buttonicon: "ui-icon-script", title: "Анкета участника",
                    onClickButton: function() { personRequest.ShowPersonsRequest("dialog-persons-request", subgrid_table_id); }
                }
            );
        }
    }

    return {
        GetColModelItem: GetColModelItem,
        AddSubTable: AddSubTable,
    };

});
