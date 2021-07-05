'use strict'

/*
 * {
 *   'detail-type': 'EC2 Instance State-change Notification',
 *   source: 'aws.ec2',
 *   account: '123456789012',
 *   time: '2019-11-11T21: 29: 54Z',
 *   region: 'us-east-1',
 *   resources: [
 *     'arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111'
 *   ],
 *   detail: {
 *     'instance-id': 'i-abcd1111',
 *     state: 'pending'
 *   }
 * }
 */
class InstanceStateChangeEvent {
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
    const instanceId = detail['instance-id']
    const state = detail.state

    const title = `*EC2 Instance ${instanceId}: ${state}*`
    const link = `<https://${region}.console.aws.amazon.com/ec2/v2/home#InstanceDetails:instanceId=${instanceId}|view details>`

    return {
      [title]: link
    }
  }

  severity () {
    const { detail } = this.event

    switch (detail.state) {
      case 'stopping':
      case 'stopped':
        return 'danger'
      default:
        return 'good'
    }
  }
}

InstanceStateChangeEvent.detailType = 'EC2 Instance State-change Notification'

module.exports = InstanceStateChangeEvent
