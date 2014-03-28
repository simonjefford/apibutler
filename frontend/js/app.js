var ajax = window.ajax = ic.ajax;

var App = window.App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

window.ajaxWithWrapperObject = function(path, klass) {
    return ajax.request(path).then(function(array) {
        return array.map(function(item) {
            return klass.create(item);
        });
    });
};

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
