import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { EventBus } from "aws-cdk-lib/aws-events";

export class EventBusStack extends Stack {
	readonly eventBus: EventBus;

	constructor(scope: Construct, id: string, props: StackProps) {
		super(scope, id, props);

		this.eventBus = new EventBus(this, "EventBus", {
			eventBusName: `${this.stackName}EventBus`,
		});
	}
}
