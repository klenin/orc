define(["utils", "datepicker/datepicker"], function(utils, datepicker) {

    function drawParam(data, event_type_id, for_saving) {
        var block;

        if (data["type"] === "region"
            || data["type"] === "district"
            || data["type"] === "city"
            || data["type"] === "street"
            || data["type"] === "building"
            || data["type"] === "text"
            || data["type"] === "password"
            || data["type"] === "phon") {
            block = $("<input/>", {type: data["type"]});
            block.attr("id", data["id"]);
            block.attr("event_type_id", event_type_id);
            block.val(data["value"]);
            block.attr("for-saving", for_saving);

        } else if (data["type"] === "textarea") {
            block = $("<textarea/>", {});
            block.attr("id", data["id"]);
            block.attr("event_type_id", event_type_id);
            block.val(data["value"]);
            block.attr("for-saving", for_saving);

        } else if (data["type"] === "date") {
            block = $("<input/>", {type: "date"});
            block.attr("id", data["id"]);
            block.attr("event_type_id", event_type_id);
            block.attr("for-saving", for_saving);
            block.val(data["value"]);
            datepicker.initDatePicker(block);
        }

        var lable = $("<label/>", {
            text: data["name"]
        });

        return $("<p/>").append(lable).append(block);
    }

    function getFormData(name) {
        var values = [];

        var $obj = $("#"+name+" [for-saving=true]");
        $obj.each(function() {
            values.push({
                "value": $(this).val(),
                "id": $(this).attr("id"),
                "event_type_id": $(this).attr("event_type_id"),
            });
        });

        return values;
    }

    function pushFormData(E, F, P, captionId, tabsId) {
        console.log("E: ", E);
        console.log("F: ", F);
        console.log("P: ", P);
        console.log("captionId: ", captionId);
        console.log("tabsId: ", tabsId);

        var F_ids = {};

        $("#" + captionId).text(E[0]["name"]);

        var div_forms = $("<div/>", {id: "event-" + E[0]["id"]});
        var ul_forms = $("<ul/>", {});

        $(div_forms).append(ul_forms);
        $(div_forms).appendTo("#" + tabsId);

        F_ids["form_id"] = []

        for (var i = 0; i < F.length; ++i) {

            var li_form = $("<li/>", {});
            var a_form = $("<a/>", {href: "#" + "form-" + F[i]["id"]}).text(F[i]["name"]);

            $(li_form).append(a_form);
            $(ul_forms).append(li_form);

            var div_tab_form = $("<div/>", {id: "form-" + F[i]["id"]});
            $(div_forms).append(div_tab_form);

            F_ids["form_id"].push(parseInt(F[i]["id"]));

            var div_params = $("<div/>", {id: "params-" + F[i]["id"]});
            $(div_tab_form).append(div_params);

            var table = $("<table/>");
            div_params.append(table);

            for (var k = 0; k < P[i].length; k++) {
                var tr = $("<tr/>").appendTo(table);
                var td_1 = $("<td/>").appendTo(tr);
                $(td_1).append(drawParam(P[i][k], F[i]["id"], true));
                tr.append($("<td/>", {id: "export-param-"+P[i][k]["id"]}));
                tr.append($("<td/>", {id: "export-val-"+P[i][k]["id"]}));
            }
        }

        $("#" + "event-" + E[0]["id"]).tabs();

        console.log("F_ids: ", F_ids);

        return F_ids;
    }

    function getListHistoryEvents(historyDiv, F_ids) {
        console.log("F_ids: ", F_ids);

        utils.postRequest(
            { "form_ids": F_ids, },
            function(response) {
                console.log("getListHistoryEvents: ", response);

                if (response["data"]) {
                    var data = response["data"];
                    for (var i = 0; i < data.length; ++i) {
                        $("#"+historyDiv+" select").append($("<option></option>")
                            .attr("value", data[i]["id"])
                            .text(data[i]["name"])
                        );
                    }
                    $("#"+historyDiv).show();
                }
            },
            "/handler/getlisthistoryevents"
        );
    }

    return {
        drawParam: drawParam,
        getFormData: getFormData,
        pushFormData: pushFormData,
        getListHistoryEvents: getListHistoryEvents,
    };

});