/*jslint browserify: true */
'use strict';

var console = require('console-browserify');

module.exports = {
    init: function (callback) {
        if (location.search) {
            callback(location.search.substring(1));
        }
    },
    update: function (params) {
        // Set URL to search query
        history.pushState(null, null, "?"+params);
}
};
