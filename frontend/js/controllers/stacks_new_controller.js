App.StacksNewController = Ember.ObjectController.extend({
    availableMiddlewares: [],

    middlewareQuery: '',

    unselectedMiddlewares: function() {
        return this.get('availableMiddlewares').filter(function(item) {
            var selected = item.get('selected');
            return Ember.isEmpty(selected) || !selected;
        });
    }.property('availableMiddlewares',
               'availableMiddlewares.isFulfilled',
               'availableMiddlewares.@each.selected'),

    filteredMiddlewares: function() {
        if (Ember.isEmpty(this.get('middlewareQuery'))) {
            return this.get('unselectedMiddlewares');
        }

        var filterExp = new RegExp(this.get('middlewareQuery'), 'i');

        return this.get('unselectedMiddlewares').filter(function(item) {
            return filterExp.test(item.get('friendlyName'));
        });
    }.property('middlewareQuery', 'unselectedMiddlewares'),

    middlewareRemaining: Ember.computed.bool('unselectedMiddlewares.length'),

    actions: {
        addToStack: function(mw) {
            mw.set('selected', true);
            this.set('middlewareQuery', '');
            var middlewares = this.get('middlewares');
            middlewares.pushObject(Ember.Object.create({
                name: mw.get('name'),
                underlying: mw,
                config: {},
                parent: middlewares
            }));
        },

        removeFromStack: function(mw) {
            mw.parent.removeObject(mw);
            mw.get('underlying').toggleProperty('selected');
        },

        configure: function(mw) {
            console.log('now configuring', mw.get('underlying.friendlyName'));
        }
    }
});
