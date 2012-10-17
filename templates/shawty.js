{{if .}}
{{if .Data.Bookmarklet}}
function shawty(data) {
	if (data.success == 1)
		prompt('Here\'s your short URL:', data.short);
	else
		alert('Unable to shorten this URL: ' + data.message);
}
{{end}}

shawty({
	'success':   {{.Data.Success}},
	'message':   '{{if .Errors}}{{index .Errors 0}}{{end}}',
	'short':     'http://{{.Data.Domain}}/{{.Data.Short}}',
	'long':      '{{.Data.Shawty.Url}}',
	'hits':      {{if .Data.Shawty}}{{.Data.Shawty.Hits}}{{else}}0{{end}}, 
	'timestamp': {{if .Data.Shawty}}{{.Data.Shawty.CreatedOn.Unix}}{{else}}0{{end}}
});
{{end}}