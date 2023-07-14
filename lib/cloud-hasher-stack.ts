import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as path from 'path';

export class CloudHasherStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // processing Lambda for the requests
    const hashRequestLambda = new lambda.Function(this, 'hashLambda', {
        runtime: lambda.Runtime.GO_1_X,
        handler: 'main',
        code: lambda.Code.fromAsset(path.join(__dirname, '../src/processorlambda/handler.zip')),
    });

    // APIGateway
    const api = new apigateway.RestApi(this, 'apigw', {
        description: 'apigateway',
        deployOptions: {
            stageName: 'dev',
        },

        // enable CORS
        defaultCorsPreflightOptions: {
            allowHeaders: [
                'Content-Type',
                'X-Amz-Date',
                'Authorization',
                'X-Api-Key',
            ],
            allowMethods: ['OPTIONS', 'POST'],
            allowCredentials: true,
            allowOrigins: ['http://localhost:3000'],
        },
    });

    api.root.addMethod("POST", new apigateway.LambdaIntegration(hashRequestLambda, {
        contentHandling: apigateway.ContentHandling.CONVERT_TO_TEXT,
        credentialsPassThrough: true,
    }));

    // Requests that failed to be processed will be sent to DLQ
    const queue = new sqs.Queue(this, 'CloudHasherDLQ', {
        visibilityTimeout: cdk.Duration.seconds(300)
    });


  }
}
