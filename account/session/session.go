package session

import (
	"context"
	"fmt"
	"github.com/cantara/cantara-annual-christmasbeer/crypto"
	"github.com/cantara/gober/eventmap"
	"github.com/cantara/gober/store/inmemory"
	"github.com/gofrs/uuid"
	"time"

	log "github.com/cantara/bragi"
)

const (
	ExpiresInSeconds = 7200
)

type sessionService struct {
	sessions eventmap.EventMap[AccessTokenSession, sessionEventMetadata]
}

var cryptKey = "2Q82oDggY6CwBs6QHFu3brYjt8JqFILnn68FDN/eTcU="

type sessionEventMetadata struct {
	Expires time.Time `json:"expires"`
}

func Init(ctx context.Context) (s sessionService, err error) {
	store, err := inmemory.Init()
	if err != nil {
		return
	}
	edSession, err := eventmap.Init[AccessTokenSession, sessionEventMetadata](store, "session", "1.0.0", "authentication_session", func(key string) string {
		return cryptKey
	}, ctx)
	if err != nil {
		return
	}
	s = sessionService{
		sessions: edSession,
	}

	go func() {
		tickChan := time.Tick(15 * time.Second)
		for {
			for range tickChan {
				now := time.Now()
				for _, key := range s.sessions.Keys() {
					session, _, err := s.sessions.Get(key)
					if err != nil {
						log.Println(err)
						continue
					}
					if session.AccessToken.Expires.After(now) {
						continue
					}
					s.sessions.Delete(key)
				}
			}
		}
	}()
	return
}

func create() (accessToken AccessToken, err error) {
	token, err := crypto.GenToken()
	if err != nil {
		return
	}
	accessToken = AccessToken{
		Token:     token,
		ExpiresIn: ExpiresInSeconds,
		Expires:   time.Now().Add(ExpiresInSeconds * time.Second),
	}
	return
}

func (s sessionService) Create(accountId uuid.UUID) (accessToken AccessToken, err error) {
	accessToken, err = create()
	if err != nil {
		return
	}
	ats := AccessTokenSession{
		AccountId:   accountId,
		AccessToken: accessToken,
	}
	sem := sessionEventMetadata{
		Expires: accessToken.Expires,
	}
	err = s.sessions.Set(accessToken.Token, ats, sem)
	return
}

func (s sessionService) Renew(token string) (accessToken AccessToken, err error) {
	/*
		This should just extend the time of the token instead of actually creating a new token.
	*/
	session, _, err := s.sessions.Get(token)
	if err != nil {
		if err == eventmap.ERROR_KEY_NOT_FOUND {
			err = fmt.Errorf("sessions does not exist")
		}
		return
	}
	log.Println("SESSION: ", session)
	if session.AccessToken.Expires.Before(time.Now()) {
		err = fmt.Errorf("sessions has expired")
		return
	}
	accessToken, err = create()
	if err != nil {
		return
	}
	defer s.sessions.Delete(token)
	accessToken.Token = token
	ats := AccessTokenSession{
		AccountId:   session.AccountId,
		AccessToken: accessToken,
	}
	sem := sessionEventMetadata{
		Expires: accessToken.Expires,
	}
	err = s.sessions.Set(accessToken.Token, ats, sem)
	return
}

func (s sessionService) Validate(token string) (accessToken AccessToken, accountId uuid.UUID, err error) {
	session, _, err := s.sessions.Get(token)
	if err != nil {
		if err == eventmap.ERROR_KEY_NOT_FOUND {
			err = fmt.Errorf("sessions does not exist")
		}
		return
	}
	accessToken = session.AccessToken
	accountId = session.AccountId
	return
}
