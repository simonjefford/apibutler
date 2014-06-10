var Middleware = DS.Model.extend({
    friendlyName: DS.attr('string'),
    schema: DS.attr(),
    name: DS.attr('string'),

    needsConfiguration: Ember.computed.bool('schema.length')
});

export default Middleware;
