# envsecrets

[Quickstart](https://docs.envsecrets.com/platform/quickstart) &nbsp;&nbsp;•&nbsp;&nbsp; [Homepage](https://envsecrets.com) &nbsp;&nbsp;•&nbsp;&nbsp; [Login](https://app.envsecrets.com) &nbsp;&nbsp;•&nbsp;&nbsp; [Community](https://join.slack.com/t/envsecrets/shared_invite/zt-24djrpzpd-RegbTvnw~f__tFCx5GsqRA) &nbsp;&nbsp;•&nbsp;&nbsp; [Twitter](https://twitter.com/envsecrets)

envsecrets is an open-source free-forever cloud account to store your environment secrets and synchronize them with third-party services.

This tool is for you if you:

- Are currently hardcoding your secrets in your code.
- Are sharing `.env` files over Slack or WhatsApp.
- Are consuming the same set of secrets in multiple services/locations.
- Do not have any access control setup for your secrets.
- Need to version your secrets.

## Security

Read our detailed [data model](https://docs.envsecrets.com/security) to understand how we keep your secrets secure.

### By Design

- **End-to-End Encryption** <br />
    You are protected with public-key cryptography. Secrets are encrypted and decrypted on client side only. Never on our servers. 
- **Zero-Knowledge Architecture** <br />
    No one can see your secrets. Not even us. If our database gets hacked/leaked, attackers will never be able to decrypt your secrets.
- **Multi-Factor Authentication** <br />
    You can enable Temporal One Time Passwords on the [platform](https://app.envsecrets.com) and scan the QR in any authenticator app like Google Authenticator or Authy.

### By Promise

- **Open Source Codebase** <br />
    Feel free to scan our code to establish trust.

## Core Features

Amongst many hidden gems, the platform's core features include:

- 🔐 **Role-Based Access Control** <br />
    Never let your interns get access to production secrets.
- 🚀 **Deployment Platform Integrations** - Vercel, Docker, etc. <br />
    Push your secrets to the third-partry services where you consume them.
- 📕 **Versioning** <br />
    Want to bring back a previous value? Rollback to an older version of your secret.
- 🔑 **Services Tokens / API Keys** <br />
    Securely export and consume your secrets in places where you cannot authenticate with your account password.
- 🏗️ **CI/CD Integrations** - Github Actions, Circle CI, etc. <br />
    Push your secrets to the third-partry services where you are consuming them.
- 🛡️ **Multi-Factor Authentication** <br />
    Activate TOTP based MFA in your account. Prevent attackers from accessing your secrets just because they got your password.

## Getting Started

### Installation

Install the CLI in your system.

**MacOS**

```
brew install envsecrets/tap/envs
```

**Linux**

```
snap install envs
```

**Windows Or Any Other OS**
Download the release binary [from here](https://github.com/envsecrets/cli/releases).


### Using w/ Local Environment

- Change directory to the root of your project.
    
    ```
    cd project_root/
    ```

- Set your first secret locally.
    
    ```
    envs set first=first
    ```
    This will save your key-value pair locally **without** encrypting it.

- Get the value of a particular key.
    
    ```
    envs get first
    ```
    This should ideally print the value of `first`.

- List your locally available keys.
    
    ```
    envs set first=first
    ```

### Using w/ Remote Environment

1. Login to your envsecrets [cloud account](https://app.envsecrets.com).
1. Create a new project from your dashboard.
1. Login to your cloud account from the CLI.
    ```
    envs login
    ```
1. Now simply using the `--env` flag will run the `get/set/ls` commands on remote environments instead of your local one. To list your keys in a remote environment called `prod`, simply run:
    ```
    envs ls -e prod
    ```
1. Similarly, to get the value for key `FIRST` in the second version of your `prod` environment secret, simply run:
    ```
    envs get FIRST -v 2 -e prod
    ```
    
### Syncing w/ Third-Party Services From CLI

1. Go to the [integrations catalog](https://app.envsecrets.com/integrations/catalog) on the [platform](https://app.envsecrets.com).
1. Choose any integration and go through the setup procedure described on the platform.
1. Activate your connected integration on the `prod` environment of any project in your organisation from the [integrations page](https://app.envsecrets.com/integrations).
1. Run the following command on your terminal:
    ```
    envs sync -e prod
    ```
1. Out of the options presented to you by the CLI, select the preferred service you want to push your secrets to.
1. That's it! Go and check your service to see if the latest values have been updated.

**[Here](https://docs.envsecrets.com/integrations/overview) is the detailed documentation on how to connect and activate every individual integration.**

## Need Help?

We particularly recommend joining our [community](https://join.slack.com/t/envsecrets/shared_invite/zt-24djrpzpd-RegbTvnw~f__tFCx5GsqRA) to remain updated on best practices and bug fixes. 

- If you are stuck anywhere, ask our team in the community.
- Read the [official documentation](https://docs.envsecrets.com) for tutorials and specifications.
- If it is something specifically related to the CLI, here are the [CLI docs](https://docs.envsecrets.com/cli).
- In case of anything confidential or legal, [email us](mailto:wahal@envsecrets.com).

## Feature Requests

To request enhancements or new features, you can do either of the following:

- Text us in the [community](https://join.slack.com/t/envsecrets/shared_invite/zt-24djrpzpd-RegbTvnw~f__tFCx5GsqRA).
- Open an issue in this repostory and label it `enhancement` or `feature`. Properly decribe your requirements in the issue.
