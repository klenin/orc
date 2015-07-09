define(["utils", "grid_lib", "blank"], function(utils, gridLib, blank) {

    function Register(dialogId, groupId, events) {
        var teamEvent = false;
        $("#"+dialogId).empty();
        $("#"+dialogId)
            .append(
                $("<p/>", { style: "color:red;" }).append(
                    "<strong>"
                    +"Внимание! После регистрации группы в мероприятии "
                    +"Вы не сможете удалять или добавлять участников группы "
                    +"во избежание потери информации. "
                    +"Участники, которые не подтвердили запрос для присоединения к группе, "
                    +"не будут зарегистрированны в мепроприятии."
                    +"</strong>")
            ).append(
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
            ).append(
                $("<p/>")
                    .append("<table id=\"events-table\"></table>")
                    .append("<div id=\"events-table-pager\"></div>")
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

        $("#"+dialogId+" #events-table").jqGrid({
            url: "/gridcontroller/load/events",
            datatype: "json",
            mtype: "POST",
            treeGrid: false,
            colNames: events.ColNames,
            colModel: gridLib.SetPrimitive(events.ColModel),
            pager: "#events-table-pager",
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
            sortorder: "asc",
            loadError: function(jqXHR, textStatus, errorThrown) {
                alert('HTTP status code: '+jqXHR.status+'\n'
                    +'textStatus: '+textStatus+'\n'
                    +'errorThrown: '+errorThrown);
                alert('HTTP message body: '+jqXHR.responseText);
            },
        });

        $("#"+dialogId+" #events-table").jqGrid("hideCol", ["id"]);

        $("#"+dialogId+" #events-table").navGrid(
           "#events-table-pager",
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

        $("#"+dialogId).dialog({
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Участвовать в мероприятии": function() {
                    var eventId = $("#"+dialogId+" #events-table").jqGrid("getGridParam", "selrow");
                    var eventType = $("#"+dialogId+" #events-table").jqGrid("getCell", eventId, "team") === "true" ? true : false;
                    console.log("eventType: ", eventType, "  teamEvent: ", teamEvent);
                    if (!eventId) {
                        console.log("Выберите запись");
                        return false;
                    } else if (eventType != teamEvent) {
                        console.log("Режим регистрации не соответсвует типу мероприятия");
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
        block.append(lf).append(f);

        var li = $("<label/>", { "text": "Имя" });
        var i = $("<input/>", { "id": 6, "for-saving": true, "required": true });
        block.append(li).append(i);

        var lo = $("<label/>", { "text": "Отчество" });
        var o = $("<input/>", { "id": 7, "for-saving": true, "required": true });
        block.append(lo).append(o);

        var le = $("<label/>", { "text": "Email" });
        var e = $("<input/>", { "id": 4, "for-saving": true, "required": true });
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
                    if (!values) {
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
