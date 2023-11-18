package account

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cantara/cantara-annual-christmasbeer/account/types"
	"github.com/cantara/cantara-annual-christmasbeer/crypto"

	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

type service struct {
	store   StoreService
	admin   AdminService
	session SessionService
}

func InitService(store StoreService, admin AdminService, session SessionService, ctx context.Context) (s service, err error) {
	s = service{
		store:   store,
		admin:   admin,
		session: session,
	}
	return
}

func (s service) Accounts() (accounts []types.Account, err error) {
	return s.store.Accounts()
}

func (s service) GetById(id uuid.UUID) (acc types.Account, err error) {
	return s.store.GetById(id)
}

func (s service) GetByUsername(username string) (acc types.Account, err error) {
	return s.store.GetByUsername(username)
}

func (s service) GetLogin(username string) (login InternalLoggin, err error) {
	data, err := s.store.GetLogin(username)
	if err != nil {
		return
	}
	err = json.Unmarshal(data.Data, &login)
	return
}

func (s service) Register(acc AccountRegister) (token session.AccessToken, err error) {
	a := types.Account{
		Id:        uuid.Must(uuid.NewV7()),
		Email:     acc.Email,
		FirstName: acc.FirstName,
		LastName:  acc.LastName,
		Number:    acc.Number,
		VerificationCodes: types.Codes{
			Email:  crypto.GenRandBase32String(8),
			Mobile: crypto.GenRandBase32String(8),
			MFA:    crypto.GenRandBase32String(16),
		},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	err = s.store.Register(a)
	if err != nil {
		return
	}
	login := InternalLoggin{
		AccountId: a.Id,
		Username:  acc.Username,
		Salt:      make([]byte, PW_SALT_BYTES),
	}
	_, err = io.ReadFull(rand.Reader, login.Salt)
	if err != nil {
		return
	}
	login.Password, err = hashPassword([]byte(acc.Password), login.Salt)
	if err != nil {
		return
	}
	data, err := json.Marshal(&login)
	if err != nil {
		return
	}
	err = s.store.Link(types.Login{
		Id:        acc.Username,
		AccountId: login.AccountId,
		Type:      types.INTERNAL,
		Data:      data,
	})
	if err != nil {
		return
	}
	return s.session.Create(a.Id)
}

func (s service) Login(username, password string) (token session.AccessToken, err error) {
	login, err := s.GetLogin(username)
	if err != nil {
		return
	}
	hashPass, err := hashPassword([]byte(password), login.Salt)
	if err != nil {
		return
	}
	if len(login.Password) != len(hashPass) {
		err = fmt.Errorf("wrong password")
		return
	}
	for i := range login.Password {
		if login.Password[i] != hashPass[i] {
			err = fmt.Errorf("wrong password")
			return
		}
	}
	return s.session.Create(login.AccountId)
}

func (s service) Renew(token string) (tokenOut session.AccessToken, err error) {
	return s.session.Renew(token)
}

func (s service) Validate(token string) (tokenOut session.AccessToken, accountId uuid.UUID, err error) {
	return s.session.Validate(token)
}

func (s service) RegisterAdmin(accountId uuid.UUID, r Rights) (err error) {
	return s.admin.Register(accountId, r)
}

func (s service) IsAdmin(accountId uuid.UUID) bool {
	p, err := s.admin.Rights(accountId)
	if err != nil {
		return false
	}
	return p.Rights.Admin
}

func (s service) IsNewbie(accountId uuid.UUID) bool {
	p, err := s.admin.Rights(accountId)
	if err != nil {
		return false
	}
	return p.Rights.Weight < 1
}

func hashPassword(password, salt []byte) ([]byte, error) {
	if len(salt) < PW_SALT_BYTES {
		return nil, fmt.Errorf("too week salt")
	}
	//hash, err := scrypt.Key(password, salt, 1<<14, 8, 1, PW_HASH_BYTES)
	hash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return hash, nil
}
