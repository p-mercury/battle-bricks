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
import { ITable } from "aws-cdk-lib/aws-dynamodb";
import { IEventBus } from "aws-cdk-lib/aws-events";
import { HttpNamespace } from "aws-cdk-lib/aws-servicediscovery";

export interface PostConfirmationHandlerProps {
	readonly table: ITable;
	readonly eventBus: IEventBus;
	readonly namespace: HttpNamespace;
}

export class PostConfirmationHandler extends GoFunction {
	readonly lambda: IFunction;

	constructor(
		scope: Construct,
		id: string,
		props: PostConfirmationHandlerProps,
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
				SERVICE_NAME: "Identity",
				NAMESPACE: props.namespace.namespaceName,
				TABLE_NAME: props.table.tableName,
				EVENT_BUS_NAME: props.eventBus.eventBusName,
			},
		});

		props.table.grantReadWriteData(this);
		props.eventBus.grantPutEventsTo(this);
	}
}
