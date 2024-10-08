package opengpg_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const rsaConfig = `
resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = [
    var.opengpg_public_key_rsa,
  ]
}

variable "opengpg_public_key_rsa" {
  description = "A public-key of type RSA 4096, using the SHA256 hashing algorithm"
  default = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=lZG6
-----END PGP PUBLIC KEY BLOCK-----
EOF
}
`

const ecc25519Config = `
resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = [
    var.opengpg_public_key_ecc25519,
  ]
}

variable "opengpg_public_key_ecc25519" {
  description = "A public-key of type ECC 25519, using the SHA512 hashing algorithm"
  default = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEZumPqxYJKwYBBAHaRw8BAQdAbEfcyIa1K25/DMwIocm+MfYYAF3jlq8+GxjY
7FjzZ9S0LGZvb2Jhci1lY2MyNTUxOSAoZm9vYmFyKSA8Zm9vQGJhci1jdXJ2ZS5j
b20+iJMEExYKADsWIQT3olI2/t6HX2MIvmYnB22SxES8hwUCZumPqwIbAwULCQgH
AgIiAgYVCgkICwIEFgIDAQIeBwIXgAAKCRAnB22SxES8hyrnAQCtqpxMtfX6XEbd
W5Ao9sfBDs3q3ajL+UOCrV/iQG3dQQEA5jbFcyju/LSL4Dkb4JF8zKiWa19hzdGW
rAlC9eYHcAm4OARm6Y+rEgorBgEEAZdVAQUBAQdAA77h3XlxlSlYygtVs/mwPXyb
szkpBnI3TlJQqUeLaTYDAQgHiHgEGBYKACAWIQT3olI2/t6HX2MIvmYnB22SxES8
hwUCZumPqwIbDAAKCRAnB22SxES8h/ErAQDlnDX+BRfsGyPR+WzhnTCV+fUvaWsG
wCnk1/Lh1fpGhAEAhFokVxfOaontUAnC/dDsxSZ7KdLVgOOuwZskhidIagk=
=1j0l
-----END PGP PUBLIC KEY BLOCK-----
EOF
}
`

const badPublicKey = `
resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = [
		"not valid message",
  ]
}
`

const badPublicKeyPEMEncoded = `
resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
	public_keys = [
		<<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

bm9wZQo=
-----END PGP PUBLIC KEY BLOCK-----
EOF
		,
	]
}

`

const noPublicKeys = `
resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = []
}
`

func TestGPGEncryptedMessageRSA(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: rsaConfig,
			},
			{
				Config:             rsaConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:      badPublicKey,
				ExpectError: regexp.MustCompile(`decoding public key #0: decoding public key: gopenpgp: error in reading key ring: openpgp: invalid data: tag byte does not have MSB set`),
			},
			{
				Config:      badPublicKeyPEMEncoded,
				ExpectError: regexp.MustCompile(`decoding public key #0: decoding public key: gopenpgp: error in reading key ring: openpgp: invalid data: tag byte does not have MSB set`),
			},
		},
	})
}

func TestGPGEncryptedMessageECC25519(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ecc25519Config,
			},
			{
				Config:             ecc25519Config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestGPGEncryptedMessageBadArguments(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      noPublicKeys,
				ExpectError: regexp.MustCompile(regexSpaceOrNewline(`Attribute public_keys requires 1 item minimum, but config has only 0 declared.`)),
				Destroy:     false,
			},
		},
	})
}

// regexSpaceOrNewline allows a string's spaces to be either a normal space or a newline. This is to allow less brittle tests when terraform changes their output
func regexSpaceOrNewline(str string) string {
	return strings.ReplaceAll(str, " ", "[\\ \\n]")
}
