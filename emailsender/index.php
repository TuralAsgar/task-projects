<?php require_once "db.php"; ?>

<?php include "header.php" ?>
<hr>
<div class="row">
    <div class="col-sm-3 col-md-2">
        <ul class="nav nav-pills nav-stacked">
            <li><a href="#"> Inbox </a>
            </li>
            <li><a href="#">Sent Mail</a></li>
        </ul>
    </div>
    <div class="col-sm-9 col-md-10">
        <form action="delete.php" method="post" id="delete">
            <a class="btn btn-info" href="compose.php">Compose</a>
            <button type="submit" class="btn btn-danger" id="deleteButton">Delete</button>
            <table class="table" width="100%">
                <thead>
                <tr>
                    <th><input type="checkbox" id="select-all"></th>
                    <th><a id="receiver">Receiver</a></th>
                    <th><a>Subject</a></th>
                    <th><a id="date">Date</a></th>
                </tr>
                </thead>
                <tbody id="data"></tbody>
            </table>
        </form>

    </div>
</div>
<?php include "footer.php" ?>
