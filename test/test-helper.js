/* eslint-env mocha */
'use strict'

process.env.STAGE = 'test'
process.env.SLACK_URL = 'https://slack.example.com/test'
process.env.SLACK_CHANNEL_ALERT = 'channel-alert'
process.env.SLACK_CHANNEL_INFO = 'channel-info'

const nock = require('nock')

function mockSlack () {
  nock(process.env.SLACK_URL)
    .post('')
    .reply(201, 'Success')
}

module.exports = {
  mockSlack
}
