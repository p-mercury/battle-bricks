package dynamo

type File struct {
	Pk   string
	Type string

	UserId        *string
	Bucket        string
	Key           string
	Name          string
	ContentType   *string
	ContentLength *uint64
	Checksum      string
	Status        int32
	Ttl           int64
}
