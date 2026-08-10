import { Duration, Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Certificate, ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import {
	CfnHealthCheck,
	IPublicHostedZone,
	PublicHostedZone,
} from "aws-cdk-lib/aws-route53";
import { CfnEndpoint, EventBus } from "aws-cdk-lib/aws-events";
import {
	Alarm,
	ComparisonOperator,
	Metric,
	TreatMissingData,
} from "aws-cdk-lib/aws-cloudwatch";
import {
	PolicyDocument,
	PolicyStatement,
	Role,
	ServicePrincipal,
} from "aws-cdk-lib/aws-iam";

import { EventBusStack } from "./event-bus-stack.js";

export interface BackboneStackProps extends StackProps {
	readonly stage: "DEVELOPMENT" | "PRODUCTION";
	readonly hostedZoneName: string;
	readonly certificateId: string;
	readonly hostedZoneId: string;
	readonly namespace: string;
	readonly hostname: string;
	readonly primaryEventBusStack: EventBusStack;
	readonly secondaryEventBusStack: EventBusStack;
	readonly activeRegions: string[];
}

export class BackboneStack extends Stack {
	readonly stage: "DEVELOPMENT" | "PRODUCTION";
	readonly certificate: ICertificate;
	readonly hostname: string;
	readonly namespace: string;
	readonly hostedZone: IPublicHostedZone;
	readonly eventBusGlobalEndpoint: CfnEndpoint;
	readonly primaryEventBus: EventBus;
	readonly secondaryEventBus: EventBus;
	readonly activeRegions: string[];

	constructor(scope: Construct, id: string, props: BackboneStackProps) {
		const {
			stage,
			hostname,
			namespace,
			primaryEventBusStack,
			secondaryEventBusStack,
			activeRegions,
			...stackProps
		} = props;

		super(scope, id, stackProps);

		this.stage = stage;
		this.hostname = hostname;
		this.namespace = namespace;
		this.primaryEventBus = primaryEventBusStack.eventBus;
		this.secondaryEventBus = secondaryEventBusStack.eventBus;
		this.activeRegions = activeRegions;

		const alarm = new Alarm(this, "EventBusLatencyAlarm", {
			alarmDescription: "High EventBridge latency",
			metric: new Metric({
				metricName: "IngestionToInvocationStartLatency",
				namespace: "AWS/Events",
				statistic: "Average",
				period: Duration.seconds(60),
			}),
			evaluationPeriods: 5,
			threshold: 30000,
			comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
			treatMissingData: TreatMissingData.MISSING,
		});

		const healthCheck = new CfnHealthCheck(this, "EventBusHealthCheck", {
			healthCheckTags: [
				{
					key: "Name",
					value: "LatencyFailuresHealthCheck",
				},
			],
			healthCheckConfig: {
				type: "CLOUDWATCH_METRIC",
				alarmIdentifier: {
					name: alarm.alarmName,
					region: this.region,
				},
				insufficientDataHealthStatus: "Healthy",
			},
		});

		this.eventBusGlobalEndpoint = new CfnEndpoint(
			this,
			"EventBusGlobalEndpoint",
			{
				eventBuses: [
					{ eventBusArn: props.primaryEventBusStack.eventBus.eventBusArn },
					{ eventBusArn: props.secondaryEventBusStack.eventBus.eventBusArn },
				],
				routingConfig: {
					failoverConfig: {
						primary: {
							healthCheck: `arn:aws:route53:::healthcheck/${healthCheck.attrHealthCheckId}`,
						},
						secondary: { route: props.secondaryEventBusStack.region },
					},
				},
				replicationConfig: {
					state: "ENABLED",
				},
				roleArn: new Role(this, "EndpointReplicationRole", {
					assumedBy: new ServicePrincipal("events.amazonaws.com"),
					inlinePolicies: {
						Replication: new PolicyDocument({
							statements: [
								new PolicyStatement({
									actions: [
										"events:PutRule",
										"events:PutTargets",
										"events:RemoveTargets",
										"events:DeleteRule",
										"events:DescribeRule",
										"events:ListTargetsByRule",
										"events:TagResource",
									],
									resources: [`arn:aws:events:*:${this.account}:rule/*`],
								}),
								new PolicyStatement({
									actions: ["events:PutEvents"],
									resources: [
										props.primaryEventBusStack.eventBus.eventBusArn,
										props.secondaryEventBusStack.eventBus.eventBusArn,
									],
								}),
								new PolicyStatement({
									actions: ["iam:PassRole"],
									resources: ["*"],
									conditions: {
										StringEquals: {
											"iam:PassedToService": "events.amazonaws.com",
										},
									},
								}),
							],
						}),
					},
				}).roleArn,
			},
		);

		this.hostedZone = PublicHostedZone.fromPublicHostedZoneAttributes(
			this,
			"HostedZone",
			{
				zoneName: props.hostedZoneName,
				hostedZoneId: props.hostedZoneId,
			},
		);

		this.certificate = Certificate.fromCertificateArn(
			this,
			"CustomDomainCertificate",
			`arn:aws:acm:us-east-1:${this.account}:certificate/${
				props.certificateId
			}`,
		);
	}
}
