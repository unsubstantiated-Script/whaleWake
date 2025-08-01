package token

import (
	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"log"
	"time"
	"whaleWake/util"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey // Use the correct type for the key
	implicit     []byte
}

// NewPasetoMaker creates a new PasetoMaker instance with the provided symmetric key.
// It returns an error if the key is not valid.
func NewPasetoMaker() (Maker, error) {

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Unable to load config:", err)
	}

	key := config.PASETO_SYMMETRIC_KEY

	if key == "" {
		return nil, ErrMissingPasetoEnvVariable
	}

	symmetricKey, err := paseto.V4SymmetricKeyFromHex(key)

	if err != nil {
		return nil, ErrFailedHexToSymmetricKeyConversion
	}

	return &PasetoMaker{
		symmetricKey, []byte{},
	}, nil
}

// CreateToken create a new token for an specific user and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// Crete paseto token
	token := paseto.NewToken()

	// create uuid for token ID
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	token.Set("id", tokenID.String())
	token.Set("username", username)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	return token.V4Encrypt(maker.symmetricKey, maker.implicit), nil
}

// VerifyToken verifies a token and returns the user ID or an error.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)

	if err != nil {
		return nil, ErrExpiredToken
	}

	payload, err := getPayloadFromToken(parsedToken)

	if err != nil {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func getPayloadFromToken(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}

	username, err := t.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}

	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}

// RefreshToken refreshes a token and returns the new token or an error.
func (maker *PasetoMaker) RefreshToken(token string) (string, error) {
	payload, err := maker.VerifyToken(token)
	if err != nil {
		return "", err
	}

	// Create a new token with the same user ID and a new expiration time
	newToken, err := maker.CreateToken(payload.Username, time.Minute*15) // Example: 24 hours duration
	if err != nil {
		return "", err
	}

	return newToken, nil
}
