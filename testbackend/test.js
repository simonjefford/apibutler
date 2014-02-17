var express = require('express');

var app = express();

app.get('/', function(req, res) {
  res.send({
    response: {
      foo: 42,
      bar: [1,2,3],
      baz: { a: "b"}
    }
  });
});

app.get('/foo', function(req, res) {
  res.send({
    response: {
      endpoint: "new"
    }
  });
});

app.listen(3000);
