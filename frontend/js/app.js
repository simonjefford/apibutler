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
});
