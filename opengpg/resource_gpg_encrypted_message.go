package opengpg

import (
	"crypto/sha256"
	"fmt"

	"github.com/coopnorge/terraform-provider-opengpg/encryption"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
					StateFunc: func(val any) string {
						publicKey, ok := val.(string)
						if !ok {
							return "MALFORMED KEY"
						}
						recipient, err := encryption.GetRecipient(publicKey)
						if err != nil {
							// We only keep KeyId in state, as we want to keep it small and also
							// we always read public keys anyway. If public key is malformed,
							// creation of resource will fail anyway, so it's fine to set it here.
							return "MALFORMED KEY"
						}

						// Instead of full ASCII-armored key, write only KeyId to state.
						return recipient.GetKeyID()
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

func getRecipients(data *schema.ResourceData) ([]*encryption.Recipient, error) {
	// Iterate over public keys, decode, parse, collect their IDs and add to recipients list.
	publicKeysAny, ok := data.Get("public_keys").([]any)
	if !ok {
		return nil, fmt.Errorf("expected type %T on key %q, got %T", []any{}, "public_keys", data.Get("public_keys"))
	}

	publicKeys := make([]string, 0, len(publicKeysAny))
	for i, v := range publicKeysAny {
		pk, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected type string on public key (idx %d), got %T", i, v)
		}
		publicKeys = append(publicKeys, pk)
	}

	return encryption.GetRecipients(publicKeys)
}

func savePublicKeys(data *schema.ResourceData, recipients []*encryption.Recipient) error {
	// Store ID of each public key, to store them in state (StateFunc does not work for TypeList for some reason).
	pksIDs := []string{}

	for _, recipient := range recipients {
		pksIDs = append(pksIDs, recipient.GetKeyID())
	}

	if err := data.Set("public_keys", pksIDs); err != nil {
		return fmt.Errorf("setting %q property: %w", "public_keys", err)
	}

	return nil
}

func resourceGPGEncryptedMessageCreate(data *schema.ResourceData, _ any) error {
	recipients, err := getRecipients(data)
	if err != nil {
		return fmt.Errorf("getting recipients: %w", err)
	}

	if err := savePublicKeys(data, recipients); err != nil {
		return fmt.Errorf("saving public keys: %w", err)
	}

	plaintextMessage, ok := data.Get("content").(string)
	if !ok {
		return fmt.Errorf("data in property %q was not a string", "content")
	}

	encryptedMessage, err := encryption.EncryptAndEncodeMessage(recipients, plaintextMessage)
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

func resourceGPGEncryptedMessageRead(_ *schema.ResourceData, _ any) error {
	return nil
}

func resourceGPGEncryptedMessageDelete(d *schema.ResourceData, _ any) error {
	d.SetId("")

	return nil
}

func sha256sum(data any) string {
	bytes, ok := data.(string)
	if !ok {
		// There is no way to handle this gracefully with existing SDK.
		panic(fmt.Sprintf("Expected state data to be of type %T, got %T", "", data))
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(bytes)))
}
