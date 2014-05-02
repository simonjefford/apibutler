App.MiddlewareSerializer = DS.RESTSerializer.extend({
    normalizeHash: {
        middlewares: function(hash) {
            if (Ember.isArray(hash.configItems)) {
                hash.configItems = hash.configItems.map(function(item) {
                    return Ember.Object.create(item);
                });
            }
            return hash;
        }
    },
});
