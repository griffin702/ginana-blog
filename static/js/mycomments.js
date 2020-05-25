//评论回复表单按钮
function show_reply_form(event) {
    event.preventDefault();
    let $this = $(this);
    let comment_id = $this.data('comment-id');
    let parent_id = $this.data('parent-id');
    $('#id_reply_pk').val(comment_id);
    $('#id_reply_fk').val(parent_id);
    $('#form-comment').appendTo($this.closest('.media-body'));
    $('#cancel_reply').show();
}

//评论关闭按钮
function cancel_reply_form(event) {
    event.preventDefault();
    $('#comment_content').val('');
    $('#id_reply_pk').val('0');
    $('#id_reply_fk').val('0');
    $('#form-comment').appendTo($('#wrap-form-comment'));
    $('#cancel_reply').hide();
}

function comment_submit(event) {
    let $this = $(this);
    let islogin = $this.data('islogin');
    if (!islogin) {
        return
    }
    event.preventDefault();
    let url = "/admin/comment/add";
    let object_pk = $("#id_object_pk").val();
    let object_pk_type = $("#id_object_pk_type").val();
    if (!object_pk_type) {
        object_pk_type = 0
    }
    let reply_pk = $('#id_reply_pk').val();
    let reply_fk = $('#id_reply_fk').val();
    let comment_content = $("#comment_content").val();
    if (comment_content !== '') {
        $("#comment_content").val('').focus();
        let timestamp = (new Date).getTime() + parseInt(10 * Math.random(), 10);
        let security_hash = hex_md5(reply_pk + timestamp + "@YO!r52w!D2*I%Ov");
        // alert((new Date).getTime());
        $.ajax({
            type: 'POST',
            data: JSON.stringify({
                'obj_pk': parseInt(object_pk, 10),
                'obj_pk_type': parseInt(object_pk_type, 10),
                'reply_pk': parseInt(reply_pk, 10),
                'reply_fk': parseInt(reply_fk, 10),
                'content': comment_content,
            }),
            url: url,
            success: function (data) {
                if (!data.data) {
                    alert("评论失败,原因:" + data.message)
                } else {
                    setTimeout(function () {
                        $.ajax({
                            type: 'GET',
                            data: {},
                            url: location.pathname,
                            // cache: true,
                            dataType: "html",
                            success: function (data) {
                                $('#form-comment').appendTo($('#wrap-form-comment'));
                                $('#cancel_reply').hide();
                                $("#wrap-comments-list").html($(data).find("#comments-list"));
                                initcommentslist();
                                $(".comments_length").html($(data).find(".comments_length p"));
                                $("#id_reply_pk").val('0');
                                $("#id_reply_fk").val('0');
                            }
                        });
                    }, 300);
                }
            }
        });
    }
}

$(document).ready(function () {
    $('#cancel_reply').hide();
    let parent_all = $('#wy-delegate-all');
    parent_all.on('click', '#comment_reply_link', show_reply_form);
    parent_all.on('click', '#cancel_reply', cancel_reply_form);
    parent_all.on('click', '#comment_submit', comment_submit);
});
