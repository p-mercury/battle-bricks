import { Construct } from "constructs";
import { Duration, Stack, StackProps } from "aws-cdk-lib";
import { HttpMethod } from "aws-cdk-lib/aws-events";
import type { BackboneRegionalStack } from "@battle-bricks/backbone";
import {
	Internal,
	External,
} from "@battle-bricks/contracts/notification/v1/service_pb";
import { HttpLambdaIntegration } from "aws-cdk-lib/aws-apigatewayv2-integrations";
import { HttpRoute, HttpRouteKey } from "aws-cdk-lib/aws-apigatewayv2";
import { HttpIamAuthorizer } from "aws-cdk-lib/aws-apigatewayv2-authorizers";

import { ApiHandler } from "./api-handler/index.js";
import { NotificationServiceDataStack } from "./data-stack.js";
import { Projector } from "./projector/index.js";
import { Notifier } from "./notifier/index.js";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

export interface NotificationServiceRegionalStackProps extends StackProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: NotificationServiceDataStack;
}

export class NotificationServiceRegionalStack extends Stack {
	constructor(
		scope: Construct,
		id: string,
		props: NotificationServiceRegionalStackProps,
	) {
		super(scope, id, props);

		const regionalTable: ITable =
			this.region === props.dataStack.region
				? props.dataStack.table
				: props.dataStack.table.replica(this.region);

		new Projector(this, "Projector", {
			backboneRegionalStack: props.backboneRegionalStack,
			dataStack: props.dataStack,
			regionalTable,
			timeout: Duration.seconds(20),
			concurrency: 2,
		});

		new Notifier(this, "Notifier", {
			backboneRegionalStack: props.backboneRegionalStack,
			dataStack: props.dataStack,
			regionalTable,
			timeout: Duration.seconds(20),
			concurrency: 2,
		});

		const internalIntegration = new HttpLambdaIntegration(
			"InternalApiHandler",
			new ApiHandler(this, "InternalApiHandler", {
				backboneRegionalStack: props.backboneRegionalStack,
				dataStack: props.dataStack,
				regionalTable,
			}),
		);
		new HttpRoute(this, `${Internal.typeName}.Get`, {
			httpApi: props.backboneRegionalStack.apiGateway,
			routeKey: HttpRouteKey.with(
				`/${Internal.typeName}/{proxy+}`,
				HttpMethod.GET,
			),
			integration: internalIntegration,
			authorizer: new HttpIamAuthorizer(),
		});
		new HttpRoute(this, `${Internal.typeName}.Post`, {
			httpApi: props.backboneRegionalStack.apiGateway,
			routeKey: HttpRouteKey.with(
				`/${Internal.typeName}/{proxy+}`,
				HttpMethod.POST,
			),
			integration: internalIntegration,
			authorizer: new HttpIamAuthorizer(),
		});

		const externalIntegration = new HttpLambdaIntegration(
			"ExternalApiHandler",
			new ApiHandler(this, "ExternalApiHandler", {
				backboneRegionalStack: props.backboneRegionalStack,
				dataStack: props.dataStack,
				regionalTable,
				reservedConcurrentExecutions: 4,
			}),
		);
		new HttpRoute(this, `${External.typeName}.Get`, {
			httpApi: props.backboneRegionalStack.apiGateway,
			routeKey: HttpRouteKey.with(
				`/${External.typeName}/{proxy+}`,
				HttpMethod.GET,
			),
			integration: externalIntegration,
			authorizer: props.backboneRegionalStack.apiGatewayAuthorizer,
		});
		new HttpRoute(this, `${External.typeName}.Post`, {
			httpApi: props.backboneRegionalStack.apiGateway,
			routeKey: HttpRouteKey.with(
				`/${External.typeName}/{proxy+}`,
				HttpMethod.POST,
			),
			integration: externalIntegration,
			authorizer: props.backboneRegionalStack.apiGatewayAuthorizer,
		});
	}
}
