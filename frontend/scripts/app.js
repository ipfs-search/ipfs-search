/*jslint browserify: true */
'use strict';

var $ = require('jquery'),
    console = require("console-browserify"),
    scrollfix = require('./scrollfix'),
    search = require('./search');

$(function () {
  console.log('init');
  scrollfix.init();
  search.init();
});



