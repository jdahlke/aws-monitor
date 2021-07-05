'use strict'

const createEvent = require('./create-event')
const SlackMessage = require('./slack/message')

async function alert (event, context, callback) {
  console.log('Alarm: ', JSON.stringify(event, null, 2))

  const alarmEvent = createEvent(event)

  if (alarmEvent.autoScaling && alarmEvent.autoScaling()) {
    callback(null)
    return
  }

  try {
    await notify({
      event: alarmEvent,
      slack: {
        url: process.env.SLACK_URL,
        channel: process.env.SLACK_CHANNEL_ALERT
      }
    })
    callback(null)
  } catch (error) {
    callback(error)
  }
}

async function info (event, context, callback) {
  console.log('Info: ', JSON.stringify(event, null, 2))

  const infoEvent = createEvent(event)

  try {
    await notify({
      event: infoEvent,
      slack: {
        url: process.env.SLACK_URL,
        channel: process.env.SLACK_CHANNEL_INFO
      }
    })
    callback(null)
  } catch (error) {
    callback(error)
  }
}

async function notify ({ event, slack }) {
  const message = new SlackMessage({
    subject: event.subject(),
    message: event.message(),
    details: event.details(),
    severity: event.severity()
  })
  const response = await message.post({
    url: slack.url,
    channel: slack.channel
  })

  return response.text().then(text => {
    console.log('Response', {
      status: response.status,
      body: text
    })
  })
}

module.exports = {
  alert,
  info
}
