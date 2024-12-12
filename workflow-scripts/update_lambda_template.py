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

with open(file_path, 'w') as file:
    yaml.dump(data, file)

print("Template file updated successfully.")