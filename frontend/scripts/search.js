/*jslint browserify: true */
'use strict';

var $ = require('jquery'),
    console = require('console-browserify'),
    after = require('./after'),
    result_template = require('./templates/results'),
    SearchHistory = require('./searchhistory'),
    blocker = require('./blocker'),
    blocker_wait_time = 300;

module.exports = {
  init: function() {
    console.log('Initializing search');

    var search_form = $('#search-form'),
        result_container = $('#result-container'),
        page_number = $('#page-number');

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

    SearchHistory.init(get_results);

    function submit_form() {
      console.log('Form submit requested.');

      var serialized_form = search_form.serialize();
      SearchHistory.update(serialized_form);
      get_results(serialized_form);

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
