package dynamo

type CustomerUser struct {
	Pk           string
	Gsi1pk       string
	Gsi1sk       string
	Type         string
	Status       int32
	Version      uint64
	Organisation *struct {
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
}
