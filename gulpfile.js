/**
 *  Created by WebLss on 2018/8/3
 */
var gulp = require('gulp')
var minifyCss = require('gulp-minify-css')
var babel = require('gulp-babel')
var uglify = require('gulp-uglify')
var clearnHtml = require('gulp-cleanhtml')
var copy = require('gulp-contrib-copy')
var browserSync = require('browser-sync').create()
var rev = require('gulp-rev');
var minimist = require('minimist');
const rm = require('rimraf')
var reload = browserSync.reload

// var env = options.env === 'production'
// 定义源代码的目录和编译压缩后的目录



var options = minimist(process.argv.slice(2));
console.log("appname:"+options.appname)
console.log("env:"+options.env)

var src = 'static/' + options.appname
var dist = 'static/' + options.appname + "_dis"
var viewSrc = 'views/' + options.appname
var viewDist = 'views/' + options.appname + "_dis"

gulp.task('cleanres', function (cb) {
  rm(dist, cb);
});

gulp.task('cleanview', function (cb) {
  rm(viewDist, cb);
});

// 编译全部scss 并压缩
gulp.task('css', function () {
  gulp.src(src + '/style/*.css')
    // .pipe(rev())
    .pipe(minifyCss())
    // .pipe(rev.manifest({
    //   base: 'build/assets',
    //   merge: true // merge with the existing manifest (if one exists)
    // }))
    .pipe(gulp.dest(dist))
})

// 编译全部js 并压缩
gulp.task('js', function () {
  gulp.src(src + '/js/*.js')
    // .pipe(rev())
    .pipe(babel({
      presets: ['es2015']
    }))
    .pipe(uglify())
    // .pipe(rev.manifest({
    //   base: 'build/assets',
    //   merge: true // merge with the existing manifest (if one exists)
    // }))
    .pipe(gulp.dest(dist))
})
// 压缩全部html
gulp.task('html', function () {
  // 编译视图文件
  gulp.src(viewSrc + '/*.html')
    .pipe(clearnHtml())

    .pipe(gulp.dest(viewDist))
})

// 其他不编译的文件直接copy
// gulp.task('copy', function () {
//   gulp.src([src + '/images/**',src + '/fonts/**'])
//     .pipe(copy())
//     .pipe(gulp.dest(dist))
// })

// // 自动刷新
// gulp.task('server', function () {
//   browserSync.init({
//     //proxy: '', // 指定代理url
//     notify: false // 刷新不弹出提示
//   })
//   // 监听其他不编译的文件 有变化直接copy
//   gulp.watch(src + '/**/*.!(css|js|html)', ['copy'])
//   // 监听html文件变化后刷新页面
//   gulp.watch(src + '/**/*.js', ['js']).on('change', reload)
//   // 监听html文件变化后刷新页面
//   gulp.watch(src + '/**/*.+(html)', ['html']).on('change', reload)
//   // 监听css文件变化后刷新页面
//   gulp.watch(dist + '/**/*.css').on('change', reload)
// })
// 监听事件
gulp.task('default', ["cleanres","cleanview",'css', 'js', 'html'])