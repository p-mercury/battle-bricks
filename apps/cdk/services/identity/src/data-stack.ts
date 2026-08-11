import { Construct } from "constructs";
import {
	AttributeType,
	Billing,
	StreamViewType,
	TableV2,
} from "aws-cdk-lib/aws-dynamodb";
import { Duration, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { BackboneStack } from "@battle-bricks/backbone";
import { BucketMesh } from "@flit/cdk-bucket-mesh";
import {
	ClientAttributes,
	FeaturePlan,
	ManagedLoginVersion,
	OAuthScope,
	StringAttribute,
	UserPool,
	UserPoolClient,
	UserPoolClientIdentityProvider,
	UserPoolDomain,
	UserPoolOperation,
} from "aws-cdk-lib/aws-cognito";
import { ARecord, RecordTarget } from "aws-cdk-lib/aws-route53";
import { UserPoolDomainTarget } from "aws-cdk-lib/aws-route53-targets";

import { IdentityServiceDataRegionalStack } from "./data-regional-stack.js";
import { PostConfirmationHandler } from "./congito-triggers/post-confirmation/index.js";

export interface IdentityServiceDataStackProps extends StackProps {
	readonly backboneStack: BackboneStack;
	readonly regionalStacks: IdentityServiceDataRegionalStack[];
}

export class IdentityServiceDataStack extends Stack {
	readonly table: TableV2;
	readonly hostname: string;
	readonly userPool: UserPool;
	readonly userPoolDomain: UserPoolDomain;
	readonly userPoolClient: UserPoolClient;

	constructor(
		scope: Construct,
		id: string,
		props: IdentityServiceDataStackProps,
	) {
		super(scope, id, props);

		this.hostname = `${props.backboneStack.hostname}`;

		this.table = new TableV2(this, "Table", {
			deletionProtection: true,
			removalPolicy: RemovalPolicy.RETAIN,
			partitionKey: { name: "pk", type: AttributeType.STRING },
			sortKey: { name: "sk", type: AttributeType.STRING },
			timeToLiveAttribute: "ttl",
			billing: Billing.onDemand(),
			dynamoStream: StreamViewType.NEW_AND_OLD_IMAGES,
			replicas: props.backboneStack.activeRegions
				.filter((r) => r !== this.region)
				.map((r) => ({ region: r })),
			pointInTimeRecoverySpecification: {
				pointInTimeRecoveryEnabled: true,
			},
		});

		this.table.addGlobalSecondaryIndex({
			indexName: "gsi1",
			partitionKey: {
				name: "gsi1pk",
				type: AttributeType.STRING,
			},
			sortKey: {
				name: "gsi1sk",
				type: AttributeType.STRING,
			},
		});

		new BucketMesh(this, "BucketMesh", {
			buckets: props.regionalStacks.map(({ bucket }) => bucket),
		});

		this.userPool = new UserPool(this, "UserPool", {
			deletionProtection: true,
			removalPolicy: RemovalPolicy.RETAIN,
			signInCaseSensitive: false,
			featurePlan: FeaturePlan.PLUS,
			selfSignUpEnabled: true,
			signInAliases: { email: true },
			autoVerify: { email: true },
			keepOriginal: { email: true },
			customAttributes: {
				userId: new StringAttribute({ mutable: true, minLen: 12, maxLen: 12 }),
			},
			standardAttributes: {
				fullname: { required: true, mutable: true },
			},
		});

		this.userPool.addTrigger(
			UserPoolOperation.POST_CONFIRMATION,
			new PostConfirmationHandler(this, "PostConfirmationHandler", {
				backboneStack: props.backboneStack,
				table: this.table,
			}),
		);

		this.userPoolDomain = this.userPool.addDomain("CustomDomain", {
			managedLoginVersion: ManagedLoginVersion.NEWER_MANAGED_LOGIN,
			customDomain: {
				domainName: `auth.${props.backboneStack.hostname}`,
				certificate: props.backboneStack.certificate,
			},
		});

		new ARecord(this, "AuthARecord", {
			recordName: `${this.userPoolDomain.domainName}.`,
			zone: props.backboneStack.hostedZone,
			target: RecordTarget.fromAlias(
				new UserPoolDomainTarget(this.userPoolDomain),
			),
		});

		this.userPoolClient = this.userPool.addClient("UserPoolClient", {
			generateSecret: true,
			preventUserExistenceErrors: true,
			refreshTokenValidity: Duration.days(8),
			accessTokenValidity: Duration.hours(2),
			readAttributes: new ClientAttributes()
				.withStandardAttributes({ email: true })
				.withCustomAttributes("userId"),
			writeAttributes: new ClientAttributes(),
			enableTokenRevocation: true,
			enablePropagateAdditionalUserContextData: true,
			refreshTokenRotationGracePeriod: Duration.seconds(0),
			oAuth: {
				flows: {
					authorizationCodeGrant: true,
				},
				scopes: [OAuthScope.OPENID],
				callbackUrls: [
					"https://localhost:5173/auth/callback",
					`https://${props.backboneStack.hostname}/auth/callback`,
				],
				logoutUrls: [
					"https://localhost:5173",
					`https://${props.backboneStack.hostname}`,
				],
			},
			supportedIdentityProviders: [UserPoolClientIdentityProvider.COGNITO],
		});
	}
}
