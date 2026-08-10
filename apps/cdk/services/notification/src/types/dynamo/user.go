package dynamo

import (
	"contracts/dist/common/v1"
	"contracts/dist/identity/v1"
)

type User struct {
	Pk                   string
	Name                 string
	EmailAddress         string
	Language             common.Language
	Version              uint64
	NotificationSettings *struct {
		Newsletter                NotificationSettingsTopic
		PurchasingQuotationPriced *NotificationSettingsTopic
		PurchasingOrderConfirmed  *NotificationSettingsTopic
	}
}

type NotificationSettingsTopic struct {
	Enabled bool
	Methods []identity.UserNotificationSettings_Topic_Method
}
