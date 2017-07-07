<?php

require_once "db.php";

$orderBy = "";

switch ($_GET['q']) {
    case 1 :
        $orderBy = "receiver";
        break;
    case 2 :
        $orderBy = "date";
        break;
    default:
        $orderBy = "id";
}

$array = $conn->query("SELECT * FROM email ORDER BY {$orderBy}")->fetchAll(PDO::FETCH_ASSOC);
echo json_encode($array);
