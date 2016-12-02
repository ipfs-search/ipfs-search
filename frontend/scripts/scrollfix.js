/*jslint browserify: true */
'use strict';

var $ = require('jquery'),
    console = require("console-browserify");

module.exports = {
  init: function() {
    var fixTop = $('.header-wrapper').height();
    console.debug(fixTop);

    $(window).scroll(function() {
      var currentScroll = $(window).scrollTop();
      if (currentScroll >= fixTop) {
        $('.fix').addClass('fixed');
        console.debug('added Class fixed');
      } else {
        $('.fix').removeClass('fixed');
        console.debug('removed Class fixed');
      }
    });
  }
};
