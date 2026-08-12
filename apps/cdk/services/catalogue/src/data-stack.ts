import { Construct } from "constructs";
import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { BackboneStack } from "@battle-bricks/backbone";
import { AttributeType, Billing, TableV2 } from "aws-cdk-lib/aws-dynamodb";

export interface CatalogueServiceDataStackProps extends StackProps {
	readonly backboneStack: BackboneStack;
}

export class CatalogueServiceDataStack extends Stack {
	readonly table: TableV2;

	constructor(
		scope: Construct,
		id: string,
		props: CatalogueServiceDataStackProps,
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
	}
}
