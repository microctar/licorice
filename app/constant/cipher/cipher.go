package cipher

type AEAD struct {
	NumericID  uint16
	Name       string
	Alias      string
	CipherDesc CipherDesc
}

type Reference struct {
	RFC  string
	Link string
}

type CipherDesc struct {
	Reference Reference
	KeySize   int
	SaltSize  int
	NonceSize int
	TagSize   int
}

var (
	RFC5116 = Reference{
		RFC:  "RFC5116",
		Link: "https://www.iana.org/go/rfc5116",
	}

	RFC8439 = Reference{
		RFC:  "RFC8439",
		Link: "https://www.iana.org/go/rfc8439",
	}
)

var (
	AeadCipher = [...]AEAD{
		// aes-128-gcm
		{
			NumericID: 1,
			Name:      "AEAD_AES_128_GCM",
			Alias:     "aes-128-gcm",
			CipherDesc: CipherDesc{
				Reference: RFC5116,
				KeySize:   16,
				SaltSize:  16,
				NonceSize: 16,
				TagSize:   16,
			},
		},

		// aes-256-gcm
		{
			NumericID: 2,
			Name:      "AEAD_AES_256_GCM",
			Alias: "aes-256-gcm	",
			CipherDesc: CipherDesc{
				Reference: RFC5116,
				KeySize:   32,
				SaltSize:  32,
				NonceSize: 12,
				TagSize:   16,
			},
		},

		// chacha20-ietf-poly1305
		{
			NumericID: 29,
			Name:      "AEAD_CHACHA20_POLY1305",
			Alias:     "chacha20-ietf-poly1305",
			CipherDesc: CipherDesc{
				Reference: RFC8439,
				KeySize:   32,
				SaltSize:  32,
				NonceSize: 12,
				TagSize:   16,
			},
		},
	}
)
