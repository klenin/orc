define(["utils", "grid-utils"], function(utils, gridUtils) {

    function Registration(dialogId, groupId, eventTableId) {
        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Участвовать в мероприятии": function() {
                    var eventId = $("#"+eventTableId).jqGrid("getGridParam", "selrow");
                    if (!eventId) return false;
                    var data = { "group_id": groupId, "event_id": eventId };
                    console.log("Registration group: ", data);
                    utils.postRequest(
                        data,
                        function(response) {
                            gridUtils.showServerPromtInDialog($("#"+dialogId), response["result"]);
                            if (response["result"] === "ok") {
                                window.location.reload();
                            }
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

    return {
        Registration: Registration,
    };

});
