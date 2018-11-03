const express = require('express')
const app = express()
const port = 9616

const validator = new require('./validator').Validator()
const searcher = new require('./searcher').Searcher()

function invalid_requst(res, msg) {
    if (!msg) msg = 'Invalid request'

    res.status(400).json({
        error: msg
    })
}

app.get('/metadata/:cid/', function (req, res) {
    if (!req.params['cid']) invalid_requst(res)

    var cid = req.params['cid']

    if (!validator.Validate(cid)) {
        invalid_requst(res, 'Invalid cid')
    }

    return searcher.Search(cid).then(function (result) {
        res.json(result)
    }).catch(function (error) {
        if (error.status) {
            res.status(error.status)
        } else {
            res.status(500)
        }

        res.json({error: error})
    })
})

app.listen(port, () => console.log(`ipfs-search metadata API listening on port ${port}!`))
