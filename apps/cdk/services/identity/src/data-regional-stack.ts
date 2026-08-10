import { Construct } from "constructs";
import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { Bucket } from "aws-cdk-lib/aws-s3";

export interface IdentityServiceDataRegionalStackProps extends StackProps {}

export class IdentityServiceDataRegionalStack extends Stack {
	readonly bucket: Bucket;

	constructor(
		scope: Construct,
		id: string,
		props: IdentityServiceDataRegionalStackProps,
	) {
		super(scope, id, props);

		this.bucket = new Bucket(this, "Bucket", {
			removalPolicy: RemovalPolicy.RETAIN,
			versioned: true,
		});
	}
}
