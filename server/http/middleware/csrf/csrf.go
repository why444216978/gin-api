package csrf

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	wraphh "github.com/turtlemonvh/gin-wraphh"
)

func CSRF(authKey []byte, domain string) gin.HandlerFunc {
	f := csrf.Protect(authKey,
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.Domain(domain),
		csrf.ErrorHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusForbidden)
			_, _ = rw.Write([]byte(`{"message": "CSRF token invalid"}`))
		})),
	)
	return wraphh.WrapHH(f)
}
