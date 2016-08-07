define(["jquery", "utils", "grid_lib", "blank"], function($, utils, gridLib, blank) {

    function Register(dialogId, groupId, eventId, groups, events) {
        if ((!groupId && !eventId) || (!eventId && !events) || (!groupId && !groups)) {
            return false;
        }

        var teamEvent = false;
        $("#"+dialogId).empty();
        $("#"+dialogId).append(
            $("<p/>", { style: "color:red;" }).append(
                "<strong>"
                +"Внимание!<br/>После регистрации группы в мероприятии "
                +"Вы не сможете удалять или добавлять участников группы "
                +"во избежание потери информации.<br/>"
                +"Участники, которые не подтвердили запрос для присоединения к группе, "
                +"не будут зарегистрированны в мепроприятии."
                +"</strong>"
            )
        );

        if (!groupId) {
            $("#"+dialogId).append(
                $("<p/>")
                    .append("<table id=\"group-reg-groups-table\"></table>")
                    .append("<div id=\"group-reg-groups-table-pager\"></div>")
            );

            $("#group-reg-groups-table").jqGrid({
                url: "/handler/usergroupsload",
                datatype: "json",
                mtype: "POST",
                treeGrid: false,
                colNames: groups.ColNames,
                colModel: gridLib.SetPrimitive(groups.ColModel),
                pager: "#group-reg-groups-table-pager",
                gridview: true,
                sortname: "id",
                viewrecords: true,
                height: "100%",
                width: "auto",
                rowNum: 5,
                rownumbers: true,
                rownumWidth: 20,
                rowList: [5, 10, 20, 50],
                caption: "Мои группы",
                sortname: "id",
                sortorder: "asc",
                editurl: "/gridcontroller/editgridrow/"+groups.TableName,

                subGrid: groups.Sub,
                subGridOptions: {
                    "plusicon": "ui-icon-triangle-1-e",
                    "minusicon": "ui-icon-triangle-1-s",
                    "openicon": "ui-icon-arrowreturn-1-e",
                    "reloadOnExpand": true,
                    "selectOnExpand": true
                },
                subGridRowExpanded: function(subgrid_id, group_id) {
                    $("#"+subgrid_id).append("<table id='"+subgrid_id+"_t"+"' class='scroll'></table>"
                        +"<div id='"+subgrid_id+"_p"+"' class='scroll'></div></br>");

                    var addDelFlag = false;
                    utils.postRequest(
                        { "group_id": group_id },
                        function(data) {
                            console.log("Is Reg Group?");
                            console.log(data["result"]);
                            addDelFlag = data["addDelFlag"];
                        },
                        "/groupcontroller/isreggroup"
                    );

                    $("#"+subgrid_id+"_t").jqGrid({
                        url: "/handler/"+groups.SubTableName.replace(/_/g, "")+"load/"+group_id,
                        datatype: "json",
                        mtype: "POST",
                        colNames: groups.SubColNames,
                        colModel: gridLib.SetPrimitive(groups.SubColModel),
                        rowNum: 5,
                        rowList: [5, 10, 20, 50],
                        pager: "#"+subgrid_id+"_p",
                        caption: groups.SubCaption,
                        sortname: "num",
                        sortorder: "asc",
                        height: "100%",
                        width: $("#group-reg-groups-table").width()-65,
                        editurl: "/gridcontroller/editgridrow/"+groups.SubTableName,
                        gridComplete: function() {
                            var rows = $("#"+subgrid_id+"_t").getDataIDs();
                            for (var i = 0; i < rows.length; i++) {
                                var status = $("#"+subgrid_id+"_t").getCell(rows[i], "status");
                                if (status === "true") {
                                    $("#"+subgrid_id+"_t").jqGrid('setRowData', rows[i], false, "row-green");
                                }
                            }
                        },
                    });

                    $("#"+subgrid_id+"_t").jqGrid("hideCol", ["id"]);

                    $(window).bind("resize", function() {
                        $("#"+subgrid_id+"_t").setGridWidth($("#group-reg-groups-table").width(), true);
                    }).trigger("resize");
                }
            });

            $("#group-reg-groups-table").jqGrid("hideCol", ["id"]);
            $("#group-reg-groups-table").jqGrid("hideCol", ["face_id"]);

            $(window).bind("resize", function() {
                $("#group-reg-groups-table").setGridWidth($(window).width()-100, true);
            }).trigger("resize");
        }

        if (!eventId) {
            $("#"+dialogId).append(
                $("<p/>").append(
                    $("<input/>", {
                        type: "radio",
                        name: "group-registration-type",
                        width: "auto",
                        value: "one",
                        checked: "checked",
                    })
                ).append(
                    "Регистрация в мероприятии каждого участника группы</br>"
                ).append(
                    $("<input/>", {
                        type: "radio",
                        name: "group-registration-type",
                        width: "auto",
                        value: "many",
                    })
                ).append(
                    "Регистрация в мероприятии группы в качестве команды</br>"
                )
            );

            $("#"+dialogId+" input:radio[name=group-registration-type]").change(function() {
                if (this.value == 'one') {
                    console.log("one");
                    teamEvent = false;
                } else if (this.value == 'many') {
                    console.log("many");
                    teamEvent = true;
                }
            });

            $("#"+dialogId).append(
                $("<p/>")
                    .append("<table id=\"group-reg-events-table\"></table>")
                    .append("<div id=\"group-reg-events-table-pager\"></div>")
            );

            $("#"+dialogId+" #group-reg-events-table").jqGrid({
                url: "/gridcontroller/load/events",
                datatype: "json",
                mtype: "POST",
                treeGrid: false,
                colNames: events.ColNames,
                colModel: gridLib.SetPrimitive(events.ColModel),
                pager: "#group-reg-events-table-pager",
                gridview: true,
                sortname: "id",
                viewrecords: true,
                height: "100%",
                width: "auto",
                rowNum: 5,
                rownumbers: true,
                rownumWidth: 20,
                rowList: [5, 10, 20, 50],
                caption: events.Caption,
                sortname: "id",
                sortorder: "asc"
            });

            $("#"+dialogId+" #group-reg-events-table").jqGrid("hideCol", ["id"]);

            $("#"+dialogId+" #group-reg-events-table").navGrid(
               "#group-reg-events-table-pager",
                {   // buttons
                    edit: false,
                    add: false,
                    del: false,
                    refresh: false,
                    view: true,
                    search: true
                }, {}, {}, {},
                {   // search
                    multipleGroup: true,
                    closeOnEscape: true,
                    multipleSearch: true,
                    closeAfterSearch: true,
                    showQuery: true,
                }
            );
        }

        $("#"+dialogId).dialog({
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Участвовать в мероприятии": function() {
                    if (!eventId) {
                        eventId = $("#"+dialogId+" #group-reg-events-table").jqGrid("getGridParam", "selrow");
                        var eventType = $("#"+dialogId+" #group-reg-events-table").jqGrid("getCell", eventId, "team") === "true" ? true : false;
                        console.log("eventType: ", eventType, "  teamEvent: ", teamEvent);
                        if (!eventId) {
                            console.log("Выберите запись");
                            return false;
                        } else if (eventType != teamEvent) {
                            console.log("Режим регистрации не соответсвует типу мероприятия");
                            return false;
                        }
                    }

                    if (!groupId) {
                        groupId = $("#"+dialogId+" #group-reg-groups-table").jqGrid("getGridParam", "selrow");
                    }

                    if (!eventId || !groupId) {
                        console.log("Так не пойдет. ", "eventId: ", eventId, "  groupId: ", groupId);
                        return false;
                    }

                    var data = { "group_id": groupId, "event_id": eventId };
                    console.log("Register group: ", data);

                    utils.postRequest(
                        data,
                        function(response) {
                            gridLib.showServerPromtInDialog($("#"+dialogId), response["result"]);
                        },
                        "/groupcontroller/register"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    function AddPerson(dialogId, groupId) {
        $("#"+dialogId).empty();

        var block = $("<div/>");

        var lf = $("<label/>", { "text": "Фамилия" });
        var f = $("<input/>", { "id": 5, "for-saving": true, "required": true });
        f.change(function() {
            f.attr("wasChanged", true);
        });
        block.append(lf).append(f);

        var li = $("<label/>", { "text": "Имя" });
        var i = $("<input/>", { "id": 6, "for-saving": true, "required": true });
        i.change(function() {
            i.attr("wasChanged", true);
        });
        block.append(li).append(i);

        var lo = $("<label/>", { "text": "Отчество" });
        var o = $("<input/>", { "id": 7, "for-saving": true, "required": true });
        o.change(function() {
            o.attr("wasChanged", true);
        });
        block.append(lo).append(o);

        var le = $("<label/>", { "text": "Email" });
        var e = $("<input/>", { "id": 4, "for-saving": true, "required": true });
        e.change(function() {
            e.attr("wasChanged", true);
        });
        block.append(le).append(e);

        $("#"+dialogId).append(block);

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Добавить участника": function() {
                    var values = blank.getFormData(dialogId);
                    if (values.length < 4) {
                        return false;
                    }
                    var data = {"group_id": groupId, "data": values };
                    console.log("AddPerson: ", data);
                    utils.postRequest(
                        data,
                        function(response) {
                            gridLib.showServerPromtInDialog($("#"+dialogId), response["result"]);
                        },
                        "/groupcontroller/addperson"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        Register: Register,
        AddPerson: AddPerson,
    };

});
