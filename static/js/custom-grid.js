define(["utils", "datepicker/datepicker", "blank", "grid-utils"],
function(utils, datepicker, blank, gridUtils) {

    function timePicker(e) {
        $(e).timepicker({"timeFormat": "HH:mm:ss"});
    }

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
        return cellvalue != undefined ?
                cellvalue.slice(0, 10)+" "
                +cellvalue.slice(11, 19)
            :
                "";
    }

    function SetPrimitive(colModel) {
        for (i = 0; i < colModel.length; ++i) {
            if (colModel[i].type != undefined && colModel[i].type === "date") {
                colModel[i]["editoptions"]["dataInit"] = datepicker.initDatePicker;
                colModel[i]["searchoptions"]["dataInit"] = datepicker.initDatePicker;
                colModel[i]["formatter"] = dateFormat;
            } else if (colModel[i].type != undefined && colModel[i].type === "time") {
                colModel[i]["editrules"]["custom_func"] = timeValidator;
                colModel[i]["editoptions"]["dataInit"] = timePicker;
                colModel[i]["searchoptions"]["dataInit"] = timePicker;
                colModel[i]["formatter"] = timeFormat;
            } else if (colModel[i].type != undefined && colModel[i].type === "timestamp") {
                // datetimepicker
                colModel[i]["editoptions"]["dataInit"] = datepicker.initDatePicker;
                colModel[i]["searchoptions"]["dataInit"] = datepicker.initDatePicker;
                colModel[i]["formatter"] = timeStampFormat;
            }
            continue;
        }
        return colModel;
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

        function collbackSUB(data) {
            console.log("collbackSUB: ", data);

            subTableCaption = data["caption"];
            subTableName    = data["name"];
            subColNames     = data["colnames"];
            subColumns      = data["columns"];
            subColModel     = SetPrimitive(data["colmodel"]);
        }

        utils.postRequest(
            { "table": tableName, "id": row_id, "index": index },
            collbackSUB,
            "/gridhandler/getsubtable"
        );

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
            multiselectWidth: 20,
            multiboxonly: true,
            editurl:     "/gridhandler/editgridrow/" + subTableName,
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
            $("#" + subTId).jqGrid(
                "navButtonAdd",
                "#" + subPId,
                {
                    caption: "", buttonicon: "ui-icon-contact", title: "Редактировать анкету участника группы",
                    onClickButton: function() {
                        var personId = gridUtils.getCurrRowId(subTId);
                        if (!personId) return false;
                        var faceId = $("#" + subTId).jqGrid("getCell", personId, "face_id");
                        blank.ShowPersonBlankFromGroup(row_id, faceId, "dialog-group-person-request");
                    }
                }
            );
        }
    }

    return {
        AddSubTable: AddSubTable,
        SetPrimitive: SetPrimitive,
    };

});
