var ajax = ic.ajax;

var App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource('paths', function() {
        this.route('new');
    });
    this.resource('applications', function() {
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
        return App.Path.findAll();
    }
});

App.Path = Ember.Object.extend({
    objectForSaving: function() {
        return {
            fragment: this.get('fragment'),
            limit: parseInt(this.get('limit'), 10),
            seconds:  parseInt(this.get('seconds'), 10)
        };
    }.property('fragment', 'limit', 'seconds'),

    save: function() {
        return ajax('paths', {
            data: JSON.stringify(this.get('objectForSaving')),
            type: 'POST',
            dataType: 'json'
        });
    }
});

App.Path.reopenClass({
    findAll: function() {
        return ajaxWithWrapperObject('/paths', App.Path);
    }
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
            model.save().then(function(result) {
                self.controllerFor('paths').addObject(model);
                self.transitionTo('paths.index');
                return result;
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
