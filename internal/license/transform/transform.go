package transform

import (
	"github.com/wenzhenxi/gorsa"
)

type Transform func([]byte) ([]byte, error)

type RSATransform interface {
	Encrypt() Transform
	Decrypt() Transform
}

type KeyPair struct {
	PriKey string
	PubKey string
}

type rsa struct {
	key KeyPair
}

func (r *rsa) Encrypt() Transform {
	return func(bytes []byte) ([]byte, error) {
		err := gorsa.RSA.SetPrivateKey(r.key.PriKey)
		if err != nil {
			return nil, err
		}
		data, err := gorsa.RSA.PriKeyENCTYPT(bytes)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (r *rsa) Decrypt() Transform {
	return func(bytes []byte) ([]byte, error) {
		err := gorsa.RSA.SetPublicKey(r.key.PubKey)
		if err != nil {
			return nil, err
		}
		data, err := gorsa.RSA.PubKeyDECRYPT(bytes)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func RSA() RSATransform {
	r := &rsa{
		key: KeyPair{
			PriKey: `-----BEGIN PRIVATE KEY-----
PROC_TYPE: Java

MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKXtd5XAlN3nWjmT
KQR45Lifq/ETYuhJ9YwOkUtjvEd4HMbuA0C6vqLc1b0OI/3Tnw4jweFgZ2CJ/HUS
alrTU2YGKLvtBzVRSVc2sITyZ2teDUcFMBKJl4MFzk/b9tOV2UUj3sVFOPQLfJW2
jcyMMPl/cag7CE9j7RtnpSdsZmpfAgMBAAECgYEAjcgcJxooGnVV41yb7/ZdemT1
x0mJenO4HbVU8daHS4qXDGTU4rqvqvqIqMMsffgWMT7crHhz3UoLLv5NYs6ws1zV
ahoKmqVayTpICkjT+wOgfDnJp3+h6ys/wPJESxsjlPEt/J8sXZqD/od+bq9LlggS
qWo0pbZjeazPpXIaSAECQQDPS66fMH+vInBBgO/4lFr39aju9X/4f7RvwFDPicLT
xT2OtHUwm0cdYGuH3gEItVRr72qvY8Oi5kdMkKLaYt8vAkEAzOmZ0GCP9wvFzQZO
+WS/P6PCicSDmmmIBChrEkDmMQq9d5oohfI+ToFHrfga4k1k6ccPoEckknc/spbr
1w/b0QI/L8ZBeG60/qfxNyeAJsoKLRtw06HA3ISSES9BcJNPU38hsMHmQE2JFjwi
jC2eD2O7ESUccU+Mxv5LcFnlLm+bAkAtp9a3kOxCtQLLXZ52/rWF7mzH2VshKmY9
1uuUU5V2U9hHL7fbsE+lmjRoVKFYzrmvRMT8hx1k7ODqX6oIbuYhAkB6Quitqyh3
DF+caa4UYANAnHI7AOeokRLBKBale0JIVOT8qQra69I/GMlL1r9S3XZih2Hng5Ii
MpmPApBFQ3tx
-----END PRIVATE KEY-----
`,
			PubKey: `-----BEGIN PUBLIC KEY-----
PROC_TYPE: Java

MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCl7XeVwJTd51o5kykEeOS4n6vx
E2LoSfWMDpFLY7xHeBzG7gNAur6i3NW9DiP9058OI8HhYGdgifx1Empa01NmBii7
7Qc1UUlXNrCE8mdrXg1HBTASiZeDBc5P2/bTldlFI97FRTj0C3yVto3MjDD5f3Go
OwhPY+0bZ6UnbGZqXwIDAQAB
-----END PUBLIC KEY-----
`,
		},
	}
	return r
}
