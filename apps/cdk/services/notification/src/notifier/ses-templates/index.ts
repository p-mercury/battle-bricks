import { Construct } from "constructs";
import { CfnTemplate } from "aws-cdk-lib/aws-ses";
import { crush } from "html-crush";
import { fileURLToPath } from "url";
import { readFileSync } from "fs";
import { Duration, RemovalPolicy } from "aws-cdk-lib";
import {
	BlockPublicAccess,
	Bucket,
	BucketAccessControl,
	HttpMethods,
	ObjectOwnership,
} from "aws-cdk-lib/aws-s3";
import {
	BucketDeployment,
	CacheControl,
	Source,
} from "aws-cdk-lib/aws-s3-deployment";

export class SesTemplates extends Construct {
	readonly info;
	readonly singleAction;
	readonly purchasingOrder;

	constructor(scope: Construct, id: string) {
		super(scope, id);

		const bucket = new Bucket(this, "AssetBucket", {
			removalPolicy: RemovalPolicy.DESTROY,
			autoDeleteObjects: true,
			websiteIndexDocument: "index.html",
			publicReadAccess: true,
			objectOwnership: ObjectOwnership.OBJECT_WRITER,
			blockPublicAccess: new BlockPublicAccess({
				blockPublicAcls: false,
				blockPublicPolicy: false,
				ignorePublicAcls: false,
				restrictPublicBuckets: false,
			}),
			cors: [
				{
					allowedMethods: [HttpMethods.GET, HttpMethods.HEAD],
					allowedOrigins: ["*"],
				},
			],
		});

		new BucketDeployment(this, "TemplateBucketDeployment", {
			destinationBucket: bucket,
			accessControl: BucketAccessControl.PUBLIC_READ,
			sources: [
				Source.asset(
					fileURLToPath(new URL("./assets", import.meta.url).href).replace(
						"/dist/",
						"/src/",
					),
				),
			],
			cacheControl: [
				CacheControl.setPublic(),
				CacheControl.maxAge(Duration.days(2)),
				CacheControl.sMaxAge(Duration.days(2)),
				CacheControl.fromString("immutable"),
			],
		});

		this.info = new CfnTemplate(scope, "Info", {
			template: {
				subjectPart: "{{subject}}",
				htmlPart: crush(
					readFileSync(
						fileURLToPath(new URL("./info.html", import.meta.url).href).replace(
							"/dist/",
							"/src/",
						),
					)
						.toString()
						.split("assets/")
						.join(`${bucket.bucketWebsiteUrl}/`),
					{
						removeHTMLComments: true,
						removeCSSComments: true,
						removeIndentations: true,
						removeLineBreaks: true,
					},
				).result,
			},
		});

		this.singleAction = new CfnTemplate(scope, "SingleActionTemplate", {
			template: {
				subjectPart: "{{subject}}",
				htmlPart: crush(
					readFileSync(
						fileURLToPath(
							new URL("./single-action.html", import.meta.url).href,
						).replace("/dist/", "/src/"),
					)
						.toString()
						.split("assets/")
						.join(`${bucket.bucketWebsiteUrl}/`),
					{
						removeHTMLComments: true,
						removeCSSComments: true,
						removeIndentations: true,
						removeLineBreaks: true,
					},
				).result,
			},
		});

		this.purchasingOrder = new CfnTemplate(scope, "PurchasingOrderTemplate", {
			template: {
				subjectPart: "{{subject}}",
				htmlPart: crush(
					readFileSync(
						fileURLToPath(
							new URL("./purchasing-order.html", import.meta.url).href,
						).replace("/dist/", "/src/"),
					)
						.toString()
						.split("assets/")
						.join(`${bucket.bucketWebsiteUrl}/`),
					{
						removeHTMLComments: true,
						removeCSSComments: true,
						removeIndentations: true,
						removeLineBreaks: true,
					},
				).result,
			},
		});
	}
}
