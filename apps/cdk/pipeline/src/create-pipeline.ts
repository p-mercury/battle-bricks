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
import { BackboneStack } from "@battle-bricks/backbone";
import { FrontEndStack } from "@battle-bricks/front-end";

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

	environmentVariables?: Record<string, any>;
};

export function createPipeline(scope: App, props: CreatePipelineProps) {
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
	});

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
						version: 0.2,
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
