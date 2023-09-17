package util

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifyMessage(message string, signedMessage string) (string, error) {
	// hash the unsigned message using EIP191
	hashedMessage := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message))
	hash := crypto.Keccak256Hash(hashedMessage)

	// get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)

	// handles cases where EIP-115 is not implemented
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	// recover a public key from the signed message
	sigPublickKeyESCSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublickKeyESCSA == nil {
		err = errors.New("could not get a public key from the message signature")
	}
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(*sigPublickKeyESCSA).String(), nil
}
