'use strict'

const AlarmEvent = require('./events/alarm')
const InstanceStateChangeEvent = require('./events/instance-state-change')
const PipelineEvent = require('./events/pipeline')

function createEvent (event) {
  const detailType = event['detail-type']

  switch (detailType) {
    case AlarmEvent.detailType:
      return new AlarmEvent(event)
    case InstanceStateChangeEvent.detailType:
      return new InstanceStateChangeEvent(event)
    case PipelineEvent.detailType:
      return new PipelineEvent(event)
    default:
      throw new Error(`Unknown event.detail-type: ${detailType}`)
  }
}

module.exports = createEvent
