{{if .}}
{{if .Bookmarklet}}
function shawty(data) {
	if (data.success == 1)
		prompt('Here\'s your short URL:', data.short);
	else
		alert('Unable to shorten this URL: ' + data.message);
}
{{end}}

shawty({
	'success':   {{.Success}},
	'message':  '{{.Message}}',
	'short':     '{{.Short}}',
	'long':      '{{.Long}}',
	'hits':      {{.Hits}}, 
	'timestamp': {{.Timestamp.Unix}}
});
{{end}}