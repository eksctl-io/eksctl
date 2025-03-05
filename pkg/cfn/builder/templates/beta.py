import copy
import json
import boto3
import os
import logging
import time
import cfnresponse  # AWS-provided helper for sending responses

# Set up logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)


def init_eks_client():
    """
    Initializes and returns an EKS client using the endpoint URL from the environment variable.
    """
    eks_endpoint = os.environ.get('AWS_ENDPOINT_URL_EKS')
    if not eks_endpoint:
        raise ValueError("AWS_ENDPOINT_URL_EKS environment variable is not set or is empty")
    return boto3.client('eks', endpoint_url=eks_endpoint)


def validate_input(event):
    """
    Validate that all required fields are present in the event.
    """
    required_fields = ['Name', 'RoleArn', 'ResourcesVpcConfig',
                       'IAMPrincipalArn', 'STSRoleArn']
    for field in required_fields:
        if field not in event['ResourceProperties']:
            raise ValueError(f"Missing required field: {field}")


def convert_keys_to_lowercase_first_letter(d):
    """
    Convert the first character of dictionary keys to lowercase.
    """
    if not isinstance(d, dict):
        return d  # Return as-is if it's not a dictionary

    new_dict = {}
    for key, value in d.items():
        # Convert the first character of the key to lowercase
        new_key = key[:1].lower() + key[1:] if key else key

        # Recursively process nested dictionaries
        if isinstance(value, dict):
            new_dict[new_key] = convert_keys_to_lowercase_first_letter(value)
        # Process lists only if they contain dictionaries
        elif isinstance(value, list):
            new_dict[new_key] = [
                convert_keys_to_lowercase_first_letter(item) if isinstance(item, dict) else item
                for item in value
            ]
        else:
            new_dict[new_key] = value
    return new_dict


def replace_boolean_strings(d):
    """
    Replace string representations of booleans with actual booleans.
    """
    if isinstance(d, (dict, list)):
        iterable = d.items() if isinstance(d, dict) else enumerate(d)
        for k, v in iterable:
            if isinstance(v, str) and v.lower() in {"true", "false"}:
                d[k] = v.lower() == "true"
            elif isinstance(v, (dict, list)):
                replace_boolean_strings(v)


def replace_integer_strings(d):
    """
    Replace string representations of integers with actual integers.
    """
    if isinstance(d, (dict, list)):
        iterable = d.items() if isinstance(d, dict) else enumerate(d)
        for k, v in iterable:
            if isinstance(v, str) and v.isdigit():
                d[k] = int(v)
            elif isinstance(v, (dict, list)):
                replace_integer_strings(v)


def create_access_entry(eks_client, principal_arn, username, cluster_name, entry_type):
    """
    Create an access entry for an IAM principal in an EKS cluster.
    """
    logger.info(f"Creating access entry in EKS cluster: {cluster_name}")
    params = {
        'clusterName': cluster_name,
        'principalArn': principal_arn,
    }
    if username is not None:
        params['username'] = username
    if entry_type is not None:
        params['type'] = entry_type

    response1 = eks_client.create_access_entry(**params)
    logger.info("Access entry called successfully:")
    logger.info("Access entry response: " + json.dumps(response1, default=str))

    if entry_type is not None and entry_type == "STANDARD":
        # Associate the admin access policy
        policy_arn = 'arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy'
        params = {
            'clusterName': cluster_name,
            'principalArn': principal_arn,
            'policyArn': policy_arn,
            'accessScope': {
                'type': 'cluster',
                'namespaces': []
            }
        }
        response2 = eks_client.associate_access_policy(**params)
        logger.info("Associate Access Policy called successfully:")
        logger.info("Associate Access Policy response: " + json.dumps(response2, default=str))
    return response1


def get_stack_tags(event):
    """
    Extracts and returns the tags of a CloudFormation stack from the given event.
    """
    stack_name = event['StackId'].split('/')[1]
    cf = boto3.client('cloudformation')
    response = cf.describe_stacks(StackName=stack_name)
    stack_tags = response['Stacks'][0].get('Tags', [])
    tag_dict = {tag['Key']: tag['Value'] for tag in stack_tags}
    return tag_dict


def update_payload_tags(payload, event):
    """
    Updates the payload with tags from the CloudFormation stack associated with the given event.
    """
    # get the stack level tags
    stack_tags = get_stack_tags(event)
    logger.info("Stack Tags: " + json.dumps(stack_tags, default=str))
    # Add stack tags to the create_cluster_payload tags
    if 'tags' not in payload:
        payload['tags'] = {}
    logger.info("Cluster Tags: " + json.dumps(payload['tags'], default=str))
    # Ensure 'tags' is a dictionary
    if isinstance(payload['tags'], list):
        tags_dict = {item['key']: item['value'] for item in payload['tags']}
        payload['tags'] = tags_dict
    payload['tags'].update(stack_tags)
    logger.info("Final Tags: " + json.dumps(payload['tags'], default=str))


def handler(event, context):
    """
    Lambda function handler for CloudFormation custom resource.
    """
    logger.info("Received event: " + json.dumps(event, default=str))

    # Validate that the invocation is from CloudFormation
    if 'RequestType' not in event:
        raise ValueError("Invalid invocation source. This Lambda function can only be invoked by CloudFormation.")

    if event['ResourceType'] == 'Custom::EksCluster':
        return cluster_handler(event, context)

    if event['ResourceType'] == 'Custom::EksManagedNodeGroup':
        return nodegroup_handler(event, context)

    if event['ResourceType'] == 'Custom::EksAccessEntry':
        return access_entry_handler(event, context)

    raise ValueError(f"Invalid resource type {event['ResourceType']}")


def cluster_handler(event, context):
    """
    Handles the creation, update, and deletion of an EKS cluster as a CloudFormation custom resource.
    """
    try:
        eks_client = init_eks_client()

        cluster_name = event['ResourceProperties']['Name']

        # Handle Delete event
        if event['RequestType'] == 'Delete':
            delete_cluster(eks_client, cluster_name)
            cfnresponse.send(event, context, cfnresponse.SUCCESS, {"Message": "Resource deleted"})
            return {
                'PhysicalResourceId': event['PhysicalResourceId']
            }

        # Validate input
        validate_input(event)

        iam_principal_arn = event['ResourceProperties']['IAMPrincipalArn']
        sts_role_arn = event['ResourceProperties']['STSRoleArn']

        # Prepare the create_cluster request payload from the custom resource properties
        create_cluster_payload = convert_keys_to_lowercase_first_letter(
            copy.deepcopy(event['ResourceProperties']))

        # cleanup create_cluster_payload as not all fields can be sent to EKS create_cluster API
        del create_cluster_payload['serviceToken']
        del create_cluster_payload['iAMPrincipalArn']
        del create_cluster_payload['sTSRoleArn']
        replace_boolean_strings(create_cluster_payload)

        update_payload_tags(create_cluster_payload, event)

        # create and wait for the eks cluster
        cluster_details, response = create_cluster(eks_client, cluster_name, create_cluster_payload)

        # Create an access entry for the EKS cluster
        create_access_entry(eks_client, iam_principal_arn, sts_role_arn, cluster_name, "STANDARD")

        # Extract required attributes
        eventData = {
            "Arn": cluster_details['arn'],
            "PhysicalResourceId": cluster_details['arn'],
            "ClusterName": cluster_details['name'],
            "ClusterSecurityGroupId": cluster_details['resourcesVpcConfig']['clusterSecurityGroupId'],
            "CertificateAuthorityData": cluster_details['certificateAuthority']['data'],
            "Endpoint": cluster_details['endpoint'],
        }

        cfnresponse.send(event, context, cfnresponse.SUCCESS,
                         eventData)
        # Return the cluster ARN as the PhysicalResourceId
        return {
            'PhysicalResourceId': response['cluster']['arn'],
            "ClusterName": response['cluster']['name'],
            'Data': json.dumps(response, default=str)
        }
    except Exception as e:
        logger.error("Error: " + str(e))
        cfnresponse.send(event, context, cfnresponse.FAILED, {"Message": str(e)})
        return None


def create_cluster(eks_client, cluster_name, create_cluster_payload):
    """
    Create the EKS cluster
    """
    try:
        # Check if the cluster already exists
        logger.info(f"Checking if EKS cluster {cluster_name} already exists")
        response = eks_client.describe_cluster(name=cluster_name)
    except eks_client.exceptions.ResourceNotFoundException:
        logger.info("Creating EKS cluster with payload: " + json.dumps(create_cluster_payload, default=str))
        response = eks_client.create_cluster(**create_cluster_payload)
        logger.info("EKS cluster created: " + json.dumps(response, default=str))

    # Wait for the cluster to become ACTIVE
    cluster_details = wait_for_cluster_creation(eks_client, cluster_name)
    return cluster_details, response


def wait_for_cluster_creation(eks_client, cluster_name):
    """
    Wait for the EKS cluster to become ACTIVE.
    """
    while True:
        response = eks_client.describe_cluster(name=cluster_name)
        status = response['cluster']['status']
        if status == 'ACTIVE':
            logger.info(f"EKS cluster {cluster_name} is now ACTIVE.")
            return response['cluster']
        elif status == 'FAILED':
            raise Exception(f"EKS cluster {cluster_name} creation failed.")
        else:
            logger.info(f"EKS cluster {cluster_name} status: {status}. Waiting...")
            time.sleep(10)  # Wait 10 seconds before polling again


def delete_cluster(eks_client, cluster_name):
    """
    Delete an EKS cluster.
    """
    logger.info(f"Deleting EKS cluster: {cluster_name}")
    eks_client.delete_cluster(name=cluster_name)
    logger.info(f"EKS cluster deleted: {cluster_name}")


def nodegroup_handler(event, context):
    """
    Handles the creation, update, and deletion of an EKS managed node group as a CloudFormation custom resource.
    """
    try:
        eks_client = init_eks_client()
        iam_client = boto3.client('iam')

        cluster_name = event['ResourceProperties']['ClusterName']
        logger.info(f"cluster name : {cluster_name}")

        # Handle Delete event
        if event['RequestType'] == 'Delete':
            nodegroup_name = event['ResourceProperties']['NodegroupName']
            logger.info(f"nodegroup name : {nodegroup_name}")

            response = eks_client.describe_nodegroup(
                clusterName=cluster_name,
                nodegroupName=nodegroup_name
            )
            nodegroup_role_arn = response['nodegroup']['nodeRole']
            print(f"IAM Role ARN associated with the EKS node group: {nodegroup_role_arn}")

            role_name = nodegroup_role_arn.split('/')[-1]
            response = iam_client.list_instance_profiles_for_role(RoleName=role_name)
            instance_profiles = response['InstanceProfiles']
            for profile in instance_profiles:
                # Process each profile
                instance_profile_name = profile['InstanceProfileName']
                logger.info(f"Processing instance profile: {instance_profile_name}")

                iam_client.remove_role_from_instance_profile(
                    InstanceProfileName=instance_profile_name,
                    RoleName=role_name
                )
                print(f"Removed role {role_name} from instance profile: {instance_profile_name}")

            eks_client.delete_nodegroup(clusterName=cluster_name, nodegroupName=nodegroup_name)
            cfnresponse.send(event, context, cfnresponse.SUCCESS, {"Message": "Resource deleted"})
            return {
                'PhysicalResourceId': event['PhysicalResourceId']
            }

        # Prepare the nodegroup request payload from the custom resource properties
        nodegroup_payload = convert_keys_to_lowercase_first_letter(
            copy.deepcopy(event['ResourceProperties']))
        del nodegroup_payload['serviceToken']
        replace_boolean_strings(nodegroup_payload)
        replace_integer_strings(nodegroup_payload)

        update_payload_tags(nodegroup_payload, event)

        logger.info("EKS nodegroup with payload: " + json.dumps(nodegroup_payload, default=str))
        response = eks_client.create_nodegroup(**nodegroup_payload)
        logger.info("EKS nodegroup created: " + json.dumps(response, default=str))

        eventData = {
            "Arn": response['nodegroup']['nodegroupArn'],
            "PhysicalResourceId": response['nodegroup']['nodegroupArn'],
        }

        cfnresponse.send(event, context, cfnresponse.SUCCESS,
                         eventData)
        # Return the cluster ARN as the PhysicalResourceId
        return {
            'PhysicalResourceId': response['nodegroup']['nodegroupArn'],
            'Data': json.dumps(response, default=str)
        }

    except Exception as e:
        logger.error("Error: " + str(e))
        cfnresponse.send(event, context, cfnresponse.FAILED, {"Message": str(e)})
        return None


def access_entry_handler(event, context):
    """
    Handles the creation, update, and deletion of an EKS access entry as a CloudFormation custom resource.
    """
    try:
        eks_client = init_eks_client()

        cluster_name = event['ResourceProperties']['ClusterName']
        logger.info(f"cluster name : {cluster_name}")

        principal_arn = event['ResourceProperties']['PrincipalArn']
        logger.info(f"principal arn : {principal_arn}")

        # Handle Delete event
        if event['RequestType'] == 'Delete':
            logger.info(f"access entry principal arn : {principal_arn}")

            eks_client.delete_access_entry(clusterName=cluster_name, principalArn=principal_arn)
            cfnresponse.send(event, context, cfnresponse.SUCCESS, {"Message": "Resource deleted"})
            return {
                'PhysicalResourceId': event['PhysicalResourceId']
            }

        username = event['ResourceProperties']['Username'] if 'Username' in event['ResourceProperties'] else None
        logger.info(f"username : {username}")

        entry_type = event['ResourceProperties']['Type']
        logger.info(f"entry type : {entry_type}")

        response = create_access_entry(eks_client, principal_arn, username, cluster_name, entry_type)
        logger.info("EKS access entry created: " + json.dumps(response, default=str))

        eventData = {
            "Arn": response['accessEntry']['accessEntryArn'],
            "PhysicalResourceId": response['accessEntry']['accessEntryArn'],
        }
        cfnresponse.send(event, context, cfnresponse.SUCCESS, eventData)
        # Return the cluster ARN as the PhysicalResourceId
        return {
            'PhysicalResourceId': response['accessEntry']['accessEntryArn'],
            'Data': json.dumps(response, default=str)
        }
    except Exception as e:
        logger.error("Error: " + str(e))
        cfnresponse.send(event, context, cfnresponse.FAILED, {"Message": str(e)})
        return None
