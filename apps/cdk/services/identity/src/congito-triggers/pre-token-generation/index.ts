import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import { fileURLToPath } from "url";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import {
	Architecture,
	IFunction,
	Runtime,
	Tracing,
} from "aws-cdk-lib/aws-lambda";
import { BackboneStack } from "@battle-bricks/backbone";
import { TableV2 } from "aws-cdk-lib/aws-dynamodb";

export interface PreTokenGenerationHandlerProps {
	readonly backboneStack: BackboneStack;
	readonly table: TableV2;
}

export class PreTokenGenerationHandler extends GoFunction {
	readonly lambda: IFunction;

	constructor(
		scope: Construct,
		id: string,
		props: PreTokenGenerationHandlerProps,
	) {
		super(scope, id, {
			entry: fileURLToPath(new URL("./handler", import.meta.url).href).replace(
				"/dist/",
				"/src/",
			),
			architecture: Architecture.ARM_64,
			runtime: Runtime.PROVIDED_AL2023,
			timeout: Duration.seconds(20),
			tracing: Tracing.ACTIVE,
			logRetention: RetentionDays.ONE_WEEK,
			memorySize: 256,
			environment: {
				STACK_NAME: Stack.of(scope).stackName,
				NAMESPACE: props.backboneStack.namespace,
				HOSTNAME: props.backboneStack.hostname,
				EVENT_BUS_NAME: props.backboneStack.primaryEventBus.eventBusName,
				EVENT_BUS_ENDPOINT_ID:
					props.backboneStack.eventBusGlobalEndpoint.attrEndpointId,
				TABLE_NAME: props.table.tableName,
				TABLE_WRITE_REGION: Stack.of(scope).region,
			},
		});

		props.backboneStack.primaryEventBus.grantPutEventsTo(this);
		props.backboneStack.secondaryEventBus.grantPutEventsTo(this);
		props.table.grantReadWriteData(this);
	}
}
