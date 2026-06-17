import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { cronExpr, weeklyCronExpr } from './cron.js'

describe('cron helpers', () => {
  it('builds daily cron from HH:MM', () => {
    assert.equal(cronExpr('09:00'), '0 9 * * *')
    assert.equal(cronExpr('22:30'), '30 22 * * *')
  })

  it('builds weekly cron from day and time', () => {
    assert.equal(weeklyCronExpr('sunday', '10:00'), '0 10 * * 0')
    assert.equal(weeklyCronExpr('friday', '18:45'), '45 18 * * 5')
  })
})
