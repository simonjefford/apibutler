import startApp from 'apibutler/tests/helpers/start-app';

var App;

module('Acceptance Tests - Stack Creation', {
    setup: function() {
        App = startApp();
    },
    teardown: function() {
        Ember.run(App, 'destroy');
    }
});

test('Stack creation pane', function() {
    visit('/stacks/new').then(function() {
        var stackPaneTitle = find('.new_stack .title');
        equal(stackPaneTitle.text(), 'New stack', 'Stack creation pane title');

        var stacks = find('.new_stack .stack_items');
        equal(stacks.length, 0, 'Stack creation pane is empty');
    });
});
