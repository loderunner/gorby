'use strict'

const merge = require('webpack-merge')
const prodWebpackConfig = require('./webpack.prod.conf')

const webpackConfig = merge(prodWebpackConfig, {
  watch: true
})

module.exports = webpackConfig