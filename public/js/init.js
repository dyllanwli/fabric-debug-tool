window.onload = function(){
    var p = document.createElement("p");
    p.innerHTML = "this is created dynamically"
    document.body.appendChild(p);

    var btn = document.getElementById("terminal");
    
    btn.onclick = function(){
        window.open('http://localhost:8088/terminal');
    };
}

