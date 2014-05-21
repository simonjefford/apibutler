import { test, moduleFor } from 'ember-qunit';

moduleFor('controller:apis', "Unit - ApisController", {
    setup: function() {},
    teardown: function() {}
});

test('it exists', function() {
    ok(this.subject(), 'can be created');
});

test('headers for table', function() {
    var headers = this.subject().get('headers');
    deepEqual(['Path', 'App'], headers);
});

test('apisToShow', function() {
    var controller = this.subject();
    controller.set('content', [{
        name: 'stock',
        isNew: true
    },{
        name: 'ingredients',
        isNew: false
    }]);

    var toShow = controller.get('apisToShow');
    ok(toShow.get('length') === 1, 'only shows loaded apis');
});
