import { Duration, Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import {
	Certificate,
	CertificateValidation,
	ICertificate,
} from "aws-cdk-lib/aws-certificatemanager";
import {
	AaaaRecord,
	ARecord,
	IPublicHostedZone,
	RecordTarget,
} from "aws-cdk-lib/aws-route53";
import { HttpNamespace } from "aws-cdk-lib/aws-servicediscovery";
import { CfnEndpoint, EventBus } from "aws-cdk-lib/aws-events";

import { BackboneStack } from "./backbone-stack.js";
import {
	CorsHttpMethod,
	DomainName,
	HttpApi,
	IHttpRouteAuthorizer,
} from "aws-cdk-lib/aws-apigatewayv2";
import { ApiGatewayv2DomainProperties } from "aws-cdk-lib/aws-route53-targets";
import {
	HttpLambdaAuthorizer,
	HttpLambdaResponseType,
} from "aws-cdk-lib/aws-apigatewayv2-authorizers";
import { ApiAuthorizer } from "./api-authorizer/index.js";

export interface BackboneRegionalStackProps extends StackProps {
	readonly backboneStack: BackboneStack;
}

export class BackboneRegionalStack extends Stack {
	readonly stage: "DEVELOPMENT" | "PRODUCTION";
	readonly certificate: ICertificate;
	readonly hostname: string;
	readonly namespace: HttpNamespace;
	readonly hostedZone: IPublicHostedZone;
	readonly eventBusGlobalEndpoint: CfnEndpoint;
	readonly primaryEventBus: EventBus;
	readonly secondaryEventBus: EventBus;
	readonly activeRegions: string[];
	readonly apiGateway: HttpApi;
	readonly apiGatewayAuthorizer: IHttpRouteAuthorizer;

	constructor(scope: Construct, id: string, props: BackboneRegionalStackProps) {
		super(scope, id, { ...props, crossRegionReferences: true });

		this.stage = props.backboneStack.stage;
		this.certificate = props.backboneStack.certificate;
		this.hostname = props.backboneStack.hostname;
		this.hostedZone = props.backboneStack.hostedZone;
		this.eventBusGlobalEndpoint = props.backboneStack.eventBusGlobalEndpoint;
		this.primaryEventBus = props.backboneStack.primaryEventBus;
		this.secondaryEventBus = props.backboneStack.secondaryEventBus;
		this.activeRegions = props.backboneStack.activeRegions;

		this.namespace = new HttpNamespace(this, "Namespace", {
			name: props.backboneStack.namespace,
		});

		const domainName = new DomainName(this, "HttpApiGatewayCustomDomain", {
			domainName: `api.${props.backboneStack.hostname}`,
			certificate: new Certificate(
				this,
				"HttpApiGatewayCustomDomainCertificate",
				{
					domainName: props.backboneStack.hostname,
					validation: CertificateValidation.fromDns(
						props.backboneStack.hostedZone,
					),
					subjectAlternativeNames: [`*.${props.backboneStack.hostname}`],
				},
			),
		});

		this.apiGateway = new HttpApi(this, "HttpApiGateway", {
			corsPreflight: {
				allowMethods: [CorsHttpMethod.GET, CorsHttpMethod.POST],
				allowCredentials: true,
				allowHeaders: [
					"content-type",
					"connect-protocol-version",
					"connect-timeout-ms",
					"grpc-timeout",
					"x-grpc-web",
					"x-user-agent",
					"authorization",
				],
				exposeHeaders: [
					"grpc-status",
					"grpc-message",
					"grpc-status-details-bin",
				],
				allowOrigins: [`https://${props.backboneStack.hostname}`],
				maxAge: Duration.hours(1),
			},
			createDefaultStage: true,
			defaultDomainMapping: { domainName },
		});

		this.apiGatewayAuthorizer = new HttpLambdaAuthorizer(
			"HttpApiGatewayAuthorizer",
			new ApiAuthorizer(this, "HttpApiGatewayLambdaAuthorizer", {
				apiGateway: this.apiGateway,
				reservedConcurrentExecutions: 10,
			}),
			{
				identitySource: ["$request.header.Cookie"],
				responseTypes: [HttpLambdaResponseType.SIMPLE],
				resultsCacheTtl: Duration.seconds(0),
			},
		);

		new ARecord(this, "HttpApiGatewayARecord", {
			recordName: `api.${props.backboneStack.hostname}.`,
			zone: props.backboneStack.hostedZone,
			region: this.region,
			setIdentifier: `${this.region}-a`,
			target: RecordTarget.fromAlias(
				new ApiGatewayv2DomainProperties(
					domainName.regionalDomainName,
					domainName.regionalHostedZoneId,
				),
			),
		});

		new AaaaRecord(this, "HttpApiGatewayAaaaRecord", {
			recordName: `api.${props.backboneStack.hostname}.`,
			zone: props.backboneStack.hostedZone,
			region: this.region,
			setIdentifier: `${this.region}-aaaa`,
			target: RecordTarget.fromAlias(
				new ApiGatewayv2DomainProperties(
					domainName.regionalDomainName,
					domainName.regionalHostedZoneId,
				),
			),
		});
	}
}
