from ruamel.yaml import YAML

yaml = YAML()

file_path = 'lambda-template.yaml'

# Load the existing YAML content
with open(file_path, 'r') as file:
    data = yaml.load(file)

# Navigate to the desired structure and perform the update
resources = data.get('Resources', {})
nr_log_forwarder = resources.get('NewRelicLogsServerlessLogForwarder', {})
properties = nr_log_forwarder.get('Properties', {})

# Update the CodeUri with the exact desired structure
properties['CodeUri'] = {
    'Bucket': yaml.load("!FindInMap [ RegionToS3Bucket, !Ref 'AWS::Region', BucketArn ]"),
    'Key': "new-relic-log-forwarder-folder/new-relic-log-forwarder.zip"
}

# Write the updated data back to the YAML file
with open(file_path, 'w') as file:
    yaml.dump(data, file)

print("Template file updated successfully.")