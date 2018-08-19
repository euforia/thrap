package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"

	homedir "github.com/mitchellh/go-homedir"
)

func VerifySignature(pubkey, data, signature []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(rawPubKey, data, &r, &s)
}

func GenerateECDSAKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// GenerateECDSAKeyPair generates a keypair writing the private key at filename
// and a file with .pub appended with the public key
func GenerateECDSAKeyPairFile(filename string, curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	// c := elliptic.P256()
	kp, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	priv, pub, err := encodeECDSA(kp, filename)
	if err == nil {
		err = writePem(priv, pub, filename)
	}
	return kp, err
}

func encodeECDSA(privateKey *ecdsa.PrivateKey, filename string) ([]byte, []byte, error) {

	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// pemEncoded := pem.Encode(privH, &pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncoded, pemEncodedPub, nil
}

func writePem(priv, pub []byte, filename string) error {
	err := ioutil.WriteFile(filename, priv, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename+".pub", pub, 0600)
}

func LoadECDSAKeyPair(filename string) (*ecdsa.PrivateKey, error) {
	if filename[0] == '~' {
		cfile, err := homedir.Expand(filename)
		if err != nil {
			return nil, err
		}
		filename = cfile
	}

	priv, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(priv)
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}

	pub, err := ioutil.ReadFile(filename + ".pub")
	if err != nil {
		return nil, err
	}
	blockPub, _ := pem.Decode(pub)
	// blockPub, _ := pem.Decode(pemEncodedPub)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	privateKey.PublicKey = *publicKey

	return privateKey, nil
}
