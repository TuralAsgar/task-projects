/**
 * Created by Tural on 7/7/2017.
 */
$(document).ready(function () {
    loadContent("");

    // Listen for click on toggle checkbox
    $('#select-all').click(function (event) {
        if (this.checked) {
            // Iterate each checkbox
            $(':checkbox').each(function () {
                this.checked = true;
            });
        } else {
            // Iterate each checkbox
            $(':checkbox').each(function () {
                this.checked = false;
            });
        }
    });

    $('#receiver').click(function (e) {
        e.preventDefault();
        loadContent(1);
    });

    $('#date').click(function (e) {
        e.preventDefault();
        loadContent(2);
    });

    function loadContent(q) {
        $("#data").empty();
        $.ajax({
            url: "emails.php?q=" + q,
            success: function (response) {
                var obj = jQuery.parseJSON(response);
                $.each(obj, function (key, value) {
                    $("#data").append('<tr><td><input type="checkbox" name="emailId[]" value="' + value['id'] + '"></td> <td><a href="email.php?id=' + value['id'] + '">' + value['receiver'] + '</a></td> <td>' + value['subject'] + '</td> <td>' + value['date'] + '</td></tr>')
                });
            }
        });
    }

});