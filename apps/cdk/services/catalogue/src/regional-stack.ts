import { Construct } from "constructs";
import { Stack, StackProps } from "aws-cdk-lib";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import {
	Internal,
	External,
} from "@battle-bricks/contracts/catalogue/v1/service_pb";
import {
	HttpMethod,
	HttpRoute,
	HttpRouteKey,
} from "aws-cdk-lib/aws-apigatewayv2";
import { HttpIamAuthorizer } from "aws-cdk-lib/aws-apigatewayv2-authorizers";
import { HttpLambdaIntegration } from "aws-cdk-lib/aws-apigatewayv2-integrations";

import { Api } from "./api/index.js";
import { CatalogueServiceDataStack } from "./data-stack.js";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

export interface CatalogueServiceRegionalStackProps extends StackProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: CatalogueServiceDataStack;
}

export class CatalogueServiceRegionalStack extends Stack {
	constructor(
		scope: Construct,
		id: string,
		props: CatalogueServiceRegionalStackProps,
	) {
		super(scope, id, props);

		const regionalTable: ITable =
			this.region === props.dataStack.region
				? props.dataStack.table
				: props.dataStack.table.replica(this.region);

		const internalIntegration = new HttpLambdaIntegration(
			"InternalApiHandler",
			new Api(this, "InternalApiHandler", {
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
			new Api(this, "ExternalApiHandler", {
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
