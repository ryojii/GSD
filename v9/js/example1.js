// specify the columns
var columnDefs = [
    {headerName: "Id", field: "idexec"},
    {headerName: "Name", field: "name"},
    {headerName: "Status", field: "status"},
    {headerName: "Trace", field: "trace"},
];


function ajaxGet(target){
    var Httpreq = new XMLHttpRequest(); // a new request
    Httpreq.open("GET",target, false);
    Httpreq.send(null);
    return Httpreq.responseText;          
}

// specify the data
var rowData = JSON.parse(ajaxGet("/execs"));
console.log("this is the author name: "+rowData.name);

// let the grid know which columns and what data to use
var gridOptions = {
    columnDefs: columnDefs,
    rowData: rowData
};

// wait for the document to be loaded, otherwise ag-Grid will not find the div in the document.
document.addEventListener("DOMContentLoaded", function() {

    // lookup the container we want the Grid to use
    var eGridDiv = document.querySelector('#myGrid');

    // create the grid passing in the div to use together with the columns & data we want to use
    new agGrid.Grid(eGridDiv, gridOptions);
});
