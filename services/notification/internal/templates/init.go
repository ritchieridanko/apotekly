package templates

import "embed"

//go:embed emails/*.html.tmpl emails/partials/*.html.tmpl
var Email embed.FS
