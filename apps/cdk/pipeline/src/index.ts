import { App } from "aws-cdk-lib";
import { createPipeline } from "./create-pipeline.js";

const APP = new App();

createPipeline(APP, {
	pipelineName: "battle-bricks",
	buildCommands: ["pnpm run forcebuild"],

	stage: "PRODUCTION",

	account: "670794226643",
	stackPrefix: "Production",
	tags: {
		Environment: "PRODUCTION",
	},

	hostedZoneName: "battlebricks.games",
	hostedZoneId: "Z074846013QDC2OD6B2LQ",
	hostname: "battlebricks.games",
	certificateId: "8b41fce0-bfe3-451d-ac2e-18f57aa243ac",
	namespace: "production.battlebricks",
	activeRegions: ["eu-central-1", "eu-west-1"],
});
