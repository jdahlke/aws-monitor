## AWS Monitor

[![GitHub Actions Test Status](https://github.com/jdahlke/aws-monitor/workflows/Tests/badge.svg?branch=develop)](https://github.com/jdahlke/aws-monitor/actions)

Post AWS events to Slack. CloudWatch alarms are posted to an alert channel,
while other events (e.g. Code Pipeline) are posted to an info channel.

It's build with the [serverless](https://serverless.com) using AWS Lambda and AWS EventBridge.


### Installation

1. Install dependencies
```
mkdir node_modules
npm install
```

2. Create `.env`


### Deployment

```
AWS_PROFILE=PROFILE sls deploy -s STAGE
```


### Test and Lint

```
npm run test
```

```
npm run lint
```
