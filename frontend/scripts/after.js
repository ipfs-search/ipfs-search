/*jslint browserify: true */
'use strict';

var $ = require('jquery');

function after(timeout) {
  var delay_promise = $.Deferred();

  setTimeout(function () {
    delay_promise.resolve();
  }, timeout);

  return delay_promise;
}

module.exports = after;
