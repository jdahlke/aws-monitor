'use strict'

/*
 *  {
 *    version: '0',
 *    id: 'c4c1c1c9-6542-e61b-6ef0-8c4d36933a92',
 *    'detail-type': 'CloudWatch Alarm State Change',
 *    source: 'aws.cloudwatch',
 *    account: '123456789012',
 *    time: '2019-10-02T17:04:40Z',
 *    region: 'us-east-1',
 *    resources: [
 *      'arn:aws:cloudwatch:us-east-1:123456789012:alarm:ServerCpuTooHigh'
 *    ],
 *    detail: {
 *      alarmName: 'ServerCpuTooHigh',
 *      configuration: {
 *        description: 'Goes into alarm when server CPU utilization is too high!',
 *        metrics: [
 *          {
 *            id: '30b6c6b2-a864-43a2-4877-c09a1afc3b87',
 *            metricStat: {
 *              metric: {
 *                dimensions: {
 *                  InstanceId: 'i-12345678901234567'
 *                },
 *                name: 'CPUUtilization',
 *                namespace: 'AWS/EC2'
 *              },
 *              period: 300,
 *              stat: 'Average'
 *            },
 *            returnData: true
 *          }
 *        ]
 *      },
 *      previousState: {
 *        reason: 'Threshold Crossed: 1 out of the last 1 datapoints [0.0666851903306472 (01/10/19 13:46:00)] was not greater than the threshold (50.0) (minimum 1 datapoint for ALARM -> OK transition).',
 *        reasonData: '{"version":"1.0","queryDate":"2019-10-01T13:56:40.985+0000","startDate":"2019-10-01T13:46:00.000+0000","statistic":"Average","period":300,"recentDatapoints":[0.0666851903306472],"threshold":50.0}',
 *        timestamp: '2019-10-01T13:56:40.987+0000',
 *        value: 'OK'
 *      },
 *      state: {
 *        reason: 'Threshold Crossed: 1 out of the last 1 datapoints [99.50160229693434 (02/10/19 16:59:00)] was greater than the threshold (50.0) (minimum 1 datapoint for OK -> ALARM transition).',
 *        reasonData: '{"version":"1.0","queryDate":"2019-10-02T17:04:40.985+0000","startDate":"2019-10-02T16:59:00.000+0000","statistic":"Average","period":300,"recentDatapoints":[99.50160229693434],"threshold":50.0}',
 *        timestamp: '2019-10-02T17:04:40.989+0000',
 *        value: 'ALARM'
 *      }
 *    }
 *  }
 */
class AlarmEvent {
  constructor (event) {
    this.event = event
  }

  subject () {
    return ''
  }

  message () {
    return ''
  }

  details () {
    const { detail, region } = this.event
    const { alarmName, state } = detail

    const title = `*Alarm ${alarmName} changed to ${state.value}*`
    const description = detail.configuration.description
    const link = `<https://${region}.console.aws.amazon.com/cloudwatch/home#alarmsV2:alarm/${alarmName}|view details>`

    return {
      [title]: `${description} - ${link}`
    }
  }

  severity () {
    const { detail } = this.event

    if (detail.state.value === 'ALARM') {
      return 'danger'
    }

    return 'good'
  }

  autoScaling () {
    const { detail } = this.event
    const description = detail.configuration.description

    if (description.includes('TargetTrackingScaling policy')) {
      return true
    }

    return false
  }
}

AlarmEvent.detailType = 'CloudWatch Alarm State Change'

module.exports = AlarmEvent
