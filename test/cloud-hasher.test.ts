import { App } from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import { CloudHasherStack } from '../lib/cloud-hasher-stack';

describe('CloudHasherStack', () => {
    let stack: CloudHasherStack;

    beforeEach(() => {
      stack = new CloudHasherStack(new App(), 'CloudHasherStack');
    });

    test('SQS Queue Created', () => {
        const template = Template.fromStack(stack);

        template.hasResourceProperties('AWS::SQS::Queue', {
            VisibilityTimeout: 300
        });
    });

    test('Lambda Created', () => {
        const template = Template.fromStack(stack);

        template.hasResourceProperties('AWS::Lambda::Function', {
            Runtime: "go1.x",
            Handler: "main",
        });
    });

    test('APIGateway Created', () => {
        const template = Template.fromStack(stack);

        template.hasResourceProperties('AWS::ApiGateway::RestApi', {
            Name: 'cloudHasherRestAPI'
        });
    });
});