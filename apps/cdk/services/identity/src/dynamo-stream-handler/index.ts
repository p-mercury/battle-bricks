import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import { fileURLToPath } from "url";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import {
	Architecture,
	Runtime,
	StartingPosition,
	Tracing,
} from "aws-cdk-lib/aws-lambda";
import { TreatMissingData } from "aws-cdk-lib/aws-cloudwatch";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import {
	DynamoEventSource,
	SqsDlq,
} from "aws-cdk-lib/aws-lambda-event-sources";
import { Queue } from "aws-cdk-lib/aws-sqs";

import { ITable } from "aws-cdk-lib/aws-dynamodb";

import { IdentityServiceDataStack } from "../data-stack.js";

export interface DynamoStreamHandlerProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: IdentityServiceDataStack;
	readonly regionalTable: ITable;
	readonly timeout: Duration;
	readonly concurrency: number;
}

export class DynamoStreamHandler extends GoFunction {
	constructor(scope: Construct, id: string, props: DynamoStreamHandlerProps) {
		super(scope, id, {
			entry: fileURLToPath(new URL("./handler", import.meta.url).href).replace(
				"/dist/",
				"/src/",
			),
			architecture: Architecture.ARM_64,
			runtime: Runtime.PROVIDED_AL2023,
			timeout: props.timeout,
			tracing: Tracing.ACTIVE,
			memorySize: 128,
			reservedConcurrentExecutions: props.concurrency * 2,
			bundling: {
				goBuildFlags: ["-trimpath", `-ldflags="-s -w"`],
			},
			environment: {
				STACK_NAME: Stack.of(scope).stackName,
				NAMESPACE: props.backboneRegionalStack.namespace.namespaceName,
				EVENT_BUS_NAME:
					props.backboneRegionalStack.primaryEventBus.eventBusName,
				EVENT_BUS_ENDPOINT_ID:
					props.backboneRegionalStack.eventBusGlobalEndpoint.attrEndpointId,
				TABLE_NAME: props.dataStack.table.tableName,
				TABLE_WRITE_REGION: props.dataStack.region,
			},
		});

		const deadLetterQueue = new Queue(this, "Dlq", {
			retentionPeriod: Duration.days(14),
		});

		deadLetterQueue
			.metricApproximateNumberOfMessagesVisible()
			.createAlarm(this, "DlqAlarm", {
				threshold: 1,
				evaluationPeriods: 1,
				treatMissingData: TreatMissingData.NOT_BREACHING,
			});

		this.addEventSource(
			new DynamoEventSource(props.dataStack.table, {
				startingPosition: StartingPosition.TRIM_HORIZON,
				reportBatchItemFailures: true,
				maxBatchingWindow: Duration.seconds(2),
				retryAttempts: 2,
				batchSize: 10,
				onFailure: new SqsDlq(deadLetterQueue),
				parallelizationFactor: props.concurrency,
			}),
		);

		props.backboneRegionalStack.primaryEventBus.grantPutEventsTo(this);
		props.backboneRegionalStack.secondaryEventBus.grantPutEventsTo(this);
		props.dataStack.table.grantReadWriteData(this);
		props.regionalTable.grantReadData(this);
	}
}
