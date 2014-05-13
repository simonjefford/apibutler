window.App = Ember.Application.create({
    LOG_TRANSITIONS: true
});


Ember.Application.initializer({
    name: 'store-debugger',
    after: 'store',

    initialize: function(container) {
        window.Store = container.lookup('store:main');
    }
});
