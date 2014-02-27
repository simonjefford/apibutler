var express = require('express');

var app = express();

var logger = function(prefix) {
    return express.logger({
        format: '[' + prefix + '] [:date] ":method :url :status" :res[content-length]'
    });
};

app.use(logger('backend1'));

app.get('/recipes', function(req, res) {
    res.send({
        response: {
            foo: 42,
            bar: [1,2,3],
            baz: { a: 'b'}
        }
    });
});

app.get('/recipes/foo', function(req, res) {
    res.send({
        response: {
            endpoint: 'new'
        }
    });
});


app.listen(3000);

var anotherapp = express();

anotherapp.use(logger('backend2'));

anotherapp.get('/recipes/other', function(req, res) {
    res.send({
        response: {
            backend: 'new'
        }
    });
});

anotherapp.listen(3001);

console.log('running');
