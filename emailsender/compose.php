<?php include "header.php" ?>
    <a href="index.php">Back to Inbox</a>
    <div class="clearfix"></div>
    <div class="col-xs-5 col-xs-offset-3">
        <div class="form-area">
            <form role="form" method="post" action="#">
                <h3 style="margin-bottom: 25px;" class="text-center">Write email</h3>

                <div class="form-group">
                    <input type="text" class="form-control" id="email" name="receiver" placeholder="Email" required>
                </div>
                <div class="form-group">
                    <input type="text" class="form-control" id="subject" name="subject" placeholder="Subject" required>
                </div>
                <div class="form-group">
                    <textarea class="form-control" name="body" placeholder="Message" rows="7"></textarea>
                </div>

                <button type="submit" id="submit" name="submit" class="btn btn-primary pull-right">Send</button>
            </form>
        </div>
    </div>
<?php include "footer.php" ?>

<?php

require_once "db.php";

if (isset($_POST['submit'])) {
    if ($_POST['receiver'] !== "" && $_POST['subject'] !== "" && $_POST['body'] !== "") {
        $receiver = $_POST['receiver'];
        $subject = $_POST['subject'];
        $body = $_POST['body'];
        $dateTime = date('Y-m-d H:i:s');

        mail($receiver, $subject, $body);

        $stmt = $conn->prepare("INSERT INTO email (receiver,subject,body,date) VALUES (?, ?, ?, ?)");
        $stmt->execute([$receiver, $subject, $body, $dateTime]);
        header('location:index.php');
    }
}