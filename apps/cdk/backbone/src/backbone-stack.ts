import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Certificate, ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import { IPublicHostedZone, PublicHostedZone } from "aws-cdk-lib/aws-route53";

import { CfnEndpoint, EventBus } from "aws-cdk-lib/aws-events";

export interface BackboneStackProps extends StackProps {
	readonly stage: "DEVELOPMENT" | "PRODUCTION";
	readonly hostedZoneName: string;
	readonly certificateId: string;
	readonly hostedZoneId: string;
	readonly namespace: string;
	readonly hostname: string;
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

	constructor(scope: Construct, id: string, props: BackboneStackProps) {
		const { stage, hostname, namespace, ...stackProps } = props;

		super(scope, id, stackProps);

		this.stage = stage;
		this.hostname = props.hostname;
		this.namespace = props.namespace;

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
