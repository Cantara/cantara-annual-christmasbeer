package beer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/cantara/bragi"
	"github.com/cantara/bragi/sbragi"

	"github.com/cantara/gober/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/beer/store"
)

const (
	CONTENT_TYPE      = "Content-Type"
	CONTENT_TYPE_JSON = "application/json"
	AUTHORIZATION     = "Authorization"
)

type accountService interface {
	Validate(token string) (tokenOut session.AccessToken, accountId uuid.UUID, err error)
	IsAdmin(accountId uuid.UUID) bool
}

type resource struct {
	path     string
	router   *gin.RouterGroup
	aService accountService
	service  service
}

type validator[bodyT any] struct {
	service accountService
}

func prov(key string) string {
	return "MdgKIHmlbRszXjLbS7pXnSBdvl+SR1bSejtpFTQXxro="
}

func InitResource(router *gin.RouterGroup, path string, as accountService, s service, ctx context.Context) (r resource, err error) {
	r = resource{
		path:     path,
		router:   router,
		aService: as,
		service:  s,
	}

	//(r *gin.RouterGroup, path string, acceptFunc func(c *gin.Context) bool, wsfunc WSHandler[T])

	websocket.Serve[store.Beer](r.router, r.path, func(c *gin.Context) bool {
		return true
	}, func(inn <-chan store.Beer, out chan<- websocket.Write[store.Beer], p gin.Params, ctx context.Context) {
		stream, err := s.BeerStream(ctx)
		if err != nil {
			log.AddError(err).Error("while starting beer stream")
		}
		defer close(out)
		for {
			select {
			case e := <-stream:
				sbragi.Info("read", "beer", e)
				errChan := make(chan error, 1)
				select {
				case out <- websocket.Write[store.Beer]{
					Data: e.Data,
					Err:  errChan,
				}:
					select {
					case err := <-errChan:
						sbragi.WithError(err).Trace("sent beer event")
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
	r.router.PUT(r.path+"/:beerId", r.registerHandler())
	return
}

func (res resource) registerHandler() func(c *gin.Context) {
	validate := validator[store.Beer]{service: res.aService}
	return validate.reqWAuthWBody(func(c *gin.Context, _ session.AccessToken, _ uuid.UUID, a store.Beer) {
		beerId := c.Param("beerId")
		_, err := res.service.Get(beerId)
		if err == nil {
			//TODO: Should create an event here
			errorResponse(c, "Conflict", http.StatusConflict)
			return
		}

		if 1980 > a.BrewYear || a.BrewYear > time.Now().Year() {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("brew year must be within the past 40ish years"),
				http.StatusBadRequest)
			return
		}
		if 0 > a.ABV || a.ABV > 98 {
			errorResponse(c,
				"Bad Request "+fmt.Sprintf("abv must be within a valid alcohol percentage range"),
				http.StatusBadRequest)
			return
		}

		err = res.service.Register(a)
		if err != nil {
			log.Println(err)
			errorResponse(c, "Error while registering", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, "")
	})
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
		authHeader := getAuthHeader(c)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			errorResponse(c, "Bad Request. Missing Bearer in "+AUTHORIZATION+" header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, accountId, err := v.service.Validate(tokenString)
		if err != nil {
			log.Println(err)
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

func (v validator[bodyT]) reqAdminWBody(f func(c *gin.Context, token session.AccessToken, accountId uuid.UUID, body bodyT)) func(c *gin.Context) {
	return v.reqWAuth(func(c *gin.Context, token session.AccessToken, accountId uuid.UUID) {
		if !v.service.IsAdmin(accountId) {
			errorResponse(c, "User is not a admin", http.StatusForbidden)
			return
		}
		body, err := unmarshalBody[bodyT](c.Request.Body)
		if err != nil {
			errorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
		f(c, token, accountId, body)
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
	//w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	//w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["error"] = message
	//json.NewEncoder(w).Encode(resp)
	c.JSON(httpStatusCode, resp)
}

func getAuthHeader(c *gin.Context) (header string) {
	headers := c.Request.Header[AUTHORIZATION]
	if len(headers) > 0 {
		header = headers[0]
	}
	return
}
