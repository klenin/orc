define(["utils", "grid-utils"], function(utils, gridUtils) {

    function listPersons(data) {
        console.log("listPersons: ", data)

        if (data["result"] !== "ok") {
            gridUtils.showErrorMsg(data["result"]);
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

    function listParams(dialogId, data) {
        console.log("listParams: ", data)

        if (data["result"] !== "ok") {
            gridUtils.showErrorMsg(data["result"]);
            return;
        }

        var select = $("<select/>", { multiple: "multiple" });
        for (i in data["data"]) {
            select.append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }
        $("#"+dialogId).append(select);
    }

    function GetPersons(dialogId, gridId) {
        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("event_id", id)

        $("#"+dialogId).empty();

        utils.postRequest(
            { "event_id": id },
            function(data) { listParams(dialogId, data); },
            "/gridhandler/getparamsbyeventid"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Получить список участников": function() {
                    var ids = [];
                    $("#"+dialogId+" select option:selected").each(function(i, selected) {
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

    return {
        GetPersons: GetPersons,
    };

});