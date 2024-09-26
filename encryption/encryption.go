package encryption

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	protonpgp "github.com/ProtonMail/gopenpgp/v3/crypto"
)

// Recipient is our own representation of a recipient/key.
// It does not export any underlying crypto-library's type, so that we are free to change it in the future, without making breaking changes.
type Recipient struct {
	protonKey *protonpgp.Key
}

// GetKeyID returns the key ID, hex encoded as a string.
func (r *Recipient) GetKeyID() string {
	return r.protonKey.GetHexKeyID()
}

// IsExpired returns whether the key is expired at the given point in time.
func (r *Recipient) IsExpired(t time.Time) bool {
	return r.protonKey.IsExpired(t.UTC().Unix())
}

// GetUserEmail returns the email of the primary identity, if found.
// The identity must be valid at the given point in time.
func (r *Recipient) GetUserEmail(t time.Time) (string, bool) {
	_, id := r.protonKey.GetEntity().PrimaryIdentity(t, nil)
	if id == nil || id.UserId == nil || id.UserId.Email == "" {
		return "", false
	}
	return id.UserId.Email, true
}

// GetRecipients decodes and parses a list of armor-encoded public keys.
func GetRecipients(publicKeys []string) ([]*Recipient, error) {
	// Store recipients for encryption.
	recipients := make([]*Recipient, 0, len(publicKeys))

	// Iterate over all the public keys, and decode and parse them, and collect them in a slice.
	for i, pk := range publicKeys {
		recipient, err := GetRecipient(pk)
		if err != nil {
			return nil, fmt.Errorf("decoding public key #%d: %w", i, err)
		}

		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

// GetRecipient decodes and parses an armor-encoded public key.
func GetRecipient(publicKey string) (*Recipient, error) {
	key, err := protonpgp.NewKeyFromArmored(publicKey)
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %w", err)
	}
	return &Recipient{protonKey: key}, nil
}

// EncryptAndEncodeMessage encrypts the message to all of the recipients.
// The message is encoded in the Armor-encoding.
func EncryptAndEncodeMessage(recipients []*Recipient, message string) (string, error) {
	if len(recipients) == 0 {
		return "", fmt.Errorf("no recipients")
	}

	keyring := &protonpgp.KeyRing{}
	for i, v := range recipients {
		err := keyring.AddKey(v.protonKey)
		if err != nil {
			return "", fmt.Errorf("adding key to keyring (index %d): %w", i, err)
		}
	}

	encrypter, err := protonpgp.PGP().Encryption().Recipients(keyring).New()
	if err != nil {
		return "", fmt.Errorf("creating encrypter: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	wcEncrypt, err := encrypter.EncryptingWriter(buf, protonpgp.Armor)
	if err != nil {
		return "", fmt.Errorf("encrypting message: %w", err)
	}

	if _, err := io.Copy(wcEncrypt, strings.NewReader(message)); err != nil {
		return "", fmt.Errorf("writing content to buffer: %w", err)
	}

	if err := wcEncrypt.Close(); err != nil {
		return "", fmt.Errorf("closing encrypted message: %w", err)
	}

	return buf.String(), nil
}
