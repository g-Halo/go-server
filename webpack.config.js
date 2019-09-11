module.exports = {
    devtool: 'eval-source-map',
    entry: __dirname + "/assets/js/main.js",
    output: {
        path: __dirname + "/public",
        filename: "application.js"
    },
    devServer: {
        contentBase: "./public",    //本地服务器所加载的页面所在的目录
        historyApiFallback: true,   //不跳转
        inline: true                //实时刷新
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
    mode: "development"
}