package gpg

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	protonpgp "github.com/ProtonMail/gopenpgp/v3/crypto"
)

func resourceGPGEncryptedMessage() *schema.Resource {
	return &schema.Resource{
		// TODO: Migrate to <Create/Read/Delete/Update>Context
		Create: resourceGPGEncryptedMessageCreate,
		// Those 2 functions below does nothing, but must be implemented.
		Read:   resourceGPGEncryptedMessageRead,
		Delete: resourceGPGEncryptedMessageDelete,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
				StateFunc: sha256sum,
			},
			"public_keys": {
				Type:     schema.TypeList,
				MinItems: 1,
				ForceNew: true,
				Required: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					ForceNew: true,
					StateFunc: func(val interface{}) string {
						recipient, err := entityFromString(val.(string))
						if err != nil {
							// We only keep KeyId in state, as we want to keep it small and also
							// we always read public keys anyway. If public key is malformed,
							// creation of resource will fail anyway, so it's fine to set it here.
							return "MALFORMED KEY"
						}

						// Instead of full ASCII-armored key, write only KeyId to state.
						return recipient.GetHexKeyID()
					},
				},
			},
			"result": {
				Type:      schema.TypeString,
				Computed:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func getRecipients(data *schema.ResourceData) ([]*protonpgp.Key, error) {
	// Store recipients for encryption.
	recipients := []*protonpgp.Key{}

	// Iterate over public keys, decode, parse, collect their IDs and add to recipients list.
	publicKeys, ok := data.Get("public_keys").([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected type %T on key %q, got %T", "public_keys", []interface{}{}, data.Get("public_keys"))
	}

	for i, pk := range publicKeys {
		recipient, err := entityFromString(pk.(string))
		if err != nil {
			return nil, fmt.Errorf("decoding public key #%d: %w", i, err)
		}

		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

func savePublicKeys(data *schema.ResourceData, recipients []*protonpgp.Key) error {
	// Store ID of each public key, to store them in state (StateFunc does not work for TypeList for some reason).
	pksIDs := []string{}

	for _, recipient := range recipients {
		pksIDs = append(pksIDs, recipient.GetHexKeyID())
	}

	if err := data.Set("public_keys", pksIDs); err != nil {
		return fmt.Errorf("setting %q property: %w", "public_keys", err)
	}

	return nil
}

func encryptAndEncodeMessage(recipients []*protonpgp.Key, message string) (string, error) {
	if len(recipients) == 0 {
		return "", fmt.Errorf("no recipients")
	}
	pgp := protonpgp.PGP()
	keyring, err := protonpgp.NewKeyRing(nil)
	if err != nil {
		return "", fmt.Errorf("creating keyring: %w", err)
	}

	for i, v := range recipients {
		err := keyring.AddKey(v)
		if err != nil {
			return "", fmt.Errorf("adding key to keyring (index %d): %w", i, err)
		}
	}

	encrypter, err := pgp.Encryption().Recipients(keyring).New()
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

func resourceGPGEncryptedMessageCreate(data *schema.ResourceData, _ interface{}) error {
	recipients, err := getRecipients(data)
	if err != nil {
		return fmt.Errorf("getting recipients: %w", err)
	}

	if err := savePublicKeys(data, recipients); err != nil {
		return fmt.Errorf("saving public keys: %w", err)
	}

	encryptedMessage, err := encryptAndEncodeMessage(recipients, data.Get("content").(string))
	if err != nil {
		return fmt.Errorf("encrypting message: %w", err)
	}

	if err := data.Set("result", encryptedMessage); err != nil {
		return fmt.Errorf("setting %q property: %w", "result", err)
	}

	// Calculate SHA-256 checksum of message for ID.
	data.SetId(sha256sum(encryptedMessage))

	return nil
}

func resourceGPGEncryptedMessageRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceGPGEncryptedMessageDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")

	return nil
}

func entityFromString(str string) (*protonpgp.Key, error) {
	key, err := protonpgp.NewKeyFromArmored(str)
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %w", err)
	}
	return key, nil
}

func sha256sum(data interface{}) string {
	bytes, ok := data.(string)
	if !ok {
		// There is no way to handle this gracefully with existing SDK.
		panic(fmt.Sprintf("Expected state data to be of type %T, got %T", "", data))
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(bytes)))
}
