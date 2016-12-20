/*jslint browserify: true */
'use strict';

var $ = require('jquery');
var block_layer = $('.blocker');

module.exports = {
    show: function () {
        console.log('Showing blocker');
        $(window).scrollTop(0);
        block_layer.show();

    },
    hide: function () {
        console.log('Hiding blocker');
        block_layer.hide();
    }
};
