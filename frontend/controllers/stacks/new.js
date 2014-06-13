var StacksNewController = Ember.ObjectController.extend({
    availableMiddlewares: [],

    middlewareQuery: '',

    hasMiddlewares: Ember.computed.bool('selectedMiddlewares.length'),

    middlewareConfig: Ember.Object.create({
        _changes: 0
    }),

    currentMiddleware: '',

    _configForMiddleware: function(name) {
        return this.get('middlewareConfig.' + name);
    },

    currentConfig: function() {
        return this._configForMiddleware(this.get('currentMiddleware'));
    }.property('currentMiddleware'),

    _middlewaresNeedConfig: function() {
        var self = this;
        return this.get('selectedMiddlewares').any(function(mw) {
            var needsConfig = mw.get('needsConfiguration') &&
                !mw.isValid(self._configForMiddleware(mw.get('name')));
            console.debug('%s: needsConfig=%s', mw.get('name'), needsConfig);
            return needsConfig;
        });
    }.property('selectedMiddlewares.@each', 'selectedMiddlewares.@each.config', 'middlewareConfig._changes'),

    canBeSaved: function() {
        return this.get('hasMiddlewares') && !this.get('_middlewaresNeedConfig');
    }.property('hasMiddlewares', '_middlewaresNeedConfig'),

    _unselectedMiddlewares: function() {
        return this.get('availableMiddlewares').filter(function(item) {
            var selected = item.get('selected');
            return Ember.isEmpty(selected) || !selected;
        });
    }.property('availableMiddlewares',
               'availableMiddlewares.isFulfilled',
               'availableMiddlewares.@each.selected'),

    selectedMiddlewares: function() {
        var selected = this.get('availableMiddlewares').filter(function(item) {
            var selected = item.get('selected');
            return !Ember.isEmpty(selected) && selected;
        });

        return Ember.ArrayController.create({content:selected});
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

        this.set('middlewareConfig', Ember.Object.create({_changes: 0}));
        this.set('currentMiddleware', '');
    },

    actions: {
        addToStack: function(mw) {
            mw.set('selected', true);
            var configPath = 'middlewareConfig.' + mw.get('name');
            if (mw.get('needsConfiguration') && Ember.isBlank(this.get(configPath))) {
                this.set(configPath, Ember.Object.create({}));
            }
            this.set('middlewareQuery', '');
        },

        removeFromStack: function(mw) {
            mw.set('selected', false);
        },

        configure: function(mw) {
            this.set('currentMiddleware', mw.get('name'));
        }
    }
});

export default StacksNewController;
