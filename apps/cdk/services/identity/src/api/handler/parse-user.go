package main

import (
	"contracts/dist/common/v1"
	"contracts/dist/identity/v1"
	"identity/src/types/dynamo"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseDynamoUser(item *dynamo.User) *identity.User {
	correlations := make([]*common.Correlation, 0, len(item.Correlations))
	for _, c := range item.Correlations {
		correlations = append(correlations, &common.Correlation{
			Provider: c.Provider,
			Kind:     c.Kind,
			Id:       c.Id,
		})
	}

	notificationSettings := &identity.UserNotificationSettings{
		Newsletter: &identity.UserNotificationSettings_Topic{
			Enabled: item.NotificationSettings.Newsletter.Enabled,
			Methods: item.NotificationSettings.Newsletter.Methods,
		},
		PurchasingQuotationPriced: &identity.UserNotificationSettings_Topic{
			Enabled: true,
			Methods: []identity.UserNotificationSettings_Topic_Method{
				identity.UserNotificationSettings_Topic_METHOD_EMAIL,
				identity.UserNotificationSettings_Topic_METHOD_IN_APP,
				identity.UserNotificationSettings_Topic_METHOD_PUSH,
			},
		},
		PurchasingOrderConfirmed: &identity.UserNotificationSettings_Topic{
			Enabled: true,
			Methods: []identity.UserNotificationSettings_Topic_Method{
				identity.UserNotificationSettings_Topic_METHOD_EMAIL,
				identity.UserNotificationSettings_Topic_METHOD_IN_APP,
				identity.UserNotificationSettings_Topic_METHOD_PUSH,
			},
		},
	}
	if item.NotificationSettings.PurchasingQuotationPriced != nil {
		notificationSettings.PurchasingQuotationPriced = &identity.UserNotificationSettings_Topic{
			Enabled: item.NotificationSettings.PurchasingQuotationPriced.Enabled,
			Methods: item.NotificationSettings.PurchasingQuotationPriced.Methods,
		}
	}
	if item.NotificationSettings.PurchasingOrderConfirmed != nil {
		notificationSettings.PurchasingOrderConfirmed = &identity.UserNotificationSettings_Topic{
			Enabled: item.NotificationSettings.PurchasingOrderConfirmed.Enabled,
			Methods: item.NotificationSettings.PurchasingOrderConfirmed.Methods,
		}
	}

	return &identity.User{
		Id:           item.Id,
		Correlations: correlations,
		Version:      item.Version,
		CreatedTime:  timestamppb.New(time.UnixMilli(item.CreatedTime)),
		ModifiedTime: timestamppb.New(time.UnixMilli(item.ModifiedTime)),

		EmailAddress:         item.EmailAddress,
		Status:               item.Status,
		Name:                 item.Name,
		Language:             item.Language,
		NotificationSettings: notificationSettings,
	}
}
