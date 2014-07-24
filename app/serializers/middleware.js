var MiddlewareSerializer = DS.RESTSerializer.extend({
    normalizeHash: {
        middlewares: function(hash) {
            if (Ember.isArray(hash.schema)) {
                hash.schema = hash.schema.map(function(item) {
                    return Ember.Object.create(item);
                });
            }
            return hash;
        }
    },
});

export default MiddlewareSerializer;
