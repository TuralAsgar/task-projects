<?php
require_once "db.php";

if (empty($_POST['emailId'])) {
    header('location:index.php');
}

$emails = implode(",", $_POST['emailId']);

$stmt = $conn->prepare("DELETE FROM email WHERE id IN ({$emails})");
if ($stmt->execute()) {
    header('location:index.php');
} else {
    echo "Error occured. Try again";
}