package score

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cantara/bragi"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cantara/gober/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	beerStore "github.com/cantara/cantara-annual-christmasbeer/beer/store"
	"github.com/cantara/cantara-annual-christmasbeer/score/store"
)

const (
	CONTENT_TYPE      = "Content-Type"
	CONTENT_TYPE_JSON = "application/json"
	AUTHORIZATION     = "Authorization"
)

type accountService interface {
	Validate(token string) (tokenOut session.AccessToken, accountId uuid.UUID, err error)
}

type beerService interface {
	Get(id string) (b beerStore.Beer, err error)
}

type resource struct {
	path     string
	router   *gin.RouterGroup
	aService accountService
	bService beerService
	service  service
}

type validator[bodyT any] struct {
	service accountService
}

func InitResource(router *gin.RouterGroup, path string, as accountService, bs beerService, s service, ctx context.Context) (r resource, err error) {
	r = resource{
		path:     path,
		router:   router,
		aService: as,
		bService: bs,
		service:  s,
	}

	websocket.Websocket(r.router, r.path, func(ctx context.Context, conn *ws.Conn) bool {
		conn.CloseRead(ctx)
		stream, err := s.BeerStream(ctx)
		if err != nil {
			log.AddError(err).Error("while starting beer stream")
		}
		for e := range stream {
			err = wsjson.Write(ctx, conn, e.Data)
			if err != nil {
				log.AddError(err).Warning("while writing to socket")
				return false
			}
		}
		return false
	})
	r.router.PUT(r.path+"/:scoreYear/:beerId", r.registerHandler())
	return
}

type score struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (res resource) registerHandler() func(c *gin.Context) {
	validate := validator[score]{service: res.aService}
	return validate.reqWAuthWBody(func(c *gin.Context, _ session.AccessToken, userid uuid.UUID, score score) {
		beerId := c.Param("beerId")
		beer, err := res.bService.Get(beerId)
		if err != nil {
			errorResponse(c, "Beer missind", http.StatusBadRequest)
			return
		}

		scoreYear, err := strconv.Atoi(c.Param("scoreYear"))
		if err != nil {
			errorResponse(c, "Score year needs to be four didgets", http.StatusBadRequest)
			return
		}

		if scoreYear < 1980 {
			errorResponse(c, "score year is too old. please contact admin if this is a relevant request", http.StatusBadRequest)
			return
		}

		if beer.BrewYear > scoreYear {
			errorResponse(c,
				"Can not score a beer that wasn't brewed yet",
				http.StatusBadRequest)
			return
		}

		if score.Rating < 1 || score.Rating > 6 {
			errorResponse(c, "rating must be in the range of a dice, 1 - 6", http.StatusBadRequest)
			return
		}

		s := store.Score{
			Year:    scoreYear,
			Scorer:  userid,
			Beer:    beer,
			Rating:  score.Rating,
			Comment: score.Comment,
		}
		_, err = res.service.Get(s.ToId())
		if err == nil {
			errorResponse(c, "score already added", http.StatusConflict)
			return
		}
		err = res.service.Register(s)
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
