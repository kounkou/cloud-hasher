import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import { Function, Code, Runtime } from 'aws-cdk-lib/aws-lambda';
import { RestApi, LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';
import { join } from 'path';

export class CloudHasherStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Requests that failed to be processed will be sent to DLQ
    const dlq = new Queue(this, 'CloudHasherDLQ', {
      visibilityTimeout: Duration.seconds(300),
    });

    // Processing Lambda for the requests
    const hashRequestLambda = new Function(this, 'CloudHasherLambda', {
      runtime: Runtime.GO_1_X,
      handler: 'main',
      deadLetterQueue: dlq,
      deadLetterQueueEnabled: true,
      code: Code.fromAsset(join(__dirname, '../src/processorlambda/main.zip')),
    });

    // APIGateway
    const restApi = new RestApi(this, 'CloudHasherRestAPI', {
      restApiName: 'cloudHasherRestAPI',
    });

    restApi.root.addMethod(
        'POST', 
        new LambdaIntegration(
            hashRequestLambda, {
                proxy: false,
                requestTemplates: {
                    'application/json': '$input.json("$")',
                },
            },
        ),
    );
  }
}