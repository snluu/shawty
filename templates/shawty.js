{{if .}}
shawty({
	'success':   {{.Success}},
	'message':  '{{.Message}}',
	'short':     '{{.Short}}',
	'long':      '{{.Long}}',
	'hits':      {{.Hits}}, 
	'timestamp': {{.Timestamp.Unix}}
});
{{end}}