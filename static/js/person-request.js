define(["utils", "grid-utils", "kladr/kladr"], function(utils, gridUtils, kladr) {

    function DrowPersonRequest(dialogId, data) {
        console.log("DrowPersonRequest: ", data);

        if (data["result"] !== "ok") {
            gridUtils.showErrorMsg(data["result"]);
            return;
        }

        var i;
        for (i in data["data"]) {
            var row = $("<div/>");

            var label = $("<label/>", {
                text: data["data"][i].name,
            });

            var block = $("<input/>", {
                type: data["data"][i].type,
                id: data["data"][i].id,
            }).val(data["data"][i].value);

            if (data["data"][i].type === "date") {
                datepicker.initDatePicker(block)
            }

            row.append(label).append(block);

            $("#"+dialogId).append(row);
        }

        // kladr.kladr();
    }

    function ShowPersonsRequest(dialogId, gridId) {
        var id = gridUtils.getCurrRowId(gridId);
        if (id == -1) return false;

        var face_id = $("#"+gridId).jqGrid("getCell", id, "face_id");
        console.log("face_id", face_id)

        $("#"+dialogId).empty();

        utils.postRequest(
            { "face_id": face_id },
            function(data) { DrowPersonRequest(dialogId, data); },
            "/gridhandler/getpersonrequest"
        );

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Сохранить изменения": function() {
                    var values = [];

                    var $objs = $("#"+dialogId+" :input");
                    $objs.each(function() {
                        values.push({
                            "value": $(this).val(),
                            "id": $(this).attr("id"),
                        });
                    });

                    utils.postRequest(
                        { "data": values },
                        function(data) { gridUtils.showServerPromtInDialog(dialogId, data["result"]); },
                        "/gridhandler/editparams"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        ShowPersonsRequest: ShowPersonsRequest,
    };

});
