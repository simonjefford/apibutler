var Middleware = DS.Model.extend({
    friendlyName: DS.attr('string'),
    schema: DS.attr(),
    name: DS.attr('string'),

    needsConfiguration: Ember.computed.bool('schema.length'),

    isValid: function(config) {
        var self = this;
        return this.get('schema').every(function(configItem) {
            return !Ember.isBlank(config) && !Ember.isBlank(config.get(configItem.name));
        });
    }
});

export default Middleware;
