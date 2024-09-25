# GPG Command Provider

The GPG provider allows to generate GPG encrypted message in ASCII-armored format
using Terraform.

This provider uses community-maintained fork [ProtonMail/go-crypto](https://github.com/ProtonMail/go-crypto)
to perform GPG encryption. Currently the only supported option is encrypting
message with public keys.

Managing GPG keyring or signing files is currently not implemented.

## Example Usage

```hcl
terraform {
  required_providers {
    opengpg = {
      source  = "coopnorge/opengpg"
      version = "0.2.0"
    }
  }
}

resource "opengpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = [
    var.opengpg_public_key,
  ]
}

output "opengpg_encrypted_message" {
  value = opengpg_encrypted_message.example.result
}

variable "opengpg_public_key" {
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
```

## Argument Reference

This provider currently takes no arguments.
