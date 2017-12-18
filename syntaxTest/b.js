function bb(){
    console.log("this is b")
}

function ss(obj){
    obj.name = "adfad"
    obj.j = {
        asd:"sd"
    }
}

par = new Object()
ss(par)

console.log(par)