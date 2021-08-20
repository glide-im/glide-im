package conn

type Serializable interface {
	Deserialize(data []byte) error

	Serialize() ([]byte, error)
}
