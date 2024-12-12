import yaml

def unknown_constructor(loader, tag_suffix, node):
    if isinstance(node, yaml.ScalarNode):
        return loader.construct_scalar(node)
    elif isinstance(node, yaml.SequenceNode):
        return loader.construct_sequence(node)
    elif isinstance(node, yaml.MappingNode):
        return loader.construct_mapping(node)
    return None

file_path = 'lambda-template.yaml'

with open(file_path, 'r') as stream:
    try:
        yaml.add_multi_constructor('!', unknown_constructor, Loader=yaml.FullLoader)
        data = yaml.load(stream, Loader=yaml.FullLoader)
    except yaml.YAMLError as exc:
        print(exc)

code_uri_path = data.get('Resources', {}).get('NewRelicLogsServerlessLogForwarder', {}).get('Properties', {})
code_uri_path['CodeUri'] = {
    'Bucket': "!FindInMap [ RegionToS3Bucket, !Ref 'AWS::Region', BucketArn ]",
    'Key': 'new-relic-log-forwarder-folder/new-relic-log-forwarder.zip'
}

with open(file_path, 'w') as file:
    yaml.dump(data, file, default_flow_style=False)

print("Template file updated successfully.")