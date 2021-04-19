'use strict'

const fetch = require('node-fetch')
const isEmpty = require('is-empty')

class SlackMessage {
  constructor (options) {
    this.subject = options.subject
    this.payload = this.createPayload({
      message: options.message || '',
      details: options.details || {},
      severity: options.severity || 'good'
    })
  }

  post ({ url, channel }) {
    if (isEmpty(url)) {
      throw new Error('SlackMessage#post: value for `url` cannot be empty')
    }
    if (isEmpty(channel)) {
      throw new Error('SlackMessage#post: value for `channel` cannot be empty')
    }

    const payload = this.payload
    payload.channel = channel

    return fetch(url, {
      method: 'POST',
      body: JSON.stringify(payload)
    })
  }

  // private

  createPayload ({ message, details, severity }) {
    const attachment = {
      color: severity,
      pretext: message
    }
    attachment.fields = Object.keys(details).map(key => {
      return {
        title: key,
        value: details[key],
        short: false
      }
    })

    return {
      username: 'AWS Monitor Bot',
      text: this.subject,
      icon_emoji: ':aws:',
      attachments: [attachment]
    }
  }
}

module.exports = SlackMessage
