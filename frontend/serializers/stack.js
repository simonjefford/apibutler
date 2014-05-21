var isArray = Ember.isArray;

var wrapObjects = function(array) {
    if (isArray(array)) {
        return array.map(function(item) {
            return Ember.Object.create(item);
        });
    }
};

var StackSerializer = DS.RESTSerializer.extend({
    normalizeHash: {
        stacks: function(hash) {
            hash.middlewares = wrapObjects(hash.middlewares);
        }
    }
});

export default StackSerializer;
