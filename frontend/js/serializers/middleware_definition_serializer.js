App.MiddlewareDefinitionSerializer = DS.RESTSerializer.extend({
    normalizeHash: {
        middleware_definitions: function(hash) {
            if (Ember.isArray(hash.schema)) {
                hash.schema = hash.schema.map(function(item) {
                    return Ember.Object.create(item);
                });
            }
            return hash;
        }
    },
});
