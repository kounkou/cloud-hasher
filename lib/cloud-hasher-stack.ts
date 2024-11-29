import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import { Function, Code, Runtime } from 'aws-cdk-lib/aws-lambda';
import { RestApi, LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';
import { join } from 'path';

export class CloudHasherStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    
    const dlq = new Queue(this, 'CloudHasherDLQ', {
      visibilityTimeout: Duration.seconds(300),
    });

    const hashRequestLambda = new Function(this, 'CloudHasherLambda', {
      runtime: Runtime.GO_1_X,
      handler: 'main',
      deadLetterQueue: dlq,
      deadLetterQueueEnabled: true,
      code: Code.fromAsset(join(__dirname, '../build/lambda.zip')),
    });

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
