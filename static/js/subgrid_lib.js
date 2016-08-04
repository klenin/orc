define(["jquery", "utils", "blank", "grid_lib"],
function($, utils, blank, gridLib) {

    function AddSubTable(subgrid_id, row_id, index, tableName, gridId, data) {
        console.log("AddSubTable");

        var subTId = subgrid_id + "_t";
        var subPId = subgrid_id + "_p";

        $("#" + subgrid_id).append("<table id='" + subTId + "' class='scroll'></table><div id='" + subPId + "' class='scroll'></div></br>");

        var subTableCaption = "";
        var subTableName    = "";
        var subColNames     = [];
        var subColModel     = [];
        var subColumns      = [];

        function collbackSUB(data) {
            console.log("collbackSUB: ", data);

            subTableCaption = data["caption"];
            subTableName    = data["name"];
            subColNames     = data["colnames"];
            subColumns      = data["columns"];
            subColModel     = gridLib.SetPrimitive(data["colmodel"]);
        }

        utils.postRequest(
            { "table": tableName, "id": row_id, "index": index },
            collbackSUB,
            "/gridcontroller/getsubtable"
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
            // width:       /*$("#"+gridId).width()*/"100%",
            multiselect: true,
            multiselectWidth: 20,
            multiboxonly: true,
            editurl:     "/gridcontroller/editgridrow/" + subTableName,
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
                afterShowForm: function(formId) {
                    if (subTableName !== "persons") {
                        return;
                    }

                    var groupId = -1;
                    if (tableName === "group_registrations") {
                        groupId = $("#"+gridId).jqGrid("getCell", row_id, "group_id");
                    } else if (tableName === "groups") {
                        groupId = row_id;
                    }

                    var grSelect = $($($($('#tr_group_id', formId)[0]).children()[1]).find('select'));
                    grSelect.val(groupId);

                    var faseTd2 = $($($('#tr_face_id', formId)[0]).children()[1]);
                    var ansSelect = $(faseTd2.find('select'));
                    ansSelect.empty();

                    var tr1 = $('<tr class="FormData" id="tr_AddInfoForChoosingParams">'
                        +'<td class="CaptionTD ui-widget-content">'
                        +'<p><b>Выберите параметры поиска:</b></p>'
                        +'</td></tr>').insertAfter($('#tr_face_id', formId).show());
                    var tdParams = $('<td><p>'
                        +'<table id="params-table"></table>'
                        +'<div id="params-table-pager"></div>'
                        +'</p></td>').appendTo(tr1);
                    var tr2 = $('<tr class="FormData" id="tr_AddInfoForChoosingFace">'
                        +'<td class="CaptionTD ui-widget-content">'
                        +'<p><b>Выберите участника:</b></p>'
                        +'</td></tr>').insertAfter($('#tr_AddInfoForChoosingParams', formId).show());
                    var tdFaces = $('<td><p>'
                        +'<table id="faces-table"></table>'
                        +'<div id="faces-table-pager"></div>'
                        +'</p></td>').appendTo(tr2);

                    var filter = {};

                    $("#params-table").jqGrid({
                        url: "/gridcontroller/load/"+data["params"].TableName,
                        datatype: "json",
                        mtype: "POST",
                        treeGrid: false,
                        colNames: data["params"].ColNames,
                        colModel: gridLib.SetPrimitive(data["params"].ColModel),
                        pager: "#params-table-pager",
                        gridview: true,
                        viewrecords: true,
                        height: "100%",
                        width: "auto",
                        rowNum: 1,
                        rownumWidth: 20,
                        caption: data["params"].Caption,
                        sortname: "id",
                        sortorder: "asc",
                        loadError: function (jqXHR, textStatus, errorThrown) {
                            alert('HTTP status code: '+jqXHR.status+'\n'
                                +'textStatus: '+textStatus+'\n'
                                +'errorThrown: '+errorThrown);
                            alert('HTTP message body (jqXHR.responseText): '+'\n'+jqXHR.responseText);
                        },
                        loadComplete: function() {
                            $("#faces-table").trigger('reloadGrid');
                        },
                        beforeRequest: function() {
                            filter = $("#params-table").getGridParam("postData").filters;
                        }
                    });

                    $("#params-table").jqGrid("hideCol", ["id"]);
                    $("#params-table").jqGrid("hideCol", ["date"]);

                    $("#params-table").navGrid(
                        "#params-table-pager",
                        {   // buttons
                            edit: false,
                            add: false,
                            del: false,
                            search: true,
                            refresh: false,
                            view: false,
                        }, {}, {}, {},
                        {   // search
                            multipleGroup: true,
                            closeOnEscape: true,
                            multipleSearch: true,
                            closeAfterSearch: true,
                            showQuery: true,
                    });

                    formId.bind("resize", function() {
                        $("#params-table").setGridWidth(formId.width()-100, true);
                    }).trigger("resize");

                    $("#faces-table").jqGrid({
                        url: "/gridcontroller/load/search",
                        datatype: "json",
                        mtype: "POST",
                        treeGrid: false,
                        colNames: data["faces"].ColNames,
                        colModel: gridLib.SetPrimitive(data["faces"].ColModel),
                        pager: "#faces-table-pager",
                        gridview: true,
                        viewrecords: true,
                        height: "100%",
                        width: "auto",
                        rowNum: 5,
                        rownumbers: true,
                        rownumWidth: 20,
                        rowList: [5, 10, 20, 50],
                        caption: data["faces"].Caption,
                        sortname: "id",
                        sortorder: "asc",
                        loadError: function (jqXHR, textStatus, errorThrown) {
                            alert('HTTP status code: '+jqXHR.status+'\n'
                                +'textStatus: '+textStatus+'\n'
                                +'errorThrown: '+errorThrown);
                            alert('HTTP message body: '+jqXHR.responseText);
                        },
                        beforeRequest: function() {
                            $("#faces-table").setGridParam({ postData: {
                                "filters": filter ? filter : null,
                            } });
                        },
                        onSelectRow: function(faceId) {
                            ansSelect.empty();
                            var fio = $($($("#faces-table").find('tr#'+faceId)[0]).find('td')[1]).text();
                            console.log("fio: ", fio);
                            var option = $("<option/>", { value: faceId, text:  fio })
                            ansSelect.append(option);
                        }
                    });

                    $("#faces-table").navGrid(
                        "#faces-table-pager",
                        {   // buttons
                            edit: false,
                            add: false,
                            del: false,
                            view: false,
                            search: false,
                            refresh: false,
                    });

                    formId.bind("resize", function() {
                        $("#faces-table").setGridWidth(formId.width()-100, true);
                    }).trigger("resize");
                }
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
                        var personId = gridLib.getCurrRowId(subTId);
                        if (!personId) return false;
                        var faceId = $("#" + subTId).jqGrid("getCell", personId, "face_id");
                        blank.showPersonBlankFromGroup(row_id, faceId, "dialog-group-person-request", "true");
                    }
                }
            );
        }
    }

    return {
        AddSubTable: AddSubTable,
    };

});
