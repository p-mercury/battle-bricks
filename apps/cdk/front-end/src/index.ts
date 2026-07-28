import { Duration, Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { ARecord, RecordTarget } from "aws-cdk-lib/aws-route53";
import { CloudFrontTarget } from "aws-cdk-lib/aws-route53-targets";
import { Runtime } from "aws-cdk-lib/aws-lambda";
import { HttpOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import {
	AllowedMethods,
	CachePolicy,
	FunctionCode,
	FunctionEventType,
	FunctionRuntime,
	OriginProtocolPolicy,
	OriginRequestPolicy,
	OriginSslPolicy,
	ViewerProtocolPolicy,
	Function,
} from "aws-cdk-lib/aws-cloudfront";
import { SvelteKitEdge } from "@battle-bricks/interface";
import { BackboneStack } from "@battle-bricks/backbone";

export interface FrontEndStackProps extends StackProps {
	readonly backboneStack: BackboneStack;
}

export class FrontEndStack extends Stack {
	constructor(scope: Construct, id: string, props: FrontEndStackProps) {
		super(scope, id, props);

		const interfaceMain = new SvelteKitEdge(this, "Interface", {
			domainNames: [props.backboneStack.hostname],
			certificate: props.backboneStack.certificate,
			runtime: Runtime.NODEJS_24_X,
			timeout: Duration.seconds(8),
			memorySize: 256,
			reservedConcurrentExecutions: 2,
			bundling: {
				esbuildArgs: { "--ignore-annotations": true },
				banner:
					'import { createRequire } from "module"; global.require = createRequire(import.meta.url);',
			},
		});

		interfaceMain.cloudFront.addBehavior(
			"/api/*",
			new HttpOrigin(`api.${props.backboneStack.hostname}`, {
				protocolPolicy: OriginProtocolPolicy.HTTPS_ONLY,
				originSslProtocols: [OriginSslPolicy.TLS_V1_2],
			}),
			{
				cachePolicy: CachePolicy.CACHING_DISABLED,
				allowedMethods: AllowedMethods.ALLOW_ALL,
				viewerProtocolPolicy: ViewerProtocolPolicy.HTTPS_ONLY,
				originRequestPolicy: OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
				functionAssociations: [
					{
						function: new Function(this, "InterfaceStripApiPathFunction", {
							runtime: FunctionRuntime.JS_2_0,
							code: FunctionCode.fromInline(`
								function handler(event) {
									var request = event.request;
									if (request.uri.startsWith("/api/")) {
										request.uri = request.uri.substring(4);
									}
									return request;
								}
							`),
						}),
						eventType: FunctionEventType.VIEWER_REQUEST,
					},
				],
			},
		);

		new ARecord(this, "ARecord", {
			recordName: `${props.backboneStack.hostname}.`,
			zone: props.backboneStack.hostedZone,
			target: RecordTarget.fromAlias(
				new CloudFrontTarget(interfaceMain.cloudFront),
			),
		});
	}
}
