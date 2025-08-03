package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// 客户端内置公钥（用于加密请求和验证签名）
const embeddedPublicKey = `-----BEGIN PUBLIC KEY-----
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

// GetEmbeddedPublicKey 获取内置公钥
func GetEmbeddedPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(embeddedPublicKey))
	if block == nil {
		log.Fatal("failed to decode public key")
	}
	
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("failed to parse public key:", err)
	}
	
	return pub.(*rsa.PublicKey)
}