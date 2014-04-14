/* global ajax */

// TODO come back and fix me

test('App.Path objectForSaving', function() {
    var path = App.Path.create({
        fragment: '/foo',
        limit: '10',
        seconds: '5'
    });

    var output = path.get('objectForSaving');
    console.log(output);
    strictEqual(output.fragment, '/foo');
    strictEqual(output.limit, 10);
    strictEqual(output.seconds, 5);
});

asyncTest('App.Path findAll', function() {
    ajax.defineFixture('/paths', {
        response: [{
            fragment: '/foo',
            limit: 10,
            seconds: 5
        }],
        jqXHR: {},
        textStatus: 'success'
    });

    App.Path.findAll().then(function(result) {
        start();
        strictEqual(result.length, 1);
        strictEqual(result[0].constructor, App.Path);
    });
});
