package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"httpRest/internal/validator"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	PlainText string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, timeDuration time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(timeDuration),
		Scope:  scope,
	}

	// make an array on bytes of size 16 , fill with random bytes from OS
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// encode the random bytes into human-readable string
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// make the hashed version to store in the db
	hash := sha256.Sum256([]byte(token.PlainText))
	// tuen the returned array into slice for data consistency and store it in the token struct
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlainText(v *validator.Validator, tokenplainText string) {
	v.Check(len(tokenplainText) > 0 && len(tokenplainText) <= 26, "token", "token is empty or used different hashing base")
}

type TokenModel struct {
	db *sql.DB
}

func (tokenModel TokenModel) New(userID int64, Exp time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, Exp, scope)
	if err != nil {
		return nil, err
	}
	err = tokenModel.Insert(token)
	return token, err
}

func (tokenmodel *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`
	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tokenmodel.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TokenModel) DeleteToken(scope string, userID int64) error {
	query := `DELETE FROM tokens
				WHERE scope = $1 AND user_id = $2`

	args := []interface{}{scope, userID}
	_, err := m.db.Exec(query, args...)
	return err
}
