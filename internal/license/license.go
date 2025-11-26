package license

type Raw interface {
	RawBytes() ([]byte, error)
}

type License interface {
	Raw
	Header
	Security
	Payload
	Signature
}
