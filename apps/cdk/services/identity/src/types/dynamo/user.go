package dynamo

import (
	"contracts/dist/common/v1"
	"contracts/dist/identity/v1"
)

type User struct {
	Pk     string
	Sk     string
	Gsi1pk string
	Gsi1sk string
	Type   string

	Id           string
	Correlations map[string]Correlation
	Version      uint64
	CreatedTime  int64
	ModifiedTime int64

	OrganisationId       string
	Status               identity.UserStatus
	EmailAddress         string
	Name                 *string
	JobTitle             *string
	Auth0Id              string
	Language             common.Language
	NotificationSettings struct {
		Newsletter                NotificationSettingsTopic
		PurchasingQuotationPriced *NotificationSettingsTopic
		PurchasingOrderConfirmed  *NotificationSettingsTopic
	}
}

type NotificationSettingsTopic struct {
	Enabled bool
	Methods []identity.UserNotificationSettings_Topic_Method
}
