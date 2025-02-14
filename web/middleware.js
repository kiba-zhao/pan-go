'use strict'

module.exports = (req, res, next) => {
  const _send = res.send
  res.send = function (body) {
    if (require('url').parse(req.originalUrl, true).query["_take"]) {
      try {
        const json = JSON.parse(body)
        if (Array.isArray(json)) {
          if (json.length === 1) {
            res.status(200)
            return _send.call(this, JSON.stringify(json[0]))
          } else if (json.length === 0) {
            res.status(404)
            return _send.call(this, '{}')
          }
        }
      } catch (e) {}
    }
    return _send.call(this, body)
  }
  next()
}