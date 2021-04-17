# chiabot

This is a release bot that will notify a Slack channel of a new release for the Chia blockchain.

This project is not affiliated with the Chia blockchain, and is a fan-made project.

# Pre-requisites

You need a Slack app. Please find documentation from Slack if you have never done this before. The Slack app should be configured with the `chat:write` OAuth scope.
The app will then need to be invited to the channel you want it to write messages to.
And finally, you will need to find your channel's ID. To do this, you just need to visit your Slack channel in your web browser and find the last string in the URL path (eg. For `https://app.slack.com/client/XXXXXXXXX/ZZZZZZZZZ` the Channel ID is `ZZZZZZZZZ`)

# Use

You can either start this container via `docker run`, `docker-compose`, or view the accompanying kube manifest. 

**Quick Start:** 
```
docker container run \
  -e CHECK_INTERVAL="60" \
  -e TOKEN="my_slack_app_token" \
  -e CHANNEL_ID="slack_channel_id" \
  registry.gitlab.com/brandonbutler/chiabot:latest
```

Example docker-compose.yml:

```
version: '2'
services:

 cloudflare:
   container_name: chiabot
   image: registry.gitlab.com/brandonbutler/chiabot:latest
   environment:
     - CHECK_INTERVAL="60"
     - TOKEN="my_slack_app_token"
     - CHANNEL_ID="slack_channel_id"
```

Or if in a Kubernetes environment:

```
kubectl apply -f deployment.yaml
```

# Environment Variables

| Variable | Function |
| ---- | ---- | 
|  CHECK_INTERVAL | The duration, in seconds, between checks for a newer version of Chia (ex: "60") | 
|  TOKEN | The Slack App token | 
|  CHANNEL_ID | The ID string of the Slack channel that the App was invited to | 