# Setup Your Environment

## Spot Credentials
Before using eksctl for Spot Ocean you need:
- Create Spot Ocean user
- Generate your access token
- Configure your Spot and aws credentials.

If you are not a Spot Ocean user, sign up for free [here](https://console.spotinst.com/spt/auth/signUp).\
Connect to your account and generate your access token [here](https://console.spotinst.com/settings/v2/tokens/permanent).\
Please copy the generated token value and the account id.

There several ways to configure your Spot credentials

To use environment variables, run:
```bash
export SPOTINST_TOKEN=<spotinst_token>
export SPOTINST_ACCOUNT=<spotinst_account>
```

To use credentials file, run the [spotctl configure](https://github.com/spotinst/spotctl#getting-started) command:
```bash
spotctl configure
? Enter your access token [? for help] **********************************
? Select your default account  [Use arrows to move, ? for more help]
> act-01234567 (prod)
  act-0abcdefg (dev)
```

Or, manually create an INI formatted file like this:
```ini
[default]
token   = <spotinst_token>
account = <spotinst_account>
```

and place it in:

- Unix/Linux/macOS:
```bash
~/.spotinst/credentials
```
- Windows:
```bash
%UserProfile%\.spotinst\credentials
```

## AWS Credentials

Please refer to aws official [Configuration and Credential File Settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) documentation for further details.

- Also see [Spot Ocean Cluster](./ocean/spot-ocean-cluster.md/#cluster)
