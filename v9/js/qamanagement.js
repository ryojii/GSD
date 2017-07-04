function submitSearch(){
    $.ajax({
        url: "/execs",
        method: "GET",
        data: $("#search-form").serialize(),
        success: function(data) {
            var searchResults = $("#search-results");
            searchResults.empty();
            data.forEach(function(result) {
                var row = $("<tr><td>"+ result.idexec+"</td><td>"+ result.idcampaign +"</td><td onclick='showDetail("+result.idexec+")'>"+result.name+"</td><td>"+result.status+"</td><td>"+result.trace+"</td><td>"+result.start+"</td><td>"+result.end+"</td><td><button type='button' value='delete' onclick='deleteExec("+result.idexec+")'>delete </button></td></tr>");
                searchResults.append(row);
            });
        }
    });
return false;
}

function deleteExec(id) {
    $.ajax({
        url: "/delete?id="+id,
        method: "GET"
    });
}

function showDetail(id) {
    $.ajax( {
        url: "/exec/"+id,
        method: "GET",
        data : $("detail-form").serialize(),
        datatype : "json",
        success: function(data) {
            $("#campaign").val(data.idcampaign);
            $("#name").val(data.name);
            $("#trace").val(data.trace);
            $("#reviewer").val(data.reviewe);
            $("#status").val(data.status);
            $("#forced-status").val(data.forcedStatus);
            $("#start").val(data.start);
            $("#end").val(data.end);
        }
    });
    $("#executions-list").hide();
    $("#execution-detail").show();

    return false;
}

function showList() {
    $("#executions-list").show();
    $("#execution-detail").hide();
}

function saveExec() {
    alert("ok, je vais enregistrer ça...");
}
