<?php
require_once "db.php";

if (!isset($_GET['id'])) {
    header('location:index.php');
}
$id = $_GET['id'];
$stmt = $conn->prepare("SELECT * FROM email WHERE id=?");
$stmt->execute([$id]);
$email = $stmt->fetch();

?>

<?php include "header.php" ?>
<a href="index.php">Back to Inbox</a>
<div class="clearfix"></div>
<hr>
<div class="row">
    <div class="col-xs-12">
        <div class=""><strong>Receiver: </strong><?php echo $email['receiver'] ?></div>
        <div class=""><strong>Subject: </strong><?php echo $email['subject'] ?></div>
        <div class=""><strong>Date: </strong><?php echo $email['date'] ?></div>
        <div class=""><strong>Message: </strong><?php echo $email['body'] ?></div>
    </div>
</div>
<?php include "footer.php" ?>


