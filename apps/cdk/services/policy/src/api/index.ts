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
import { Bucket } from "aws-cdk-lib/aws-s3";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

import { PolicyServiceDataStack } from "../data-stack.js";

export interface ApiProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: PolicyServiceDataStack;
	readonly bucket: Bucket;
	readonly regionalTable: ITable;
	readonly reservedConcurrentExecutions?: number;
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
			timeout: Duration.seconds(4),
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
				EVENT_BUS_NAME:
					props.backboneRegionalStack.primaryEventBus.eventBusName,
				EVENT_BUS_ENDPOINT_ID:
					props.backboneRegionalStack.eventBusGlobalEndpoint.attrEndpointId,
				TABLE_NAME: props.dataStack.table.tableName,
				TABLE_WRITE_REGION: props.dataStack.region,
				BUCKET_URL: props.bucket.urlForObject(),
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
		props.bucket.grantRead(this);
		this.addToRolePolicy(
			new PolicyStatement({
				actions: ["execute-api:Invoke"],
				resources: [props.backboneRegionalStack.apiGateway.arnForExecuteApi()],
			}),
		);
	}
}
