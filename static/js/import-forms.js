define(["utils", "grid-utils"], function(utils, gridUtils) {

    function listEventTypes(dialogId, data) {
        console.log("listEventTypes: ", data);

        if (data["result"] !== "ok") {
            gridUtils.showErrorMsg(data["result"]);
            return;
        }

        for (i in data["data"]) {
            $("#"+dialogId+" select").append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }
    }

    function ImportForms(dialogId, gridId) {
        console.log("ImportForms");

        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        console.log("event_id", id)

        $("#"+dialogId+" select").empty();

        utils.postRequest(
            { "event_id": id },
            function(data) { listEventTypes(dialogId, data); },
            "/gridhandler/geteventtypesbyeventid"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Импорт": function() {
                    var ids = [];
                    $("#"+dialogId+" select option:selected").each(function(i, selected) {
                       ids[i] = $(selected).val();
                    });
                    utils.postRequest(
                        { "event_id": id, "event_types_ids": ids },
                        function(data) { gridUtils.showServerPromtInDialog(dialogId, data["result"]); },
                        "/gridhandler/importforms"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        ImportForms: ImportForms,
    };

});