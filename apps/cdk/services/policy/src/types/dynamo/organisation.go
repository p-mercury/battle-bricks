package dynamo

type Organisation struct {
	Pk           string
	Type         string
	Id           string
	Entitlements struct {
		MaxUsers  uint32
		Platforms struct {
			Sales struct {
				Status              int32
				Integration         int32
				AllowPublicListings bool
				MaxProducts         uint32
			}
			Purchasing struct {
				Status      int32
				Integration int32
			}
		}
	}
	Version uint64
}
