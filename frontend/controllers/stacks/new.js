var StacksNewController = Ember.ObjectController.extend({
    availableMiddlewares: [],

    middlewareQuery: '',

    hasMiddlewares: Ember.computed.bool('middlewares.length'),

    _middlewaresNeedConfig: function() {
        return this.get('middlewares').any(function(mw) {
            return mw.get('underlying.needsConfiguration');
        });
    }.property('middlewares.@each', 'middlewares.@each.config'),

    canBeSaved: function() {
        return this.get('hasMiddlewares') && !this.get('_middlewaresNeedConfig');
    }.property('hasMiddlewares', 'middlewaresNeedConfig'),

    _unselectedMiddlewares: function() {
        return this.get('availableMiddlewares').filter(function(item) {
            var selected = item.get('selected');
            return Ember.isEmpty(selected) || !selected;
        });
    }.property('availableMiddlewares',
               'availableMiddlewares.isFulfilled',
               'availableMiddlewares.@each.selected'),

    filteredMiddlewares: function() {
        var unselected = this.get('_unselectedMiddlewares'),
            query = this.get('middlewareQuery');

        if (Ember.isEmpty(query)) {
            return unselected;
        }

        var filterExp = new RegExp(query, 'i');

        return unselected.filter(function(item) {
            return filterExp.test(item.get('friendlyName'));
        });
    }.property('middlewareQuery', '_unselectedMiddlewares'),

    middlewareRemaining: Ember.computed.bool('_unselectedMiddlewares.length'),

    resetSelected: function() {
        this.get('availableMiddlewares').forEach(function(mw) {
            mw.set('selected', false);
        });
    },

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
        }
    }
});

export default StacksNewController;
