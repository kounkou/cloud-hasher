import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as path from 'path';

export class CloudHasherStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create an IAM role in the target account
    const targetRole = new iam.Role(this, 'TargetRole', {
        assumedBy: new iam.AccountPrincipal('TARGET_ACCOUNT_ID'), // Replace with the target account ID
        // Other role configuration options
    });

    // Requests that failed to be processed will be sent to DLQ
    const dlq = new sqs.Queue(this, 'CloudHasherDLQ', {
      visibilityTimeout: cdk.Duration.seconds(300),
    });

    // Processing Lambda for the requests
    const hashRequestLambda = new lambda.Function(this, 'CloudHasherLambda', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'main',
      deadLetterQueue: dlq,
      deadLetterQueueEnabled: true,
      code: lambda.Code.fromAsset(path.join(__dirname, '../src/processorlambda/main.zip')),
    });

    // APIGateway
    const restApi = new apigateway.RestApi(this, 'CloudHasherRestAPI', {
      restApiName: 'cloudHasherRestAPI',
    });

    restApi.root.addMethod('GET', new apigateway.LambdaIntegration(hashRequestLambda, {
      contentHandling: apigateway.ContentHandling.CONVERT_TO_TEXT, // convert to base64
      credentialsPassthrough: false,
    }));
    restApi.root.addMethod('POST', new apigateway.LambdaIntegration(hashRequestLambda));
  }
}