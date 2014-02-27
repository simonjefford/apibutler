var express = require('express');

var app = express();

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

anotherapp.get('/recipes/other', function(req, res) {
    res.send({
        response: {
            backend: 'new'
        }
    });
});

anotherapp.listen(3001);

console.log('running');
