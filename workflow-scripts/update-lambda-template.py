import yaml

file_path = 'lambda-template.yaml'

# Modify lambda-template with production s3 bucket arn
with open(file_path, 'r') as stream:
    data = yaml.safe_load(stream)

data['Resources']['NewRelicLogsServerlessLogForwarder']['Properties']['CodeUri'] = {
    'Bucket': '!FindInMap [ RegionToS3Bucket, !Ref \'AWS::Region\', BucketArn ]',
    'Key': 'new-relic-log-forwarder-folder/new-relic-log-forwarder.zip'
}

# dump changes to lambda-template.yaml
with open(file_path, 'w') as outfile:
    yaml.dump(data, outfile, default_flow_style=False)

print("Template file updated successfully.")