module.exports = function (req, res, next) {
    if (req.method === 'POST') {
        req.method = 'GET'
    }
    if (req.method === 'PUT') {
        req.method = 'GET'
    }
    next()
}
