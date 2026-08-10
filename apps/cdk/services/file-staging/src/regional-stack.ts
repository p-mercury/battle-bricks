import { Construct } from "constructs";
import { Duration, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import type { BackboneRegionalStack } from "@battle-bricks/backbone";
import { Bucket, EventType, HttpMethods } from "aws-cdk-lib/aws-s3";
import { SqsDestination } from "aws-cdk-lib/aws-s3-notifications";
import { Queue } from "aws-cdk-lib/aws-sqs";
import { SqsEventSource } from "aws-cdk-lib/aws-lambda-event-sources";
import { CfnMalwareProtectionPlan } from "aws-cdk-lib/aws-guardduty";
import {
	CfnRole,
	PolicyDocument,
	PolicyStatement,
	Role,
	ServicePrincipal,
} from "aws-cdk-lib/aws-iam";
import {
	Internal,
	External,
} from "@battle-bricks/contracts/filestaging/v1/service_pb";
import { HttpMethod } from "aws-cdk-lib/aws-events";
import { HttpRoute, HttpRouteKey } from "aws-cdk-lib/aws-apigatewayv2";
import { HttpLambdaIntegration } from "aws-cdk-lib/aws-apigatewayv2-integrations";
import { HttpIamAuthorizer } from "aws-cdk-lib/aws-apigatewayv2-authorizers";
import { TreatMissingData } from "aws-cdk-lib/aws-cloudwatch";

import { Api } from "./api/index.js";
import { FileStagingServiceDataStack } from "./data-stack.js";
import { UploadHandler } from "./upload-handler/index.js";
import { ITable } from "aws-cdk-lib/aws-dynamodb";

export interface FileStagingServiceRegionalStackProps extends StackProps {
	readonly backboneRegionalStack: BackboneRegionalStack;
	readonly dataStack: FileStagingServiceDataStack;
}

export class FileStagingServiceRegionalStack extends Stack {
	constructor(
		scope: Construct,
		id: string,
		props: FileStagingServiceRegionalStackProps,
	) {
		super(scope, id, props);

		const regionalTable: ITable =
			this.region === props.dataStack.region
				? props.dataStack.table
				: props.dataStack.table.replica(this.region);

		const bucket = new Bucket(this, "Bucket", {
			removalPolicy: RemovalPolicy.DESTROY,
			autoDeleteObjects: true,
			enforceSSL: true,
			lifecycleRules: [
				{
					abortIncompleteMultipartUploadAfter: Duration.days(1),
				},
				{
					expiration: Duration.days(31),
					expiredObjectDeleteMarker: false,
				},
			],
			cors: [
				{
					allowedHeaders: ["*"],
					allowedMethods: [HttpMethods.GET, HttpMethods.PUT, HttpMethods.POST],
					allowedOrigins: ["*"],
				},
			],
		});

		const malwareScanRole = new Role(this, "GuardDutyMalwareScanRole", {
			assumedBy: new ServicePrincipal(
				"malware-protection-plan.guardduty.amazonaws.com",
			),
			inlinePolicies: {
				S3Policy: new PolicyDocument({
					statements: [
						new PolicyStatement({
							sid: "AllowManagedRuleToSendS3EventsToGuardDuty",
							actions: [
								"events:PutRule",
								"events:DeleteRule",
								"events:PutTargets",
								"events:RemoveTargets",
							],
							resources: [
								`arn:aws:events:${this.region}:${this.account}:rule/DO-NOT-DELETE-AmazonGuardDutyMalwareProtectionS3*`,
							],
							conditions: {
								StringLike: {
									"events:ManagedBy":
										"malware-protection-plan.guardduty.amazonaws.com",
								},
							},
						}),
						new PolicyStatement({
							sid: "AllowGuardDutyToMonitorEventBridgeManagedRule",
							actions: ["events:DescribeRule", "events:ListTargetsByRule"],
							resources: [
								`arn:aws:events:${this.region}:${this.account}:rule/DO-NOT-DELETE-AmazonGuardDutyMalwareProtectionS3*`,
							],
						}),
						new PolicyStatement({
							sid: "AllowEnableS3EventBridgeEvents",
							actions: ["s3:PutBucketNotification", "s3:GetBucketNotification"],
							resources: [bucket.bucketArn],
						}),
						new PolicyStatement({
							sid: "AllowCheckBucketOwnership",
							actions: ["s3:ListBucket"],
							resources: [bucket.bucketArn],
						}),
						new PolicyStatement({
							sid: "AllowPutValidationObject",
							actions: ["s3:PutObject"],
							resources: [
								`${bucket.bucketArn}/malware-protection-resource-validation-object`,
							],
						}),
						new PolicyStatement({
							sid: "AllowMalwareScan",
							actions: ["s3:GetObject", "s3:GetObjectVersion"],
							resources: [bucket.arnForObjects("*")],
						}),
						new PolicyStatement({
							sid: "AllowPostScanTag",
							actions: [
								"s3:PutObjectTagging",
								"s3:GetObjectTagging",
								"s3:PutObjectVersionTagging",
								"s3:GetObjectVersionTagging",
							],
							resources: [bucket.arnForObjects("*")],
						}),
					],
				}),
			},
		});

		const plan = new CfnMalwareProtectionPlan(this, "MalwareProtectionPlan", {
			protectedResource: {
				s3Bucket: {
					bucketName: bucket.bucketName,
				},
			},
			role: malwareScanRole.roleArn,
			actions: {
				tagging: {
					status: "ENABLED",
				},
			},
		});

		plan.addDependency(malwareScanRole.node.defaultChild as CfnRole);

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
			visibilityTimeout: Duration.seconds(80),
			deadLetterQueue: {
				maxReceiveCount: 5,
				queue: deadLetterQueue,
			},
		});

		bucket.addEventNotification(
			EventType.OBJECT_TAGGING_PUT,
			new SqsDestination(bufferQueue),
		);

		const uploadHandler = new UploadHandler(this, "UploadHandler", {
			backboneRegionalStack: props.backboneRegionalStack,
			dataStack: props.dataStack,
			bucket,
			regionalTable,
			timeout: Duration.seconds(40),
			reservedConcurrentExecutions: 4,
		});

		uploadHandler.addEventSource(
			new SqsEventSource(bufferQueue, {
				maxConcurrency: 4,
				batchSize: 2,
				reportBatchItemFailures: true,
			}),
		);

		const internalIntegration = new HttpLambdaIntegration(
			"InternalApiHandler",
			new Api(this, "InternalApiHandler", {
				backboneRegionalStack: props.backboneRegionalStack,
				dataStack: props.dataStack,
				bucket,
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
				bucket,
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
