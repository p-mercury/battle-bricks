import { Construct } from "constructs";
import { Bucket, BucketEncryption } from "aws-cdk-lib/aws-s3";
import { BucketDeployment, Source } from "aws-cdk-lib/aws-s3-deployment";
import { fileURLToPath } from "url";
import { Duration, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import { HttpMethod } from "aws-cdk-lib/aws-events";
import {
	Internal,
	External,
} from "@battle-bricks/contracts/policy/v1/service_pb";
import { HttpLambdaIntegration } from "aws-cdk-lib/aws-apigatewayv2-integrations";
import { HttpRoute, HttpRouteKey } from "aws-cdk-lib/aws-apigatewayv2";
import { HttpIamAuthorizer } from "aws-cdk-lib/aws-apigatewayv2-authorizers";

import { Api } from "./api/index.js";
import { PolicyServiceDataStack } from "./data-stack.js";
import { Projector } from "./projector/index.js";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

export interface PolicyServiceRegionalStackProps extends StackProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: PolicyServiceDataStack;
}

export class PolicyServiceRegionalStack extends Stack {
	constructor(
		scope: Construct,
		id: string,
		props: PolicyServiceRegionalStackProps,
	) {
		super(scope, id, props);

		const regionalTable: ITable =
			this.region === props.dataStack.region
				? props.dataStack.table
				: props.dataStack.table.replica(this.region);

		const bucket = new Bucket(this, "Bucket", {
			enforceSSL: true,
			encryption: BucketEncryption.S3_MANAGED,
			removalPolicy: RemovalPolicy.DESTROY,
			autoDeleteObjects: true,
		});

		new BucketDeployment(this, "DeployBundle", {
			destinationBucket: bucket,
			sources: [
				Source.asset(fileURLToPath(new URL(".", import.meta.url).href)),
			],
			exclude: ["*"],
			include: ["bundle.tar.gz"],
		});

		new Projector(this, "Projector", {
			backboneRegionalStack: props.backboneRegionalStack,
			dataStack: props.dataStack,
			regionalTable,
			timeout: Duration.seconds(20),
			concurrency: 2,
		});

		const internalIntegration = new HttpLambdaIntegration(
			"InternalApiHandler",
			new Api(this, "InternalApiHandler", {
				bucket,
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
				bucket,
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
