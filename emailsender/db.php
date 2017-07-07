<?php

$serverName = "localhost";
$userName = "root";
$password = "";
$db_name = "email";
$db_port = "3306";

try {
    $conn = new PDO("mysql:host=$serverName;port=$db_port;dbname=$db_name", $userName, $password);
    // set the PDO error mode to exception
    $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
} catch (PDOException $e) {
    echo "Connection failed: " . $e->getMessage();
}