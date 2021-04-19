'use strict'

/*
 * {
 *   'detail-type': 'CodePipeline Action Execution State Change',
 *   source: 'aws.codepipeline',
 *   account: 123456789012,
 *   time: '2020-01-24T22:03:07Z',
 *   region: 'us-east-1',
 *   resources: [
 *     'arn:aws:codepipeline:us-east-1:123456789012:myPipeline'
 *   ],
 *   detail: {
 *     pipeline: 'myPipeline',
 *     'execution-id': '12345678-1234-5678-abcd-12345678abcd',
 *     stage: 'Prod',
 *     action: 'myAction',
 *     state: 'STARTED',
 *     type: {
 *       owner: 'AWS',
 *       category: 'Deploy',
 *       provider: 'CodeDeploy',
 *       version: 1
 *     },
 *     'input-artifacts': [
 *       {
 *         name: 'SourceArtifact',
 *         s3location: {
 *           bucket: 'codepipeline-us-east-1-BUCKETEXAMPLE',
 *           key: 'myPipeline/SourceArti/KEYEXAMPLE'
 *         }
 *       }
 *     ]
 *   }
 * }
 */
class PipelineEvent {
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
    const { pipeline, action } = detail
    const execution = detail['execution-id']
    const state = detail.state.toLowerCase()

    const title = `*Pipeline ${pipeline} ${action} ${state}*`
    const link = `<https://${region}.console.aws.amazon.com/codesuite/codepipeline/pipelines/${pipeline}/executions/${execution}/timeline|view details>`

    return {
      [title]: link
    }
  }

  severity () {
    const { detail } = this.event

    if (detail.state === 'FAILED') {
      return 'danger'
    }

    return 'good'
  }
}

module.exports = PipelineEvent
