var urlBox;
var msgDiv;

$(function() {
	urlBox = $('#url');
	msgDiv = $('#msg');
	urlBox.focus();
});

function shawty(data) {
	if (data.success == 0) {
		msgDiv.text(data.message)
	}
	else {
		msgDiv.text(' ')
		urlBox.val(data.short);
		urlBox.focus();
		urlBox.select();
	}			
}

function shorten() {
	var urlStr = urlBox.val();
	$.getScript('shawty.js?url=' + escape(urlStr));

	return false;
}