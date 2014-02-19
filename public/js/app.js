var ajax = ic.ajax;

var App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource('paths', function() {
        this.route('new');
    });
});

var ajaxWithWrapperObject = function(path, klass) {
    return ajax(path).then(function(array) {
        return array.map(function(item) {
            return klass.create(item);
        });
    });
};

App.PathsRoute = Ember.Route.extend({
    model: function() {
        return ajaxWithWrapperObject('/paths', App.Path);
    }
});

App.Path = Ember.Object.extend({
    objectForSaving: function() {
        return {
            fragment: this.get('fragment'),
            limit: parseInt(this.get('limit'), 10),
            seconds:  parseInt(this.get('seconds'), 10)
        };
    }.property('fragment', 'limit', 'seconds')
});

App.NavbarLinkComponent = Ember.Component.extend({
    tagName: ''
});

App.PathsNewRoute = Ember.Route.extend({
    model: function() {
        return App.Path.create();
    },

    actions: {
        save: function(model) {
            var self = this;
            console.log('Now saving %o', JSON.stringify(model));
            ajax('/paths', {
                data: JSON.stringify(model.get('objectForSaving')),
                type: 'POST',
                dataType: 'json'
            }).then(function(result) {
                console.log('in then');
                self.controllerFor('paths').addObject(model);
                self.transitionTo('paths.index');
                return result;
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
