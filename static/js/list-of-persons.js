define(["utils", "grid-utils"], function(utils, gridUtils) {

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
                    var url ="/gridhandler/getpersonsbyeventid?event="+id+"&params=";

                    $("#"+dialogId+" select option:selected").each(function(i, selected) {
                       url += $(selected).val() + ",";
                    });

                    url = url.slice(0, url.length-1);

                    location.href = url;
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