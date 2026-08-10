import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import { fileURLToPath } from "url";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import {
	Architecture,
	LayerVersion,
	Runtime,
	Tracing,
} from "aws-cdk-lib/aws-lambda";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import {
	AnomalyDetectionAlarm,
	ComparisonOperator,
	TreatMissingData,
} from "aws-cdk-lib/aws-cloudwatch";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import { ITable } from "aws-cdk-lib/aws-dynamodb";
import {
	UserPool,
	UserPoolClient,
	UserPoolDomain,
} from "aws-cdk-lib/aws-cognito";

import { IdentityServiceDataStack } from "../data-stack.js";

export interface ApiProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: IdentityServiceDataStack;
	readonly regionalTable: ITable;
	readonly reservedConcurrentExecutions?: number;
	readonly userPool: UserPool;
	readonly userPoolClient: UserPoolClient;
	readonly userPoolDomain: UserPoolDomain;
}

export class Api extends GoFunction {
	constructor(scope: Construct, id: string, props: ApiProps) {
		super(scope, id, {
			entry: fileURLToPath(new URL("./handler", import.meta.url).href).replace(
				"/dist/",
				"/src/",
			),
			architecture: Architecture.ARM_64,
			runtime: Runtime.PROVIDED_AL2023,
			timeout: Duration.seconds(20),
			tracing: Tracing.ACTIVE,
			memorySize: 512,
			reservedConcurrentExecutions: props.reservedConcurrentExecutions,
			bundling: {
				goBuildFlags: ["-trimpath", `-ldflags="-s -w"`],
			},
			layers: [
				LayerVersion.fromLayerVersionArn(
					scope,
					`${id}LambdaAdapterLayer`,
					`arn:aws:lambda:${
						Stack.of(scope).region
					}:753240598075:layer:LambdaAdapterLayerArm64:28`,
				),
			],
			environment: {
				STACK_NAME: Stack.of(scope).stackName,
				NAMESPACE: props.backboneRegionalStack.namespace.namespaceName,
				API_URL: props.backboneRegionalStack.apiGateway.url!,
				HOSTNAME: props.backboneRegionalStack.hostname,
				EVENT_BUS_NAME:
					props.backboneRegionalStack.primaryEventBus.eventBusName,
				EVENT_BUS_ENDPOINT_ID:
					props.backboneRegionalStack.eventBusGlobalEndpoint.attrEndpointId,
				TABLE_NAME: props.dataStack.table.tableName,
				TABLE_WRITE_REGION: props.dataStack.region,

				USER_POOL_ID: props.userPool.userPoolId,
				USER_POOL_URL: props.userPoolDomain.baseUrl(),
				USER_POOL_PROVIDER_URL: props.userPool.userPoolProviderUrl,
				USER_POOL_CLIENT_ID: props.userPoolClient.userPoolClientId,
				USER_POOL_CLIENT_SECRET:
					props.userPoolClient.userPoolClientSecret.unsafeUnwrap(),
			},
		});

		new AnomalyDetectionAlarm(this, "ErrorAnomalyAlarm", {
			alarmDescription: "Anomalous lambda error rate",
			metric: this.metricErrors(),
			evaluationPeriods: 1,
			comparisonOperator: ComparisonOperator.GREATER_THAN_UPPER_THRESHOLD,
			treatMissingData: TreatMissingData.NOT_BREACHING,
		});

		props.backboneRegionalStack.primaryEventBus.grantPutEventsTo(this);
		props.backboneRegionalStack.secondaryEventBus.grantPutEventsTo(this);
		props.dataStack.table.grantReadWriteData(this);
		props.regionalTable.grantReadData(this);
		this.addToRolePolicy(
			new PolicyStatement({
				actions: ["execute-api:Invoke"],
				resources: [props.backboneRegionalStack.apiGateway.arnForExecuteApi()],
			}),
		);
	}
}
