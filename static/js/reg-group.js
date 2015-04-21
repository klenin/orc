define(["utils", "grid-utils"], function(utils, gridUtils) {

    function listEvents(dialogId, data) {
        console.log("listEvents: ", data);

        if (data["result"] !== "ok") {
            gridUtils.showErrorMsg(data["result"]);
            return;
        }

        var select = $("<select/>");

        for (i in data["data"]) {
            select.append($("<option/>", {
                value: data["data"][i]["id"],
                text: data["data"][i]["name"],
            }));
        }

        $("#"+dialogId).append(select);
    }

    function Registration(dialogId, gridId) {
        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        $("#"+dialogId).empty();

        utils.postRequest(
            { "table": "events", "fields": ["id", "name"] },
            function(data) { listEvents(dialogId, data); },
            "/handler/getlist"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Участвовать в мероприятии": function() {
                    var event_id = $("#"+dialogId+" select").find(":selected").attr("value");

                    utils.postRequest(
                        { "group_id": id, "event_id": event_id },
                        function(data) { gridUtils.showServerPromtInDialog(dialogId, data["result"]); },
                        "/gridhandler/reggroup"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        Registration: Registration,
    };

});
