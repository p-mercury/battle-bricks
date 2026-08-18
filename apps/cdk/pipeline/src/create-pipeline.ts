import { Annotations, App } from "aws-cdk-lib";
import {
	Pipeline,
	PipelineSegment,
	StackSegment,
	Artifact,
	CodeStarSourceSegment,
} from "@flit/cdk-pipeline";
import {
	BuildSpec,
	ComputeType,
	LinuxArmBuildImage,
} from "aws-cdk-lib/aws-codebuild";
import {
	BackboneRegionalStack,
	BackboneStack,
	EventBusStack,
} from "@battle-bricks/backbone";
import { FrontEndStack } from "@battle-bricks/front-end";
import {
	NotificationServiceDataStack,
	NotificationServiceRegionalStack,
} from "@battle-bricks/notification-service";
import { IdentityServiceDataRegionalStack } from "@battle-bricks/identity-service/dist/data-regional-stack";
import {
	IdentityServiceDataStack,
	IdentityServiceRegionalStack,
} from "@battle-bricks/identity-service";
import {
	CatalogueServiceDataStack,
	CatalogueServiceRegionalStack,
} from "@battle-bricks/catalogue-service";
import {
	FileStagingServiceDataStack,
	FileStagingServiceRegionalStack,
} from "@battle-bricks/file-staging-service";
import {
	PolicyServiceDataStack,
	PolicyServiceRegionalStack,
} from "@battle-bricks/policy-service";

export type CreatePipelineProps = {
	pipelineName: string;
	buildCommands: string[];

	stage: "PRODUCTION" | "DEVELOPMENT";

	account: string;
	stackPrefix: string;
	tags: { [key: string]: string };

	hostedZoneName: string;
	hostedZoneId: string;
	hostname: string;
	certificateId: string;
	namespace: string;
	activeRegions: string[];

	environmentVariables?: Record<string, any>;
};

export function createPipeline(scope: App, props: CreatePipelineProps) {
	const EVENT_BUS_EU_CENTRAL = new EventBusStack(
		scope,
		`${props.stackPrefix}EventBusEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}EventBus`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
		},
	);

	const EVENT_BUS_EU_WEST = new EventBusStack(
		scope,
		`${props.stackPrefix}EventBusEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}EventBus`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
		},
	);

	const BACKBONE = new BackboneStack(scope, `${props.stackPrefix}Backbone`, {
		stackName: `BattleBricks${props.stackPrefix}Backbone`,
		terminationProtection: true,
		crossRegionReferences: true,
		tags: props.tags,
		env: {
			account: props.account,
			region: "eu-central-1",
		},
		stage: props.stage,
		hostedZoneName: props.hostedZoneName,
		hostedZoneId: props.hostedZoneId,
		hostname: props.hostname,
		certificateId: props.certificateId,
		namespace: props.namespace,
		activeRegions: props.activeRegions,

		primaryEventBusStack: EVENT_BUS_EU_CENTRAL,
		secondaryEventBusStack: EVENT_BUS_EU_WEST,
	});

	const BACKBONE_EU_CENTRAL = new BackboneRegionalStack(
		scope,
		`${props.stackPrefix}BackboneEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}BackboneRegional`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const BACKBONE_EU_WEST = new BackboneRegionalStack(
		scope,
		`${props.stackPrefix}BackboneEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}BackboneRegional`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const CATALOGUE_SERVICE_DATA = new CatalogueServiceDataStack(
		scope,
		`${props.stackPrefix}CatalogueServiceData`,
		{
			stackName: `BattleBricks${props.stackPrefix}CatalogueServiceData`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const CATALOGUE_SERVICE_EU_CENTRAL = new CatalogueServiceRegionalStack(
		scope,
		`${props.stackPrefix}CatalogueServiceEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}CatalogueServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneRegionalStack: BACKBONE_EU_CENTRAL,
			dataStack: CATALOGUE_SERVICE_DATA,
		},
	);

	const CATALOGUE_SERVICE_EU_WEST = new CatalogueServiceRegionalStack(
		scope,
		`${props.stackPrefix}CatalogueServiceEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}CatalogueServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneRegionalStack: BACKBONE_EU_WEST,
			dataStack: CATALOGUE_SERVICE_DATA,
		},
	);

	const FILE_STAGING_SERVICE_DATA = new FileStagingServiceDataStack(
		scope,
		`${props.stackPrefix}FileStagingServiceData`,
		{
			stackName: `BattleBricks${props.stackPrefix}FileStagingServiceData`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const FILE_STAGING_SERVICE_EU_CENTRAL = new FileStagingServiceRegionalStack(
		scope,
		`${props.stackPrefix}FileStagingServiceEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}FileStagingServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneRegionalStack: BACKBONE_EU_CENTRAL,
			dataStack: FILE_STAGING_SERVICE_DATA,
		},
	);

	const FILE_STAGING_SERVICE_EU_WEST = new FileStagingServiceRegionalStack(
		scope,
		`${props.stackPrefix}FileStagingServiceEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}FileStagingServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneRegionalStack: BACKBONE_EU_WEST,
			dataStack: FILE_STAGING_SERVICE_DATA,
		},
	);

	const IDENTITY_SERVICE_DATA_EU_CENTRAL = new IdentityServiceDataRegionalStack(
		scope,
		`${props.stackPrefix}IdentityServiceDataEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}IdentityServiceDataRegional`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
		},
	);

	const IDENTITY_SERVICE_DATA_EU_WEST = new IdentityServiceDataRegionalStack(
		scope,
		`${props.stackPrefix}IdentityServiceDataEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}IdentityServiceDataRegional`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
		},
	);

	const IDENTITY_SERVICE_DATA = new IdentityServiceDataStack(
		scope,
		`${props.stackPrefix}IdentityServiceData`,
		{
			stackName: `BattleBricks${props.stackPrefix}IdentityServiceData`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
			regionalStacks: [
				IDENTITY_SERVICE_DATA_EU_CENTRAL,
				IDENTITY_SERVICE_DATA_EU_WEST,
			],
		},
	);

	const IDENTITY_SERVICE_EU_CENTRAL = new IdentityServiceRegionalStack(
		scope,
		`${props.stackPrefix}IdentityServiceEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}IdentityServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneRegionalStack: BACKBONE_EU_CENTRAL,
			dataStack: IDENTITY_SERVICE_DATA,
		},
	);

	const IDENTITY_SERVICE_EU_WEST = new IdentityServiceRegionalStack(
		scope,
		`${props.stackPrefix}IdentityServiceEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}IdentityServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneRegionalStack: BACKBONE_EU_WEST,
			dataStack: IDENTITY_SERVICE_DATA,
		},
	);

	const NOTIFICATION_SERVICE_DATA = new NotificationServiceDataStack(
		scope,
		`${props.stackPrefix}NotificationServiceData`,
		{
			stackName: `BattleBricks${props.stackPrefix}NotificationServiceData`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const NOTIFICATION_SERVICE_EU_CENTRAL = new NotificationServiceRegionalStack(
		scope,
		`${props.stackPrefix}NotificationServiceEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}NotificationServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneRegionalStack: BACKBONE_EU_CENTRAL,
			dataStack: NOTIFICATION_SERVICE_DATA,
		},
	);

	const NOTIFICATION_SERVICE_EU_WEST = new NotificationServiceRegionalStack(
		scope,
		`${props.stackPrefix}NotificationServiceEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}NotificationServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneRegionalStack: BACKBONE_EU_WEST,
			dataStack: NOTIFICATION_SERVICE_DATA,
		},
	);

	const POLICY_SERVICE_DATA = new PolicyServiceDataStack(
		scope,
		`${props.stackPrefix}PolicyServiceData`,
		{
			stackName: `BattleBricks${props.stackPrefix}PolicyServiceData`,
			terminationProtection: true,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneStack: BACKBONE,
		},
	);

	const POLICY_SERVICE_EU_CENTRAL = new PolicyServiceRegionalStack(
		scope,
		`${props.stackPrefix}PolicyServiceEuCentral1`,
		{
			stackName: `BattleBricks${props.stackPrefix}PolicyServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-central-1",
			},
			backboneRegionalStack: BACKBONE_EU_CENTRAL,
			dataStack: POLICY_SERVICE_DATA,
		},
	);

	const POLICY_SERVICE_EU_WEST = new PolicyServiceRegionalStack(
		scope,
		`${props.stackPrefix}PolicyServiceEuWest1`,
		{
			stackName: `BattleBricks${props.stackPrefix}PolicyServiceRegional`,
			crossRegionReferences: true,
			tags: props.tags,
			env: {
				account: props.account,
				region: "eu-west-1",
			},
			backboneRegionalStack: BACKBONE_EU_WEST,
			dataStack: POLICY_SERVICE_DATA,
		},
	);

	const FRONT_END = new FrontEndStack(scope, `${props.stackPrefix}FrontEnd`, {
		stackName: `BattleBricks${props.stackPrefix}FrontEnd`,
		crossRegionReferences: true,
		tags: props.tags,
		env: {
			account: props.account,
			region: "us-east-1",
		},
		backboneStack: BACKBONE,
	});

	const SOURCE_ARTIFACT = new Artifact();
	const BUILD_ARTIFACT = new Artifact();

	new Pipeline(scope, `BattleBricks${props.stackPrefix}Pipeline`, {
		pipelineName: props.pipelineName,
		rootDir: "apps/cdk/pipeline",
		terminationProtection: true,
		tags: props.tags,
		env: {
			account: props.account,
			region: "eu-central-1",
		},
		segments: [
			new CodeStarSourceSegment({
				output: SOURCE_ARTIFACT,
				connectionArn:
					"arn:aws:codeconnections:eu-central-1:670794226643:connection/c524f512-96f4-4ec5-b6a5-d7110e7b35f9",
				owner: "p-mercury",
				repository: "battle-bricks",
				branch: "main",
			}),
			new PipelineSegment({
				input: SOURCE_ARTIFACT,
				output: BUILD_ARTIFACT,
				project: {
					environment: {
						computeType: ComputeType.LARGE,
						buildImage: LinuxArmBuildImage.AMAZON_LINUX_2023_STANDARD_3_0,
						privileged: true,
						environmentVariables: props.environmentVariables,
					},
					buildSpec: BuildSpec.fromObject({
						Version: 0.2,
						phases: {
							install: {
								"runtime-versions": {
									nodejs: "24",
									golang: "1.26",
								},
								commands: [
									"corepack enable pnpm",
									"pnpm ci",
									"curl -L -o /usr/local/bin/opa https://openpolicyagent.org/downloads/latest/opa_linux_arm64_static",
									"chmod 755 /usr/local/bin/opa",
								],
							},
							build: {
								commands: props.buildCommands,
							},
						},
					}),
				},
			}),
			{
				stageName: "EventBusStacks",
				segments: [
					new StackSegment({
						stack: EVENT_BUS_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: EVENT_BUS_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
				],
			},
			new StackSegment({
				stack: BACKBONE,
				input: BUILD_ARTIFACT,
			}),
			{
				stageName: "ServiceDataRegionalStacks",
				segments: [
					new StackSegment({
						stack: IDENTITY_SERVICE_DATA_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: IDENTITY_SERVICE_DATA_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
				],
			},
			{
				stageName: "ServiceDataStacks",
				segments: [
					new StackSegment({
						stack: CATALOGUE_SERVICE_DATA,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: FILE_STAGING_SERVICE_DATA,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: IDENTITY_SERVICE_DATA,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: NOTIFICATION_SERVICE_DATA,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: POLICY_SERVICE_DATA,
						input: BUILD_ARTIFACT,
					}),
				],
			},
			new StackSegment({
				stack: BACKBONE_EU_CENTRAL,
				input: BUILD_ARTIFACT,
			}),
			{
				stageName: "ServiceStacksEuCentral1",
				segments: [
					new StackSegment({
						stack: CATALOGUE_SERVICE_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: FILE_STAGING_SERVICE_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: IDENTITY_SERVICE_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: NOTIFICATION_SERVICE_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: POLICY_SERVICE_EU_CENTRAL,
						input: BUILD_ARTIFACT,
					}),
				],
			},
			new StackSegment({
				stack: BACKBONE_EU_WEST,
				input: BUILD_ARTIFACT,
			}),
			{
				stageName: "ServiceStacksEuWest1",
				segments: [
					new StackSegment({
						stack: CATALOGUE_SERVICE_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: FILE_STAGING_SERVICE_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: IDENTITY_SERVICE_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: NOTIFICATION_SERVICE_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
					new StackSegment({
						stack: POLICY_SERVICE_EU_WEST,
						input: BUILD_ARTIFACT,
					}),
				],
			},
			new StackSegment({
				stack: FRONT_END,
				input: BUILD_ARTIFACT,
			}),
		],
	});

	Annotations.of(scope).acknowledgeWarning(
		"@aws-cdk/aws-lambda-go-alpha:goBuildFlagsSecurityWarning",
		"goBuildFlags are defined in-repo and reviewed; no untrusted input.",
	);
}
