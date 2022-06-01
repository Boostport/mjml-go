const path = require("path");
const webpack = require("webpack");

module.exports = {
  mode: "production",
  target: "es2019",
  optimization: {
    sideEffects: true,
  },
  resolve: {
    extensions: [".js"],
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
    filename: "mjml.js",
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
