## AWS Monitor

[![GitHub Actions Status](https://github.com/jdahlke/aws-monitor/actions/workflows/tests.yml/badge.svg)](https://github.com/jdahlke/aws-monitor/actions/workflows/tests.yml)

Post AWS events to Slack. CloudWatch alarms are posted to an alert channel,
while other events (e.g. Code Pipeline) are posted to an info channel.

#### Test

Run Go tests as normal

```
go test ./...
```

#### Deployment

Deploy new application code

```
AWS_PROFILE=placeholder bin/deploy staging|production
```

#### Setup

Use the CloudFormation template in `aws/template.yml` to create a new Stack
`aws-monitor`.
The Lambda function only needs basic permissions.
