
const webpack = require('webpack');
const path = require('path');
const CleanWebpackPlugin = require("clean-webpack-plugin");
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const StyleExtHtmlWebpackPlugin = require("style-ext-html-webpack-plugin");
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');
// const CompressionPlugin = require('compression-webpack-plugin');
const ConcatPlugin = require('webpack-concat-plugin');
const HtmlWebpackPlugin = require("html-webpack-plugin");
const PreloadWebpackPlugin = require("preload-webpack-plugin");
const HtmlWebpackExcludeAssetsPlugin = require('html-webpack-exclude-assets-plugin');

const sourcePath = path.join(__dirname, 'src');
const stylePath = path.join(sourcePath,"styles");
const outputPath = path.join(__dirname,"build");

// const htmlIndexPath = path.join(__dirname,"/");




module.exports = env =>  {

  const NODE_ENV = (env.NODE_ENV===undefined)?"development":env.NODE_ENV;
  const isDevelopment = NODE_ENV === "development";

  const sourceMap = (isDevelopment)? "inline-source-map":undefined;
  // const publicPath = (isDevelopment) ? "http://localhost:8080/":"http://192.168.0.5:8080/"
  const publicPath = path.join(__dirname,'./')
  
  const uglifyjsPlugin = (isDevelopment)? 
    ()=>{}:
    new UglifyJSPlugin({
      sourceMap: false,
      uglifyOptions:{
        output: {
          comments: false
        }
      },
      parallel: true,
 
      
    });

  const extractLess = new ExtractTextPlugin({
    filename: "[name].[contenthash].css",
    allChunks :true,
    disable: isDevelopment
  });

  const excludeStyleJs = (isDevelopment)? "":/style.*.js/;

  const commonsChunkPlugin = new webpack.optimize.CommonsChunkPlugin({
    name:"vendor",
    filename: "vendor.[chunkhash].js",
    minChunks: function(module, count) { 
      let validation = (/node_modules/.test(module.resource));
      return validation;
    },
  });



  return{    
    entry: {
      styles: path.resolve(stylePath,"main.less"),      
      app: path.resolve(sourcePath,"index.tsx"),
      
      // app:"./src/init.tsx"
      // init:"./src/init.tsx"
    },
    output: {
        path: outputPath,
        filename: "[name].[chunkhash].js",
        // publicPath: outputPath,

        //window.location.href
        
    },
    resolve: {
      extensions: ['.webpack-loader.js', '.web-loader.js', '.loader.js','d.ts','.ts', '.tsx', '.js'],
      modules: [
        sourcePath,
        path.resolve(__dirname, 'node_modules')
    ]
    },


    plugins: [
      new CleanWebpackPlugin([outputPath],{exclude:["sdk.js"]}),
      commonsChunkPlugin,
        // extract inline css into separate 'styles.css'
        // new ExtractTextPlugin("[name].css"),
        new webpack.IgnorePlugin(/^\.\/locale$/,/moment$/),
        new webpack.DefinePlugin({
          "process.env":{
            "NODE": JSON.stringify(NODE_ENV)
          }
        }),
        uglifyjsPlugin,
        new HtmlWebpackPlugin({
          template: path.join(sourcePath, "template.html"),
          excludeChunks: ["base"],
          excludeAssets: excludeStyleJs,
          filename: path.join(outputPath,"index.html"),
          minify:{
            collapseWhitespace: true,
            collapseInlineTagWhitespace: true,
            removeComments: true,
            removeRedundantAttrubutes: true
          },
        }),
        new HtmlWebpackExcludeAssetsPlugin(),
        new webpack.HashedModuleIdsPlugin(),
        new PreloadWebpackPlugin({
          rel: "preload",
          as: "script",
          include: "all",
          fileBlacklist: [/\.(css|map)$/, /base?.+/,/style?./,/app?./],
        }),
        extractLess,
        // new StyleExtHtmlWebpackPlugin({
        //   minify: true
        // })
        
    ],


    module: {
      rules: [
        {
              test: /\.less$/,
              include: stylePath,
              use: extractLess.extract({
                use: [{
                  // include: stylePath,
                  loader: "css-loader",options: { minimize: true } 
                  }, {
                  // include: stylePath,                    
                  loader: "less-loader",options: { minimize: true } 
                }],
                fallback: "style-loader"
              })
              },
        { 
          test: /\.tsx?$/, 
          include: [sourcePath,stylePath],
          loader: 'ts-loader',
          options: {
            transpileOnly: true
          } 
        }
      ],

    },
    devtool: sourceMap,

    devServer: {
      contentBase: __dirname,
      inline: true,
      stats: "errors-only"
    },
  };
};
