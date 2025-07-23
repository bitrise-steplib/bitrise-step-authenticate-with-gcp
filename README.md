# Authenticate with Google Cloud Platform (GCP)

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-authenticate-wth-gcp?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-authenticate-wth-gcp/releases)

The step authenticates with GCP using a service account key.

<details>
<summary>Description</summary>

This step authenticates with Google Cloud Platform (GCP) using a service account key.

It generates a Google auth token and sets the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

The step can also log in to Docker to the supplied Artifact Registry locations.
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `service_account_key` | The service account key in JSON format.  You can generate a service account key in the Google Cloud Console. Make sure to grant the necessary permissions to the service account.  The step will use this key to authenticate with GCP and generate a Google auth token. | required, sensitive |  |
| `docker_login` | Performs Docker login with an auth token.  The step will log in to the Artifact Registry locations specified in the `artifact_registry_locations` input.  It is supported only on the Linux stacks. | required | `false` |
| `artifact_registry_locations` | A newline (`\n`) separated list of Artifact Registry locations to log in to.  The step will log in to the specified locations using the Google auth token generated from the service account key.  These locations can point ot either the deprecated Container Registry or the new Artifact Registry. |  |  |
| `verbose` | Enable logging additional information for debugging. | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `GOOGLE_AUTH_TOKEN` | The generated Google auth token.  This token can be used to authenticate with GCP services. |
| `GOOGLE_APPLICATION_CREDENTIALS` | The path to the service account key file.  This file is used by GCP SDKs and libraries to authenticate with GCP services. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-authenticate-wth-gcp/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-authenticate-wth-gcp/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Note: this step's end-to-end tests (defined in e2e/bitrise.yml) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
