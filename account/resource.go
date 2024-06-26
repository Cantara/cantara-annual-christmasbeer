package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	log "github.com/cantara/bragi"
	"github.com/cantara/bragi/sbragi"

	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/oidc"
	"github.com/gin-gonic/gin"
)

const (
	CONTENT_TYPE      = "Content-Type"
	CONTENT_TYPE_JSON = "application/json"
	AUTHORIZATION     = "Authorization"
)

type resource struct {
	path             string
	router           *gin.RouterGroup
	service          service
	requireFirstName bool
	requireLastName  bool
	requireEmail     bool
	requireNumber    bool
}

type validator[bodyT any] struct {
	service service
}

func InitResource(router *gin.RouterGroup, path string, s service) (r resource, err error) {
	r = resource{
		path:             path,
		router:           router,
		service:          s,
		requireFirstName: true,
	}
	// These endpoints are actions not objects, this goes against REST.
	r.router.GET(r.path, r.userHandler())
	r.router.GET(r.path+"/valid", r.validateHandler())
	r.router.GET(r.path+"/renew", r.renewHandler())
	r.router.GET(r.path+"/logins", r.loginsHandler())
	r.router.GET(r.path+"/accounts", r.accountsHandler())
	r.router.PUT(r.path+"/privilege/:account_id", r.registerAdminHandler())
	r.router.GET(r.path+"/admin", r.adminHandler())
	return
}

func (r *resource) InitResourceInternal() error {
	r.router.PUT(r.path+"/:username", r.registerHandler())
	r.router.POST(r.path+"/login", r.loginHandler())
	return nil
}

func (res resource) registerHandler() func(c *gin.Context) {
	validate := validator[AccountRegister]{service: res.service}
	return validate.reqWBody(func(c *gin.Context, a AccountRegister) {
		username := c.Param("username")
		_, err := res.service.GetByUsername(username)
		if err == nil {
			//TODO: Should create an event here
			errorResponse(c, "Conflict", http.StatusConflict)
			return
		}

		if a.Username != username {
			errorResponse(c, "Bad Request. Missmatch usernames", http.StatusBadRequest)
			return
		}
		var msisdn = regexp.MustCompile(`^\+[1-9][0-9]{9,14}$`)
		var number = regexp.MustCompile(`[0-9]`)
		var lower = regexp.MustCompile(`[a-z]`)
		var upper = regexp.MustCompile(`[A-Z]`)
		if len(a.Password) < 8 {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Password is not atleast 8 characters long."),
				http.StatusBadRequest)
			return
		}
		if !(number.MatchString(a.Password) && lower.MatchString(a.Password) && upper.MatchString(a.Password)) {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Password does not meet requrement of number, lower and upercase characters."),
				http.StatusBadRequest)
			return
		}
		if res.requireFirstName && (len(a.FirstName) <= 2 || number.MatchString(a.FirstName)) {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Firstname needs to be atleast 2 chars long and not contain numbers %s", a.FirstName),
				http.StatusBadRequest)
			return
		}
		if res.requireLastName && (len(a.LastName) <= 2 || number.MatchString(a.LastName)) {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Lastname needs to be atleast 2 chars long and not contain numbers %s", a.LastName),
				http.StatusBadRequest)
			return
		}
		var email = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
		if res.requireEmail && !email.MatchString(a.Email) {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Your email does not match email format"),
				http.StatusBadRequest)
			return
		}
		if res.requireNumber && !msisdn.MatchString(a.Number) {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("Your phone number does not match MSISDN format"),
				http.StatusBadRequest)
			return
		}
		a.Username = strings.ToLower(a.Username)

		token, err := res.service.Register(a)
		if err != nil {
			log.Println(err)
			errorResponse(c, "Error while registring", http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:     "token",
			Value:    token.Token,
			Expires:  token.Expires,
			MaxAge:   token.ExpiresIn,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
			HttpOnly: true,
			Domain:   "localhost:3030",
			Path:     "/",
		}
		http.SetCookie(c.Writer, &cookie)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "localhost:3030")
		c.Writer.WriteHeader(http.StatusNoContent)
		//c.JSON(http.StatusOK, token)
	})
}

func (res resource) userHandler() func(c *gin.Context) {
	validate := validator[any]{service: res.service}
	return validate.reqWAuth(func(c *gin.Context, token session.AccessToken, accountId uuid.UUID) {
		acc, err := res.service.GetById(accountId)
		if err != nil {
			errorResponse(c, "Not found", http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, acc)
	})
}

func (res resource) loginsHandler() func(c *gin.Context) {
	validate := validator[any]{service: res.service}
	return validate.req(func(c *gin.Context) {
		logins := make(map[string]string)
		if os.Getenv("account.internal.enable") == "true" {
			logins["internal"] = fmt.Sprintf("localhost:3000%s/login", res.path)
		}
		if os.Getenv("account.external.enable") == "true" {
			if os.Getenv("twitch.enable") == "true" {
				logins[string(oidc.TWITCH)] = fmt.Sprintf("localhost:3000%s/external/session/%s", res.path, oidc.TWITCH)
			}
			if os.Getenv("youtube.enable") == "true" {
				logins[string(oidc.YOUTUBE)] = fmt.Sprintf("localhost:3000%s/external/session/%s", res.path, oidc.YOUTUBE)
			}
		}
		c.JSON(http.StatusOK, logins)
	})
}

func (res resource) accountsHandler() func(c *gin.Context) {
	validate := validator[any]{service: res.service}
	return validate.reqWAdmin(func(c *gin.Context, _ session.AccessToken, _ uuid.UUID) {
		accounts, err := res.service.Accounts()
		if err != nil {
			errorResponse(c, "unable to get all accounts", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, accounts)
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (res resource) loginHandler() func(c *gin.Context) {
	validate := validator[loginRequest]{service: res.service}
	return validate.reqWBody(func(c *gin.Context, lr loginRequest) {
		lr.Username = strings.ToLower(lr.Username)
		token, err := res.service.Login(lr.Username, lr.Password)
		if err != nil {
			log.Println(err)
			errorResponse(c, "Error while loggin inn: "+err.Error(), http.StatusForbidden)
			return
		}
		//w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
		//json.NewEncoder(w).Encode(&token)
		cookie := http.Cookie{
			Name:     "token",
			Value:    token.Token,
			Expires:  token.Expires,
			MaxAge:   token.ExpiresIn,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
			Domain:   "localhost:3030",
			Path:     "/",
		}
		http.SetCookie(c.Writer, &cookie)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "localhost:3030")
		c.Writer.WriteHeader(http.StatusNoContent)
	})
}

func (res resource) validateHandler() func(c *gin.Context) {
	validate := validator[any]{service: res.service}
	return validate.reqWAuth(func(c *gin.Context, _ session.AccessToken, _ uuid.UUID) {
		errorResponse(c, "valid", http.StatusOK)
	})
}

func (res resource) renewHandler() func(c *gin.Context) {
	validate := validator[any]{service: res.service}
	return validate.reqWAuth(func(c *gin.Context, token session.AccessToken, _ uuid.UUID) {
		tokenNew, err := res.service.Renew(token.Token)
		if err != nil {
			log.Println(err)
			errorResponse(c, "Error while renewing token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, tokenNew)
	})
}

func (res resource) registerAdminHandler() func(c *gin.Context) {
	validate := validator[Rights]{service: res.service}
	return validate.reqWAdminWBody(func(c *gin.Context, _ session.AccessToken, _ uuid.UUID, right Rights) {
		accountId, err := uuid.FromString(c.Param("account_id"))
		if err != nil {
			errorResponse(c, "account id must be a uuid", http.StatusBadRequest)
			return
		}
		account, err := res.service.GetById(accountId)
		if err != nil {
			errorResponse(c, "User does not exist", http.StatusBadRequest)
			return
		}
		if res.service.IsAdmin(account.Id) {
			errorResponse(c, "Account is already admin", http.StatusConflict)
			return
		}
		err = res.service.RegisterAdmin(account.Id, right)
		if err != nil {
			errorResponse(c, "Could not register account as admin", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, "")
	})
}

func (res resource) adminHandler() func(c *gin.Context) {
	validate := validator[AccountRegister]{service: res.service}
	return validate.reqWAdmin(func(c *gin.Context, _ session.AccessToken, _ uuid.UUID) {
		c.JSON(http.StatusOK, "")
	})
}

func (res resource) Account(c *gin.Context) (accountId uuid.UUID) {
	headers := c.Request.Header[AUTHORIZATION]
	var tokenString string
	if len(headers) > 0 {
		if strings.HasPrefix(headers[0], "Bearer ") {
			tokenString = strings.TrimPrefix(headers[0], "Bearer ")
		}
	}
	if tokenString == "" {
		tokenString, _ = c.Cookie("token")
	}
	if tokenString == "" {
		return
	}
	token, accountId, err := res.service.Validate(tokenString)
	if err != nil {
		sbragi.WithError(err).Warning("while validating token")
		return
	}
	if time.Now().Add(time.Hour).After(token.Expires) {
		token, err = res.service.Renew(tokenString)
		if err != nil {
			sbragi.WithError(err).Warning("while renewing validated token")
			return
		}
		c.SetCookie("token", token.Token, token.ExpiresIn, "/", "localhost", true, true)
	}
	return
}

func (v validator[bodyT]) req(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Header[CONTENT_TYPE][0] != CONTENT_TYPE_JSON {
			errorResponse(c, "Content Type is not "+CONTENT_TYPE_JSON, http.StatusUnsupportedMediaType)
			return
		}
		f(c)
	}
}

func (v validator[bodyT]) reqWAuth(f func(c *gin.Context, token session.AccessToken, accountId uuid.UUID)) func(c *gin.Context) {
	return v.req(func(c *gin.Context) {
		headers := c.Request.Header[AUTHORIZATION]
		var tokenString string
		if len(headers) > 0 {
			if !strings.HasPrefix(headers[0], "Bearer ") {
				errorResponse(c, "Bad Request. Missing Bearer in "+AUTHORIZATION+" header", http.StatusUnauthorized)
				return
			}
			tokenString = strings.TrimPrefix(headers[0], "Bearer ")
		} else {
			tokenString, _ = c.Cookie("token")
		}
		token, accountId, err := v.service.Validate(tokenString)
		if err != nil {
			log.Println(err)
			errorResponse(c, "Forbidden", http.StatusForbidden)
			return
		}
		f(c, token, accountId)
	})
}

func (v validator[bodyT]) reqWAdmin(f func(c *gin.Context, token session.AccessToken, accountId uuid.UUID)) func(c *gin.Context) {
	return v.reqWAuth(func(c *gin.Context, token session.AccessToken, accountId uuid.UUID) {
		if !v.service.IsAdmin(accountId) {
			errorResponse(c, "Forbidden", http.StatusForbidden)
			return
		}
		f(c, token, accountId)
	})
}

func (v validator[bodyT]) reqWBody(f func(c *gin.Context, body bodyT)) func(c *gin.Context) {
	return v.req(func(c *gin.Context) {
		body, err := unmarshalBody[bodyT](c.Request.Body)
		if err != nil {
			errorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
		f(c, body)
	})
}

func (v validator[bodyT]) reqWAuthWBody(f func(c *gin.Context, token session.AccessToken, accountId uuid.UUID, body bodyT)) func(c *gin.Context) {
	return v.reqWAuth(func(c *gin.Context, token session.AccessToken, accountId uuid.UUID) {
		body, err := unmarshalBody[bodyT](c.Request.Body)
		if err != nil {
			errorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
		f(c, token, accountId, body)
	})
}

func (v validator[bodyT]) reqWAdminWBody(f func(c *gin.Context, token session.AccessToken, accountId uuid.UUID, body bodyT)) func(c *gin.Context) {
	return v.reqWAdmin(func(c *gin.Context, token session.AccessToken, accountId uuid.UUID) {
		body, err := unmarshalBody[bodyT](c.Request.Body)
		if err != nil {
			errorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
		f(c, token, accountId, body)
	})
}

func unmarshalBody[bodyT any](body io.ReadCloser) (v bodyT, err error) {
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&v)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			err = fmt.Errorf("Bad Request. Wrong Type provided for field %s", unmarshalErr.Field)
		} else {
			err = fmt.Errorf("Bad Request %v", err)
		}
		return
	}
	return
}

func errorResponse(c *gin.Context, message string, httpStatusCode int) {
	resp := make(map[string]string)
	resp["error"] = message
	c.JSON(httpStatusCode, resp)
}

/*
func getToken(c *gin.Context) (header string) {
	return
}
*/
