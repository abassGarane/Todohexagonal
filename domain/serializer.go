package domain

type Serializer interface {
	Decode(input []byte) (*Todo, error)
	Encode(input *Todo) ([]byte, error)
}
