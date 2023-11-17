package score

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/cantara/bragi"
	"github.com/cantara/bragi/sbragi"
	"github.com/cantara/cantara-annual-christmasbeer/account/types"

	"github.com/cantara/gober/stream/event"
	"github.com/cantara/gober/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

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
	GetById(id uuid.UUID) (user types.Account, err error)
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

type calculated struct {
	Beer beerStore.Beer
	sum  float32
	Avg  int
	Num  int
}

func InitResource(router *gin.RouterGroup, path string, as accountService, bs beerService, s service, ctx context.Context) (r resource, err error) {
	r = resource{
		path:     path,
		router:   router,
		aService: as,
		bService: bs,
		service:  s,
	}

	scoreStream, err := s.ScoreStream(ctx)
	if err != nil {
		log.AddError(err).Error("while starting global score stream")
	}

	signalChan := make(chan struct{}, 0)
	var cache string
	go func(scoreStream <-chan event.Event[store.Score], scores []store.Score) {
		p := sync.Pool{
			New: func() any {
				return []byte{}
			},
		}
		for score := range scoreStream {
			scores = append(scores, score.Data)

			ant := 5
			start := len(scores) - ant
			if start < 0 {
				start = 0
			}

			calcs := make(map[string]calculated)
			for _, score := range scores {
				c, ok := calcs[score.Beer.ToId()]
				if !ok {
					c = calculated{
						Beer: score.Beer,
					}
				}
				c.Num++
				c.sum += score.Rating
				calcs[score.Beer.ToId()] = c
			}

			if len(calcs) < ant {
				ant = len(calcs)
			}
			high := make([]calculated, ant)
			most := make([]calculated, ant)
			for _, v := range calcs {
				v.Avg = int(v.sum / float32(v.Num))
				insert(high, func(v1, v2 calculated) bool { return v1.Avg < v2.Avg }, v)
				insert(most, func(v1, v2 calculated) bool { return v1.Num < v2.Num }, v)
			}
			slices.Reverse[[]calculated](high)
			slices.Reverse[[]calculated](most)

			b := p.Get().([]byte)
			buf := bytes.NewBuffer(b)
			sumary(scores[start:], high, most).Render(ctx, buf)
			cache = buf.String()
			p.Put(b)
			close(signalChan)
			signalChan = make(chan struct{}, 0)
		}
	}(scoreStream, make([]store.Score, 0, 1024))

	websocket.Serve[string](r.router, r.path+"/sumary", func(c *gin.Context) bool {
		return true
	}, func(_ <-chan string, out chan<- websocket.Write[string], p gin.Params, ctx context.Context) {
		defer close(out)
		out <- websocket.Write[string]{Data: cache}
		for {
			select {
			case <-signalChan: //This logic is flawd and will spam!!
				errChan := make(chan error, 1)
				select {
				case out <- websocket.Write[string]{
					Data: cache,
					Err:  errChan,
				}:
					select {
					case err := <-errChan:
						sbragi.WithError(err).Trace("sent sumary")
					case <-ctx.Done():
						return
					}
					time.Sleep(time.Second) //Adding this to reduce spam frequency
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
	websocket.Serve[store.Score](r.router, r.path, func(c *gin.Context) bool {
		return true
	}, func(inn <-chan store.Score, out chan<- websocket.Write[store.Score], p gin.Params, ctx context.Context) {
		stream, err := s.ScoreStream(ctx)
		if err != nil {
			log.AddError(err).Error("while starting session score stream")
		}
		defer close(out)
		for {
			select {
			case e := <-stream:
				errChan := make(chan error, 1)
				select {
				case out <- websocket.Write[store.Score]{
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
	r.router.POST(r.path+"/:scoreYear/:beerId", r.registerHandler())
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
		user, err := res.aService.GetById(userid)
		if err != nil {
			log.AddError(err).Error("while getting user info for score")
			errorResponse(c, "while getting user info", http.StatusInternalServerError)
			return
		}
		s := store.Score{
			Year:       scoreYear,
			ScorerId:   userid,
			Scorer:     user.FirstName,
			Beer:       beer,
			RatingBase: score.Rating,
			Comment:    score.Comment,
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

func insert(arr []calculated, f func(v1, v2 calculated) bool, v calculated) {
	for i := len(arr) - 1; i >= 0; i-- {
		if !f(arr[i], v) {
			continue
		}
		v, arr[i] = arr[i], v
	}
}
