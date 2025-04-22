package gate

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/glide-im/glide/pkg/hash"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
)

type CredentialCrypto interface {
	EncryptCredentials(c *ClientAuthCredentials) ([]byte, error)

	DecryptCredentials(src []byte) (*ClientAuthCredentials, error)
}

// AesCBCCrypto cbc mode PKCS7 padding
type AesCBCCrypto struct {
	Key []byte
}

func NewAesCBCCrypto(key []byte) *AesCBCCrypto {
	keyLen := len(key)
	count := 0
	switch true {
	case keyLen <= 16:
		count = 16 - keyLen
	case keyLen <= 24:
		count = 24 - keyLen
	case keyLen <= 32:
		count = 32 - keyLen
	default:
		key = key[:32]
	}
	if count != 0 {
		key = append(key, bytes.Repeat([]byte{0}, count)...)
	}
	return &AesCBCCrypto{Key: key}
}

func (a *AesCBCCrypto) EncryptCredentials(c *ClientAuthCredentials) ([]byte, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	// generate random iv
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	encryptBody, err := a.Encrypt(jsonBytes, iv)
	if err != nil {
		return nil, err
	}

	// NOTE: append iv
	var encrypt []byte
	encrypt = append(encrypt, iv...)
	encrypt = append(encrypt, encryptBody...)

	// base64 encoding encrypted json credentials
	b64Bytes := make([]byte, base64.RawStdEncoding.EncodedLen(len(encrypt)))
	base64.RawStdEncoding.Encode(b64Bytes, encrypt)
	return b64Bytes, nil
}

func (a *AesCBCCrypto) DecryptCredentials(src []byte) (*ClientAuthCredentials, error) {

	encrypt := make([]byte, base64.RawStdEncoding.DecodedLen(len(src)))
	_, err := base64.RawStdEncoding.Decode(encrypt, src)
	if err != nil {
		return nil, err
	}
	var iv []byte
	iv = append(iv, encrypt[:aes.BlockSize]...)
	var encryptBody []byte
	encryptBody = append(encryptBody, encrypt[aes.BlockSize:]...)

	jsonBytes, err := a.Decrypt(encryptBody, iv)
	if err != nil {
		return nil, err
	}

	credentials := ClientAuthCredentials{}
	err = json.Unmarshal(jsonBytes, &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (a *AesCBCCrypto) Encrypt(src, iv []byte) ([]byte, error) {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	// padding
	blockSize := block.BlockSize()
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	src = append(src, padtext...)

	encryptData := make([]byte, len(src))

	if len(iv) != block.BlockSize() {
		iv = a.cbcIVPending(iv, blockSize)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptData, src)

	return encryptData, nil
}

func (a *AesCBCCrypto) Decrypt(src, iv []byte) ([]byte, error) {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		iv = a.cbcIVPending(iv, blockSize)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, src)

	length := len(dst)
	if length == 0 {
		return nil, errors.New("unpadding")
	}
	unpadding := int(dst[length-1])
	if length < unpadding {
		return nil, errors.New("unpadding")
	}
	res := dst[:(length - unpadding)]

	return res, nil
}

func (a *AesCBCCrypto) cbcIVPending(iv []byte, blockSize int) []byte {
	k := len(iv)
	if k < blockSize {
		return append(iv, bytes.Repeat([]byte{0}, blockSize-k)...)
	} else if k > blockSize {
		return iv[0:blockSize]
	}
	return iv
}

// Authenticator handle client authentication message
type Authenticator struct {
	credentialCrypto CredentialCrypto
	gateway          DefaultGateway
}

func NewAuthenticator(gateway DefaultGateway, key string) *Authenticator {
	k := sha512.Sum512([]byte(key))
	return &Authenticator{
		credentialCrypto: NewAesCBCCrypto(k[:]),
		gateway:          gateway,
	}
}

func (a *Authenticator) MessageInterceptor(dc DefaultClient, msg *messages.GlideMessage) bool {

	if dc.GetCredentials() == nil {
		return false
	}
	switch msg.Action {
	case messages.ActionGroupMessage, messages.ActionChatMessage, messages.ActionChatMessageResend:
		break
	default:
		return false
	}

	if dc.GetCredentials() == nil || dc.GetCredentials().Secrets == nil {
		_ = a.gateway.EnqueueMessage(dc.GetInfo().ID, messages.NewMessage(msg.GetSeq(), messages.ActionNotifyForbidden, "no credentials"))
		return true
	}

	secret := dc.GetCredentials().Secrets.MessageDeliverSecret
	if secret == "" {
		_ = a.gateway.EnqueueMessage(dc.GetInfo().ID, messages.NewMessage(msg.GetSeq(), messages.ActionNotifyForbidden, "no message deliver secret"))
		return true
	}

	var ticket = msg.Ticket
	// sha1 hash
	if len(ticket) != 40 {
		_ = a.gateway.EnqueueMessage(dc.GetInfo().ID, messages.NewMessage(msg.GetSeq(), messages.ActionNotifyForbidden, "invalid ticket"))
		return true
	}
	sum1 := hash.SHA1(secret + msg.To)
	id := dc.GetInfo().ID
	expectTicket := hash.SHA1(secret + id.UID() + sum1)

	if strings.ToUpper(ticket) != strings.ToUpper(expectTicket) {
		logger.I("invalid ticket, expected=%s, actually=%s, secret=%s, to=%s, from=%s", expectTicket, ticket, secret, msg.To, id.UID())
		// invalid ticket
		_ = a.gateway.EnqueueMessage(dc.GetInfo().ID, messages.NewMessage(msg.GetSeq(), messages.ActionNotifyForbidden, "ticket expired"))
		return true
	}
	return false
}

func (a *Authenticator) ClientAuthMessageInterceptor(dc DefaultClient, msg *messages.GlideMessage) (intercept bool) {
	if msg.Action != messages.ActionAuthenticate {
		return false
	}

	intercept = true

	var err error
	var errMsg string
	var newId ID
	var span int64
	var authCredentials *ClientAuthCredentials

	credential := EncryptedCredential{}
	err = msg.Data.Deserialize(&credential)
	if err != nil {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	if len(credential.Credential) < 5 {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	authCredentials, err = a.credentialCrypto.DecryptCredentials([]byte(credential.Credential))
	if err != nil {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	span = time.Now().UnixMilli() - authCredentials.Timestamp
	if span > 1500*1000 {
		errMsg = "credential expired"
		goto DONE
	}

	newId, err = a.updateClient(dc, authCredentials)

DONE:

	ac, _ := json.Marshal(authCredentials)
	logger.D("credential: %s", string(ac))

	logger.D("client auth message intercepted %s, %v", dc.GetInfo().ID, err)

	if err != nil || errMsg != "" {
		_ = a.gateway.EnqueueMessage(dc.GetInfo().ID, messages.NewMessage(msg.GetSeq(), messages.ActionNotifyError, errMsg))
	} else {
		_ = a.gateway.EnqueueMessage(newId, messages.NewMessage(msg.GetSeq(), messages.ActionNotifySuccess, nil))
	}
	return
}

func (a *Authenticator) updateClient(dc DefaultClient, authCredentials *ClientAuthCredentials) (ID, error) {

	dc.SetCredentials(authCredentials)

	oldID := dc.GetInfo().ID
	newID := NewID2(authCredentials.UserID)
	err := a.gateway.SetClientID(oldID, newID)
	if IsIDAlreadyExist(err) {
		if newID.Equals(oldID) {
			// already authenticated
			return newID, nil
		}
		tempID, _ := GenTempID("")
		err = a.gateway.SetClientID(newID, tempID)
		if err != nil {
			return "", err
		}
		kickOut := messages.NewMessage(0, messages.ActionNotifyKickOut, &messages.KickOutNotify{
			DeviceName: authCredentials.DeviceName,
			DeviceId:   authCredentials.DeviceID,
		})
		_ = a.gateway.EnqueueMessage(tempID, kickOut)
		err = a.gateway.SetClientID(oldID, newID)
		if err != nil {
			return "", err
		}
	}
	return newID, err
}
