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

$(document).ready(function () {
    $('#cancel_reply').hide();
    let parent_all = $('#wy-delegate-all');
    parent_all.on('click', '#comment_reply_link', show_reply_form);
    parent_all.on('click', '#cancel_reply', cancel_reply_form);
});
