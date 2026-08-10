import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import { fileURLToPath } from "url";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import {
	Architecture,
	IFunction,
	Runtime,
	Tracing,
} from "aws-cdk-lib/aws-lambda";
import {
	AnomalyDetectionAlarm,
	ComparisonOperator,
	TreatMissingData,
} from "aws-cdk-lib/aws-cloudwatch";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

import { FileStagingServiceDataStack } from "../data-stack.js";

export interface UploadHandlerProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: FileStagingServiceDataStack;
	readonly bucket: Bucket;
	readonly regionalTable: ITable;
	readonly timeout?: Duration;
	readonly reservedConcurrentExecutions?: number;
}

export class UploadHandler extends GoFunction {
	readonly lambda: IFunction;

	constructor(scope: Construct, id: string, props: UploadHandlerProps) {
		super(scope, id, {
			entry: fileURLToPath(new URL("./handler", import.meta.url).href).replace(
				"/dist/",
				"/src/",
			),
			architecture: Architecture.ARM_64,
			runtime: Runtime.PROVIDED_AL2023,
			timeout: props.timeout,
			tracing: Tracing.ACTIVE,
			memorySize: 512,
			reservedConcurrentExecutions: props.reservedConcurrentExecutions,
			environment: {
				STACK_NAME: Stack.of(scope).stackName,
				NAMESPACE: props.backboneRegionalStack.namespace.namespaceName,
				EVENT_BUS_NAME:
					props.backboneRegionalStack.primaryEventBus.eventBusName,
				EVENT_BUS_ENDPOINT_ID:
					props.backboneRegionalStack.eventBusGlobalEndpoint.attrEndpointId,
				TABLE_NAME: props.dataStack.table.tableName,
				TABLE_WRITE_REGION: props.dataStack.region,
				BUCKET_NAME: props.bucket.bucketName,
			},
		});

		new AnomalyDetectionAlarm(this, "ErrorAnomalyAlarm", {
			alarmDescription: "Anomalous lambda error rate",
			metric: this.metricErrors(),
			evaluationPeriods: 1,
			comparisonOperator: ComparisonOperator.GREATER_THAN_UPPER_THRESHOLD,
			treatMissingData: TreatMissingData.NOT_BREACHING,
		});

		props.bucket.grantReadWrite(this);
		props.backboneRegionalStack.primaryEventBus.grantPutEventsTo(this);
		props.backboneRegionalStack.secondaryEventBus.grantPutEventsTo(this);
		props.regionalTable.grantReadWriteData(this);
	}
}
