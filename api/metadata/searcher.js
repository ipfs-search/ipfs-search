const elasticsearch = require('elasticsearch');

function process_results(results) {
    var new_results = results["_source"]

    // Add some metadata to the metadata
    new_results["version"] = results._version
    new_results["type"] = results._type

    return new_results
}

function search(client, cid) {
    // Perform search on cid, returns promise with resulting metadata
    return new Promise(
        // The resolver function is called with the ability to resolve or
        // reject the promise
        function(resolve, reject) {
            // Just to keep code readable. *ughJS*
            function success(results) {
                resolve(process_results(results))
            }

            function fail(error) {
                console.trace(error)
                reject(error)
            }

            // Perform the actual search
            client.get({
                index: 'ipfs',
                id: cid,
                type: 'file',
                realtime: false,
                _source: 'metadata',
            }).then(success).catch(fail)
        }
    )
}

function make_searcher() {
    // Don't *do* anything at module loadnpm
    var client = new elasticsearch.Client({
        host: 'localhost:9200',
        log: 'info'
    });

    return {
        Search: function (cid) { return search(client, cid) }
    }

}

module.exports = {
    Searcher: make_searcher
}

