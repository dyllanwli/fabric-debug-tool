var express = require('express');
var router = express.Router();
var mysql = require('mysql');
var conf = require('../config.js');

//定义pool池
var pool = mysql.createPool(
    {
        host        : conf.dbMysql.host,
        user        : conf.dbMysql.user,
        password    : conf.dbMysql.password,
        database    : conf.dbMysql.database,
        // port        : conf.dbMysql.port
    }
);

router.get('/', function(req, res) {
    var selectSites = "select *, date_format(do_time, '%Y-%m-%d %H:%i:%s') as time from siteinfo order by id";
    pool.getConnection(function(err, connection) {
        if (err) throw err;
        connection.query(selectSites, function(err, rows) {
            if (err) throw  err;
            res.render('sites', {title : '站点分布', results : rows});
            //回收pool
            connection.release();
        });
    });
});

module.exports = router;