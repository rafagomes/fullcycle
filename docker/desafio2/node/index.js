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

const sql = 'INSERT INTO people(name) values("Rafa")';
connection.query(sql);
connection.end();

app.get("/", (req, res) => {
  res.send("<h1>Full Cycle Rocks!</h1>");
});

app.listen(3000, () => console.log("Server is up and running"));
