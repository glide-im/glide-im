package gate

import (
	"crypto/sha512"
	"testing"
	"time"

	"github.com/glide-im/glide/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id := NewID("gate", "uid", "dev")
	assert.Equal(t, "gate_uid_dev", string(id))
}

func TestID_SetGateway(t *testing.T) {
	id := NewID2("empty-uid")
	suc := id.SetGateway("gateway")

	assert.True(t, suc)
	assert.Equal(t, "gateway_empty-uid_", string(id))
}

func TestID_SetDevice(t *testing.T) {
	id := NewID2("empty-uid")
	suc := id.SetDevice("device")

	assert.True(t, suc)
	assert.Equal(t, "_empty-uid_device", string(id))
}

func TestID_IsTemp(t *testing.T) {
	id := NewID2(tempIdPrefix + "temp-uid")
	assert.True(t, id.IsTemp())

	id = NewID2("uid")
	assert.False(t, id.IsTemp())
}

func TestID_UID(t *testing.T) {
	id := NewID2("uid")
	assert.Equal(t, "uid", id.UID())
}

func TestID_Gateway(t *testing.T) {
	id := NewID("gate", "uid", "dev")
	assert.Equal(t, "gate", id.Gateway())

	id = NewID2("uid")
	assert.Equal(t, "", id.Gateway())
}

func TestAesCBC_Decrypt(t *testing.T) {

	key := sha512.Sum512([]byte("secret_key"))
	cbcCrypto := NewAesCBCCrypto(key[:])

	credentials := ClientAuthCredentials{
		Type:       1,
		UserID:     "",
		DeviceID:   "1",
		DeviceName: "iPhone 6s",
		Secrets: &ClientSecrets{
			MessageDeliverSecret: "secret",
		},
		ConnectionID: "1",
		Timestamp:    time.Now().UnixMilli(),
	}
	encryptCredentials, err := cbcCrypto.EncryptCredentials(&credentials)
	assert.NoError(t, err)

	t.Log(string(encryptCredentials))
	decryptCredentials, err := cbcCrypto.DecryptCredentials(encryptCredentials)
	assert.NoError(t, err)

	assert.Equal(t, decryptCredentials.UserID, credentials.UserID)
}

func TestGenerateTicket(t *testing.T) {

	secret := "secret"
	sum1 := hash.SHA1(secret + "to")
	expectTicket := hash.SHA1(secret + "from" + sum1)

	t.Log(expectTicket)
}
