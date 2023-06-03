package user

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

// 1st responsibility
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

// 2nd responsibility
func (u *User) Save(ctx context.Context) error {
	conn, err := sql.Open("sqlite3", "file:foo.db?cache=shared&mode=memory")
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	query := "INSERT INTO users VALUES (?, ?, ?)"
	_, err = conn.ExecContext(ctx, query, u.ID, u.Email, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to persist user: %w", err)
	}

	return nil
}

// 3rd responsibility
func (u *User) SetPassword(password string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	var iterations uint32 = 3
	var memory uint32 = 64 * 1024
	var parallelism uint8 = 2
	var keyLength uint32 = 32

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m%d, t=%d, p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash,
	)

	u.PasswordHash = encodedHash

	return nil
}

// 4th responsibility
func (u *User) SendEmail(ctx context.Context, subject string, body string) error {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sesClient := ses.NewFromConfig(conf)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{u.Email},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &types.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(u.Email),
	}

	if _, err = sesClient.SendEmail(ctx, input); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
