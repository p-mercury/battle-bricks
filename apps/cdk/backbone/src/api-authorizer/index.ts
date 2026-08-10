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
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import {
	AnomalyDetectionAlarm,
	ComparisonOperator,
	TreatMissingData,
} from "aws-cdk-lib/aws-cloudwatch";
import { HttpApi } from "aws-cdk-lib/aws-apigatewayv2";

export interface ApiAuthorizerProps {
	readonly apiGateway: HttpApi;
	readonly reservedConcurrentExecutions?: number;
}

export class ApiAuthorizer extends GoFunction {
	readonly lambda: IFunction;

	constructor(scope: Construct, id: string, props: ApiAuthorizerProps) {
		super(scope, id, {
			entry: fileURLToPath(new URL("./handler", import.meta.url).href).replace(
				"/dist/",
				"/src/",
			),
			architecture: Architecture.ARM_64,
			runtime: Runtime.PROVIDED_AL2023,
			timeout: Duration.seconds(6),
			tracing: Tracing.ACTIVE,
			memorySize: 256,
			reservedConcurrentExecutions: props.reservedConcurrentExecutions,
			bundling: {
				goBuildFlags: ["-trimpath", `-ldflags="-s -w"`],
			},
			environment: {
				STACK_NAME: Stack.of(scope).stackName,
				API_URL: props.apiGateway.url!,
			},
		});

		new AnomalyDetectionAlarm(this, "ErrorAnomalyAlarm", {
			alarmDescription: "Anomalous lambda error rate",
			metric: this.metricErrors(),
			evaluationPeriods: 1,
			comparisonOperator: ComparisonOperator.GREATER_THAN_UPPER_THRESHOLD,
			treatMissingData: TreatMissingData.NOT_BREACHING,
		});

		this.addToRolePolicy(
			new PolicyStatement({
				actions: ["execute-api:Invoke"],
				resources: [props.apiGateway.arnForExecuteApi()],
			}),
		);
	}
}
