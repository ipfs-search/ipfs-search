const elasticsearch = require('elasticsearch');
const http = require('http');
const url = require('url');
const port = 9615;

var client = new elasticsearch.Client({
  host: 'localhost:9201',
  log: 'trace'
});

function query(q) {
  var body = {
      "query": {
          "query_string": {
              "query": q,
              "default_operator": "AND"
          }
      },
/*
      "highlight": {
          "order" : "score",
          "require_field_match": false,
          "encoder": "html",
          "fields": {
              "*": {
                  "number_of_fragments" : 1
              }
          }
      },
*/
      "_source": [
        "metadata.title", "metadata.name", "metadata.description",
        "references"
      ]
  }

  return client.search({
    index: 'ipfs',
    body: body
  });
}

function error_response(response, code, error) {
  console.trace(code+": "+error.message);

  response.writeHead(code, {"Content-Type": "application/json"});
  response.write(JSON.stringify({
    "error": error
  }));
  response.end();
}

console.info("Starting server on http://localhost:"+port+"/");

http.createServer(function(request, response) {
  var parsed_url;

  try {
    console.log(request.url);

    try {
      parsed_url = url.parse(request.url, true);
    } catch(err) {
      error_response(response, 400, err.message);
    }

    if (parsed_url.pathname === "/search") {
      if (!"q" in parsed_url.query) {
        error_response(response, 422, "query argument missing");
      }

      query(parsed_url.query['q']).then(function (body) {
        console.info("200: Returning "+body.hits.hits.length+" results");

        response.writeHead(200, {"Content-Type": "application/json"});
        response.write(JSON.stringify(body.hits));
        response.end();
      }, function (error) {
        throw(error);
      });

    } else {
        error_response(response, 404, "file not found");
    }

  } catch(err) {
    // Catch generic errors
    error_response(response, 500, err.message);
  } finally {

  }
}).listen(port);

