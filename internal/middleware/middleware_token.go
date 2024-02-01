package middleware

//
//import (
//	"github.com/Jackalgit/BuildShortURL/internal/handlers"
//	"github.com/Jackalgit/BuildShortURL/internal/util"
//	"log"
//	"net/http"
//)
//
//func (s *handlers.ShortURL) TokenMiddleware(next http.Handler) http.Handler {
//	tokenFn := func(w http.ResponseWriter, r *http.Request) {
//
//		cookie, err := r.Cookie("token")
//		if err == http.ErrNoCookie {
//			SetCookie(w)
//			next.ServeHTTP(w, r)
//		}
//
//		cookieStr := cookie.Value
//		userId := util.GetUserID(cookieStr)
//
//		if _, ok := a.cache[userId]; !ok {
//			SetCookie(w)
//			next.ServeHTTP(w, r)
//		}
//
//		next.ServeHTTP(w, r)
//
//	}
//	return http.HandlerFunc(tokenFn)
//}
//
//func SetCookie(w http.ResponseWriter) {
//	tokenString, err := util.BuildJWTString()
//	if err != nil {
//		log.Printf("[BuildJWTString] %q", err)
//	}
//	cookie := http.Cookie{Name: "token", Value: tokenString}
//	http.SetCookie(w, &cookie)
//}
