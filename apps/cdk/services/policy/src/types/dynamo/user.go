package dynamo

type User struct {
	Pk      string
	Gsi1pk  string
	Gsi1sk  string
	Type    string
	Status  int32
	Version uint64
}
