var mysql = require("mysql")

var connection = mysql.createConnection({
    host: 'localhost',
    user: 'root',
    password: '1234',
    database: 'test'
})

// check connection
connection.connect(function (err) {
    if (err) throw err;
    console.log("Connected!");
});

// clean table
connection.connect(function (err) {
    if (err) throw err;
    var sql = "DROP TABLE users";
    con.query(sql, function (err, result) {
        if (err) {
            console.log('[DROP ERROR] - ', err.message)
            console.log('[WARNNING] - Check if the table exiting.')
            return
        };
        console.log("Table deleted");
    });
})

// create table
connection.connect(function (err) {
    if (err) throw err;
    console.log("Connected!");
    var sql = "CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), organization VARCHAR(255),token VARCHAR(255))";
    con.query(sql, function (err, result) {
        if (err) {
            console.log('[DROP ERROR] - ', err.message)
            return
        };
        console.log("Table created");
    });
});

// sql
var addsql = 'INSERT INTO users(Id, name,organization,token) VALUES(0,?,?,?)'
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

var query = 'SELECT * FROM user'
connection.query(query, function (err, result) {
    if (err) {
        console.log('[SELECT ERROR] - ', err.message);
        return
    }
    console.log('---------------------------SELECT * --------------------')
    console.log(result)
    console.log('----------------------------------------------')
})

connection.end();