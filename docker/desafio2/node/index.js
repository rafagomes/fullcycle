const express = require("express");
const app = express();

const config = {
  host: "db",
  user: "root",
  password: "root",
  database: "desafio2",
};

const mysql = require("mysql2");
const connection = mysql.createConnection(config);

const sqlCreateTable = `CREATE TABLE IF NOT EXISTS people (id int auto_increment, name varchar(255), primary key(id))`;
connection.query(sqlCreateTable);

const sql = 'INSERT INTO people(name) values("Rafa")';
connection.query(sql);
connection.end();

app.get("/", (req, res) => {
  res.send("<h1>Full Cycle Rocks!</h1>");
});

app.listen(3000, () => console.log("Server is up and running"));
