import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import { fileURLToPath } from "url";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { Architecture, Runtime, Tracing } from "aws-cdk-lib/aws-lambda";
import * as IdentityEvent from "@battle-bricks/contracts/identity/events/registry";
import { BackboneRegionalStack } from "@battle-bricks/backbone";
import { SqsEventSource } from "aws-cdk-lib/aws-lambda-event-sources";
import { Rule } from "aws-cdk-lib/aws-events";
import { Queue } from "aws-cdk-lib/aws-sqs";
import { SqsQueue } from "aws-cdk-lib/aws-events-targets";
import { TreatMissingData } from "aws-cdk-lib/aws-cloudwatch";

import { ITable } from "aws-cdk-lib/aws-dynamodb";

import { PolicyServiceDataStack } from "../data-stack.js";

export interface ProjectorProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: PolicyServiceDataStack;
	readonly regionalTable: ITable;
	readonly timeout: Duration;
	readonly concurrency: number;
}

export class Projector extends Construct {
	constructor(scope: Construct, id: string, props: ProjectorProps) {
		super(scope, id);

		[
			props.backboneRegionalStack.primaryEventBus,
			props.backboneRegionalStack.secondaryEventBus,
		].forEach((eventBus) => {
			if (eventBus.stack.region !== Stack.of(this).region) return;

			const handler = new GoFunction(this, "Handler", {
				entry: fileURLToPath(
					new URL("./handler", import.meta.url).href,
				).replace("/dist/", "/src/"),
				architecture: Architecture.ARM_64,
				runtime: Runtime.PROVIDED_AL2023,
				timeout: props.timeout,
				tracing: Tracing.ACTIVE,
				memorySize: 128,
				reservedConcurrentExecutions: props.concurrency,
				bundling: {
					goBuildFlags: ["-trimpath", `-ldflags="-s -w"`],
				},
				environment: {
					STACK_NAME: Stack.of(scope).stackName,
					NAMESPACE: props.backboneRegionalStack.namespace.namespaceName,
					TABLE_NAME: props.dataStack.table.tableName,
					TABLE_WRITE_REGION: props.dataStack.region,
				},
			});

			props.dataStack.table.grantReadWriteData(handler);
			props.regionalTable.grantReadData(handler);

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

			const bufferQueue = new Queue(this, "Queue", {
				retentionPeriod: Duration.days(2),
				visibilityTimeout: Duration.seconds(6 * props.timeout.toSeconds()),
				deadLetterQueue: {
					maxReceiveCount: 5,
					queue: deadLetterQueue,
				},
			});

			handler.addEventSource(
				new SqsEventSource(bufferQueue, {
					maxConcurrency: props.concurrency,
					batchSize: 10,
					reportBatchItemFailures: true,
				}),
			);

			new Rule(this, "IdentityRule", {
				eventBus,
				eventPattern: {
					source: [
						IdentityEvent.getSource(
							props.backboneRegionalStack.namespace.namespaceName,
						),
					],
					detailType: [
						IdentityEvent.UserCreatedDetailType,
						IdentityEvent.UserUpdatedDetailType,
						IdentityEvent.UserDeletedDetailType,
					],
				},
				targets: [new SqsQueue(bufferQueue)],
			});
		});
	}
}
