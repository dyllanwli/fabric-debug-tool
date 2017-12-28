var mysql = require("mysql")

var connection = mysql.createConnection({
    host: 'localhost',
    user: 'root',
    password: '1234',
    database: 'test'
})

err = connection.connect();

// check connection
var check = function(err) {
    if (err) throw err;
    console.log("Connected!");
}


// clean table 
var clean = function(err) {
    if (err) {
        console.log("CLEAN ERROR: ", err.message)
        return
    }
    var sql = "DROP TABLE users";
    connection.query(sql, function (err, result) {
        if (err) {
            console.log('[DROP ERROR] - ', err.message)
            console.log('[WARNNING] - Check if the table exiting.')
            return
        };
        console.log("Table deleted");
    });
}

// create table
var create = function(err) {
    var sql = "CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), organization VARCHAR(255),token VARCHAR(255))";
    connection.query(sql, function (err, result) {
        if (err) {
            console.log('[DROP ERROR] - ', err.message)
            return
        };
        console.log("Table created");
    });
}

// sql
var addsql = 'INSERT INTO users (Id, name,organization,token) VALUES(0,?,?,?)'
var addparameter = ['diya', 'org1', '1234123']
connection.query(addsql, addparameter, function (err, result) {
    if (err) {
        console.log('[INSERT ERROR] - ', err.message);
        return;
    }

    console.log('--------------------------INSERT----------------------------');
    //console.log('INSERT ID:',result.insertId);        
    console.log('INSERT ID:', result);
    console.log('-----------------------------------------------------------------\n\n');
});

var query = 'SELECT * FROM users'
connection.query(query, function (err, result) {
    if (err) {
        console.log('[SELECT ERROR] - ', err.message);
        return
    }
    console.log('---------------------------SELECT * --------------------')
    console.log(result)
    console.log('----------------------------------------------')
})

check(err);
clean(err);
create(err);

connection.end();
