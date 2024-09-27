[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/telemetry-manager)](https://api.reuse.software/info/github.com/kyma-project/telemetry-manager)

# slack-bot

Slack bot copying reference to a ping to notifications channel

The bot is designed to listen on messages containing specific tag, copy reference to a designated notifications channel and reply to the original message an acknowledgement. For configuration and deployment, use [helm chart](./chart)

General information about Slack tokens can be found <https://api.slack.com/authentication/token-types>.
The `gopher-bot` slack profile <https://api.slack.com/apps/A035YCW3PAP>.

## Release

The `post-slack-bot-build` workflow creates a `slack-bot` Docker image in the [registry](https://console.cloud.google.com/artifacts/docker/kyma-project/europe/prod/slack-bot)
when there is a push to the main branch. The image is tagged with the date and commit SHA. Once the Docker image is built and pushed to the registry,
replace the current Docker image in deployment with the newly created image.

```bash
kubectl set image deployment/gopher-slack-bot bot=europe-docker.pkg.dev/kyma-project/prod/slack-bot:<tag>
```
