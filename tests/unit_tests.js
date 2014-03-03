/* global App */
test('App.path objectForSaving', function() {
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
