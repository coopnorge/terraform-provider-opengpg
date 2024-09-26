package encryption

import (
	"testing"
	"time"

	protonpgp "github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptMessageSuccessful(t *testing.T) {
	testCases := []struct {
		name       string
		publicKeys []string
	}{
		{name: "rsa", publicKeys: []string{publicKeyRSA}},
		{name: "curve", publicKeys: []string{publicKeyCurve}},
		{name: "rsa+curve", publicKeys: []string{publicKeyRSA, publicKeyCurve}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recipients, err := GetRecipients(tc.publicKeys)
			require.NoError(t, err)
			message := "hello world"
			result, err := EncryptAndEncodeMessage(recipients, message)
			require.NoError(t, err)
			assert.True(t, protonpgp.IsPGPMessage(result), "encrypted messages is not PGP message")
		})
	}
}

func TestEncryptMessageExpired(t *testing.T) {
	testCases := []struct {
		name       string
		publicKeys []string
	}{
		{name: "expired rsa", publicKeys: []string{publicKeyRSAExpired}},
		{name: "expired curve", publicKeys: []string{publicKeyCurveExpired}},
		{name: "expired rsa + expired curve", publicKeys: []string{publicKeyRSAExpired, publicKeyCurveExpired}},
		{name: "valid rsa + expired curve", publicKeys: []string{publicKeyRSA, publicKeyCurveExpired}},
		{name: "expired rsa + valid curve", publicKeys: []string{publicKeyRSAExpired, publicKeyCurve}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recipients, err := GetRecipients(tc.publicKeys)
			require.NoError(t, err)
			message := "hello world"
			result, err := EncryptAndEncodeMessage(recipients, message)
			require.ErrorContains(t, err, "cannot encrypt a message")
			require.ErrorContains(t, err, "no valid encryption keys")
			assert.Empty(t, result)
		})
	}
}

func TestGetKeyID(t *testing.T) {
	testCases := []struct {
		name          string
		publicKey     string
		expectedKeyID string
	}{
		{name: "rsa", publicKey: publicKeyRSA, expectedKeyID: "4f54663daabdbaff"},
		{name: "curve", publicKey: publicKeyCurve, expectedKeyID: "27076d92c444bc87"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recipient, err := GetRecipient(tc.publicKey)
			require.NoError(t, err)

			keyID := recipient.GetKeyID()
			assert.Equal(t, tc.expectedKeyID, keyID)
		})
	}
}

func TestIsExpired(t *testing.T) {
	testCases := []struct {
		name            string
		publicKey       string
		time            time.Time
		expectedExpired bool
	}{
		{name: "rsa (non-expiring)", publicKey: publicKeyRSA, time: time.Now(), expectedExpired: false},
		{name: "curve (non-expiring)", publicKey: publicKeyCurve, time: time.Now(), expectedExpired: false},
		{name: "rsa (before expiry-date)", publicKey: publicKeyRSAExpired, time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), expectedExpired: false},
		{name: "rsa (after expiry-date)", publicKey: publicKeyRSAExpired, time: time.Date(2024, 9, 25, 16, 0, 0, 0, time.UTC), expectedExpired: true},
		{name: "curve (before expiry-date)", publicKey: publicKeyCurveExpired, time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), expectedExpired: false},
		{name: "curve (after expiry-date)", publicKey: publicKeyCurveExpired, time: time.Date(2024, 9, 25, 16, 0, 0, 0, time.UTC), expectedExpired: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recipient, err := GetRecipient(tc.publicKey)
			require.NoError(t, err)

			isExpired := recipient.IsExpired(tc.time)
			assert.Equal(t, tc.expectedExpired, isExpired)
		})
	}
}

func TestGetUserEmail(t *testing.T) {
	testCases := []struct {
		name          string
		publicKey     string
		time          time.Time
		expectedEmail string
	}{
		{name: "rsa (non-expiring)", publicKey: publicKeyRSA, time: time.Now(), expectedEmail: "bar@foo.com"},
		{name: "curve (non-expiring)", publicKey: publicKeyCurve, time: time.Now(), expectedEmail: "foo@bar-curve.com"},
		{name: "rsa (before expiry-date)", publicKey: publicKeyRSAExpired, time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), expectedEmail: "foo@coop.no"},
		{name: "rsa (after expiry-date)", publicKey: publicKeyRSAExpired, time: time.Date(2024, 9, 25, 16, 0, 0, 0, time.UTC), expectedEmail: "foo@coop.no"},
		{name: "curve (before expiry-date)", publicKey: publicKeyCurveExpired, time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), expectedEmail: "foo@coop.no"},
		{name: "curve (after expiry-date)", publicKey: publicKeyCurveExpired, time: time.Date(2024, 9, 25, 16, 0, 0, 0, time.UTC), expectedEmail: "foo@coop.no"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recipient, err := GetRecipient(tc.publicKey)
			require.NoError(t, err)

			email, ok := recipient.GetUserEmail(tc.time)
			assert.True(t, ok)
			assert.Equal(t, tc.expectedEmail, email)
		})
	}
}

// The following RSA key does not have an expiry-date
var publicKeyRSA = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=lZG6
-----END PGP PUBLIC KEY BLOCK-----`

// The following Curve key does not have an expiry-date
var publicKeyCurve = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
-----END PGP PUBLIC KEY BLOCK-----`

// The following RSA key is set to expire at 2024-09-25, around 12 o'clock (yes, that is in the past, even when writing that)
var publicKeyRSAExpired = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
2wtG2A==
=ezHw
-----END PGP PUBLIC KEY BLOCK-----`

// The following Curve key is set to expire at 2024-09-25, around 12 o'clock (yes, that is in the past, even when writing that)
var publicKeyCurveExpired = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEZvKFHhYJKwYBBAHaRw8BAQdAfuy7802mdT5uZgS7FZ+Y+sjlAxNk19eFHO9j
X2MEPva0H2V4cGlyZWQgY3VydmUga2V5IDxmb29AY29vcC5ubz6ImQQTFgoAQRYh
BNip/li+plk5FFLOPZ7bP9GBou6fBQJm8oUeAhsDBQkAAVGABQsJCAcCAiICBhUK
CQgLAgQWAgMBAh4HAheAAAoJEJ7bP9GBou6fYmYBAJ9dNjhsQTymabwBLA0Db4Nx
ekfwu0pEipM5kgAcs5y7AQCa3c6pwEjzY52E+GsfgfATdawCupO8TPOri++K2HGN
Brg4BGbyhR4SCisGAQQBl1UBBQEBB0DeYzjTp0tVao6VfAHp28L/IS1gxupXA7uV
utG/6KCGAgMBCAeIfgQYFgoAJhYhBNip/li+plk5FFLOPZ7bP9GBou6fBQJm8oUe
AhsMBQkAAVGAAAoJEJ7bP9GBou6fR4wA/j3zghlADAmyJaNXNDFPqA7D31N/VGiP
1OcotwikB6g1AP9kOhgTm6AWgHj2RU2p1VH/LuTLgr69DH5w/7//MqEKDw==
=S5pq
-----END PGP PUBLIC KEY BLOCK-----`
