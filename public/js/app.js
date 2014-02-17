var ajax = ic.ajax;

var App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource('paths', function() {
        this.route('new');
    });
});

App.PathsRoute = Ember.Route.extend({
    model: function() {
        return ajax('/paths').then(function(paths) {
            return paths.map(function(path) {
                return App.Path.create(path);
            });
        });
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

App.PathsNewController = Ember.ObjectController.extend({
});
