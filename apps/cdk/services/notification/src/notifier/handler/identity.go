package main

import (
	"context"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"
	"dynamokit/dynamolease"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func handlerIdentityEvent(ctx context.Context, event events.EventBridgeEvent) error {
	switch event.DetailType {
	case identityevents.UserCreatedDetailType:
		var details identityevents.UserUpserted
		if err := protojson.Unmarshal(event.Detail, &details); err != nil {
			Logger.Error("Error parsing UserUpserted event", slog.Any("error", err))
			return err
		}

		incomingLease, err := Leaser.Acquire(ctx, &dynamolease.AcquireInput{
			Key:     "EVENT#" + event.Source + "." + event.DetailType + "." + details.Id + "." + strconv.FormatUint(details.Version, 10),
			Version: details.Version,
			Ttl:     new(time.Now().Add(time.Hour * 120)),
		})
		defer incomingLease.Cancel(ctx)
		if incomingLease == nil {
			if err != nil {
				return fmt.Errorf("acquiring incoming lease: %w", err)
			}
			return err
		}

		var passwordChangeUrl string
		{
			resp, err := IdentityService.CreatePasswordChangeSession(ctx, connect.NewRequest(&identity.CreatePasswordChangeSessionRequest{
				Id: &details.Id,
			}))
			if err != nil {
				Logger.Error("Error creating password change session", slog.Any("error", err))
				return err
			}
			passwordChangeUrl = resp.Msg.Session.Url
		}

		var templateData string
		{
			a := map[string]any{
				"subject": "Invitation",
				"preActionParagraphs": []map[string]string{
					{
						"content": "Hi " + details.Name + ",",
					},
					{
						"content": "Welcome to the new jumper.de, your one-stop shop for all things backend testing. Full Rasco spare part support and more features are coming in Q2 2026. For now, we support change kits and sockets for any and all handlers you have on site.",
					},
					{
						"content": "The jumper.de platform is now part of Vierroth GmbH. Everything you're used to stays the same, just with more functionality than before. Going forward, all queries from Jumper Systems and Vierroth GmbH will run through this platform to streamline your workflow and simplify purchasing.",
					},
					{
						"content": "To keep your access, please take a moment to reauthenticate your account:",
					},
				},
				"action": map[string]string{
					"label": "Reauthenticate your account",
					"href":  passwordChangeUrl,
				},
				"postActionParagraphs": []map[string]string{
					{
						"content": "If you have any questions, just reply to this email.",
					},
				},
				"signature": map[string]string{
					"thanks": "Best regards,",
					"name":   "Your jumper.de team",
				},
			}
			b, err := json.Marshal(a)
			if err != nil {
				Logger.Error("Error marshaling template data", slog.Any("error", err))
				return err
			}
			templateData = string(b)
		}

		_, err = Ses.SendTemplatedEmail(ctx, &ses.SendTemplatedEmailInput{
			Template: SingleActionSesTemplateName,
			Source:   new("jumper.de <support@jumper.de>"),
			Destination: &sesTypes.Destination{
				ToAddresses: []string{details.EmailAddress},
			},
			TemplateData: &templateData,
		})
		if err != nil {
			Logger.Error("Error sending email through ses", slog.Any("error", err))
			return err
		}

		incomingLease.Lock(ctx, &dynamolease.LockLeaseInput{
			Version: details.Version,
			Ttl:     new(time.Now().Add(time.Hour * 120)),
		})
	}

	return nil
}
