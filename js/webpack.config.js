const path = require("path");
const webpack = require("webpack");
const merge = require("lodash.merge");

const config = {
  mode: "production",
  target: "es2019",
  optimization: {
    sideEffects: true,
  },
  context: path.resolve(__dirname, "src"),
  resolve: {
    fallback: {
      fs: false,
      https: false,
      http: false,
      os: false,
      path: false,
      url: false,
    },
    alias: {
      "uglify-js": path.resolve(__dirname, "shims/uglify-js"), // We need to do this because uglify-js can't be built with webpack
    },
  },
  output: {
    globalObject: "this",
    path: "/tmp",
    libraryTarget: "umd",
    library: "Suborbital",
    chunkFormat: "array-push",
  },
  plugins: [
    new webpack.ProvidePlugin({
      process: "process/browser",
      window: "global/window",
    }),
  ],
  performance: {
    hints: false,
  },
};

const wasmConfig = merge({}, config, {
  entry: "./index.js",
  output: {
    filename: "mjml.js",
  },
});

const testServerConfig = merge({}, config, {
  target: "node18",
  entry: "./server.js",
  output: {
    filename: "server.js",
  },
});

module.exports = [wasmConfig, testServerConfig];
