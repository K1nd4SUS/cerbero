const interval = setInterval(function() {
    check();
}, 5000); 
 
clearInterval(interval);

function check(){
    readTextFile("./assets/info.json", function(text){
        var data = JSON.parse(text);
        if(data != localStorage.getItem('config')){
            memConfig(JSON.stringify(data));
            refreshList(data);
            window.location.reload();
        }
    });
}

function readTextFile(file, callback) {
    var rawFile = new XMLHttpRequest();
    rawFile.overrideMimeType("application/json");
    rawFile.open("GET", file, true);
    rawFile.onreadystatechange = function() {
        if (rawFile.readyState === 4 && rawFile.status == "200") {
            callback(rawFile.responseText);
        }
    }
    rawFile.send(null);
}

function change(txt) {

    var x=document.getElementById("olTest");
    newLI = document.createElement("li");

    newText = document.createTextNode(txt);
    newLI.appendChild(newText);
    x.appendChild(newLI);

}


function retrieve_rules(){
    var r1 = document.getElementById("s1").value;
    var r2 = document.getElementById("s2").value;
    var r3 = document.getElementById("s3").value;
    var r4 = document.getElementById("s4").value;
    send_rules("http://127.0.0.1:9090", { 's1': r1, 's2': r2, 's3': r3, 's4': r4}, 'post')
}

function send_rules(servURL, params, method) {
    method = method || "post";
    var form = document.createElement("form");
    form.setAttribute("method", method);
    form.setAttribute("action", servURL);
    for(var key in params) {
        var hiddenField = document.createElement("input");
        hiddenField.setAttribute("type", "hidden");
        hiddenField.setAttribute("name", key);
        hiddenField.setAttribute("value", params[key]);
        form.appendChild(hiddenField);
    }
    document.body.appendChild(form);
    form.submit();
}

document.getElementById("buttonTest").addEventListener('click', function(){
    retrieve_rules();
});