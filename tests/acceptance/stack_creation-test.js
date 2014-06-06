import startApp from 'apibutler/tests/helpers/start-app';
import Middleware from 'apibutler/models/middleware';

var App;

module('Acceptance Tests - Stack Creation', {
    setup: function() {
        Middleware.reopenClass({
            FIXTURES: [
                {
                    friendlyName: 'foo',
                    id: 'foo'
                },
                {
                    friendlyName: 'bar',
                    id: 'bar'
                }
            ]
        });
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

        var stacks = find('.new_stack .stack_item');
        equal(stacks.length, 0, 'Stack creation pane is empty');
    });
});

test('Available middleware pane', function() {
    visit('/stacks/new').then(function() {
        var stackPaneTitle = find('.available .title');
        equal(stackPaneTitle.text(), 'Available middleware', 'Available middleware pane title');

        var stacks = find('.available .stack_item');
        equal(stacks.length, 2, 'Available middleware pane has some middleware');
    });
});

test('Adding some middleware to the stack', function() {
    expect(3);
    visit('/stacks/new').then(function() {
        click('.foo');
    }).then(function() {
        var stacks = find('.new_stack .stack_item');
        equal(stacks.length, 1, 'Added a middleware to the stack. Now navigating away.');
        visit('/');
    }).then(function() {
        visit('/stacks/new');
    }).then(function() {
        var stacks = find('.new_stack .stack_item');
        equal(stacks.length, 0, 'After coming back, stack creation pane is empty');
        var available = find('.available .stack_item');
        equal(available.length, 2, 'After coming back, available middleware pane has all middleware');
    });
});
