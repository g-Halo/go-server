var path = require('path');

module.exports = {
    entry: __dirname + "/assets/js/main.js",
    devtool: 'cheap-module-eval-source-map',
    output: {
        path: path.resolve(__dirname, 'public/'),
        filename: "application.js",
        publicPath: '/public/',
        sourceMapFilename: '[name].map',
        chunkFilename: '[id].chunk.js',
    },
    devServer: {
        inline: true,
        contentBase: path.join(__dirname, 'public'),
        writeToDisk: true,
        port: 4399,
    },
    resolve: {
        extensions: ['*', '.js', '.jsx']
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ['babel-loader']
            },
            {
                test: /\.css$/, // Only .css files
                loader: 'css-loader' // Run both loaders
            },
            {
                test: /\.scss$/,
                loader: 'style-loader!css-loader!sass-loader'
            }
        ]
    },
    // plugins: [HTMLWebpackPluginConfig],
    mode: "development"
}