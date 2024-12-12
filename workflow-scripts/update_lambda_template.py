import re

file_path = 'lambda-template.yaml'

with open(file_path, 'r') as file:
    content = file.read()

# CodeUri: src/ pattern is identified and replaced with production s3 bucket arn
pattern = re.compile(
    r'CodeUri:\s*src/\s*'
)

replacement = (
    "CodeUri:\n"
    "    Bucket: !FindInMap [ RegionToS3Bucket, !Ref 'AWS::Region', BucketArn ]\n"
    "    Key: 'new-relic-log-forwarder-folder/new-relic-log-forwarder.zip'\n"
)

new_content = pattern.sub(replacement, content)

with open(file_path, 'w') as file:
    file.write(new_content)

print("Template file updated successfully.")