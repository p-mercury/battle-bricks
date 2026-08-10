import { Construct } from "constructs";
import {
	AttributeType,
	Billing,
	ProjectionType,
	TableV2,
} from "aws-cdk-lib/aws-dynamodb";
import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { BackboneStack } from "@battle-bricks/backbone";

export interface PolicyServiceDataStackProps extends StackProps {
	readonly backboneStack: BackboneStack;
}

export class PolicyServiceDataStack extends Stack {
	readonly table: TableV2;

	constructor(
		scope: Construct,
		id: string,
		props: PolicyServiceDataStackProps,
	) {
		super(scope, id, props);

		this.table = new TableV2(this, "Table", {
			deletionProtection: true,
			removalPolicy: RemovalPolicy.RETAIN,
			partitionKey: { name: "pk", type: AttributeType.STRING },
			timeToLiveAttribute: "ttl",
			billing: Billing.onDemand(),
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
			projectionType: ProjectionType.KEYS_ONLY,
		});
	}
}
