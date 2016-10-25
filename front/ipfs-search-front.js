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

function get_title(result) {
  // Get title from result

  var src = result._source;
  var titles = [];

  // Try metadata
  if ("metadata" in src) {
    if (src.metadata.title) {
      titles.push(src.metadata.title[0]);
    }

    if (src.metadata.name) {
      titles.push(src.metadata.name[0]);
    }
  }

  // // Try references (return the first)
  src.references.forEach(function (item) {
    if (item.name) {
      titles.push(item.name);
    }
  });

  // // Pick longest title
  if (titles.length > 0) {
    titles.sort(function (a, b) { return b.length - a.length });

    return titles[0];
  } else {
    // Fallback to id
    return result._id;
  }
}

function transform_results(results) {
  // Sanitize search results into a list like this:
  // {
  //   "hash": <>
  //   "title": <>
  //   "description": <>
  // }
  var hits = [];

  results.hits.forEach(function (item) {
    var title = get_title(item);
    var description;

    hits.push({
      "hash": item._id,
      "title": title,
      "description": description
    })
  });

  // Overwrite existing list of hits
  results.hits = hits;
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

        transform_results(body.hits);

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

