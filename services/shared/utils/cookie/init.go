package cookie

import "github.com/gin-gonic/gin"

type Cookie struct {
	domain   string
	isSecure bool
}

func Init(env, domain string) *Cookie {
	isProd := env == "prod"
	if !isProd {
		domain = ""
	}
	return &Cookie{
		domain:   domain,
		isSecure: isProd,
	}
}

func (c *Cookie) Set(ctx *gin.Context, name, value, path string, duration int) {
	ctx.SetCookie(name, value, duration, path, c.domain, c.isSecure, true)
}

func (c *Cookie) Unset(ctx *gin.Context, name, path string) {
	ctx.SetCookie(name, "", -1, path, c.domain, c.isSecure, true)
}
