import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import { Function, Code, Runtime } from 'aws-cdk-lib/aws-lambda';
import { RestApi, MethodLoggingLevel, LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';
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
      memorySize: 512,
      timeout: Duration.seconds(10),
      deadLetterQueue: dlq,
      deadLetterQueueEnabled: true,
      code: Code.fromAsset(join(__dirname, '../src/processorlambda'), {
        bundling: {
          image: Runtime.GO_1_X.bundlingImage,
          command: [
            'bash', '-c',
            [
              'mkdir -p /tmp/go-cache /tmp/go-mod',
              'cd /asset-input',
              'ls -al',
              'GOCACHE=/tmp/go-cache GOMODCACHE=/tmp/go-mod GOOS=linux GOARCH=amd64 go build -o /asset-output/main main.go'
            ].join(' && ')
          ],
        },
      }),
    });

    const restApi = new RestApi(this, 'CloudHasherRestAPI', {
      restApiName: 'cloudHasherRestAPI',
    });

    restApi.root.addMethod(
        'POST',
        new LambdaIntegration(
            hashRequestLambda, {
                proxy: true,
            },
        ),
    );
  }
}
