/*jslint node: true */

var gulp        = require('gulp');
var browserify  = require('browserify');
var browserSync = require('browser-sync').create();
var less        = require('gulp-less');
var buffer      = require('vinyl-buffer');
var gutil       = require('gulp-util');
var reload      = browserSync.reload;
var sourcemaps  = require('gulp-sourcemaps');
var source      = require('vinyl-source-stream');
var uglify      = require('gulp-uglify');
var handlebars  = require('gulp-handlebars');
var defineModule = require('gulp-define-module');

// Use browserify
gulp.task('browserify', function() {
    return browserify('./scripts/app.js', { debug: true })
        .bundle()
        .on('error', gutil.log.bind(gutil, 'Browserify Error'))
        // Convert it to streaming vinyl file object
        .pipe(source('bundle.js'))
        .pipe(buffer())
        .pipe(sourcemaps.init({loadMaps: true})) // loads map from browserify file
        .pipe(uglify())
        .pipe(sourcemaps.write('./')) // writes .map file
        // Start piping stream to tasks!
        .pipe(gulp.dest('./app/'));
});


// Compile less
gulp.task('less', function () {
  return gulp.src('less/style.less')
    .pipe(less())
    .pipe(gulp.dest('./app/style'))
    .pipe(reload({stream: true}));
});

// Templates
gulp.task('templates', function(){
  gulp.src('templates/*.hbs')
    .pipe(handlebars({
        handlebars: require('handlebars')
    }))
    .pipe(defineModule('commonjs'))
    .pipe(gulp.dest('scripts/templates/'));
});


// Watch Files For Changes
gulp.task('watch', function() {
    gulp.watch("./less/**/*.less", ['less']);
    gulp.watch("./scripts/**/*.js", ['browserify']);
    gulp.watch("./templates/**/*.hbs", ['templates']);
});

// Static server
gulp.task('browser-sync', function() {
    browserSync.init({
        server: {
            baseDir: "./app/"
        }
    });
    gulp.watch("./app/index.html").on('change', reload);
    gulp.watch("./app/bundle.js").on('change', reload);
});


gulp.task('build', ['templates', 'browserify', 'less']);
gulp.task('default', ['build', 'watch', 'browser-sync']);
