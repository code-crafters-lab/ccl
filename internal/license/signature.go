package license

type Signature interface {
	Raw
	//Sign(payload Payload, header Header) []byte
}

type signature struct {
	header  Raw
	payload Raw
	// 签名算法
	signatureAlgorithm string
	// 原始签名值
	rawSignatureData []byte
}

func (s *signature) RawBytes() ([]byte, error) {
	return nil, nil
}
