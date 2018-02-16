/*jslint browserify: true */
'use strict';

var $ = require('jquery'),
    console = require('console-browserify'),
    after = require('./after'),
    result_template = require('./templates/results'),
    form_history = require('./form_history'),
    blocker = require('./blocker'),
    blocker_wait_time = 300;

// Ugly shizzle right here!
window.jQuery = $;
require('./jquery.deserialize');

module.exports = {
  init: function() {
    console.log('Initializing search');

    var search_form = $('#search-form'),
        result_container = $('#result-container'),
        page_number = $('#page-number');

    search_form.data('ipfs-gateway', 'https://gateway.ipfs.io');

    var gateway_promise = $.get({'url': 'http://localhost:8080/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv/ping', 'cache': false});

    gateway_promise.done(function (result) {
      if (result === 'ipfs') {
        search_form.data('ipfs-gateway', 'http://localhost:8080');
      }
    });

    function get_results(params) {
      // Show blocker for longer waits

      // This is silly, but sad. Promises can't be used with timeouts for
      // some reason. Maybe switch to 'proper' promises?
      var done = false;

      var result_promise = $.get(
        search_form.attr('action'),
        params
      );

      result_promise.done(function (results) {

        results.gateway = search_form.data('ipfs-gateway');

        result_container.html(result_template(results));

        done = true;
        blocker.hide();

        // Wait for re-render
        setTimeout(function () {
          $(window).scrollTop($('.header-wrapper').height());
        }, 100);
      });


      after(blocker_wait_time).done(function () {
        if (!done) {
          blocker.show();
        }
      });
    }

    form_history.init(function (params) {
      $("#search-form").deserialize(params);
      get_results(params);
    });

    function submit_form() {
      console.log('Form submit requested.');

      var serialized_form = search_form.serialize();
      get_results(serialized_form);
      form_history.update(serialized_form);

      return false;
    }

    search_form.submit(function () {
      // Reset page number
      page_number.val(0);
      return submit_form();
    });

    window.next_page = function () {
      console.log('Page increase');
      page_number.val(parseInt(page_number.val()) + 1);
      submit_form();
    };

    window.prev_page = function () {
      console.log('Page decrease');
      page_number.val(parseInt(page_number.val()) - 1);
      submit_form();
    };

  }
};
