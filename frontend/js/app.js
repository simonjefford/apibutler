window.App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

require('js/components/*');
require('js/controllers/*');
require('js/models/*');
require('js/routes/*');

App.Router.map(function() {
    this.resource('apis', function() {
        this.route('new');
    });
    this.resource('applications', function() {
        this.route('new');
    });
    this.resource('middlewares', function() {
        this.route('middleware', { path: '/:middleware_id'});
    });
});

App.Middlewares = [
    Ember.Object.create({
        name: 'Authorisation',
        count: 10,
        id: 1
    }),
    Ember.Object.create({
        name: 'Throttling',
        count: 50,
        configSettings: [
            'timeInterval',
            'callCount'
        ],
        id: 2
    })
];
