function getAll(){
    $("#search-results").empty();
    $.ajax({
        url: "/execs",
        method: "GET",
        data: $("#search-form").serialize(),
        success: insertResults
    });
    showList();
    return false;
}

function submitSearch(){
    $("#search-results").empty();
    $.ajax({
        url: "/search/"+$("#searchMethod").find(":selected").val()+"/"+$("#search").val(),
        method: "GET",
        data: $("#search-form").serialize(),
        success: insertResults
    });
    return false;
}

function deleteExec(id) {
    $.ajax({
        url: "/delete?id="+id,
        method: "GET"
    });
    return false;
}

function showDetail(id) {
    $.ajax( {
        url: "/exec/"+id,
        method: "GET",
        datatype : "json",
        success: function(data) {
            $("#campaign").text(data.idcampaign);
            $("#name").text(data.name);
            $("#trace").text(data.trace);
            $("#reviewer").val(data.reviewe);
            $("#status").text(data.status);
            $("#forced-status").val(data.forcedStatus);
            $("#start").text(data.start);
            $("#end").text(data.end);
        }
    });
    $("#reviewer").on( "change", function() {
        updateReviewer(id, this.value);
    });
    $("#forced-status").on( "change", function() {
        updateStatus(id, this.value);
    });
    $("#executions-list").hide();
    $("#execution-detail").show();

    return false;
}

function showList() {
    $("#executions-list").show();
    $("#forced-status").off();
    $("#reviewer").off();
    $("#execution-detail").hide();
}

function updateStatus(id, value) {
    $.ajax({
        url: "/update/"+id+"/status/"+value,
        method: "PUT",
        success: function(id, value) {
            console.log(id +"-"+ value);
        }
    });
    $("#forced-status").off();
    $("#forced-status").on("change", function() {
        updateStatus(id, this.value);
    });
}

function updateReviewer(id, value) {
    $.ajax({
        url: "/update/"+id+"/reviewer/"+value,
        method: "PUT"
        // en cas de succes je dois mettre à jour le tableau
    });
    $("#reviewer").off();
    $("#reviewer").on("change", function() {
        updateReviewer(id, this.value);
    });
}

function insertResults(data) {
    var searchResults = $("#search-results");
    searchResults.empty();
    data.forEach(function(result) {
        var row = $("<tr id='"+result.idexec+"'><td>"+ result.idexec+"</td><td>"+ result.idcampaign +"</td><td onclick='showDetail("+result.idexec+")'>"+result.name+"</td><td>"+result.status+"</td><td>"+result.fstatus+"</td><td>"+result.reviewer+"</td><td>"+result.trace+"</td><td>"+result.start+"</td><td>"+result.end+"</td><td><button type='button' value='delete' onclick='deleteExec("+result.idexec+")'>delete </button></td></tr>");
        searchResults.append(row);
    });
}
