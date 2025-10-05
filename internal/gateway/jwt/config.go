package jwt

type SignerConfig struct {
	ActiveKeyPairName string
	KeyPairs          []RSAKeyBytesPair
}

type VerifierConfig struct {
	KeyPairs []RSAKeyBytesPair
}

type RSAKeyBytesPair struct {
	Name    string
	Private []byte
	Public  []byte
}
