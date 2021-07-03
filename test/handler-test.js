/* eslint-env mocha */
'use strict'

const LambdaTester = require('lambda-tester')

const { mockSlack } = require('./test-helper')
const handler = require('../src/handler')

describe('handler', () => {
  beforeEach(() => {
    mockSlack()
  })

  describe('.alert', () => {
    const event = {
      'detail-type': 'CloudWatch Alarm State Change',
      source: 'aws.cloudwatch',
      account: '123456789012',
      time: '2019-10-02T17:04:40Z',
      region: 'us-east-1',
      resources: [
        'arn:aws:cloudwatch:us-east-1:123456789012:alarm:ServerCpuTooHigh'
      ],
      detail: {
        alarmName: 'ServerCpuTooHigh',
        configuration: {
          description: 'Goes into alarm when server CPU utilization is too high!',
          metrics: [
            {
              id: '30b6c6b2-a864-43a2-4877-c09a1afc3b87',
              metricStat: {
                metric: {
                  dimensions: {
                    InstanceId: 'i-12345678901234567'
                  },
                  name: 'CPUUtilization',
                  namespace: 'AWS/EC2'
                },
                period: 300,
                stat: 'Average'
              },
              returnData: true
            }
          ]
        },
        previousState: {
          reason: 'Threshold Crossed: 1 out of the last 1 datapoints [0.0666851903306472 (01/10/19 13:46:00)] was not greater than the threshold (50.0) (minimum 1 datapoint for ALARM -> OK transition).',
          reasonData: '{"version":"1.0","queryDate":"2019-10-01T13:56:40.985+0000","startDate":"2019-10-01T13:46:00.000+0000","statistic":"Average","period":300,"recentDatapoints":[0.0666851903306472],"threshold":50.0}',
          timestamp: '2019-10-01T13:56:40.987+0000',
          value: 'OK'
        },
        state: {
          reason: 'Threshold Crossed: 1 out of the last 1 datapoints [99.50160229693434 (02/10/19 16:59:00)] was greater than the threshold (50.0) (minimum 1 datapoint for OK -> ALARM transition).',
          reasonData: '{"version":"1.0","queryDate":"2019-10-02T17:04:40.985+0000","startDate":"2019-10-02T16:59:00.000+0000","statistic":"Average","period":300,"recentDatapoints":[99.50160229693434],"threshold":50.0}',
          timestamp: '2019-10-02T17:04:40.989+0000',
          value: 'ALARM'
        }
      }
    }

    it('returns success', () => {
      return LambdaTester(handler.alert)
        .event(event)
        .expectResult()
    })
  })

  describe('.info', () => {
    describe('Code Pipeline event', () => {
      const event = {
        'detail-type': 'CodePipeline Action Execution State Change',
        source: 'aws.codepipeline',
        account: 123456789012,
        time: '2020-01-24T22:03:07Z',
        region: 'us-east-1',
        resources: [
          'arn:aws:codepipeline:us-east-1:123456789012:myPipeline'
        ],
        detail: {
          pipeline: 'myPipeline',
          'execution-id': '12345678-1234-5678-abcd-12345678abcd',
          stage: 'Prod',
          action: 'myAction',
          state: 'STARTED',
          type: {
            owner: 'AWS',
            category: 'Deploy',
            provider: 'CodeDeploy',
            version: 1
          }
        }
      }

      it('returns success', () => {
        return LambdaTester(handler.info)
          .event(event)
          .expectResult()
      })
    })

    describe('EC2 Instance state change event', () => {
      const event = {
        'detail-type': 'EC2 Instance State-change Notification',
        source: 'aws.ec2',
        account: '123456789012',
        time: '2019-11-11T21: 29: 54Z',
        region: 'us-east-1',
        resources: [
          'arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111'
        ],
        detail: {
          'instance-id': 'i-abcd1111',
          state: 'pending'
        }
      }

      it('returns success', () => {
        return LambdaTester(handler.info)
          .event(event)
          .expectResult()
      })
    })
  })
})
