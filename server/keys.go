package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// 内置RSA私钥
const privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

// 内置RSA公钥
const publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsW2m+fxeWHTcDl4LBHVI
sTbLJyOG7xJm9lhit9AWaMAz3XIXM4WF9hT6VO3E9nJbcTL5ts56nhxOFg4AToze
FDWg1sXk2pTfcBiTKNLQvAc+5t01a0gFupmXhDs0Z79l5UemwVDAThwJ1yOciN0k
ZaMaXs3i9lnI/ZUSRPPB9TdjCCoy6XLiMM/6CGF2+73vz9mVyg1FwzUxxUJNKp6l
vkLf/Z0ybYQRxsLOx03IfAImnyYVTMS/LwgCBodRGHDyWUChmdPcAk4BcpM8onH9
sSIgdtJR3Nym6bPQSqdf6J4qH0EKb610tk+hWJ13FI8ikKmIOD1KUHNp8GKgN7i3
gqvGCNYEuDrKuxYB3ktH/WHz5Ou/w6xOPbTJB4de5g34MRoiybcnaVOrXmQdx38F
fBqrlw9YR0QMLeE1nWQBu9cZ3TsQVzbWSh7iIqeSt+QJgfjiw1vvHA1P7vypjYj4
3mI95ocxc0o/g1avjGzEvsZJEz00njwqJqdbgv2l7SYWdTS2ohcrk9XzaM/FDv4j
i4kpUwHaaZHOmpmDFNr67RJe3WjddBPVGnzMwKg9INJy4hRRDSYt77NKlgOkMouA
jXJ1PNPefl3r5hrpFiwlolTUqaoju5T84MXHYpFlN3a+xIORoLePBHGIGlxmCTtz
55kTNDhieF1Y9TnCvrVZsj8CAwEAAQ==
-----END PUBLIC KEY-----`

// GetPrivateKey 获取内置私钥
func GetPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		log.Fatal("failed to decode private key")
	}
	
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("failed to parse private key:", err)
	}
	
	return key
}

// GetPublicKey 获取内置公钥
func GetPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		log.Fatal("failed to decode public key")
	}
	
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("failed to parse public key:", err)
	}
	
	return pub.(*rsa.PublicKey)
}