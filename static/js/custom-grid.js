define(["utils", "datepicker/datepicker", "blank"],
function(utils, datepicker, blank) {

    function GetColModelItem(refData, refFields, field) {
        console.log("GetColModelItem");

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
        data["searchoptions"] = {};

        // if ((field.indexOf("id") > -1)) {
        if (field == "id") {
            data["editable"] = false;

        // } else if (field.indexOf("id") > -1) {
        //     data["editrules"].edithidden = true;
        //     data["hidden"] = true;
        //     data["editable"] = false;

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
        data["searchoptions"]["sort"] = ["eq","ne","bw","cn"];

        return data;
    }

    function AddSubTable(subgrid_id, row_id, index, tableName, gridId) {
        console.log("AddSubTable");

        var subTId = subgrid_id + "_t";
        var subPId = subgrid_id + "_p";

        $("#" + subgrid_id).append("<table id='" + subTId + "' class='scroll'></table><div id='" + subPId + "' class='scroll'></div></br>");

        var subTableCaption = "";
        var subTableName    = "";
        var subColNames     = [];
        var subColModel     = [];
        var subData         = [];
        var subColumns      = [];
        var subRefData      = [];

        function collbackSUB(data) {
            console.log("collbackSUB: ", data);

            subTableCaption = data["caption"];
            subTableName    = data["name"];
            subColNames     = data["colnames"];
            subColumns      = data["columns"];
            subRefData      = data["refdata"];
            subRefFields    = data["reffields"];
        }

        utils.postRequest(
            { "table": tableName, "id": row_id, "index": index },
            collbackSUB,
            "/gridhandler/getsubtable"
        );

        for (var i = 0; i < subColumns.length; ++i) {
            subColModel.push(GetColModelItem(subRefData, subRefFields, subColumns[i]));
        }

        var url = "/handler/"+subTableName.replace(/_/g, '')+"load";
        if (tableName == "group_registrations") {
            var group_id = $("#"+gridId).jqGrid("getCell", row_id, "group_id");
            url += "/"+group_id;

        } else if (tableName == "groups") {
            url += "/"+row_id;
        }


        $("#" + subTId).jqGrid({
            url:         url,
            datatype:    "json",
            mtype:       "POST",
            colNames:    subColNames,
            colModel:    subColModel,
            rowNum:      5,
            rowList:     [5, 10, 20, 50],
            pager:       subPId,
            caption:     subTableCaption,
            sortname:    "num",
            sortorder:   "asc",
            height:      "100%",
            width:       $("#"+gridId).width()-65,
            multiselect: true,
            editurl:     "/gridhandler/edit/" + subTableName,
        });

        $("#" + subTId).navGrid(
            "#" + subPId,
            {
                edit:    true,
                add:     true,
                del:     true,
                refresh: false,
                view:    false,
                search:  false
            },
            {
                width: "100%",
                recreateForm: true,
                closeAfterEdit: true,
            },
            {
                width: "100%",
                recreateForm: true,
                clearAfterAdd: true,
                closeAfterAdd: true,
                addedrow: "last",
            },
            {
                closeAfterAdd: true,
                viewPagerButtons: false
            }
        );

        if (tableName == "group_registrations" && subTableName == "persons") {
            var event_id = $("#"+gridId).jqGrid("getCell", row_id, "event_id");
            $("#" + subTId).jqGrid(
                "navButtonAdd",
                "#" + subPId,
                {
                    caption: "", buttonicon: "ui-icon-pencil", title: "Редактировать анкету участника группы",
                    onClickButton: function() {
                        blank.ShowPersonBlankFromGroup(row_id, event_id, "dialog-group-person-request", subTId);
                    }
                }
            );
        }
    }

    return {
        GetColModelItem: GetColModelItem,
        AddSubTable: AddSubTable,
    };

});
