const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const outputDirectory = path.resolve(__dirname, 'build');
const htmlTemplate = path.resolve(__dirname, 'public/index.html');
const serviceAddress = process.env.SERVICE_ADDRESS || 'http://localhost:9123';

module.exports = {
  entry: './src/index.js',
  output: {
    path: outputDirectory,
    publicPath: '/',
    filename: 'main.bundle.js'
  },
  module: {
    rules: [
      {
        test: /\.(jsx|js)$/,
        exclude: /node_modules/,
        use: ['babel-loader']
      },
      {
        test: /\.(jpe?g|gif|png|svg|woff|ttf)$/,
        use: [
          {
            loader: 'file-loader',
            options: {},
          }
        ],
      },
      {
        test: /\.css$/i,
        use: ['css-loader']
      },
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: htmlTemplate
    })
  ],
  devServer: {
    contentBase: outputDirectory,
    port: 9999,
    historyApiFallback: true,
    proxy: {
      '/api/*': {
        'target': serviceAddress,
        'secure': false,
        'logLevel': 'debug'
      }
    }
  },
};
