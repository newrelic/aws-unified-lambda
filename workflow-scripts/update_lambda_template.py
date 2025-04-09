from ruamel.yaml import YAML

yaml = YAML()

file_path = 'lambda-template.yaml'

with open(file_path, 'r') as file:
    data = yaml.load(file)

resources = data.get('Resources', {})
nr_log_forwarder = resources.get('NewRelicLogsServerlessLogForwarder', {})
properties = nr_log_forwarder.get('Properties', {})

properties['CodeUri'] = {
    'Bucket': yaml.load("!FindInMap [ RegionToS3Bucket, !Ref 'AWS::Region', BucketArn ]"),
    'Key': "new-relic-log-forwarder-folder/new-relic-log-forwarder.zip"
}

mappings = {
    'RegionToS3Bucket': {
        'us-east-1': {'BucketArn': 'unified-logging-lambda-code-us-east-1'},
        'us-east-2': {'BucketArn': 'unified-logging-lambda-code-us-east-2'},
        'eu-west-1': {'BucketArn': 'unified-logging-lambda-code-eu-west-1'},
        'eu-west-2': {'BucketArn': 'unified-logging-lambda-code-eu-west-2'},
        'us-west-1': {'BucketArn': 'unified-logging-lambda-code-us-west-1'},
        'us-west-2': {'BucketArn': 'unified-logging-lambda-code-us-west-2'},
        'af-south-1': {'BucketArn': 'unified-logging-lambda-code-af-south-1'},
        'ap-south-1': {'BucketArn': 'unified-logging-lambda-code-ap-south-1'},
        'ap-northeast-3': {'BucketArn': 'unified-logging-lambda-code-ap-northeast-3'},
        'ap-northeast-2': {'BucketArn': 'unified-logging-lambda-code-ap-northeast-2'},
        'ap-southeast-1': {'BucketArn': 'unified-logging-lambda-code-ap-southeast-1'},
        'ap-southeast-2': {'BucketArn': 'unified-logging-lambda-code-ap-southeast-2'},
        'ap-northeast-1': {'BucketArn': 'unified-logging-lambda-code-ap-northeast-1'},
        'ca-central-1': {'BucketArn': 'unified-logging-lambda-code-ca-central-1'},
        'ca-west-1': {'BucketArn': 'unified-logging-lambda-code-ca-west-1'},
        'eu-central-1': {'BucketArn': 'unified-logging-lambda-code-eu-central-1'},
        'eu-south-1': {'BucketArn': 'unified-logging-lambda-code-eu-south-1'},
        'eu-west-3': {'BucketArn': 'unified-logging-lambda-code-eu-west-3'},
        'eu-south-2': {'BucketArn': 'unified-logging-lambda-code-eu-south-2'},
        'eu-north-1': {'BucketArn': 'unified-logging-lambda-code-eu-north-1'},
        'eu-central-2': {'BucketArn': 'unified-logging-lambda-code-eu-central-2'},
        'me-south-1': {'BucketArn': 'unified-logging-lambda-code-me-south-1'},
        'me-central-1': {'BucketArn': 'unified-logging-lambda-code-me-central-1'},
        'ap-east-1': {'BucketArn': 'unified-logging-lambda-code-ap-east-1'},
        'ap-south-2': {'BucketArn': 'unified-logging-lambda-code-ap-south-2'},
        'ap-southeast-3': {'BucketArn': 'unified-logging-lambda-code-ap-southeast-3'},
        'ap-southeast-5': {'BucketArn': 'unified-logging-lambda-code-ap-southeast-5'},
        'ap-southeast-4': {'BucketArn': 'unified-logging-lambda-code-ap-southeast-4'},
        'il-central-1': {'BucketArn': 'unified-logging-lambda-code-il-central-1'},
        'sa-east-1': {'BucketArn': 'unified-logging-lambda-code-sa-east-1'},
    }
}

data['Mappings'] = mappings

with open(file_path, 'w') as file:
    yaml.dump(data, file)

print("Template file updated successfully.")