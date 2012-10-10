var urlBox;
var msgDiv;

$(function() {
	urlBox = $('#url');
	msgDiv = $('#msg');
	urlBox.focus();
});

function err(str) {
	msgDiv.text(str);
}

function shawty(data) {
	if (data.success == 0) {
		msgDiv.text(data.message)
	}
	else {
		err(' ');
		urlBox.val(data.short);
	}

	urlBox.focus();
	urlBox.select();
}

function shorten() {
	var urlStr = urlBox.val();
	urlStr = $.trim(urlStr);
	urlBox.val(urlStr);
	if (urlStr != '')
		$.getScript('shawty.js?url=' + escape(urlStr));
	else 
		err('a valid url has to start with http:// or https://');

	return false;
}