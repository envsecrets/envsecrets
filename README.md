# Backend

## How to develop locally?

1. Clone the `main` branch of this repository with the following command:

    ```
    git clone git@github.com:envsecrets/envsecrets.git
    ```

1. Launch the docker engine in your system. Download [from here](https://www.docker.com/products/docker-desktop/) if you don't already have it.

1. Install the Nhost CLI in your system using the following command:

    ```
    sudo curl -L https://raw.githubusercontent.com/nhost/cli/main/get.sh | bash
    ```

1. Run the local environment using:

    ```
    nhost up
    ```

## How to contribute?

1. Pull latest changes from the `main` branch.
1. Checkout your feature branch.
1. Make changes/migrations from Hasura console or directly in the Go code.
1. Open a PR to the `main` branch.
1. Make sure your PR is thoroughly reviewed and approved before being merged.
1. Always follow **squash and merge** to keep the commit history clean. Please don't merge the commit history to the `main` branch.

**NEVER** push directly to the `main` branch because it deploys to production.
